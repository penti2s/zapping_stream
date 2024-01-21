package chat

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)

var mutex = &sync.Mutex{}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Registrar la nueva conexión.
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// Eliminar la conexión cuando se cierra o cuando ocurre un error.
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			return
		}
		broadcastMessage(messageType, p)
	}
}

func broadcastMessage(messageType int, p []byte) {
	mutex.Lock()
	defer mutex.Unlock()
	for client := range clients {
		if err := client.WriteMessage(messageType, p); err != nil {
			// Eliminar la conexión si hay un error al enviar el mensaje.
			delete(clients, client)
			client.Close()
		}
	}
}
