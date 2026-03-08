package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BeilerOl/HobbyManager/backend/internal/repository/sqlite"
	"github.com/BeilerOl/HobbyManager/backend/internal/server"
)

func main() {
	dataSource := os.Getenv("DB_PATH")
	if dataSource == "" {
		dataSource = "file:hobby.db"
	}
	db, err := sqlite.NewDB(dataSource)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := server.NewMux(db)
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	// CORS for local frontend dev (e.g. Vue on localhost:5173)
	withCORS := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}
