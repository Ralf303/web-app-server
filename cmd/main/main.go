package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/myapp/internal/database"
	"example.com/myapp/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.prod")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	router := server.Routes(db)

	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		log.Fatalf("Certificate file not found: %v", err)
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		log.Fatalf("Key file not found: %v", err)
	}

	fmt.Println("Server is running on port 8080 at", time.Now())
	log.Fatal(http.ListenAndServeTLS(":8080", certPath, keyPath, router))
}
