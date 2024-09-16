package Accounts

import (
	"candlelight-ruleengine/Engine"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// A user with the password field omitted
type SafeUser struct {
	Username string `json:"username"`
}

const (
	ModuleLogPrefix  = "CANDLELIGHT-RULEENGINE"
	PackageLogPrefix = "Accounts"
)

var RDB = Engine.RDB

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Logs the error in the format of "[funcLogPrefix] ERROR! [err]"
func LogError(funcLogPrefix string, err error) {
	Engine.LogError(funcLogPrefix, err)
}

func GenerateId() string {
	return Engine.GenerateId()
}
