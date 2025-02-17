package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {

	port := os.Getenv("PORT")

	router := http.NewServeMux()

	router.HandleFunc("GET /hello", hello)

	fmt.Printf("Starting server at port: %s\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	connectionDB()
	fmt.Fprintf(w, "HELLO!")
}

func connectionDB() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)

	query := `
	CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL);`

	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create table: %v\n", err)
		os.Exit(1)
	}

	// Execute Insert
	var insertedID int
	err = conn.QueryRow(context.Background(), `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id;`, "John Doe", "john@example.com").Scan(&insertedID)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}

	fmt.Println()
	fmt.Printf("Inserted user with ID: %d\n", insertedID)
}
