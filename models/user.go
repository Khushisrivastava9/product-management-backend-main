package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) CreateUser(db *sql.DB) error {
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
	err := db.QueryRow(query, u.Username, u.Email, u.Password).Scan(&u.ID)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func (u *User) GetUserByID(db *sql.DB, id int) error {
	query := "SELECT id, username, email, password FROM users WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return fmt.Errorf("error getting user by id: %v", err)
	}
	return nil
}

func (u *User) UpdateUser(db *sql.DB) error {
	query := "UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4"
	_, err := db.Exec(query, u.Username, u.Email, u.Password, u.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}

func (u *User) DeleteUser(db *sql.DB) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := db.Exec(query, u.ID)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}
