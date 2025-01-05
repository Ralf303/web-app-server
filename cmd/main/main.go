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
	err := godotenv.Load(".env.prod")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	router := server.Routes(db)

	fmt.Println("Server is running on port 8000 at", time.Now())

	log.Fatal(http.ListenAndServe(":8000", router))
}
