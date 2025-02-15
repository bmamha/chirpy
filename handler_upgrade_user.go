package main

import (
	"encoding/json"
	"net/http"

	"github.com/bmamha/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) userUpgradeChirpyRed(w http.ResponseWriter, r *http.Request) {
	api_key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	if api_key != cfg.POLKA_KEY {
		w.WriteHeader(401)
		return
	}
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	user_id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	err = cfg.db.UpgradeUserChirpyRedById(r.Context(), user_id)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
}
