package backend

import "github.com/google/uuid"

var howManyAI = map[int]int{
	1: 2,
	2: 3,
	3: 4,
}

type Game struct {
	users []User
}

func NewGame(users []User) *Game {
	// Initalize game, add AI roles to Users
	numOfHuman := len(users)
	numOfAI, ok := howManyAI[numOfHuman]
	if !ok {
		numOfAI = 2 * numOfHuman
	}

	for i := 0; i < numOfAI; i++ {
		// TODO: Make custom AI names
		users = append(users, CreateUser("AI", uuid.New().String(), AI))
	}

	return &Game{
		users: users,
	}
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
