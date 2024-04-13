package backend

import "strings"

var howManyAI = map[int]int{
	1: 2,
	2: 3,
	3: 4,
}

type Player interface {
	UUID() string
	Name() string

	Eliminated() bool
	Eliminate()
}

type Game struct {
	UUIDToPlayers map[string]Player
}

func NewGame(users []User) *Game {
	g := Game{UUIDToPlayers: map[string]Player{}}

	// Add Players
	for _, user := range users {
		g.AddPlayer(&user)
	}

	// Add AIs
	numOfHumans := len(users)
	numOfAIs, ok := howManyAI[numOfHumans]
	if !ok {
		numOfAIs = 2 * numOfHumans
	}

	for i := 0; i < numOfAIs; i++ {
		// TODO: Create AIs
	}

	return &g
}

func (g *Game) AddPlayer(user *User) {
	g.UUIDToPlayers[user.UUID()] = user
}

func (g *Game) BeginRound() {
	// Get AI responses to prompt

}

func (g *Game) ProcessResponse() {
	// Process responses from users

}

func (g *Game) FinalizeResponses() {
	// Once all responses have been given, send them back to the users
}

func (g *Game) ProcessElimination(UUID string) {
	// Process elimination of user, do we go onto the next round

}

func (g *Game) CalculateLeaderboard() {
	// Calculate the leaderboard
}

func (g *Game) sanitizeResponse(resp string) string {
	// make responses lowercase and remove trailing period
	respLower := strings.ToLower(resp)
	respNoPeriod := strings.TrimSuffix(respLower, ".")
	return respNoPeriod
}
