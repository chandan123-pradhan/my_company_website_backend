package routers

import (
	"log"

	"blog_project.com/controllers"
	"github.com/gorilla/mux"

	"net/http"
)

func NewRouter() *mux.Router{
	r := mux.NewRouter()
	r.HandleFunc("/register",controllers.CreateUser).Methods("POST")
	r.HandleFunc("/login",controllers.LoginUser).Methods("POST")
	r.HandleFunc("/profile",controllers.GetUserProfile).Methods("GET")
	log.Println("Server starting on port :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err) // Log any errors that occur while starting the server
	}
	return r;
}