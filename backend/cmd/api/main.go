package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"rwiratama.com/m/internal/database"
	"rwiratama.com/m/internal/handlers"
	czm "rwiratama.com/m/internal/middleware"
	"rwiratama.com/m/internal/utils"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	// Load .env file from project root
	// if err := godotenv.Load("../../../.env"); err != nil {
	// 	log.Printf("No .env file found, using environment variables: %v", err)
	// }

	secret := utils.GetEnv("JWT_SECRET") //os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	tokenAuth = utils.GetJWT(secret) //jwtauth.New("HS256", []byte(secret), nil)
}

func main() {
	ctx := context.Background()

	// Load configuration from environment
	databaseURL := utils.GetEnv("DATABASE_URL") //os.Getenv("DATABASE_URL")
	port := utils.GetEnv("PORT")                //os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Resolve uploads dir & default-image source. Both are optional in
	// development; if UPLOADS_DIR is empty the upload routes will return 500
	// when hit, which is the desired behaviour for misconfigured deployments.
	uploadsDir := utils.GetEnv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "/var/lib/api/uploads"
	}
	defaultImageSrc := utils.GetEnv("DEFAULT_IMAGE_SOURCE")
	if defaultImageSrc == "" {
		// Resolved relative to the working directory of the running binary.
		defaultImageSrc = "assets/default_image.jpg"
	}
	if err := handlers.EnsureDefaultImage(uploadsDir, defaultImageSrc); err != nil {
		log.Printf("warning: failed to seed default image into uploads dir: %v", err)
	}

	// Cookie + secure flag for the auth handler.
	cookieSecure, _ := strconv.ParseBool(utils.GetEnv("AUTH_COOKIE_SECURE"))

	// Initialize database connection
	pool, err := database.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Run migrations
	if err := runMigrations(databaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize and create database indexes
	indexManager := database.NewIndexManager(pool)
	if err := indexManager.CreateProductIndexes(ctx); err != nil {
		log.Fatalf("Failed to create product indexes: %v", err)
	}
	log.Println("Database indexes initialized successfully")

	// Initialize handlers
	userHandler := handlers.NewUserHandler(pool)
	articleHandler := handlers.NewArticleHandler(pool)
	productHandler := handlers.NewProductHandler(pool)
	authHandler := handlers.NewAuthHandler(pool, tokenAuth, handlers.AuthHandlerConfig{
		TokenTTL:   24 * time.Hour,
		CookieName: czm.AuthCookieName,
		Secure:     cookieSecure,
	})
	uploadHandler := handlers.NewUploadHandler(uploadsDir)
	logHandler := handlers.NewLogHandler()

	// Initialize Coraza WAF
	waf, err := czm.InitializeWAF()
	if err != nil {
		log.Fatalf("Failed to initialize WAF: %v", err)
	}

	// Setup routes with chi router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(czm.WAFMiddleware(waf))

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: utils.GetAllowedOrigins(), // Use this to allow specific origin hosts from environment variable
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Route("/api", func(r chi.Router) {

		// /api/auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/logout", authHandler.Logout)
			r.Group(func(r chi.Router) {
				r.Use(czm.JWTVerifier(tokenAuth))
				r.Use(jwtauth.Authenticator(tokenAuth))
				r.Get("/me", authHandler.Me)
			})
		})

		// /api/logs - frontend error reporting (open; rate-limited at nginx).
		r.Post("/logs", logHandler.Receive)

		// /api/assets/{filename} - serves uploaded images (public).
		r.Get("/assets/{filename}", uploadHandler.ServeAsset)

		// /api/uploads - admin-only image upload.
		r.Route("/uploads", func(r chi.Router) {
			r.Use(czm.JWTVerifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Use(czm.RequireAdmin)
			r.Post("/images", uploadHandler.UploadImage)
		})

		// /api/users
		r.Route("/users", func(r chi.Router) {
			// Public
			r.Group(func(r chi.Router) {
				r.Post("/", userHandler.CreateUser)
				r.Get("/", userHandler.GetAllUsers)
				r.Get("/{uid}", userHandler.GetUser)
				r.Get("/{uid}/articles", articleHandler.GetUserArticles)
			})
			// Protected
			r.Group(func(r chi.Router) {
				r.Use(czm.JWTVerifier(tokenAuth))
				r.Use(jwtauth.Authenticator(tokenAuth))
				r.Put("/{uid}", userHandler.UpdateUser)
				r.Delete("/{uid}", userHandler.DeleteUser)
			})
		})

		// /api/articles
		r.Route("/articles", func(r chi.Router) {
			// Public
			r.Group(func(r chi.Router) {
				r.With(czm.Paginate).Get("/", articleHandler.GetAllArticles)
				r.Get("/{id}", articleHandler.GetArticle)
			})
			// Protected
			r.Group(func(r chi.Router) {
				r.Use(czm.JWTVerifier(tokenAuth))
				r.Use(jwtauth.Authenticator(tokenAuth))
				r.Post("/", articleHandler.CreateArticle)
				r.Put("/{id}", articleHandler.UpdateArticle)
				r.Delete("/{id}", articleHandler.DeleteArticle)
			})
		})

		// /api/products
		r.Route("/products", func(r chi.Router) {
			// Public
			r.Group(func(r chi.Router) {
				r.With(czm.Paginate).Get("/", productHandler.GetAllProducts)
				r.Get("/{id}", productHandler.GetProductByID)
			})
			// Protected
			r.Group(func(r chi.Router) {
				r.Use(czm.JWTVerifier(tokenAuth))
				r.Use(jwtauth.Authenticator(tokenAuth))
				r.Post("/", productHandler.CreateProduct)
				r.Put("/{id}", productHandler.UpdateProduct)
				r.Delete("/{id}", productHandler.DeleteProduct)
			})
		})
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("health: failed to write response: %v", err)
		}
	})

	// Start server
	fmt.Printf("Server listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func runMigrations(databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("failed to close migration db: %v", cerr)
		}
	}()

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
