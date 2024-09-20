package CreationStudio

import (
	"candlelight-api/LogUtil"
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-ruleengine/Engine"
	"encoding/json"
	"fmt"
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

// Generates the dummy game and inserts it into the local DB. Useful for testing
func GenerateJSON(w http.ResponseWriter, r *http.Request) {
	sharedViewId := Engine.GenerateId()
	player1ViewId := Engine.GenerateId()
	player2ViewId := Engine.GenerateId()
	player3ViewId := Engine.GenerateId()
	player4ViewId := Engine.GenerateId()
	game := Game.Game{
		Id:         "game123",
		Name:       "HAPPY FUN GAME!!!!",
		Genre:      "Happy and fun",
		Author:     "Ryan",
		MaxPlayers: 4,
		Resources:  []Game.GameResource{},
		Views: []Game.View{
			{
				Id:                sharedViewId,
				OwnerPlayerNumber: 0,
				Pieces: Pieces.PieceSet{
					Decks: []Pieces.Deck{
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: sharedViewId,
								Name:       "Ancient Deck",
								Tags:       map[string]string{},
								Position:   Pieces.Position{X: 0, Y: 0},
								Style:      Pieces.Style{},
							},
							Cards: []Pieces.Card{
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Earth",
										Tags: map[string]string{
											"color": "#186e21",
										},
										ParentView: sharedViewId,
									},
									Value: 1,
								},
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Wind",
										Tags: map[string]string{
											"color": "#bad9e8",
										},
										ParentView: sharedViewId,
									},
									Value: 2,
								},
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Fire",
										Tags: map[string]string{
											"color": "#ed8f6f",
										},
										ParentView: sharedViewId,
									},
									Value: 3,
								},
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Ice",
										Tags: map[string]string{
											"color": "#aaf7fa",
										},
										ParentView: sharedViewId,
									},
									Value: 4,
								},
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Air",
										Tags: map[string]string{
											"color": "#d5eaeb",
										},
										ParentView: sharedViewId,
									},
									Value: 5,
								},
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Thunder",
										Tags: map[string]string{
											"color": "#ffed91",
										},
										ParentView: sharedViewId,
									},
									Value: 6,
								},
							},
						},
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: sharedViewId,
								Name:       "Primitive Deck",
								Tags:       map[string]string{},
								Position:   Pieces.Position{X: 0, Y: 0},
								Style:      Pieces.Style{},
							},
							Cards: []Pieces.Card{
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Lava",
										Tags: map[string]string{
											"color": "#a15d55",
										},
										ParentView: sharedViewId,
									},
									Value: 1,
								},
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Plasma",
										Tags: map[string]string{
											"color": "#d1b4de",
										},
										ParentView: sharedViewId,
									},
									Value: 2,
								},
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Ribonucleic Acid",
										Tags: map[string]string{
											"color": "#b4deb8",
										},
										ParentView: sharedViewId,
									},
									Value: 3,
								},
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Brown",
										Tags: map[string]string{
											"color": "#594233",
										},
										ParentView: sharedViewId,
									},
									Value: 4,
								},
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Sparks",
										Tags: map[string]string{
											"color": "#fce481",
										},
										ParentView: sharedViewId,
									},
									Value: 5,
								},
								{
									GamePiece: Pieces.GamePiece{
										Id:   Engine.GenerateId(),
										Name: "Storm",
										Tags: map[string]string{
											"color": "#c1c2de",
										},
										ParentView: sharedViewId,
									},
									Value: 6,
								},
							},
						},
					},
					CardPlaces: []Pieces.CardPlace{
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: sharedViewId,
								Name:       "Prehistoric World",
								Tags:       map[string]string{},
								Position:   Pieces.Position{X: 0, Y: 0},
								Style:      Pieces.Style{},
							},
							PlacedCards: []Pieces.Card{},
						},
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: sharedViewId,
								Name:       "Old World",
								Tags:       map[string]string{},
								Position:   Pieces.Position{X: 0, Y: 0},
								Style:      Pieces.Style{},
							},
							PlacedCards: []Pieces.Card{},
						},
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: sharedViewId,
								Name:       "Modern World",
								Tags:       map[string]string{},
								Position:   Pieces.Position{X: 0, Y: 0},
								Style:      Pieces.Style{},
							},
							PlacedCards: []Pieces.Card{},
						},
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: sharedViewId,
								Name:       "Space World",
								Tags:       map[string]string{},
								Position:   Pieces.Position{X: 0, Y: 0},
								Style:      Pieces.Style{},
							},
							PlacedCards: []Pieces.Card{},
						},
					},
					Orphans: Pieces.Deck{
						GamePiece: Pieces.GamePiece{
							Id:         Engine.GenerateId(),
							ParentView: sharedViewId,
						},
						Cards: []Pieces.Card{},
					},
				},
			},
			{
				Id:                player1ViewId,
				OwnerPlayerNumber: 1,
				Pieces: Pieces.PieceSet{
					Decks: []Pieces.Deck{
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: player1ViewId,
							},
							Cards: []Pieces.Card{},
						},
					},
					CardPlaces: []Pieces.CardPlace{
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: player1ViewId,
							},
							PlacedCards: []Pieces.Card{},
						},
					},
					Orphans: Pieces.Deck{
						GamePiece: Pieces.GamePiece{
							Id:         Engine.GenerateId(),
							ParentView: player1ViewId,
						},
						Cards: []Pieces.Card{},
					},
				},
			},
			{
				Id:                player2ViewId,
				OwnerPlayerNumber: 2,
				Pieces: Pieces.PieceSet{
					Decks: []Pieces.Deck{
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: player2ViewId,
							},
							Cards: []Pieces.Card{},
						},
					},
					CardPlaces: []Pieces.CardPlace{
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: player2ViewId,
							},
							PlacedCards: []Pieces.Card{},
						},
					},
					Orphans: Pieces.Deck{
						GamePiece: Pieces.GamePiece{
							Id:         Engine.GenerateId(),
							ParentView: player2ViewId,
						},
						Cards: []Pieces.Card{},
					},
				},
			},
			{
				Id:                player3ViewId,
				OwnerPlayerNumber: 3,
				Pieces: Pieces.PieceSet{
					Decks: []Pieces.Deck{
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: player3ViewId,
							},
							Cards: []Pieces.Card{},
						},
					},
					CardPlaces: []Pieces.CardPlace{
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: player3ViewId,
							},
							PlacedCards: []Pieces.Card{},
						},
					},
					Orphans: Pieces.Deck{
						GamePiece: Pieces.GamePiece{
							Id:         Engine.GenerateId(),
							ParentView: player3ViewId,
						},
						Cards: []Pieces.Card{},
					},
				},
			},
			{
				Id:                player4ViewId,
				OwnerPlayerNumber: 4,
				Pieces: Pieces.PieceSet{
					Decks: []Pieces.Deck{
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: player4ViewId,
							},
							Cards: []Pieces.Card{},
						},
					},
					CardPlaces: []Pieces.CardPlace{
						{
							GamePiece: Pieces.GamePiece{
								Id:         Engine.GenerateId(),
								ParentView: player4ViewId,
							},
							PlacedCards: []Pieces.Card{},
						},
					},
					Orphans: Pieces.Deck{
						GamePiece: Pieces.GamePiece{
							Id:         Engine.GenerateId(),
							ParentView: player4ViewId,
						},
						Cards: []Pieces.Card{},
					},
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
