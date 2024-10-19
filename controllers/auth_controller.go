package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

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

func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Limit the size of the request body
    r.ParseMultipartForm(10 << 20) // 10MB max file size

    // Parse the multipart form
    fullName := r.FormValue("full_name")
    email := r.FormValue("email")
    password := r.FormValue("password")


	if fullName == "" {
        respondWithError(w, http.StatusBadRequest, "Full name is mandatory")
        return
    }
    if email == "" {
        respondWithError(w, http.StatusBadRequest, "Email is mandatory")
        return
    }
    if password == "" {
        respondWithError(w, http.StatusBadRequest, "Password is mandatory")
        return
    }
    // Get the file from the form input
    file, handler, err := r.FormFile("profile_pic")
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Profile picture is mandatory")
        return
    }
    defer file.Close()

    // Create unique file name with timestamp
    fileName := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), handler.Filename)

    // Save the file to the server
    outFile, err := os.Create(fileName)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to save profile picture")
        return
    }
    defer outFile.Close()

    // Copy the uploaded file's content to the new file
    _, err = io.Copy(outFile, file)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to store profile picture")
        return
    }

    // Hash the password
    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
        return
    }

    // Insert user into the database
    result, err := db.Exec("INSERT INTO users (full_name, email, password, profile_pic) VALUES (?, ?, ?, ?)",
        fullName, email, hashedPassword, fileName)
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

    // Generate a token for the user
    token, err := utils.GenerateToken(int(userId), email)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to generate token, please try again")
        return
    }

    // Build success response
    successResponse := models.Response{
        Status:  true,
        Message: "User registration successful",
        Data: map[string]interface{}{
            "id":          userId,
            "full_name":   fullName,
            "email":       email,
            "profile_pic": fileName,
        },
        Token: token,
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

    // Construct the full URL or file path for the profile picture
    profilePicURL := fmt.Sprintf("/uploads/%s", dbUser.ProfilePic) // Adjust if needed for full URL

    // Prepare login response with profile_pic URL
    loginResponse := models.LoginResponse{
        ID:         dbUser.ID,
        FullName:   dbUser.FullName,
        Email:      dbUser.Email,
        ProfilePic: profilePicURL,  // Return the profile picture URL or file path
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
