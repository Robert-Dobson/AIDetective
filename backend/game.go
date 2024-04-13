package backend

import (
	"ai-detective/llm"
	"strings"
)

var howManyAI = map[int]int{
	1: 5,
	2: 7,
	3: 9,
}

type Player interface {
	UUID() string
	Name() string

	Eliminated() bool
	Eliminate()
}

type Game struct {
	UUIDToPlayers    map[string]Player
	PlayerToResponse map[Player]string
}

func NewGame(users []User, ais []llm.AI) *Game {
	g := Game{
		UUIDToPlayers:    map[string]Player{},
		PlayerToResponse: map[Player]string{},
	}

	// Add Players
	for _, user := range users {
		g.AddPlayer(&user)
	}

	// Add AIs
	for _, ai := range ais {
		g.AddPlayer(&ai)
	}

	return &g
}

func (g *Game) AddPlayer(player Player) {
	g.UUIDToPlayers[player.UUID()] = player
}

func (g *Game) BeginRound(prompt string) {
	// Reset responses
	g.PlayerToResponse = map[Player]string{}
}

func (g *Game) ProcessResponse(player Player, response string) {
	if !player.Eliminated() {
		g.PlayerToResponse[player] = sanitizeResponse(response)
	}
}

func (g *Game) EveryoneResponded() bool {
	return len(g.PlayerToResponse) >= g.GetNumberOfActivePlayers()
}

func (g *Game) ProcessElimination(UUID string) {
	// Process elimination of user, do we go onto the next round
	player, ok := g.UUIDToPlayers[UUID]
	if !ok {
		//TODO: Handle error
	}
	player.Eliminate()

	// Check if the game is over
	numActivePlayers := g.GetNumberOfActivePlayers()
	numActiveHumans := g.GetNumberOfActiveHumans()
	numActiveAIs := numActivePlayers - numActiveHumans

	if numActiveHumans == 0 {
		// Detective win condition
		// TODO: Send Detective win message
	} else if numActiveAIs < 3 {
		// Human win condition
		// TODO: Send Human win message
	} else {
		// Continue to next round
		// TODO: Send next round message
	}
}

func (g *Game) GetNumberOfActivePlayers() int {
	activePlayers := 0
	for _, player := range g.UUIDToPlayers {
		if !player.Eliminated() {
			activePlayers++
		}
	}
	return activePlayers
}

func (g *Game) GetNumberOfActiveHumans() int {
	activeHumans := 0
	for _, player := range g.UUIDToPlayers {
		if user, ok := player.(*User); ok && !user.Eliminated() && user.role == Human {
			activeHumans++
		}
	}
	return activeHumans
}

func (g *Game) CalculateLeaderboard() {
	// Calculate the leaderboard
}

func sanitizeResponse(resp string) string {
	// make responses lowercase and remove trailing punctuation
	respLower := strings.ToLower(resp)
	respNoPeriod := strings.TrimSuffix(respLower, ".")
	respNoExclamation := strings.TrimSuffix(respNoPeriod, "!")
	return respNoExclamation
}
