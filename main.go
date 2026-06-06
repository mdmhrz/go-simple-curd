package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

var users = []User{
	{
		Id:    1,
		Name:  "Kalam",
		Age:   25,
		Email: "kalam@example.com",
	},
	{
		Id:    2,
		Name:  "Salam",
		Age:   29,
		Email: "salam@example.com",
	},
	{
		Id:    3,
		Name:  "Alam",
		Age:   20,
		Email: "alam@example.com",
	},
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("POST /create-user", createUserHandler)
	mux.HandleFunc("GET /users", getUsers)
	mux.HandleFunc("GET /users/{id}", getUserById)
	mux.HandleFunc("PUT /users/{id}", updateUserById)
	mux.HandleFunc("DELETE /users/{id}", deleteUserById)

	fmt.Println("Server is running at port 5000")
	err := http.ListenAndServe(":5000", mux)

	if err != nil {
		fmt.Println("Server error", err)
	}
}

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

	newUser.Id = len(users) + 1
	users = append(users, newUser)

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	// users, _ := json.Marshal(users)
	// w.Write(users)

	encoder := json.NewEncoder(w)
	encoder.Encode(users)
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	idParams := r.PathValue("id")

	id, err := strconv.Atoi(idParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	for _, user := range users {
		if user.Id == id {
			w.Header().Set("Content-Type", "Application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "User not found")

}

func updateUserById(w http.ResponseWriter, r *http.Request) {
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

	for idx, user := range users {
		if user.Id == id {

			updatedUser.Id = id
			users[idx] = updatedUser

			w.Header().Set("Content-Type", "Application/json")
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "User not found")
}

func deleteUserById(w http.ResponseWriter, r *http.Request) {
	idParams := r.PathValue("id")

	id, error := strconv.Atoi(idParams)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	for idx, user := range users {
		if user.Id == id {
			// users = append(users[:idx], users[idx+1:]...)
			users = slices.Delete(users, idx, idx+1)
			w.WriteHeader(http.StatusNoContent)
			fmt.Fprintln(w, "User information deleted")
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "User not found")

}
