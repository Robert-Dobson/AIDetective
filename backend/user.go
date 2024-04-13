package backend

import "strings"

type Role int

const (
	Detective Role = iota
	Human
)

// User implements Player
type User struct {
	name       string
	uuid       string
	role       Role
	eliminated bool
}

func (u *User) UUID() string {
	return u.uuid
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Eliminated() bool {
	return u.eliminated
}

func (u *User) Eliminate() {
	u.eliminated = true
}

func CreateUser(name string, UUID string, role Role) User {
	return User{
		name:       name,
		uuid:       UUID,
		role:       role,
		eliminated: false,
	}
}

func GetRole(role string) Role {
	switch strings.ToLower(role) {
	case "detective":
		return Detective
	case "human":
		return Human
	default:
		// Always default to human
		return Human
	}
}
