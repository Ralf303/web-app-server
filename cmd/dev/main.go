package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/myapp/internal/database"
	"example.com/myapp/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.ConnectFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := server.Routes(db)

	fmt.Println("Server running on :8080 at", time.Now().Format(time.RFC3339))
	log.Fatal(http.ListenAndServe(":8080", router))
}
