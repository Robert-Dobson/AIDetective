package backend

// TODO: Handle what happens if the user disconnects at any time

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/olahol/melody"
)

type Server struct {
	m              *melody.Melody
	mutex          sync.Mutex
	sessionUserMap map[*melody.Session]*User
	game           *Game
	isDetectiveIn  bool
}

type MessageData struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type EliminateData struct {
	UUID string `json:"UUID"`
}

func NewServer() *Server {
	return &Server{
		m:              melody.New(),
		sessionUserMap: map[*melody.Session]*User{},
		game:           nil,
		isDetectiveIn:  false,
	}
}

func (s *Server) RunServer() {
	// Serve the frontend
	http.Handle("/", http.FileServer(http.Dir("frontend")))

	// Upgrade to websocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.m.HandleRequest(w, r)
	})

	// Initialize user on websocket connection
	s.m.HandleConnect(func(session *melody.Session) {
		name := session.Request.URL.Query().Get("name")
		UUID := session.Request.URL.Query().Get("uuid")
		roleName := session.Request.URL.Query().Get("role")

		if name == "" || UUID == "" || roleName == "" {
			session.CloseWithMsg([]byte("Request Fields are missing"))
			log.Printf("Refused connection as request fields are missing")
			return
		}

		s.mutex.Lock()
		defer s.mutex.Unlock()

		if s.game != nil {
			session.CloseWithMsg([]byte("Game already started"))
			log.Printf("Refused connection as game already started")
			return
		}

		role := GetRole(roleName)

		if role == Detective {
			if s.isDetectiveIn {
				session.CloseWithMsg([]byte("Detective already in game"))
				log.Printf("Refused connection as Detective already in game")
				return
			} else {
				s.isDetectiveIn = true
			}
		}

		user := CreateUser(name, UUID, role)
		s.sessionUserMap[&melody.Session{}] = &user
		log.Printf("Added user %s to lobby", name)
	})

	// Receive messages from clients
	s.m.HandleMessage(func(session *melody.Session, msg []byte) {
		var data MessageData
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Printf("%w", err)
			return
		}

		s.mutex.Lock()
		defer s.mutex.Unlock()

		log.Printf("Receieved message: %s", data.Type)

		switch data.Type {
		case "beginGame":
			// Initialize game
			users := getUsersFromSessionUserMap(s.sessionUserMap)
			game := NewGame(users)
			s.game = game
			log.Printf("Started game with %d players", len(users))

			// Broadcast beginGame to all players
			response, _ := json.Marshal(data)
			s.m.Broadcast(response)
			log.Printf("Broadcasted beginGame to all players")
		case "beginRound":
			// TODO: Ensure another player is in the game

			response, _ := json.Marshal(data)
			s.m.Broadcast(response)
			log.Printf("Broadcasted beginRound to all players")
		case "respond":
			// TODO
		case "eliminate":
			// Parse elimination data
			var data EliminateData
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Printf("%w", err)
				return
			}

			if s.game == nil {
				log.Printf("Game not initialized")
				return
			}

			// Process Elimination
			s.game.ProcessElimination(data.UUID)
		}
	})

	// Handle disconnection
	s.m.HandleDisconnect(func(session *melody.Session) {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		// Remove user from sessionUserMap
		user := s.sessionUserMap[session]
		delete(s.sessionUserMap, session)

		// If user is detective, set isDetectiveIn to false
		if user.role == Detective {
			s.isDetectiveIn = false

			if s.game != nil {
				// Detective disconnected, end game
				// TODO: Send End Game message
			}
		} else {
			// If game is initialized, eliminate user silently
			if s.game != nil {
				s.game.UUIDToPlayers[user.UUID()].Eliminate()
			}
		}
	})

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}

func getUsersFromSessionUserMap(m map[*melody.Session]*User) []User {
	users := []User{}
	for _, user := range m {
		users = append(users, *user)
	}
	return users
}
