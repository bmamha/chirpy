package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/bmamha/chirpy/internal/database"
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
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	log.Fatal(s.ListenAndServe())
}
