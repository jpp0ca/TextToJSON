package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	httpAdapter "GoApiProject/internal/adapters/input/http"
	"GoApiProject/internal/adapters/output/gemini"
	"GoApiProject/internal/application"

	_ "GoApiProject/docs"

	"github.com/joho/godotenv"
)

//	@title			GoApiProject
//	@version		1.0
//	@description	A text-to-JSON structuring API powered by Google Gemini.

//	@host		localhost:8080
//	@BasePath	/

func main() {
	// --- Load .env ---
	// Try current directory first, then project root (for when IDE sets CWD to cmd/api/)
	if err := godotenv.Load(); err != nil {
		if err2 := godotenv.Load("../../.env"); err2 != nil {
			log.Println("No .env file found, falling back to environment variables")
		}
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is required. Set it in .env or as an environment variable.")
	}

	// --- Driven adapters (output) ---
	llmClient := gemini.NewClient(apiKey)

	// --- Application services ---
	structurerService := application.NewStructurerService(llmClient)

	// --- Driving adapters (input) ---
	handler := httpAdapter.NewStructureHandler(structurerService)
	router := httpAdapter.NewRouter(handler)

	// --- Start HTTP server ---
	addr := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", addr)
	fmt.Printf("Swagger UI:  http://localhost%s/swagger/index.html\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
