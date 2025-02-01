package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bmamha/chirpy/internal/auth"
	"github.com/bmamha/chirpy/internal/database"
)

func (cfg *apiConfig) UserCreationHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json "password"`
		Email    string `json "email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Error decoding parameters"))
		return
	}
	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Unable to hash password"))
		return
	}
	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: passwordHash,
	}
	user, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Error creating user"))
		return
	}

	jsonResponse := fmt.Sprintf("{\"id\": \"%s\", \"created_at\": \"%s\",\"updated_at\": \"%s\",\"email\": \"%s\"}",
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(jsonResponse))
}
