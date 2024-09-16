package Accounts

import (
	"candlelight-api/LogUtil"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func SaveNewAccount(user User) (SafeUser, error) {
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	funcLogPrefix := "==SaveNewAccount==:"
	log.Printf("%s Received Account to create with username == {%s}", funcLogPrefix, user.Username)

	//Username overlap checking
	log.Printf("%s Checking if username is already taken...", funcLogPrefix)
	_, err := RDB.Get(ctx, "user:"+user.Username).Result()
	if err != redis.Nil {
		log.Printf("%s User with username == {%s} already exists!", funcLogPrefix, user.Username)
		return SafeUser{user.Username}, fmt.Errorf("Username already taken") //I know it shouldn't be capitalized, but this string goes straight to the client
	}

	//Password hashing
	log.Printf("Hashing user's password...")
	hashed, err := HashPassword(user.Password)

	if err != nil {
		LogError(funcLogPrefix, err)
		return SafeUser{user.Username}, err
	}

	user.Password = hashed

	key := "user:" + user.Username
	log.Printf("%s Password hashed. Saving User to DB with key == {%s}", funcLogPrefix, key)

	//Save to redis
	asJson, err := json.Marshal(user)
	if err != nil {
		LogError(funcLogPrefix, err)
		return SafeUser{user.Username}, err
	}

	err = RDB.Set(ctx, key, asJson, 0).Err()
	if err != nil {
		LogError(funcLogPrefix, err)
		return SafeUser{user.Username}, err
	}

	log.Printf("%s User saved with key == {%s}", funcLogPrefix, key)

	return SafeUser{user.Username}, nil
}

func ChangePassword(user User) (SafeUser, error) {
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	funcLogPrefix := "==ChangePassword==:"
	log.Printf("%s Attempting login for user with username == {%s}", funcLogPrefix, user.Username)

	//Check if user exists
	log.Printf("%s Checking if user exists...", funcLogPrefix)
	dbUserAsJson, err := RDB.Get(ctx, "user:"+user.Username).Result()
	if err == redis.Nil {
		log.Printf("%s Could not find User with username == {%s}", funcLogPrefix, user.Username)
		return SafeUser{user.Username}, fmt.Errorf("%s Could not find User with username == {%s}", funcLogPrefix, user.Username)
	}

	dbUser := User{}
	json.Unmarshal([]byte(dbUserAsJson), &dbUser)

	//Hash and set new password
	hashed, err := HashPassword(user.Password)
	if err != nil {
		LogError(funcLogPrefix, err)
		return SafeUser{user.Username}, err
	}

	dbUser.Password = hashed

	key := "user:" + dbUser.Username
	log.Printf("%s Saving User to DB with key == {%s}", funcLogPrefix, key)

	//Save to redis
	asJson, err := json.Marshal(dbUser)
	if err != nil {
		LogError(funcLogPrefix, err)
		return SafeUser{user.Username}, err
	}

	err = RDB.Set(ctx, key, asJson, 0).Err()
	if err != nil {
		LogError(funcLogPrefix, err)
		return SafeUser{user.Username}, err
	}

	log.Printf("%s User saved with key == {%s}", funcLogPrefix, key)

	return SafeUser{dbUser.Username}, nil
}

func AttemptLogin(user User) (SafeUser, error) {
	//Note: Lots of logging is commented out. THIS IS INTENTIONAL for security. Uncomment ONLY for debugging!
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	funcLogPrefix := "==AttemptLogin==:"
	log.Printf("%s Attempting login for user with username == {%s}", funcLogPrefix, user.Username)

	//Check if user exists
	//log.Printf("%s Checking if user exists...", funcLogPrefix)
	dbUserAsJson, err := RDB.Get(ctx, "user:"+user.Username).Result()
	if err == redis.Nil {
		//log.Printf("%s Could not find User with username == {%s}", funcLogPrefix, user.Username)
		return SafeUser{user.Username}, fmt.Errorf("%s Could not find User with username == {%s}", funcLogPrefix, user.Username)
	}

	dbUser := User{}
	json.Unmarshal([]byte(dbUserAsJson), &dbUser)

	//Check password
	passwordMatch := CheckPasswordHash(user.Password, dbUser.Password)

	if !passwordMatch {
		//log.Printf("%s Password did not match", funcLogPrefix)
		return SafeUser{user.Username}, fmt.Errorf("%s Password did not match", funcLogPrefix)
	}

	return SafeUser{dbUser.Username}, nil
}
