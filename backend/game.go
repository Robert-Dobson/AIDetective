package backend

import "strings"

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

func NewGame(users []User) *Game {
	g := Game{
		UUIDToPlayers:    map[string]Player{},
		PlayerToResponse: map[Player]string{},
	}

	// Add Players
	for _, user := range users {
		g.AddPlayer(&user)
	}

	// Add AIs
	numOfHumans := len(users)
	numOfAIs, ok := howManyAI[numOfHumans]
	if !ok {
		numOfAIs = 3 * numOfHumans
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
	g.PlayerToResponse = map[Player]string{}
}

func (g *Game) ProcessResponse(player Player, response string) {
	if player.Eliminated() {
		return
	}

	g.PlayerToResponse[player] = sanitizeResponse(response)

	if len(g.PlayerToResponse) >= g.GetNumberOfActivePlayers() {
		g.FinalizeResponses()
	}
}

func (g *Game) FinalizeResponses() {
	// Once all responses have been given, send them back to the users
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
