package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"blog_project.com/models"
	"blog_project.com/utils"
)

var db *sql.DB

// Initialize sets the database connection
//
// This function should be called at the start of the application
// to set up the database connection that will be used by the
// controllers for user management.
func Initialize(database *sql.DB) {
	db = database
}

// CreateUser handles user registration.
//
// It decodes the incoming request body to extract user information,
// validates the input, hashes the password, inserts the user
// into the database, generates a token, and returns a JSON response
// with user details and a success message.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterUserModel

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input fields
	errors := utils.ValidateUserInput(user, true) // true for registration
	if len(errors) > 0 {
		respondWithError(w, http.StatusBadRequest, utils.ErrorMessages(errors))
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Insert user into the database
	result, err := db.Exec("INSERT INTO users (full_name, email, profile_pic, password) VALUES (?, ?, ?, ?)",
		user.FullName, user.Email, user.ProfilePic, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusConflict, "Email ID already exists")
		return
	}

	// Retrieve the new user ID
	userId, err := result.LastInsertId()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user ID")
		return
	}
	user.ID = int(userId)

	// Generate a token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token, please try again")
		return
	}

	// Build success response
	successResponse := models.Response{
		Status:  true,
		Message: "User registration successful",
		Data:    user,
		Token:   token,
	}

	respondWithJSON(w, http.StatusOK, successResponse)
}

// LoginUser handles user login.
//
// It decodes the incoming request body to extract login credentials,
// validates the input, retrieves the user from the database,
// checks the password, generates a token, and returns a JSON
// response with user details and a success message.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.LoginUserModel
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input fields
	errors := utils.ValidateUserInput(user, false) // false for login
	if len(errors) > 0 {
		respondWithError(w, http.StatusBadRequest, utils.ErrorMessages(errors))
		return
	}

	// Retrieve user from database
	var dbUser models.RegisterUserModel
	err := db.QueryRow("SELECT id, full_name, email, profile_pic, password FROM users WHERE email = ?", user.Email).Scan(
		&dbUser.ID, &dbUser.FullName, &dbUser.Email, &dbUser.ProfilePic, &dbUser.Password,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Check password
	if err := utils.CheckPasswordHash(user.Password, dbUser.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(dbUser.ID, dbUser.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token, please try again")
		return
	}

	// Prepare login response
	loginResponse := models.LoginResponse{
		ID:         dbUser.ID,
		FullName:   dbUser.FullName,
		Email:      dbUser.Email,
		ProfilePic: dbUser.ProfilePic,
	}

	// Send success response
	successResponse := models.Response{
		Status:  true,
		Message: "Login successful",
		Data:    loginResponse,
		Token:   token,
	}
	respondWithJSON(w, http.StatusOK, successResponse)
}


// GetUserProfile retrieves the user's profile data based on the user ID
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	var userID int
	var err error

	// Extract user ID from the token
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " from tokenString and trim spaces
	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	if tokenString == "" {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}
	if tokenString != "" {
		userID, err =utils.ParseToken(tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}
	} else {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization token")
		return
	}

	// Now you can retrieve the user data based on userID
	var user models.GetUserProfileModel
	err = db.QueryRow("SELECT full_name, email, profile_pic FROM users WHERE id = ?", userID).Scan(
		&user.FullName, &user.Email, &user.ProfilePic,
	)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Prepare the response
	successResponse := models.Response{
		Status:  true,
		Message: "User profile retrieved successfully",
		Data:    user,
	}

	respondWithJSON(w, http.StatusOK, successResponse)
}

// respondWithJSON sends a JSON response.
//
// This utility function sets the Content-Type header to
// "application/json", writes the HTTP status code, and
// encodes the payload into a JSON response body.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// respondWithError sends an error response in JSON format.
//
// This utility function uses the standard response structure
// to send error messages, making it easy to maintain consistency
// across error responses.
func respondWithError(w http.ResponseWriter, code int, message string) {
	errorResponse := models.Response{
		Status:  false,
		Message: message,
		Data:    struct{}{},
	}
	respondWithJSON(w, code, errorResponse)
}
