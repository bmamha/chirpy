package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bmamha/chirpy/internal/auth"
	"github.com/bmamha/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) userUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\": Error retreiving token"))
		return

	}
	userid, err := auth.ValidateJWT(tokenString, cfg.SECRET)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\": Unauthorized"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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

	updateParameters := database.UpdateUserEmailAndPasswordParams{
		Email:          params.Email,
		HashedPassword: passwordHash,
		UpdatedAt:      time.Now(),
		ID:             userid,
	}

	user, err := cfg.db.UpdateUserEmailAndPassword(r.Context(), updateParameters)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte("\"error\": Error updating parameter in database"))
		return
	}

	type userJson struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	jsonUser := userJson{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	jsonResponse, err := json.Marshal(jsonUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte("\"error\": Error parsing Json Response"))
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonResponse)
}
