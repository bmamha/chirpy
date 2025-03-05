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
	SECRET         string
	POLKA_KEY      string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("unable to open database")
	}
	dbQueries := database.New(db)

	directoryPath := http.Dir(".")
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)

	mux := http.NewServeMux()
	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	apiCfg := &apiConfig{
		db:        dbQueries,
		PLATFORM:  os.Getenv("PLATFORM"),
		SECRET:    os.Getenv("SECRET"),
		POLKA_KEY: os.Getenv("POLKA_KEY"),
	}

	//	fsHandler := apiCfg.middleWareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(directoryPath))))

	mux.Handle("/", apiCfg.middleWareMetricsInc(http.FileServer(directoryPath)))
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetHandler)
	mux.HandleFunc("POST /api/users", apiCfg.UserCreationHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetChirpsHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.ChirpCreationHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.userUpdateHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.ChirpDeletionHandler)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.userUpgradeChirpyRed)
	log.Fatal(s.ListenAndServe())
}
