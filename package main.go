package main

import (
	"fmt"
	"gosmart/pkg/db"
	"os"
)

func main() {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		fmt.Println("PG_DSN not set")
		return
	}
	conn, err := db.NewPostgres(dsn)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	if err := conn.Ping(); err != nil {
		fmt.Println("Ping error:", err)
		return
	}
	fmt.Println("Database connection successful!")
}
