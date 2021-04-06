package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	userRepository UserRepository
	passwordHasher PasswordHasher
	infoLog        *log.Logger
)

func main() {
	// ideally things like this would be wired up with some sort of DI container
	userRepository = &MapUserRepository{}
	passwordHasher = &BcryptPasswordHasher{cost: 10}

	// would probably want a more comprehensive structured logging solution for real use
	infoLog = log.New(log.Writer(), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	r := gin.Default()

	// set CORS to be wide open; a real implementation would restrict this appropriately
	r.Use(cors.Default())

	r.POST("/login", login)

	r.Run()
}

// usually good to make sure a login function like this returns in constant time, regardless of
// whether the user exists or their password is correct, etc to prevent username enumeration attacks
func login(c *gin.Context) {
	var loginCmd LoginCommand
	c.BindJSON(&loginCmd)

	// should probably have some more error handling as well

	user, found := userRepository.FindUser(loginCmd.Username)

	if found && passwordHasher.ComparePasswords(user.PasswordHash, loginCmd.Password) {
		expectedToken := time.Now().Format("1504" /* HHMM */)

		if expectedToken == loginCmd.Token {
			c.Status(http.StatusOK)
			infoLog.Printf("Login succeeded - %v", loginCmd.Username)
			return
		}
	}

	// would be better to differentiate between different types of failures when logging
	infoLog.Printf("Login failed - %v", loginCmd.Username)

	// all failures are returned as 401, but may be better differentiate certain types of errors
	// e.g. 400 if the body was malformed, 415 for wrong content type, etc
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
