package main

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	ComparePasswords(hashedPassword string, plainPassword string) bool
	HashAndSalt(password string) string
}

type BcryptPasswordHasher struct {
	cost int
}

func (bcryptPasswordHasher BcryptPasswordHasher) ComparePasswords(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

// not actually used here, except to compute the user's password hash once,
// would normally be used when setting or changing a password
func (bcryptPasswordHasher BcryptPasswordHasher) HashAndSalt(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcryptPasswordHasher.cost)
	return string(hash)
}
