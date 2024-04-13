package backend

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
}

func NewGame(numOfHumans int) *Game {
	// Initalize game, add AI roles to Users
	numOfAI, ok := howManyAI[numOfHumans]
	if !ok {
		numOfAI = 2 * numOfHumans
	}

	for i := 0; i < numOfAI; i++ {
		// TODO: Create AIs
	}

	return &Game{}
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

func (g *Game) ProcessElimination() {
	// Process elimination of user, do we go onto the next round
}

func (g *Game) CalculateLeaderboard() {
	// Calculate the leaderboard
}
