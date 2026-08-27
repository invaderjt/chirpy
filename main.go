package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("GET /api/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /api/reset", apiCfg.handleReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func handleHealthz(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	_, err := writer.Write([]byte("OK"))
	if err != nil {
		log.Println(err)
	}
}

func (cfg *apiConfig) handleMetrics(writer http.ResponseWriter, req *http.Request) {
	hits := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	_, err := writer.Write([]byte(hits))
	if err != nil {
		log.Println(err)
	}
}

func (cfg *apiConfig) handleReset(writer http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Swap(0)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(writer, req)
	})
}
