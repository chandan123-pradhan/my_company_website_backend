package routers

import (
	"blog_project.com/controllers"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router{
	r := mux.NewRouter()
	r.HandleFunc("/register",controllers.CreateUser).Methods("POST")
	r.HandleFunc("/status", controllers.StatusHandler).Methods("GET")
	return r
}