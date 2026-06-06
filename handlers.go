package main

import (
	"context"
	"encoding/json"
	"fmt"
	"http-crud/db"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to go server")

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is up and healthy")

}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User

	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid input")
		return
	}

	// newUser.Id = len(users) + 1
	// users = append(users, newUser)

	query := `
		insert into users (username, age, email)
		values ($1, $2, $3)
		returning id
	`

	err = db.Db.QueryRow(context.Background(), query, newUser.Name, newUser.Age, newUser.Email).Scan(&newUser.Id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not create user")
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)

}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {

	query := `select * from users`

	rows, err := db.Db.Query(context.Background(), query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not retrive users")
		return
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Email)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Could not map users")
			return
		}

		users = append(users, user)
	}

	err = rows.Err()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not read users")
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	// users, _ := json.Marshal(users)
	// w.Write(users)

	encoder := json.NewEncoder(w)
	encoder.Encode(users)
}

func getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var user User

	query := `
		SELECT id, username, age, email
		FROM users
		WHERE id = $1
	`

	err = db.Db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Age,
		&user.Email,
	)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	idParams := r.PathValue("id")

	id, error := strconv.Atoi(idParams)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	var updatedUser User

	err := json.NewDecoder(r.Body).Decode(&updatedUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid input")
		return
	}

	query := `
		update users
		set username = $1, age = $2, email = $3
		where id = $4
		returning id, username, age, email
	`

	err = db.Db.QueryRow(context.Background(), query, updatedUser.Name, updatedUser.Age, updatedUser.Email, id).Scan(&updatedUser.Id, &updatedUser.Name, &updatedUser.Age, &updatedUser.Email)

	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User not found")
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not update user")
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(updatedUser)

}

func deleteUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	idParams := r.PathValue("id")

	id, error := strconv.Atoi(idParams)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	query := `delete from users where id = $1`
	cmdTag, err := db.Db.Exec(context.Background(), query, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not delete user")
		return
	}

	if cmdTag.RowsAffected() != 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w, "User deleted successfully")

}
