package backend

import (
	"net/http"

	"github.com/olahol/melody"
)

type Server struct {
	m *melody.Melody
}

func NewServer() Server {
	return Server{
		m: melody.New(),
	}
}

func (s Server) RunServer() {
	// Serve the frontend
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/chattest.html")
	})

	// Upgrade to websocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.m.HandleRequest(w, r)
	})

	// Broadcast messages to all clients
	s.m.HandleMessage(func(session *melody.Session, msg []byte) {
		s.m.Broadcast(msg)
	})

	http.ListenAndServe(":5000", nil)
}
