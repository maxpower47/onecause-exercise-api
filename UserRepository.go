package main

type UserRepository interface {
	FindUser(username string) (User, bool)
}

// simple implementation of a UserRepository based on a static map of users, indexed by username.
// real usage would be based on e.g. a database lookup
type MapUserRepository struct {
}

func (mapUserRepository MapUserRepository) FindUser(username string) (User, bool) {
	user, found := userTable[username]
	return user, found
}

var userTable = map[string]User{
	"c137@onecause.com": {Username: "c137@onecause.com", PasswordHash: "$2a$04$/quOaP3z173uvQxIA2g67ORiLBox.uvEUd3m6zuRhcQtyPBFS2ZqO"},
}
