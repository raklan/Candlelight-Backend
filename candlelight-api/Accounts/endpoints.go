package Accounts

import (
	"candlelight-api/LogUtil"
	"candlelight-ruleengine/Accounts"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	funcLogPrefix := "==CreateAccount=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(LogUtil.ModuleLogPrefix, PackagePrefix)

	log.Printf("%s received request to create account!", funcLogPrefix)

	d := json.NewDecoder(r.Body)
	req := Accounts.User{}

	err := d.Decode(&req)
	if err != nil {
		LogUtil.LogError(funcLogPrefix, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//They must supply both a Username and Password
	if !(req.Username != "" && req.Password != "") {
		log.Print("Client did not supply both Username and Password. Rejecting request")
		http.Error(w, "Please supply both Username and Password", http.StatusBadRequest)
		return
	}

	log.Printf("%s Sending user to be saved", funcLogPrefix)
	saved, err := Accounts.SaveNewAccount(req)
	if err != nil {
		LogUtil.LogError(funcLogPrefix, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%s Save successful, sending response", funcLogPrefix)

	asJson, _ := json.Marshal(saved)
	fmt.Fprint(w, string(asJson))
}

// func ChangePassword(w http.ResponseWriter, r *http.Request) { //TODO: Currently no authentication for this
// 	funcLogPrefix := "==ChangePassword=="
// 	defer LogUtil.EnsureLogPrefixIsReset()
// 	LogUtil.SetLogPrefix(LogUtil.ModuleLogPrefix, PackagePrefix)

// 	log.Printf("%s received request to change password!", funcLogPrefix)

// 	d := json.NewDecoder(r.Body)
// 	req := Accounts.User{}

// 	err := d.Decode(&req)
// 	if err != nil {
// 		LogUtil.LogError(funcLogPrefix, err)
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	saved, err := Accounts.ChangePassword(req)

// 	if err != nil {
// 		LogUtil.LogError(funcLogPrefix, err)
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	log.Printf("%s Save successful, sending response", funcLogPrefix)

// 	asJson, _ := json.Marshal(saved)
// 	fmt.Fprint(w, string(asJson))
// }

func Login(w http.ResponseWriter, r *http.Request) {
	funcLogPrefix := "==Login=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(LogUtil.ModuleLogPrefix, PackagePrefix)

	log.Printf("%s received request to login!", funcLogPrefix)

	d := json.NewDecoder(r.Body)
	req := Accounts.User{}

	err := d.Decode(&req)
	if err != nil {
		LogUtil.LogError(funcLogPrefix, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !(req.Username != "" && req.Password != "") {
		log.Print("Client did not supply both Username and Password")
		http.Error(w, "Please supply both Username and Password", http.StatusBadRequest)
		return
	}

	user, err := Accounts.AttemptLogin(req)

	if err != nil {
		//Commented out for security
		//LogUtil.LogError(funcLogPrefix, err)
		http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
		return
	}

	log.Printf("%s Login successful, sending response", funcLogPrefix)

	asJson, _ := json.Marshal(user)
	fmt.Fprint(w, string(asJson))
}
