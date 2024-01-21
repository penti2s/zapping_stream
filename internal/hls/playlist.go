package hls

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ServePlaylist(w http.ResponseWriter, r *http.Request) {
	log.Println("Solicitado playlist.m3u8")

	m3u8Path := filepath.Join("assets", "hls", "segment.m3u8")

	if _, err := os.Stat(m3u8Path); os.IsNotExist(err) {
		http.Error(w, "Archivo no encontrado.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")

	http.ServeFile(w, r, m3u8Path)
}

func ServeSegment(w http.ResponseWriter, r *http.Request) {
	log.Println("Solicitado segmento:", r.URL.Path)

	segmentName := filepath.Base(r.URL.Path)

	segmentPath := filepath.Join("assets", "hls", segmentName)

	if !strings.HasSuffix(segmentPath, ".ts") {
		http.Error(w, "Archivo no válido.", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, segmentPath)
}
