package routers

import (
	"net/http"

	"blog_project.com/controllers" // Replace with your actual package path
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// SetupRouter initializes and returns a configured router with both API and static file routes.
func SetupRouter() http.Handler {
	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// API Routes
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	apiRouter.HandleFunc("/profile", controllers.GetUserProfile).Methods("GET")
	apiRouter.HandleFunc("/add-story", controllers.AddStory).Methods("POST")
	apiRouter.HandleFunc("/get-story", controllers.GetStory).Methods("GET")

	// Static file handler for serving files from the "uploads" directory
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Adjust to your frontend's origin
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		Debug:            true, // Set to true only during development
	})

	// Wrap the router with the CORS handler
	return c.Handler(r)
}
