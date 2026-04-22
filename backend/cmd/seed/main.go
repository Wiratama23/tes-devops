package main

import (
	"context"
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

	// Create sample users
	user1, err := userRepo.Create(ctx, "alice", "alice@example.com")
	if err != nil {
		log.Printf("Error creating user1: %v", err)
	} else {
		log.Printf("Created user: %v", user1)
	}

	user2, err := userRepo.Create(ctx, "bob", "bob@example.com")
	if err != nil {
		log.Printf("Error creating user2: %v", err)
	} else {
		log.Printf("Created user: %v", user2)
	}

	// Create sample articles
	if user1 != nil {
		article1, err := articleRepo.Create(ctx, user1.UID, "Getting Started with Go", "Go is a modern programming language...")
		if err != nil {
			log.Printf("Error creating article1: %v", err)
		} else {
			log.Printf("Created article: %v", article1)
		}

		article2, err := articleRepo.Create(ctx, user1.UID, "Advanced Go Patterns", "In this post we'll explore some advanced patterns...")
		if err != nil {
			log.Printf("Error creating article2: %v", err)
		} else {
			log.Printf("Created article: %v", article2)
		}
	}

	if user2 != nil {
		article3, err := articleRepo.Create(ctx, user2.UID, "Web Development with Go", "Learn how to build web applications...")
		if err != nil {
			log.Printf("Error creating article3: %v", err)
		} else {
			log.Printf("Created article: %v", article3)
		}
	}

	// Create sample products
	if user1 != nil {
		product1, err := productRepo.Create(ctx, "SKU10001", "Coffee Beans", 100, "29.99", "10", user1.UID, "assets/coffee.jpg")
		if err != nil {
			log.Printf("Error creating product1: %v", err)
		} else {
			log.Printf("Created product: %v", product1)
		}

		product2, err := productRepo.Create(ctx, "SKU05001", "Programming Book", 50, "49.99", "05", user1.UID, "assets/book.jpg")
		if err != nil {
			log.Printf("Error creating product2: %v", err)
		} else {
			log.Printf("Created product: %v", product2)
		}
	}

	if user2 != nil {
		product3, err := productRepo.Create(ctx, "SKU10002", "Tea Leaves", 75, "19.99", "10", user2.UID, "assets/tea.jpg")
		if err != nil {
			log.Printf("Error creating product3: %v", err)
		} else {
			log.Printf("Created product: %v", product3)
		}
	}

	log.Println("Seeding completed successfully!")
}
