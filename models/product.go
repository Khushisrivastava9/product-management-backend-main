package models

import (
	"database/sql"
	"fmt"
)

type Product struct {
	ID                    int      `json:"id"`
	UserID                int      `json:"user_id"`
	ProductName           string   `json:"product_name"`
	ProductDescription    string   `json:"product_description"`
	ProductImages         []string `json:"product_images"`
	ProductPrice          float64  `json:"product_price"`
	CompressedProductImages []string `json:"compressed_product_images"`
}

func (p *Product) Create(db *sql.DB) error {
	query := `INSERT INTO products (user_id, product_name, product_description, product_images, product_price, compressed_product_images) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := db.QueryRow(query, p.UserID, p.ProductName, p.ProductDescription, pq.Array(p.ProductImages), p.ProductPrice, pq.Array(p.CompressedProductImages)).Scan(&p.ID)
	if err != nil {
		return fmt.Errorf("could not create product: %v", err)
	}
	return nil
}

func (p *Product) GetByID(db *sql.DB, id int) error {
	query := `SELECT id, user_id, product_name, product_description, product_images, product_price, compressed_product_images 
			  FROM products WHERE id = $1`
	row := db.QueryRow(query, id)
	err := row.Scan(&p.ID, &p.UserID, &p.ProductName, &p.ProductDescription, pq.Array(&p.ProductImages), &p.ProductPrice, pq.Array(&p.CompressedProductImages))
	if err != nil {
		return fmt.Errorf("could not get product by id: %v", err)
	}
	return nil
}

func (p *Product) Update(db *sql.DB) error {
	query := `UPDATE products SET user_id = $1, product_name = $2, product_description = $3, product_images = $4, product_price = $5, compressed_product_images = $6 
			  WHERE id = $7`
	_, err := db.Exec(query, p.UserID, p.ProductName, p.ProductDescription, pq.Array(p.ProductImages), p.ProductPrice, pq.Array(p.CompressedProductImages), p.ID)
	if err != nil {
		return fmt.Errorf("could not update product: %v", err)
	}
	return nil
}

func (p *Product) Delete(db *sql.DB) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := db.Exec(query, p.ID)
	if err != nil {
		return fmt.Errorf("could not delete product: %v", err)
	}
	return nil
}

func GetAllProducts(db *sql.DB, userID int, minPrice, maxPrice float64, productName string) ([]Product, error) {
	query := `SELECT id, user_id, product_name, product_description, product_images, product_price, compressed_product_images 
			  FROM products WHERE user_id = $1`
	args := []interface{}{userID}

	if minPrice > 0 {
		query += " AND product_price >= $2"
		args = append(args, minPrice)
	}
	if maxPrice > 0 {
		query += " AND product_price <= $3"
		args = append(args, maxPrice)
	}
	if productName != "" {
		query += " AND product_name ILIKE $4"
		args = append(args, "%"+productName+"%")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not get all products: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.UserID, &p.ProductName, &p.ProductDescription, pq.Array(&p.ProductImages), &p.ProductPrice, pq.Array(&p.CompressedProductImages))
		if err != nil {
			return nil, fmt.Errorf("could not scan product: %v", err)
		}
		products = append(products, p)
	}

	return products, nil
}
