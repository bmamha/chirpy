package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bmamha/chirpy/internal/auth"
	"github.com/bmamha/chirpy/internal/database"
)

func (cfg *apiConfig) ChirpCreationHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\": Error obtaining token"))
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.SECRET)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte("\"error\": Unauthorized"))
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Error decoding parameters"))
		return
	}

	if len(params.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("\"error\": Chirp is too long"))
		return
	}

	cleaned_body := badWordsHandler(params.Body)

	chirpParameters := database.CreateChirpsParams{
		Body:   cleaned_body,
		UserID: userid,
	}
	chirp, err := cfg.db.CreateChirps(r.Context(), chirpParameters)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Error creating user"))
		return
	}
	jsonResponse := fmt.Sprintf("{\"id\": \"%s\", \"created_at\": \"%s\",\"updated_at\": \"%s\",\"body\": \"%s\",\"user_id\": \"%s\"}",
		chirp.ID, chirp.CreatedAt, chirp.UpdatedAt, chirp.Body, chirp.UserID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(jsonResponse))
}
