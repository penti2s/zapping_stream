package main

import (
	"log"
	"net/http"
	"zapping_stream/internal/hls"
)

func main() {
	hls.SetupUpdater()

	http.HandleFunc("/hls/playlist.m3u8", hls.ServePlaylist)
	http.HandleFunc("/hls/segments/", hls.ServeSegment)
	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
