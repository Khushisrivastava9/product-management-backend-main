package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/yourproject/models"
	"github.com/yourusername/yourproject/services"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = product.Create(services.DB)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	err = services.Queue.AddProductImages(product.ProductImages)
	if err != nil {
		http.Error(w, "Failed to add product images to queue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := services.Cache.GetProductByID(id)
	if err == nil {
		json.NewEncoder(w).Encode(product)
		return
	}

	var product models.Product
	err = product.GetByID(services.DB, id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	services.Cache.SetProductByID(id, product)

	json.NewEncoder(w).Encode(product)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	minPrice, _ := strconv.ParseFloat(r.URL.Query().Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(r.URL.Query().Get("max_price"), 64)
	productName := r.URL.Query().Get("product_name")

	products, err := models.GetAllProducts(services.DB, userID, minPrice, maxPrice, productName)
	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}
