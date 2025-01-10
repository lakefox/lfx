package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// Init initializes the database and creates necessary tables.
func Init(file string) *sql.DB {
	var err error
	db, err = sql.Open("sqlite", file)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	createTables()
	return db
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		text TEXT NOT NULL,
		score INTEGER DEFAULT 0,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		username TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		joined DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (username) REFERENCES users(username)
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}

// GetDB returns the database instance.
func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call Init() first.")
	}
	return db
}
