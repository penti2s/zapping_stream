package hls

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func SetupUpdater() {
	go updatePlaylistPeriodically()
}

func updatePlaylistPeriodically() {

	fmt.Println("Esperando...")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := updatePlaylist()
			if err != nil {
				// Manejar el error, posiblemente con un log
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func updatePlaylist() error {
	m3u8Path := "assets/hls/segment.m3u8"

	// Leer el archivo .m3u8 actual
	file, err := os.ReadFile(m3u8Path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(file), "\n")
	updatedLines, _ := modifyPlaylist(lines)

	return os.WriteFile(m3u8Path, []byte(strings.Join(updatedLines, "\n")), 0666)
}

func modifyPlaylist(lines []string) ([]string, int) {
	var header []string
	var segments []string
	var mediaSequence int
	var targetDuration string

	for _, line := range lines {
		if strings.HasPrefix(line, "#EXT-X-TARGETDURATION:") {
			targetDuration = strings.Split(line, ":")[1]
			header = append(header, line)
		} else if strings.HasPrefix(line, "#EXTINF:") || strings.HasPrefix(line, "segment") {
			segments = append(segments, line)
		} else if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:") {
			mediaSequence, _ = strconv.Atoi(strings.Split(line, ":")[1])
			header = append(header, line)
		} else {
			header = append(header, line)
		}
	}

	mediaSequence++
	nextSegmentFile := fmt.Sprintf("segment%d.ts", mediaSequence-1)
	nextSegmentPath := fmt.Sprintf("assets/hls/%s", nextSegmentFile)

	if !segmentExists(nextSegmentPath) {
		mediaSequence = 1
		nextSegmentFile = "segment0.ts"
	}

	nextSegment := fmt.Sprintf("#EXTINF:%s,\n%s", targetDuration, nextSegmentFile)
	segments = append(segments, nextSegment)

	if len(segments) > 6 {
		segments = segments[2:]
	}

	// Actualiza la secuencia de medios en el encabezado
	for i, line := range header {
		if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:") {
			header[i] = fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d", mediaSequence)
			break
		}
	}

	return append(header, segments...), mediaSequence
}

func segmentExists(segmentPath string) bool {
	if _, err := os.Stat(segmentPath); os.IsNotExist(err) {
		return false
	}
	return true
}
