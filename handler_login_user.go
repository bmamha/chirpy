package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bmamha/chirpy/internal/auth"

	"github.com/bmamha/chirpy/internal/database"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte("\"error\": \"error decoding parameters\""))
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\": \"error finding user\""))
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\":\"incorrect email or password\""))
		return
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.SECRET)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\":\"Unable to generate token for user\""))
		return

	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\":\"Unable to generate token for user\""))
		return
	}

	jsonUser := UserResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        tokenString,
		RefreshToken: refreshToken,
	}
	expireTime := time.Now().AddDate(0, 0, 60)
	Refresh, err := cfg.db.CreateRefreshTokens(r.Context(), database.CreateRefreshTokensParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expireTime,
	})
	if err != nil {
		fmt.Println(Refresh, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Unable to create refresh token"))
		return

	}

	jsonResponse, err := json.Marshal(jsonUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Unable to parse response"))
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(jsonResponse))
}
