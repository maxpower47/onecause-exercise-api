package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var userRepository UserRepository
var passwordHasher PasswordHasher

func main() {
	// ideally things like this would be wired up with some sort of DI container
	userRepository = &MapUserRepository{}
	passwordHasher = &BcryptPasswordHasher{}

	r := gin.Default()

	// set CORS to be wide open; a real implementation would restrict this appropriately
	r.Use(cors.Default())

	r.POST("/login", login)

	r.Run()
}

func login(c *gin.Context) {
	var loginCmd LoginCommand
	c.BindJSON(&loginCmd)

	user, found := userRepository.FindUser(loginCmd.Username)

	if found && passwordHasher.ComparePasswords(user.PasswordHash, loginCmd.Password) {
		expectedToken := time.Now().Format("1504" /* HHMM */)

		if expectedToken == loginCmd.Token {
			c.Status(http.StatusOK)
			return
		}
	}

	c.Status(http.StatusUnauthorized)
}

// there's probably better ways to organize things like these value types than dumping them in the main file

type LoginCommand struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type User struct {
	Username     string
	PasswordHash string
}
