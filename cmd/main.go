package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"zapping_stream/internal/auth"
	"zapping_stream/internal/db"
	"zapping_stream/internal/hls"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	err = db.InitDB()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos: ", err)
	}

	hls.SetupUpdater()

	http.HandleFunc("/hls/playlist.m3u8", enableCors(hls.ServePlaylist))
	http.HandleFunc("/hls/", enableCors(hls.ServeSegment))

	http.HandleFunc("/api/auth/register", enableCors(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Register(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResponse := map[string]string{"token": token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonResponse)
	}))

	http.HandleFunc("/api/auth/login", enableCors(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Login(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResponse := map[string]string{"token": token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonResponse)
	}))
	http.HandleFunc("/api/auth/user", enableCors(auth.UserHandler))

	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func enableCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next(w, r)
	}
}
