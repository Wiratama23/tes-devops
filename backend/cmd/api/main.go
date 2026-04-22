package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"rwiratama.com/m/internal/database"
	"rwiratama.com/m/internal/handlers"
	czm "rwiratama.com/m/internal/middleware"
)

func main() {
	ctx := context.Background()

	// Load .env file from project root
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Printf("No .env file found, using environment variables: %v", err)
	}

	// Load configuration from environment
	databaseURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

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

	// User routes
	router.Route("/api/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.GetAllUsers)
		r.Get("/{uid}", userHandler.GetUser)
		r.Put("/{uid}", userHandler.UpdateUser)
		r.Delete("/{uid}", userHandler.DeleteUser)

		r.Route("/{uid}/articles", func(r chi.Router) {
			r.Get("/", articleHandler.GetUserArticles)
		})
	})

	// Article routes
	router.Route("/api/articles", func(r chi.Router) {
		r.Post("/", articleHandler.CreateArticle)
		r.Get("/", articleHandler.GetAllArticles)
		r.Get("/{id}", articleHandler.GetArticle)
		r.Put("/{id}", articleHandler.UpdateArticle)
		r.Delete("/{id}", articleHandler.DeleteArticle)
	})

	// Product routes
	router.Route("/api/products", func(r chi.Router) {
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetAllProducts)
		r.Get("/{id}", productHandler.GetProductByID)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

	// Wrap router with Coraza WAF
	wrappedRouter := czm.WrapHandlerWithWAF(waf, router)

	// Start server
	fmt.Printf("Server listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, wrappedRouter); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func runMigrations(databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
