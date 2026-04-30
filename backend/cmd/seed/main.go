package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"rwiratama.com/m/internal/database"
	"rwiratama.com/m/internal/repository"
)

func main() {
	ctx := context.Background()

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables: %v", err)
	}

	// Load configuration from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Initialize database connection
	pool, err := database.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Create repositories
	userRepo := repository.NewUserRepository(pool)
	articleRepo := repository.NewArticleRepository(pool)
	productRepo := repository.NewProductRepository(pool)

	// 0. Seed admin user. Password is taken from ADMIN_PASSWORD env (defaults
	// to "admin123" so docker compose works out of the box; rotate via env in
	// production).
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}
	adminHash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash admin password: %v", err)
	}
	if _, err := userRepo.CreateWithPassword(ctx, "admin", "admin@example.com", string(adminHash), true); err != nil {
		log.Printf("admin user may already exist: %v", err)
	} else {
		log.Println("Seeded admin user (username=admin)")
	}

	const seedCount = 5000

	log.Printf("Starting to seed %d records for users, articles, and products...", seedCount)

	for i := 1; i <= seedCount; i++ {
		// 1. Generate User
		username := fmt.Sprintf("user%d", i)
		email := fmt.Sprintf("user%d@example.com", i)

		user, err := userRepo.Create(ctx, username, email)
		if err != nil {
			log.Printf("Error creating user %d: %v", i, err)
			continue // Skip article/product creation for this iteration if user creation fails
		}

		// 2. Generate Article (Linked to the generated user)
		articleTitle := fmt.Sprintf("Generated Article %d", i)
		articleBody := fmt.Sprintf("This is the automated body content for article number %d.", i)

		_, err = articleRepo.Create(ctx, user.UID, articleTitle, articleBody)
		if err != nil {
			log.Printf("Error creating article %d: %v", i, err)
		}

		// 3. Generate Product (Linked to the generated user)
		sku := fmt.Sprintf("SKU%05d", i) // Formats as SKU00001, SKU00002, etc.
		productName := fmt.Sprintf("Test Product %d", i)
		qty := 100
		price := "29.99"
		categoryId := "10"
		// Match the default used by product_handler.go and the file actually
		// shipped in /assets so frontend resolves to /api/assets/default_image.jpg
		imagePath := "assets/default_image.jpg"

		_, err = productRepo.Create(ctx, sku, productName, qty, price, categoryId, user.UID, imagePath)
		if err != nil {
			log.Printf("Error creating product %d: %v", i, err)
		}

		// Log progress every 1,000 iterations to keep the console clean
		if i%1000 == 0 {
			log.Printf("Successfully seeded %d/%d records...", i, seedCount)
		}
	}

	// Refresh the planner statistics that articles/products pagination uses
	// for an instant approximate count. Without this the very first
	// /articles page hit on a fresh DB would see reltuples = -1 and report a
	// total of 0 until autovacuum eventually runs.
	for _, table := range []string{"articles", "products", "users"} {
		if _, err := pool.Exec(ctx, "ANALYZE "+table); err != nil {
			log.Printf("ANALYZE %s failed: %v", table, err)
		}
	}

	log.Println("Seeding completed successfully!")
}
