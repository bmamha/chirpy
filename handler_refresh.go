package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bmamha/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Fatalf("unable to extract token: %v", err)
	}

	token, err := cfg.db.GetRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\": Error obtaining refresh token"))
		return
	}

	if token.ExpiresAt.Before(time.Now()) || token.RevokedAt.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\": token has expired or is revoked"))
		return
	}
	accessToken, err := auth.MakeJWT(token.UserID, cfg.SECRET)
	if err != nil {
		log.Fatalf("unable to create access Token")
	}
	jsonResponse := fmt.Sprintf("{\"token\":\"%v\"}", accessToken)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(jsonResponse))
}
