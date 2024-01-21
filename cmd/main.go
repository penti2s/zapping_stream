package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
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
	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func enableCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Puedes ajustar "*" para ser más específico dependiendo de tus necesidades de seguridad.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next(w, r)
	}
}
