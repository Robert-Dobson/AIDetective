package backend

type Game struct {
	Users []User
}

func NewGame(users []User) Game {
	// Reset/Initalize game, add AI roles to Users
	return Game{
		Users: users,
	}
}

func (g Game) BeginRound() {
	// Get AI responses to prompt

}

func (g Game) ProcessResponse() {
	// Process responses from users

}

func (g Game) FinalizeResponses() {
	// Once all responses have been given, send them back to the users
}

func (g Game) ProcessElimination() {
	// Process elimination of user, do we go onto the next round
}

func (g Game) CalculateLeaderboard() {
	// Calculate the leaderboard
}
