package main

import (
	"fmt"
	"net/http"

	"github.com/bmamha/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) ChirpDeletionHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\": Error obtaining token"))
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.SECRET)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\": Unauthorized"))
		return
	}

	chirp_id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("\"error\": Failed to parse id given"))
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirp_id)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte("\"error\": Error getting chirp"))
		return
	}

	if chirp.UserID != userId {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(403)
		w.Write([]byte("\"error\": User is not the author of the chirp"))
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirp_id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte("\"error\": Error deleting chirp"))
		return
	}

	w.WriteHeader(204)
}
