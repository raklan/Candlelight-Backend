package Engine

import (
	"candlelight-models/Game"
	"fmt"
	"log"
	"math/rand" //May want to change this to crypto/rand for better security, but for the prototype this is more than fine
	"slices"
	"time"

	"github.com/go-redis/redis/v8"
)

type Criteria struct {
	Authors []string
	Genres  []string
}

// Returns true IFF for each array in Criteria, if any values are given, the Game's matching field
// must be equal to one of the values in that array. If an array is empty, the Game automatically passes
// that part of the Criteria
func (c Criteria) Check(game Game.Game) bool {
	return (len(c.Authors) == 0 || slices.Contains(c.Authors, game.Author)) &&
		(len(c.Genres) == 0 || slices.Contains(c.Genres, game.Genre))
}

const (
	ModuleLogPrefix  = "CANDLELIGHT-RULEENGINE"
	PackageLogPrefix = "Engine"
	RedisAddress     = "redis:6379"
)

var RDB = redis.NewClient(&redis.Options{
	Addr: RedisAddress,
})

// Generates an ID for something. To ensure it's unique, I'm just using the current UNIX time in
// milliseconds with a random set of 10 characters appended to the end. Will probably need to change to something more random later
func GenerateId() string {
	const lettersAndNumbers = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 10)
	for i := range code {
		code[i] = lettersAndNumbers[rand.Intn(len(lettersAndNumbers))]
	}

	return fmt.Sprint(time.Now().UnixMilli(), string(code))
}

// Generates a random 4-character room code. For use in creating lobbies
func generateRoomCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, 4)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}

	return string(code)
}

// Logs the error in the format of "[funcLogPrefix] ERROR! [err]"
func LogError(funcLogPrefix string, err error) {
	log.Printf("%s ERROR! %s", funcLogPrefix, err)
}
