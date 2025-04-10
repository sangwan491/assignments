package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sangwan491/backend-assignments/employee-management/backend/controllers"
	router "github.com/sangwan491/backend-assignments/employee-management/backend/routes"
)

func main() {
	// Connect to MongoDB first
	if err := controllers.ConnectToMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		os.Exit(1)
	}

	r := router.SetupRouter()
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
	// This line will never be executed due to log.Fatal above
	// fmt.Println("Server started on port 8080")
}
