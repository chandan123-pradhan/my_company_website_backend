package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"blog_project.com/models"
	"blog_project.com/utils"
)

// AddStory handles adding a single story for a user.
func AddStory(w http.ResponseWriter, r *http.Request) {
	var userID int
	var err error

	// Extract and validate token from the header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization token")
		return
	}

	// Remove "Bearer " prefix and trim spaces
	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	if tokenString == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid token format")
		return
	}

	// Parse the user ID from the token
	userID, err = utils.ParseToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	// Parse the JSON story from the request body
	var req models.AddStoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Marshal the JSON story to store it as a JSON column
	storyJSON, err := json.Marshal(req.Story)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to encode story")
		return
	}

	// Insert the new story into the database with userID
	_, err = db.Exec("INSERT INTO usersStory (stories, userId) VALUES (?, ?)", storyJSON, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to store story")
		return
	}

	// Send success response
	successResponse := models.UserStoryAddSuccessModel{
		Status:  true,
		Message: "Story added successfully",
	}
	respondWithJSON(w, http.StatusOK, successResponse)
}


func GetStory(w http.ResponseWriter, r *http.Request){
	var userID int
	var err error

	// Extract and validate token from the header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization token")
		return
	}
	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	if tokenString == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid token format")
		return
	}

	// Parse the user ID from the token
	userID, err = utils.ParseToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	// Retrieve all stories for the given user ID
	rows, err := db.Query("SELECT stories FROM usersStory WHERE userId = ?", userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve stories")
		return
	}
	defer rows.Close()
	// Collect all stories in a slice
	var stories []map[string]interface{}
	for rows.Next() {
		var storyData sql.NullString
		if err := rows.Scan(&storyData); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to read story")
			return
		}

		// Unmarshal JSON data if valid
		if storyData.Valid {
			var story map[string]interface{}
			if err := json.Unmarshal([]byte(storyData.String), &story); err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to parse story data")
				return
			}
			stories = append(stories, story)
		}
	}

	// Send the response with all stories for the user
	successResponse := models.Response{
		Status:  true,
		Message: "Stories retrieved successfully",
		Data:    stories,
	}
	respondWithJSON(w, http.StatusOK, successResponse)

}