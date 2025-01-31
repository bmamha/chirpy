package main

import (
	"net/http"
)

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
