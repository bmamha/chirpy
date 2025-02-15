package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bmamha/chirpy/internal/auth"
	"github.com/bmamha/chirpy/internal/database"
	"github.com/google/uuid"
)

type UserCreationResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) UserCreationHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Error creating user"))
		return
	}
	jsonUser := UserCreationResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	jsonResponse, err := json.Marshal(jsonUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Unable to parse response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(jsonResponse))
}
