package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
		imagePath := "assets/default.jpg"

		_, err = productRepo.Create(ctx, sku, productName, qty, price, categoryId, user.UID, imagePath)
		if err != nil {
			log.Printf("Error creating product %d: %v", i, err)
		}

		// Log progress every 1,000 iterations to keep the console clean
		if i%1000 == 0 {
			log.Printf("Successfully seeded %d/%d records...", i, seedCount)
		}
	}

	log.Println("Seeding completed successfully!")
}