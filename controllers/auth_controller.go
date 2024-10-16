package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"blog_project.com/models"
	"blog_project.com/utils"
)

var db *sql.DB

func Initialize(database *sql.DB) {
	db = database
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Check the database status
	err := db.Ping()
	if err != nil {
		http.Error(w, "MySQL server is down or not accessible: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "MySQL server is up and running!")
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.FullName != "" {
		if user.Email != "" {
			if user.Password != "" {
				if user.ProfilePic != "" {
					result, err := db.Exec("INSERT INTO users (full_name, email, profile_pic,password) VALUES (?,?,?,?)",
						user.FullName, user.Email, user.ProfilePic, user.Password)

					if err != nil {
						fmt.Println("Error ", err)
						// http.Error(w, err.Error(), http.StatusInternalServerError)
						returnError(w, "Email ID already Exists")
						return
					}

					userId, err := result.LastInsertId() //get user id here..

					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					token, err :=utils.GenerateToken(user.ID, user.Email)
					if err != nil {
						http.Error(w, "Failed to generate token, please try again", http.StatusInternalServerError)
						return
					}
					user.ID = int(userId)
					successResponse := models.Response{
						Status:  true,
						Message: "user register sucessfully done",
						Data:    user,
						Token:   token,
					}
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(successResponse)
				} else {
					returnError(w, "Profile Pic is Mandentory")
				}
			} else {
				returnError(w, "Password is Mandentory")
			}
		} else {
			returnError(w, "Email is Mandentory")
		}
	} else {
		returnError(w, "Fullname is Mandentory")
	}

}



func LoginUser(w http.ResponseWriter, r *http.Request){
	fmt.Println("Login api called")
}


func returnError(w http.ResponseWriter, msg string) {
	errorResponse := models.Response{
		Status:  false,
		Message: msg,
		Data:    struct{}{},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest) // Optional: set HTTP status code to 500
	json.NewEncoder(w).Encode(errorResponse)

}
