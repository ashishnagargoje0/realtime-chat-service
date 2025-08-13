package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Postgres driver
)

// DB is a global database connection
var DB *sql.DB

// InitDB initializes the Postgres connection
func InitDB() {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "password")
	dbname := getEnv("POSTGRES_DB", "chatdb")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")

	// Create table if not exists
	createTable()
}

// createTable creates the messages table
func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create messages table: %v", err)
	}
}

// SaveMessage inserts a message into the database
func SaveMessage(username, content string) error {
	query := `INSERT INTO messages (username, content) VALUES ($1, $2)`
	_, err := DB.Exec(query, username, content)
	return err
}

// GetLastMessages fetches the last N messages
func GetLastMessages(limit int) ([]Message, error) {
	query := `SELECT id, username, content, created_at FROM messages ORDER BY id DESC LIMIT $1`
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		err := rows.Scan(&m.ID, &m.Username, &m.Content, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append([]Message{m}, messages...) // reverse order
	}
	return messages, nil
}

// Message struct
type Message struct {
	ID        int
	Username  string
	Content   string
	CreatedAt string
}

// helper to read env variables
func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
