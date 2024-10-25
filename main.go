package main

import (
	"log"
	"net/http"

	"blog_project.com/controllers"
	"blog_project.com/routers"
)

func main() {
	// Initialize the database
	controllers.InitDB()
	defer controllers.DB.Close()
	controllers.Initialize(controllers.DB)

   

	// Get the configured router with API and static file handling
	handler := routers.SetupRouter()

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
