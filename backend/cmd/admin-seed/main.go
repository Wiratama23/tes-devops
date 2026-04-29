package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"rwiratama.com/m/internal/database"
	"rwiratama.com/m/internal/repository"
)

func main() {
	ctx := context.Background()

	// Define command-line flags
	username := flag.String("username", "", "Admin username")
	email := flag.String("email", "", "Admin email address")
	password := flag.String("password", "", "Admin password (will be prompted if not provided)")
	flag.Parse()

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

	// Create user repository
	userRepo := repository.NewUserRepository(pool)

	// Get username
	if *username == "" {
		fmt.Print("Enter admin username: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		*username = strings.TrimSpace(input)
		if *username == "" {
			log.Fatal("Username cannot be empty")
		}
	}

	// Get email
	if *email == "" {
		fmt.Print("Enter admin email: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		*email = strings.TrimSpace(input)
		if *email == "" {
			log.Fatal("Email cannot be empty")
		}
	}

	// Get password
	if *password == "" {
		fmt.Print("Enter admin password: ")
		reader := bufio.NewReader(os.Stdin)
		passwordInput, _ := reader.ReadString('\n')
		*password = strings.TrimSpace(passwordInput)

		if *password == "" {
			log.Fatal("Password cannot be empty")
		}

		// Confirm password
		fmt.Print("Confirm admin password: ")
		confirmInput, _ := reader.ReadString('\n')
		confirm := strings.TrimSpace(confirmInput)

		if *password != confirm {
			log.Fatal("Passwords do not match")
		}
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create admin user
	user, err := userRepo.CreateWithPassword(ctx, *username, *email, string(passwordHash), true)
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Printf("\n✅ Admin user created successfully!\n")
	fmt.Printf("   UID:      %s\n", user.UID)
	fmt.Printf("   Username: %s\n", user.Username)
	fmt.Printf("   Email:    %s\n", user.Email)
	fmt.Printf("   Is Admin: %v\n", user.IsAdmin)
}
