package CreationStudio

import (
	"candlelight-api/LogUtil"
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-models/Sparks"
	"candlelight-ruleengine/Engine"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"net/http"
	"slices"
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

func generatePickColor(id string) []float32 {
	hasher := fnv.New32a()
	hasher.Write([]byte(id))
	hash := hasher.Sum32()
	r := float32((hash & 0xFF))
	g := float32((hash >> 8) & 0xFF)
	b := float32((hash >> 16) & 0xFF)
	return []float32{r, g, b, 255}
}

func generateCardArray(sharedViewId string) []Pieces.Card {
	cards := []Pieces.Card{}
	for _, suit := range []string{
		"Spades", "Hearts", "Diamonds", "Clubs",
	} {
		aceId := Engine.GenerateId()
		jackId := Engine.GenerateId()
		queenId := Engine.GenerateId()
		kingId := Engine.GenerateId()
		ace := Pieces.Card{
			GamePiece: Pieces.GamePiece{
				Id:         aceId,
				Name:       fmt.Sprintf("Ace of %s", suit),
				Color:      []float32{255, 255, 255, 255},
				PickColor:  generatePickColor(aceId),
				Tags:       map[string]string{},
				Text:       fmt.Sprintf("Ace of %s", suit),
				Layer:      0,
				X:          0,
				Y:          0,
				ParentView: sharedViewId,
			},
		}
		jack := Pieces.Card{
			GamePiece: Pieces.GamePiece{
				Id:         jackId,
				Name:       fmt.Sprintf("Jack of %s", suit),
				Color:      []float32{255, 255, 255, 255},
				PickColor:  generatePickColor(jackId),
				Tags:       map[string]string{},
				Text:       fmt.Sprintf("Jack of %s", suit),
				Layer:      0,
				X:          0,
				Y:          0,
				ParentView: sharedViewId,
			},
		}
		queen := Pieces.Card{
			GamePiece: Pieces.GamePiece{
				Id:         queenId,
				Name:       fmt.Sprintf("Queen of %s", suit),
				Color:      []float32{255, 255, 255, 255},
				PickColor:  generatePickColor(queenId),
				Tags:       map[string]string{},
				Text:       fmt.Sprintf("Queen of %s", suit),
				Layer:      0,
				X:          0,
				Y:          0,
				ParentView: sharedViewId,
			},
		}
		king := Pieces.Card{
			GamePiece: Pieces.GamePiece{
				Id:         kingId,
				Name:       fmt.Sprintf("King of %s", suit),
				Color:      []float32{255, 255, 255, 255},
				PickColor:  generatePickColor(kingId),
				Tags:       map[string]string{},
				Text:       fmt.Sprintf("King of %s", suit),
				Layer:      0,
				X:          0,
				Y:          0,
				ParentView: sharedViewId,
			},
		}
		for x := range 9 {
			id := Engine.GenerateId()
			card := Pieces.Card{
				GamePiece: Pieces.GamePiece{
					Id:         id,
					Name:       fmt.Sprintf("%d of %s", x+2, suit),
					Color:      []float32{255, 255, 255, 255},
					PickColor:  generatePickColor(id),
					Tags:       map[string]string{},
					Text:       fmt.Sprintf("%d of %s", x+2, suit),
					Layer:      0,
					X:          0,
					Y:          0,
					ParentView: sharedViewId,
				},
			}
			cards = append(cards, card)
		}
		cards = append(cards, ace)
		cards = append(cards, jack)
		cards = append(cards, queen)
		cards = append(cards, king)
	}
	return cards
}

// Generates the dummy game and inserts it into the local DB. Useful for testing
func GenerateJSON(w http.ResponseWriter, r *http.Request) {
	sharedViewId := Engine.GenerateId()
	player1ViewId := Engine.GenerateId()
	player2ViewId := Engine.GenerateId()
	player3ViewId := Engine.GenerateId()
	player4ViewId := Engine.GenerateId()

	cards := generateCardArray(sharedViewId)
	slices.SortFunc(cards, func(a Pieces.Card, b Pieces.Card) int {
		if rand.Intn(100)%2 == 0 {
			return 1
		} else {
			return -1
		}
	})

	//cardPlaceId := Engine.GenerateId()
	deckId := Engine.GenerateId()

	game := Game.Game{
		Id:        "game123",
		Name:      "Shuffled Deck of Cards",
		Genre:     "Card",
		Author:    "CandlelightDevTeam",
		Published: true,
		Rules: Game.GameRules{
			ShowOtherPlayerDetails: true,
			EnforceTurnOrder:       true,
		},
		Sparks: Sparks.Sparks{
			Dealer: Sparks.Dealer{
				Enabled:   true,
				NumToDeal: 5,
				DeckToUse: deckId,
			},
		},
		MaxPlayers: 4,
		Views: []Game.View{
			{
				Id:                sharedViewId,
				OwnerPlayerNumber: 0,
				Pieces: Pieces.PieceSet{
					CardPlaces: []Pieces.CardPlace{},
					Decks: []Pieces.Deck{
						{
							GamePiece: Pieces.GamePiece{
								Id:         deckId,
								Name:       "Draw Deck",
								Color:      []float32{0.5, 0.5, 0.5, 1},
								PickColor:  generatePickColor(deckId),
								Tags:       map[string]string{},
								Text:       "Draw Deck",
								Layer:      0,
								X:          50,
								Y:          50,
								ParentView: sharedViewId,
							},
							Cards: cards,
						},
					},
					Orphans: []Pieces.Card{},
				},
			},
			{
				Id:                player1ViewId,
				OwnerPlayerNumber: 1,
				Pieces: Pieces.PieceSet{
					Decks:      []Pieces.Deck{},
					CardPlaces: []Pieces.CardPlace{},
					Orphans:    []Pieces.Card{},
				},
			},
			{
				Id:                player2ViewId,
				OwnerPlayerNumber: 2,
				Pieces: Pieces.PieceSet{
					Decks:      []Pieces.Deck{},
					CardPlaces: []Pieces.CardPlace{},
					Orphans:    []Pieces.Card{},
				},
			},
			{
				Id:                player3ViewId,
				OwnerPlayerNumber: 3,
				Pieces: Pieces.PieceSet{
					Decks:      []Pieces.Deck{},
					CardPlaces: []Pieces.CardPlace{},
					Orphans:    []Pieces.Card{},
				},
			},
			{
				Id:                player4ViewId,
				OwnerPlayerNumber: 4,
				Pieces: Pieces.PieceSet{
					Decks:      []Pieces.Deck{},
					CardPlaces: []Pieces.CardPlace{},
					Orphans:    []Pieces.Card{},
				},
			},
		},
	}
	Engine.SaveGameDefToDB(game)
	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(game)
	asJson, _ := json.Marshal(game)
	fmt.Fprint(w, string(asJson))
}
