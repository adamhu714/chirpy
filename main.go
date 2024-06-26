package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileServerHits int
	jwtSecret      string
	polkaApiKey    string
}

func main() {
	godotenv.Load()
	const port = "8080"
	const fileServerPathRoot = "./public/"

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		debugMode()
	}

	apiCfg := apiConfig{
		fileServerHits: 0,
		jwtSecret:      os.Getenv("JWT_SECRET"),
		polkaApiKey:    os.Getenv("POLKA_APIKEY"),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fileServerPathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerPostChirps)
	mux.HandleFunc("GET /api/chirps", handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", handlerGetChirpsId)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCfg.handlerDeleteChirpsId)

	mux.HandleFunc("POST /api/users", handlerPostUsers)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerPutUsers)

	mux.HandleFunc("POST /api/login", apiCfg.handlerPostLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerPostRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerPostRevoke)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPostPolkaWebhooks)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html>\n\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n\n</html>\n", cfg.fileServerHits)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func debugMode() {
	os.Remove("database.json")
}
