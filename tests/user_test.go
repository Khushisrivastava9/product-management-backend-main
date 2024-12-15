package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/yourusername/yourproject/controllers"
	"github.com/yourusername/yourproject/models"
	"github.com/yourusername/yourproject/services"
)

func TestCreateUser(t *testing.T) {
	// Set up the database connection
	db, err := services.NewDB()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")

	// Create a new user
	user := models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password",
	}
	userJSON, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	var createdUser models.User
	err = json.NewDecoder(rr.Body).Decode(&createdUser)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if createdUser.Username != user.Username || createdUser.Email != user.Email {
		t.Errorf("handler returned unexpected body: got %v want %v", createdUser, user)
	}
}

func TestGetUserByID(t *testing.T) {
	// Set up the database connection
	db, err := services.NewDB()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", controllers.GetUserByID).Methods("GET")

	// Create a new user
	user := models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password",
	}
	err = user.CreateUser(db)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get the user by ID
	req, err := http.NewRequest("GET", "/users/"+strconv.Itoa(user.ID), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var retrievedUser models.User
	err = json.NewDecoder(rr.Body).Decode(&retrievedUser)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if retrievedUser.Username != user.Username || retrievedUser.Email != user.Email {
		t.Errorf("handler returned unexpected body: got %v want %v", retrievedUser, user)
	}
}
