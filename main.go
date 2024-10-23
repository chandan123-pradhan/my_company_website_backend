package main

import (
	"fmt"
	"log"
	"net/http"

	"blog_project.com/controllers"
	"blog_project.com/routers"
	"github.com/rs/cors"
)

func main() {
    fs := http.FileServer(http.Dir("./uploads"))
    http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))
    fmt.Println("Hello, world!")
    controllers.InitDB()
    defer controllers.DB.Close()
    controllers.Initialize(controllers.DB);
    

    r:= routers.NewRouter();
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"}, // Adjust this to your Flutter app's URL
        AllowCredentials: true,
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type"},
        // For debugging, log every CORS request
        Debug: true,
    })
    // Use CORS middleware
    handler := c.Handler(r)

    // Start the server
    http.ListenAndServe(":8080", handler)
    log.Println("starting server on : 8080")
    
}