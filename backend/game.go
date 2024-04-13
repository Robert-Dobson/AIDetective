package backend

type Game struct {
	Users []User
}

func BeginGame() Game {
	// Initialize game, add AI roles to Users
	return Game{
		Users: []User{},
	}
}

func BeginRound() {
	// Get AI responses to prompt

}

func ProcessResponse() {
	// Process responses from users

}

func FinalizeResponses() {
	// Once all responses have been given, send them back to the users
}

func ProcessElimination() {
	// Process elimination of user, do we go onto the next round
}

func CalculateLeaderboard() {
	// Calculate the leaderboard
}
