package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initialise la base de données SQLite
func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		return err
	}

	// Tester la connexion
	err = DB.Ping()
	if err != nil {
		return err
	}

	// Créer la table users
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("✅ Base de données SQLite initialisée (forum.db)")
	return nil
}

// CreateUser crée un nouvel utilisateur
func CreateUser(username, email, password string) (*User, error) {
	result, err := DB.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		username, email, password,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       int(id),
		Username: username,
		Email:    email,
		Password: password,
	}, nil
}

// UserExists vérifie si un utilisateur existe
func UserExists(username, email string) (bool, error) {
	var count int
	err := DB.QueryRow(
		"SELECT COUNT(*) FROM users WHERE username = ? OR email = ?",
		username, email,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserByUsername récupère un utilisateur par son nom
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, username, email, password FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail récupère un utilisateur par email
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, username, email, password FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}
	return user, nil
}
