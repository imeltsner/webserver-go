package main

import (
	"log"
	"net/http"
	"os"

	"github.com/imeltsner/webserver-go/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	jwtSecret      string
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	const port = "8080"
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("Error creating DB: %s", err)
		return
	}

	cfg := apiConfig{
		db:             db,
		fileserverHits: 0,
		jwtSecret:      jwtSecret,
	}

	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricHandler)
	mux.HandleFunc("/api/reset", cfg.resetHandler)
	mux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpHandler)
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("POST /api/login", cfg.loginUserHandler)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
