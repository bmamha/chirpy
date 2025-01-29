package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/bmamha/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	PLATFORM       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("unable to open database")
	}
	dbQueries := database.New(db)

	directoryPath := "."
	_, err = os.Stat(directoryPath)
	if os.IsNotExist(err) {
		fmt.Printf("Directory %s not found.\n", directoryPath)
	}
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	apiCfg := &apiConfig{db: dbQueries, PLATFORM: os.Getenv("PLATFORM")}

	fileServer := http.FileServer(http.Dir(directoryPath))
	handler := http.StripPrefix("/app", fileServer)

	mux.Handle("/app/", apiCfg.middleWareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetHandler)
	mux.HandleFunc("POST /api/users", apiCfg.UserCreationHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetChirpsHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.ChirpCreationHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpHandler)
	log.Fatal(s.ListenAndServe())
}

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (apiCfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", apiCfg.fileServerHits.Load())))
	if err != nil {
		log.Fatal()
	}
}

func (apiCfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if apiCfg.PLATFORM != "dev" {
		w.WriteHeader(403)
		return
	}

	err := apiCfg.db.DeleteUsers(r.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Error deleting users"))
		return
	}

	apiCfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) middleWareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) UserCreationHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json "email"`
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

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
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

func (cfg *apiConfig) ChirpCreationHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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

	userid, err := uuid.Parse(params.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("\"error\": Failed to parse id given"))
		return
	}
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

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("\"error\": Error getting chirps"))
		return
	}

	jsonchirps := []Chirp{}

	for _, chirp := range chirps {
		jsonchirps = append(jsonchirps, Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.ID})
	}

	jsonResponse, err := json.Marshal(jsonchirps)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonResponse)
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
