package utils

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
)

// loadEnvOnce ensures we only re-open and parse the .env file the first time
// any caller invokes GetEnv. The previous implementation re-loaded .env on
// every lookup, which was both inefficient under load and noisy in the logs.
var loadEnvOnce sync.Once

func GetAllowedOrigins() []string {
	// 1. Get the string from the environment
	rawOrigins := os.Getenv("ALLOWED_ORIGINS")

	// 2. Fallback for local development if the variable is missing
	if rawOrigins == "" {
		return []string{"http://localhost:4000"}
	}

	// 3. Split the string by commas
	origins := strings.Split(rawOrigins, ",")

	// 4. Clean up the strings (removes accidental spaces)
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return origins
}

func GetEnv(value string) string {
	loadEnvOnce.Do(func() {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Printf("No .env file found, using environment variables: %v", err)
		}
	})
	return os.Getenv(value)
}

func GetJWT(value string) *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(value), nil)
}
