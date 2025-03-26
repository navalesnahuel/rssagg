package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users  []User
	nextID int        = 1
	mu     sync.Mutex // Para evitar condiciones de carrera
)

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 201, users)
}

func postUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&newUser)
	if err != nil {
		respondWithError(w, 400, "error at payload")
		return
	}

	// Create the user
	mu.Lock()
	user := User{
		ID:   nextID,
		Name: newUser.Name,
	}
	users = append(users, user)
	nextID++
	mu.Unlock()

	// Respond with the created user
	respondWithJSON(w, http.StatusCreated, user)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from query parameter
	idStr := chi.URLParam(r, "id") // Get ID as a string
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Find and remove the user
	for i, user := range users {
		if user.ID == id {
			// Remove the user by slicing
			users = append(users[:i], users[i+1:]...)
			respondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted"})
			return
		}
	}

	// If user not found
	respondWithError(w, http.StatusNotFound, "User not found")
}

func putUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from query parameter
	idStr := chi.URLParam(r, "id") // Get ID as a string
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Decode the new user data
	var updatedUser User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedUser); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Find and update the user
	for i, user := range users {
		if user.ID == id {
			// Update the user's name
			users[i].Name = updatedUser.Name
			respondWithJSON(w, http.StatusOK, users[i])
			return
		}
	}

	// If user not found
	respondWithError(w, http.StatusNotFound, "User not found")
}
