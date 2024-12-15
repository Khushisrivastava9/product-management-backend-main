package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/yourusername/yourproject/controllers"
	"github.com/yourusername/yourproject/models"
	"github.com/yourusername/yourproject/services"
)

func TestCreateProduct(t *testing.T) {
	// Initialize the necessary services
	services.InitLogger()
	services.InitCache("localhost", "6379")
	services.InitDB("user=youruser dbname=yourdb sslmode=disable")

	// Create a new product
	product := models.Product{
		UserID:             1,
		ProductName:        "Test Product",
		ProductDescription: "This is a test product",
		ProductImages:      []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
		ProductPrice:       19.99,
	}

	// Convert product to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to marshal product: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(productJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new HTTP recorder
	rr := httptest.NewRecorder()

	// Create a new router and register the handler
	router := mux.NewRouter()
	router.HandleFunc("/products", controllers.CreateProduct).Methods("POST")

	// Serve the HTTP request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	var createdProduct models.Product
	err = json.NewDecoder(rr.Body).Decode(&createdProduct)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if createdProduct.ProductName != product.ProductName {
		t.Errorf("Handler returned unexpected product name: got %v want %v", createdProduct.ProductName, product.ProductName)
	}
}

func TestGetProductByID(t *testing.T) {
	// Initialize the necessary services
	services.InitLogger()
	services.InitCache("localhost", "6379")
	services.InitDB("user=youruser dbname=yourdb sslmode=disable")

	// Create a new product
	product := models.Product{
		UserID:             1,
		ProductName:        "Test Product",
		ProductDescription: "This is a test product",
		ProductImages:      []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
		ProductPrice:       19.99,
	}

	// Save the product to the database
	err := product.Create(services.DB)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/products/"+strconv.Itoa(product.ID), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new HTTP recorder
	rr := httptest.NewRecorder()

	// Create a new router and register the handler
	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", controllers.GetProductByID).Methods("GET")

	// Serve the HTTP request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var fetchedProduct models.Product
	err = json.NewDecoder(rr.Body).Decode(&fetchedProduct)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if fetchedProduct.ProductName != product.ProductName {
		t.Errorf("Handler returned unexpected product name: got %v want %v", fetchedProduct.ProductName, product.ProductName)
	}
}

func TestGetAllProducts(t *testing.T) {
	// Initialize the necessary services
	services.InitLogger()
	services.InitCache("localhost", "6379")
	services.InitDB("user=youruser dbname=yourdb sslmode=disable")

	// Create a new product
	product1 := models.Product{
		UserID:             1,
		ProductName:        "Test Product 1",
		ProductDescription: "This is a test product 1",
		ProductImages:      []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
		ProductPrice:       19.99,
	}

	product2 := models.Product{
		UserID:             1,
		ProductName:        "Test Product 2",
		ProductDescription: "This is a test product 2",
		ProductImages:      []string{"http://example.com/image3.jpg", "http://example.com/image4.jpg"},
		ProductPrice:       29.99,
	}

	// Save the products to the database
	err := product1.Create(services.DB)
	if err != nil {
		t.Fatalf("Failed to create product 1: %v", err)
	}

	err = product2.Create(services.DB)
	if err != nil {
		t.Fatalf("Failed to create product 2: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/products?user_id=1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new HTTP recorder
	rr := httptest.NewRecorder()

	// Create a new router and register the handler
	router := mux.NewRouter()
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")

	// Serve the HTTP request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var products []models.Product
	err = json.NewDecoder(rr.Body).Decode(&products)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(products) != 2 {
		t.Errorf("Handler returned unexpected number of products: got %v want %v", len(products), 2)
	}
}

func BenchmarkGetProductByID(b *testing.B) {
	// Initialize the necessary services
	services.InitLogger()
	services.InitCache("localhost", "6379")
	services.InitDB("user=youruser dbname=yourdb sslmode=disable")

	// Create a new product
	product := models.Product{
		UserID:             1,
		ProductName:        "Test Product",
		ProductDescription: "This is a test product",
		ProductImages:      []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
		ProductPrice:       19.99,
	}

	// Save the product to the database
	err := product.Create(services.DB)
	if err != nil {
		b.Fatalf("Failed to create product: %v", err)
	}

	// Create a new router and register the handler
	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", controllers.GetProductByID).Methods("GET")

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a new HTTP request
		req, err := http.NewRequest("GET", "/products/"+strconv.Itoa(product.ID), nil)
		if err != nil {
			b.Fatalf("Failed to create request: %v", err)
		}

		// Create a new HTTP recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		router.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusOK {
			b.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}
}
