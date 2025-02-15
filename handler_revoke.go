package main

import (
	"log"
	"net/http"

	"github.com/bmamha/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\": Error obtaining refresh token"))
		return
	}

	w.WriteHeader(204)
}
