package backend

import (
	"encoding/json"
	"net/http"

	"github.com/olahol/melody"
)

type Server struct {
	m              *melody.Melody
	SessionUserMap map[*melody.Session]*User
	game           *Game
	isDetectiveIn  bool
}

type MessageData struct {
	Type string      `json:"type"`
	Data interface{} `json:"message"`
}

func NewServer() Server {
	return Server{
		m:              melody.New(),
		SessionUserMap: map[*melody.Session]*User{},
		game:           nil,
		isDetectiveIn:  false,
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
		name := session.Request.URL.Query().Get("name")
		UUID := session.Request.URL.Query().Get("UUID")
		role_name := session.Request.URL.Query().Get("role")

		if name == "" || UUID == "" || role_name == "" {
			session.CloseWithMsg([]byte("Request Fields are missing"))
			return
		}

		role := GetRole(role_name)

		if role == Detective {
			if s.isDetectiveIn {
				session.CloseWithMsg([]byte("Detective already in game"))
				return
			} else {
				s.isDetectiveIn = true
			}
		}

		user := CreateUser(name, UUID, role)
		s.SessionUserMap[&melody.Session{}] = &user
	})

	// Receive messages from clients
	s.m.HandleMessage(func(session *melody.Session, msg []byte) {
		var data MessageData
		if err := json.Unmarshal(msg, &data); err != nil {
			return
		}

		if data.Type == "beginGame" {
			response, _ := json.Marshal(data)
			s.m.Broadcast(response)
			users := getUsersFromSessionUserMap(s.SessionUserMap)
			game := NewGame(users)
			s.game = &game
		} else if data.Type == "beginRound" {
			response, _ := json.Marshal(data)
			s.m.Broadcast(response)
		} else if data.Type == "respond" {
			// TODO
		} else if data.Type == "eliminate" {
			// TODO
		}
	})

	http.ListenAndServe(":5000", nil)
}

func getUsersFromSessionUserMap(m map[*melody.Session]*User) []User {
	users := []User{}
	for _, user := range m {
		users = append(users, *user)
	}
	return users
}
