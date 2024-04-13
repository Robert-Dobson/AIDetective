package backend

import "strings"

type Role int

const (
	Detective Role = iota
	Human
	AI
)

type User struct {
	Name string
	UUID string
	role Role
}

func CreateUser(name string, UUID string, role Role) User {
	return User{
		Name: name,
		UUID: UUID,
		role: role,
	}
}

func GetRole(role string) Role {
	switch strings.ToLower(role) {
	case "detective":
		return Detective
	case "human":
		return Human
	case "ai":
		return AI
	default:
		// Always default to human
		return Human
	}
}
