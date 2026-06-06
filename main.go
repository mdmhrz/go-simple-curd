package main

import (
	"context"
	"fmt"
	"http-crud/db"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	var err error

	err = godotenv.Load()

	if err != nil {
		panic("env file not found")
	}

	db.ConnectDb()

	defer db.Db.Close(context.Background())

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("POST /create-user", createUserHandler)
	mux.HandleFunc("GET /users", getUsersHandler)
	mux.HandleFunc("GET /users/{id}", getUserByIdHandler)
	mux.HandleFunc("PUT /users/{id}", updateUserByIdHandler)
	mux.HandleFunc("DELETE /users/{id}", deleteUserByIdHandler)

	fmt.Println("Server is running at port 5000")
	err = http.ListenAndServe(":5000", mux)

	if err != nil {
		fmt.Println("Server error", err)
	}
}
