package main

import (
	"candlelight-api/Accounts"
	"candlelight-api/CreationStudio"
	"candlelight-api/Lobby"
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-ruleengine/Engine"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
)

func main() {
	startServer()
}

func startServer() {
	//Go does the strangest datetime string formatting I've ever seen. You give it a specific date/time (Specifically Jan 2, 2006 3:04:05 PM GMT-7)
	//in the format you want, and it'll match whatever the object is into that format
	logName := fmt.Sprintf("./logs/%v.log", time.Now().Format("2006-01-02_15-04-05"))

	//Log file & Server startup
	log.SetPrefix("CANDLELIGHT-API: ")
	logfile, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	log.Println("Starting HTTP listener...")

	//Start the server at localhost:10000 & register all paths
	mux := http.NewServeMux()
	registerPathHandlers(mux)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	handler := c.Handler(mux)

	http.ListenAndServe(":10000", handler)
}
func registerPathHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", heartbeat)
	mux.HandleFunc("/dummy", generateJSON)

	//Gamedef-related requests
	mux.HandleFunc("/studio", CreationStudio.Studio)
	mux.HandleFunc("/allGames", CreationStudio.GetAllGames)

	//Lobby-related requests
	mux.HandleFunc("/joinLobby", Lobby.HandleJoinLobby)
	mux.HandleFunc("/hostLobby", Lobby.HostLobby)
	mux.HandleFunc("/rejoinLobby", Lobby.HandleRejoinLobby)

	//Account-related Requests
	mux.HandleFunc("/createAccount", Accounts.CreateAccount)
	mux.HandleFunc("/login", Accounts.Login)
	//mux.HandleFunc("/changePassword", Accounts.ChangePassword)
}

// Simple heartbeat endpoint to test if the server is up and running
func heartbeat(w http.ResponseWriter, r *http.Request) {
	log.Println("==Heartbeat==: Returning dummy response...")
	fmt.Fprintf(w, "Buh-dump, buh-dump")
}

// Generates the dummy game and inserts it into the local DB. Useful for testing
func generateJSON(w http.ResponseWriter, r *http.Request) {
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
