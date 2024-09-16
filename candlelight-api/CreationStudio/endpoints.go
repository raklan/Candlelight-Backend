package CreationStudio

import (
	"candlelight-api/LogUtil"
	"candlelight-models/Game"
	"candlelight-ruleengine/Engine"
	"encoding/json"
	"log"
	"net/http"
)

// Routing endpoint for the Creation Studio. GET requests are used to load a certain GameDef, POST are used to save a GameDef, DELETE to delete
func Studio(w http.ResponseWriter, r *http.Request) {
	/**
	* Reminder on REST api structuring
	* CREATE -> POST
	* READ -> GET
	* UPDATE -> PUT
	* DELETE -> DELETE
	 */
	switch r.Method {
	case "GET":
		getGame(w, r)
	case "POST":
		saveGame(w, r)
	case "DELETE":
		deleteGame(w, r)
	default:
		http.Error(w, "Method not allowed or implemented", http.StatusMethodNotAllowed)
	}
}

func saveGame(w http.ResponseWriter, r *http.Request) {
	funcLogPrefix := "==SaveGame==:"
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(LogUtil.ModuleLogPrefix, PackagePrefix)

	log.Printf("%s Received request to save a game!", funcLogPrefix)

	d := json.NewDecoder(r.Body)
	req := Game.Game{}

	err := d.Decode(&req)
	if err != nil {
		log.Printf("%s ERROR! %s", funcLogPrefix, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%s Sending game to be saved...", funcLogPrefix)
	saved, err := Engine.SaveGameDefToDB(req)
	if err != nil {
		log.Printf("%s ERROR! %s", funcLogPrefix, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%s Save successful, sending response to client...", funcLogPrefix)

	// Return the game data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saved)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	funcLogPrefix := "==GetGame==:"

	log.Printf("%s Received request to get a game!", funcLogPrefix)

	// Assuming the game ID is passed as a query parameter
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		log.Printf("%s No game ID provided", funcLogPrefix)
		http.Error(w, "No game ID provided", http.StatusBadRequest)
		return
	}

	log.Printf("%s Sending request for game from DB", funcLogPrefix)
	game, err := Engine.GetGameDefFromDB(gameID)
	if err != nil {
		log.Printf("%s ERROR: %s", funcLogPrefix, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("%s Load successful, sending response to client", funcLogPrefix)
	// Return the game data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func deleteGame(w http.ResponseWriter, r *http.Request) {
	funcLogPrefix := "==DeleteGame==:"
	LogUtil.SetLogPrefix(LogUtil.ModuleLogPrefix, Engine.PackageLogPrefix)

	log.Printf("%s Received request to delete a game!", funcLogPrefix)

	// Assuming the game ID is passed as a query parameter
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		log.Printf("%s No game ID provided", funcLogPrefix)
		http.Error(w, "No Game Definition ID provided", http.StatusBadRequest)
		return
	}

	deleted, err := Engine.DeleteGameDefFromDB(gameID)
	if err != nil {
		log.Printf("%s Error trying to delete game: %s", funcLogPrefix, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("%s Game deleted", funcLogPrefix)
	// Return the game data as JSON
	w.Header().Set("Content-Type", "application/json")
	type DeleteRepsonse struct {
		DeletedId string
	}
	json.NewEncoder(w).Encode(DeleteRepsonse{DeletedId: deleted})
}

func GetAllGames(w http.ResponseWriter, r *http.Request) {
	funcLogPrefix := "==GetAllGames=="

	log.Printf("%s Received request for ALL games", funcLogPrefix)

	authors := r.URL.Query()["author"]
	genres := r.URL.Query()["genre"]

	criteria := struct {
		Authors []string
		Genres  []string
	}{
		Authors: authors,
		Genres:  genres,
	}

	games, err := Engine.GetAllGamesFromDB(criteria)
	if err != nil {
		log.Printf("%s ERROR: %s", funcLogPrefix, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("%s Load successful, checking if client wants it slimmed or not", funcLogPrefix)
	w.Header().Set("Content-Type", "application/json")

	slimmedReq := r.URL.Query().Get("slimmed")
	if slimmedReq == "" {
		log.Printf("%s [slimmed] URL param not detected. Returning full object...", funcLogPrefix)
		json.NewEncoder(w).Encode(games)
	} else if slimmedReq == "true" {
		//TODO: This is kinda a terrible way to slim this down. Definitely will need optimization as the system gets bigger
		log.Printf("%s Client requested slimmed object. Returning array with just IDs and Names...", funcLogPrefix)
		type SlimmedGame struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}
		slimmedArr := []SlimmedGame{}

		for _, g := range games {
			new := SlimmedGame{
				Id:   g.Id,
				Name: g.Name,
			}

			slimmedArr = append(slimmedArr, new)
		}

		json.NewEncoder(w).Encode(slimmedArr)
	} else {
		log.Printf("%s [slimmed] URL param set to unrecognized value {%s}. Returning full object...", funcLogPrefix, slimmedReq)
		json.NewEncoder(w).Encode(games)
	}
}
