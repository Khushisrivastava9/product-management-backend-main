package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize logger
	logger := logrus.New()
	logger.Out = os.Stdout

	// Set up router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/products", createProduct).Methods("POST")
	router.HandleFunc("/products/{id}", getProductByID).Methods("GET")
	router.HandleFunc("/products", getAllProducts).Methods("GET")

	// Middleware for logging
	router.Use(loggingMiddleware(logger))

	// Start server
	http.ListenAndServe(":8080", router)
}

func loggingMiddleware(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
			}).Info("Request received")
			next.ServeHTTP(w, r)
		})
	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Save the product to the database (placeholder logic)
	// In a real application, you would save the product to a database
	product.ID = "12345" // Mock ID

	// Respond with the created product
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	vars := mux.Vars(r)
	productID := vars["id"]

	// Retrieve the product from the database (placeholder logic)
	// In a real application, you would retrieve the product from a database
	product := Product{
		ID:   productID,
		Name: "Sample Product",
	}

	// Respond with the product
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	// Retrieve all products from the database (placeholder logic)
	// In a real application, you would retrieve the products from a database
	products := []Product{
		{ID: "12345", Name: "Sample Product 1"},
		{ID: "67890", Name: "Sample Product 2"},
	}

	// Respond with the list of products
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Product represents a product in the system
type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
