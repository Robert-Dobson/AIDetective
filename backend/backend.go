package backend

import (
	"net/http"

	"github.com/olahol/melody"
)

type Server struct {
	m              *melody.Melody
	SessionUserMap map[*melody.Session]*User
}

func NewServer() Server {
	return Server{
		m:              melody.New(),
		SessionUserMap: map[*melody.Session]*User{},
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

	// Initialize user on websocket connection
	s.m.HandleConnect(func(session *melody.Session) {
		// TODO: Handle errors/invalid requests
		name := session.Request.URL.Query().Get("name")
		UUID := session.Request.URL.Query().Get("UUID")
		role_name := session.Request.URL.Query().Get("role")
		role := GetRole(role_name)

		user := CreateUser(name, UUID, role)
		s.SessionUserMap[&melody.Session{}] = &user
	})

	// Broadcast messages to all clients
	s.m.HandleMessage(func(session *melody.Session, msg []byte) {
		s.m.Broadcast(msg)
	})

	http.ListenAndServe(":5000", nil)
}
