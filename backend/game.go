package backend

import (
	"ai-detective/llm"
	"strings"
)

type RoundResult int

const (
	DetectiveWin RoundResult = iota
	HumanWin
	Continue
)

var howManyAI = map[int]int{
	1: 3,
	2: 5,
	3: 8,
}

type Player interface {
	UUID() string
	Name() string
	IsAi() bool

	Eliminated() bool
	Eliminate()
}

type PlayerInfo struct {
	uuid string
	name string
	isAi bool
}

type Game struct {
	UUIDToPlayers     map[string]Player
	PlayerToResponse  map[Player]string
	eliminatedPlayers []Player
}

func NewGame(users []User, ais []llm.AI) *Game {
	g := Game{
		UUIDToPlayers:     map[string]Player{},
		PlayerToResponse:  map[Player]string{},
		eliminatedPlayers: []Player{},
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
	if ok {
		player.Eliminate()
		g.eliminatedPlayers = append(g.eliminatedPlayers, player)
	}
}

func (g *Game) GetRoundResult() RoundResult {
	// Check if the game is over
	numActivePlayers := g.GetNumberOfActivePlayers()
	numActiveHumans := g.GetNumberOfActiveHumans()
	numActiveAIs := numActivePlayers - numActiveHumans

	if numActiveHumans == 0 {
		return DetectiveWin
	}
	if numActiveAIs < 3 {
		return HumanWin
	}
	return Continue
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

func sanitizeResponse(resp string) string {
	// make responses lowercase and remove trailing punctuation
	respLower := strings.ToLower(resp)
	respNoPeriod := strings.TrimSuffix(respLower, ".")
	respNoExclamation := strings.TrimSuffix(respNoPeriod, "!")
	respNoQuotes := strings.Replace(respNoExclamation, "\"", "", -1)
	return respNoQuotes
}

func numberOfAIs(numberOfHumans int) int {
	numOfAIs, ok := howManyAI[numberOfHumans]
	if !ok {
		numOfAIs = 3 * numberOfHumans
	}

	return numOfAIs
}

func (g *Game) GetPlayerInfo(uuid string) PlayerInfo {
	player := g.UUIDToPlayers[uuid]
	return PlayerInfo{name: player.Name(), uuid: player.UUID(), isAi: player.IsAi()}

}
