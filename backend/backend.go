package backend

import (
	"ai-detective/llm"
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
	llm            llm.LLM
}

type MessageData struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type StartRoundData struct {
	Prompt string `json:"prompt"`
}

type EliminateData struct {
	UUID string `json:"uuid"`
}

type RespondData struct {
	Response string `json:"response"`
}

type Response struct {
	UUID     string `json:"uuid"`
	Response string `json:"response"`
}

type AllResponseData struct {
	Responses []Response `json:"responses"`
}

type StopRoundMessage struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	IsAi         bool   `json:"isAi"`
	NumOfPlayers int    `json:"numOfPlayers"`
	NumOfHumans  int    `json:"numOfHumans"`
}

type EliminatedPlayer struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	IsAi bool   `json:"isAi"`
}

type StopGameMessage struct {
	NumOfPlayers      int                `json:"numOfPlayers"`
	NumOfHumans       int                `json:"numOfHumans"`
	DetectiveWin      bool               `json:"didDetectiveWin"`
	EliminatedPlayers []EliminatedPlayer `json:"eliminatedPlayers"`
}

func NewServer() *Server {
	return &Server{
		m:              melody.New(),
		sessionUserMap: map[*melody.Session]*User{},
		game:           nil,
		isDetectiveIn:  false,
		llm:            llm.New(),
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
		s.sessionUserMap[session] = &user
		log.Printf("Added user %s to lobby", name)
	})

	// Receive messages from clients
	s.m.HandleMessage(func(session *melody.Session, msg []byte) {
		var data MessageData
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Printf("%v", err)
			return
		}

		s.mutex.Lock()
		defer s.mutex.Unlock()

		log.Printf("Receieved message: %s", data.Type)

		switch data.Type {
		case "beginGame":
			// Initialize game
			users := getHumansFromSessionUserMap(s.sessionUserMap)

			if len(users) == 0 {
				log.Printf("No humans in game")
				// TODO: Return message to Frontend
				return
			}

			// Create AIs
			numOfHumans := len(users)
			numOfAIs := numberOfAIs(numOfHumans)
			ais := s.llm.MakeAIs(numOfAIs)

			// Initialize game with players
			game := NewGame(users, ais)
			s.game = game
			log.Printf("Started game with %d humans and %d AI", len(users), len(ais))

			// Broadcast beginGame to all players
			response, _ := json.Marshal(data)
			s.m.Broadcast(response)
			log.Printf("Broadcasted beginGame to all players")

		case "beginRound":
			if s.game == nil {
				log.Printf("Game not initialized")
				return
			}

			// Broadcast beginRound to all players
			response, _ := json.Marshal(data)
			s.m.Broadcast(response)
			log.Printf("Broadcasted beginRound to all players")

			// Request response from AI
			var startRoundData StartRoundData
			if err := json.Unmarshal(data.Data, &startRoundData); err != nil {
				log.Printf("%v", err)
				return
			}

			prompt := startRoundData.Prompt

			for _, player := range s.game.UUIDToPlayers {
				if ai, ok := player.(*llm.AI); ok {
					go func() {
						s.mutex.Lock()
						defer s.mutex.Unlock()
						response := s.llm.AskAI(prompt, ai)
						s.game.ProcessResponse(ai, response)
						if s.game.EveryoneResponded() {
							s.BroadcastResponses()
						}
					}()
				}
			}

		case "respond":
			// Parse response data
			var response RespondData
			if err := json.Unmarshal(data.Data, &response); err != nil {
				log.Printf("%v", err)
				return
			}

			if s.game == nil {
				log.Printf("Game not initialized")
				return
			}

			user, ok := s.sessionUserMap[session]
			if !ok {
				log.Printf("User not found")
				return
			}

			if user.role == Detective {
				log.Printf("Detective cannot respond")
				return
			}

			// Process human response
			s.game.ProcessResponse(user, response.Response)
			if s.game.EveryoneResponded() {
				s.BroadcastResponses()
			}

		case "eliminate":
			// Parse elimination data
			var elimination EliminateData
			if err := json.Unmarshal(data.Data, &elimination); err != nil {
				log.Printf("%v", err)
				return
			}

			if s.game == nil {
				log.Printf("Game not initialized")
				return
			}

			// Process Elimination
			s.game.ProcessElimination(elimination.UUID)

			// Get round result
			roundResult := s.game.GetRoundResult()
			if roundResult == Continue {
				// Send stopRound message
				s.BroadcastStopRound(s.game.GetPlayerInfo(elimination.UUID))
				return
			}

			// End game, someone won
			s.BroadcastStopGame(roundResult)
		}
	})

	// Handle disconnection
	s.m.HandleDisconnect(func(session *melody.Session) {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		// Remove user from sessionUserMap
		user, ok := s.sessionUserMap[session]
		if !ok {
			log.Println("Disconnected user without session")
			return
		}
		delete(s.sessionUserMap, session)

		// If user is detective, set isDetectiveIn to false
		if user.role == Detective {
			s.isDetectiveIn = false

			if s.game != nil {
				// Detective disconnected, end game
				s.BroadcastStopGame(HumanWin)
			}
		} else {
			// TODO: If game is initialized, eliminate user silently (?)
			if s.game != nil {
				s.game.UUIDToPlayers[user.UUID()].Eliminate()
			}
		}
	})

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) BroadcastResponses() {

	// Once all responses have been given, send them back to the users
	var allResponses AllResponseData
	for player, response := range s.game.PlayerToResponse {
		allResponses.Responses = append(allResponses.Responses, Response{
			UUID:     player.UUID(),
			Response: response,
		})
	}

	responseData, _ := json.Marshal(allResponses)
	data := MessageData{
		Type: "finishResponses",
		Data: responseData,
	}

	response, _ := json.Marshal(data)
	s.m.Broadcast(response)
	log.Printf("Broadcasted finishResponses to all players")
}

func (s *Server) BroadcastStopRound(eliminatedPlayer PlayerInfo) {
	numOfPlayers := s.game.GetNumberOfActivePlayers()
	numOfHumans := s.game.GetNumberOfActiveHumans()

	// Send stopRound message
	stopResultMessage := StopRoundMessage{
		UUID:         eliminatedPlayer.uuid,
		Name:         eliminatedPlayer.name,
		IsAi:         eliminatedPlayer.isAi,
		NumOfPlayers: numOfPlayers,
		NumOfHumans:  numOfHumans,
	}
	stopResultData, _ := json.Marshal(stopResultMessage)
	data := MessageData{
		Type: "stopRound",
		Data: stopResultData,
	}

	response, _ := json.Marshal(data)
	s.m.Broadcast(response)
	log.Printf("Broadcasted stopRound to all players")
}

func (s *Server) BroadcastStopGame(roundResult RoundResult) {
	numOfPlayers := s.game.GetNumberOfActivePlayers()
	numOfHumans := s.game.GetNumberOfActiveHumans()
	detectiveWin := roundResult == DetectiveWin
	eliminatedPlayers := s.game.eliminatedPlayers

	// Process eliminatedPlayers
	eliminatedPlayersInfo := []EliminatedPlayer{}
	for _, player := range eliminatedPlayers {
		eliminatedPlayersInfo = append(eliminatedPlayersInfo, EliminatedPlayer{
			UUID: player.UUID(),
			Name: player.Name(),
			IsAi: player.IsAi(),
		})
	}

	// Send stopGame message
	stopGameMessage := StopGameMessage{
		NumOfPlayers:      numOfPlayers,
		NumOfHumans:       numOfHumans,
		DetectiveWin:      detectiveWin,
		EliminatedPlayers: eliminatedPlayersInfo,
	}
	stopGameData, _ := json.Marshal(stopGameMessage)
	data := MessageData{
		Type: "stopGame",
		Data: stopGameData,
	}

	response, _ := json.Marshal(data)
	s.m.Broadcast(response)
	log.Printf("Broadcasted stopGame to all players")

	// Reset game
	s.game = nil
	s.isDetectiveIn = false
	s.sessionUserMap = map[*melody.Session]*User{}
	s.m.CloseWithMsg([]byte("Game ended"))
}

func getHumansFromSessionUserMap(m map[*melody.Session]*User) []User {
	users := []User{}
	for _, user := range m {
		if user.role == Human {
			users = append(users, *user)
		}
	}
	return users
}
