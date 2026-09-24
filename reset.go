package main

import "net/http"

func (cfg *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(writer, nil, "Forbidden", 403)
		return
	}
	cfg.fileserverHits.Swap(0)
	cfg.dbQueries.ResetUsers(req.Context())
}
