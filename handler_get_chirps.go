package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	t := r.URL.Query().Get("sort")

	if s == "" {
		chirps, err := cfg.db.GetChirps(r.Context())
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte("\"error\": Error getting chirps"))
			return
		}
		if t == "asc" || t == "" {
			sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
		} else {
			sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
		}

		jsonchirps := []Chirp{}

		for _, chirp := range chirps {
			jsonchirps = append(jsonchirps, Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID})
		}

		jsonResponse, err := json.Marshal(jsonchirps)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(jsonResponse)
	} else {
		user_id, err := uuid.Parse(s)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte("\"error\": Error getting chirps"))
			return
		}
		chirps, err := cfg.db.GetChirpsByAuthor(r.Context(), user_id)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte("\"error\": Error getting chirps"))
			return

		}
		jsonchirps := []Chirp{}

		for _, chirp := range chirps {
			jsonchirps = append(jsonchirps, Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID})
		}

		jsonResponse, err := json.Marshal(jsonchirps)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(jsonResponse)

	}
}

func (cfg *apiConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {
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

	jsonChirp := Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.ID}
	jsonResponse, err := json.Marshal(jsonChirp)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonResponse)
}
