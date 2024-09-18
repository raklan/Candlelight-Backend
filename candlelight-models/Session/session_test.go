package Session

import (
	"candlelight-models/Actions"
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-models/Player"
	"fmt"
	"maps"
	"slices"
	"testing"
)

func TestPlacementTurn_Execute(t *testing.T) {
	//TODO: Add tests for invalid player trying to submit turn
	var tests = []struct {
		name                 string
		turn                 PlacementTurn
		cardEndsInPlayerHand bool
		cardEndsInGame       bool
		shouldReturnError    bool
	}{
		{
			//Valid turns should execute
			name: "Valid Turn",
			turn: PlacementTurn{
				ActionId: "placement",
				PieceId:  "card",
				TargetId: "place",
			},
			cardEndsInPlayerHand: false,
			cardEndsInGame:       true,
			shouldReturnError:    false,
		},
		// {
		// 	//Invalid turns should make no change
		// 	name: "Invalid Turn",
		// 	turn: PlacementTurn{
		// 		ActionId: "invalid",
		// 		PieceId:  "card",
		// 		TargetId: "place",
		// 	},
		// 	cardEndsInPlayerHand: true,
		// 	cardEndsInGame:       false,
		// 	shouldReturnError:    true,
		// },
		{
			//Invalid piece should make no change
			name: "Invalid Piece",
			turn: PlacementTurn{
				ActionId: "placement",
				PieceId:  "invalid",
				TargetId: "place",
			},
			cardEndsInPlayerHand: true,
			cardEndsInGame:       false,
			shouldReturnError:    true,
		},
		{
			//Invalid target should make no change
			name: "Invalid Target",
			turn: PlacementTurn{
				ActionId: "placement",
				PieceId:  "card",
				TargetId: "invalid",
			},
			cardEndsInPlayerHand: true,
			cardEndsInGame:       false,
			shouldReturnError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Setup: Create gamestate, player, and card to move from player to game
			cardToPlay := Pieces.Card{
				GamePiece: Pieces.GamePiece{
					Id: "card",
				},
			}

			me := Player.Player{
				Id: "ryan",
				Hand: []Game.View{
					{
						Pieces: Pieces.PieceSet{
							Decks: []Pieces.Deck{
								{
									Cards: []Pieces.Card{
										cardToPlay,
									},
								},
							},
						},
					},
				},
			}

			myState := PlayerState{
				Player: me,
				AllowedActions: Actions.ActionSet{
					Placements: []Actions.Placement{
						{
							Action: Actions.Action{Id: "placement"},
						},
					},
				},
			}

			gameState := GameState{
				Views: []Game.View{
					{
						Pieces: Pieces.PieceSet{
							CardPlaces: []Pieces.CardPlace{
								{
									GamePiece: Pieces.GamePiece{
										Id: "place",
									},
								},
							},
						},
					},
				},
				PlayerStates: []PlayerState{
					myState,
				},
			}

			//Execute the turn
			changelog, err := tt.turn.Execute(&gameState, me)

			if err != nil && !tt.shouldReturnError {
				t.Errorf("%s -- Error Not nil: Got {%s}", tt.name, err)
			}

			playerDeck := myState.Player.Hand[0].Pieces.Decks[0]
			gameCardPlace := gameState.Views[0].Pieces.CardPlaces[0]

			cardInPlayerHand := cardInCollection(cardToPlay, playerDeck.Cards)
			cardInGame := cardInCollection(cardToPlay, gameCardPlace.PlacedCards)
			zeroedCardInHand := cardInCollection(Pieces.Card{GamePiece: Pieces.GamePiece{Id: ""}}, playerDeck.Cards)

			if !tt.cardEndsInPlayerHand && zeroedCardInHand {
				fmt.Println(playerDeck.Cards)
				t.Errorf("%s -- Zeroed card found in player hand", tt.name)
			}

			if tt.cardEndsInPlayerHand != cardInPlayerHand {
				t.Errorf("%s -- CardInPlayerHand: Expected {%t}, Got {%t}", tt.name, tt.cardEndsInPlayerHand, cardInPlayerHand)
			}
			if tt.cardEndsInGame != cardInGame {
				t.Errorf("%s -- CardInGame: Expected {%t}, Got {%t}", tt.name, tt.cardEndsInGame, cardInGame)
			}

			//Changelog should only have entries if no error
			if tt.shouldReturnError {
				if len(changelog.CardPlaces) != 0 || len(changelog.Decks) != 0 {
					t.Errorf("Changelog recorded changes! len(Decks) == %d, len(CardPlaces) == %d", len(changelog.Decks), len(changelog.CardPlaces))
				}
			} else {
				if !slices.ContainsFunc(changelog.Decks, func(d Pieces.Deck) bool { return d.Id == playerDeck.Id }) {
					t.Error("Deck card was removed from not found in changelog!")
				}
				if !slices.ContainsFunc(changelog.CardPlaces, func(c Pieces.CardPlace) bool { return c.Id == gameCardPlace.Id }) {
					t.Error("CardPlace where card was played not found in changelog!")
				}
			}
		})
	}
}

func TestMoveTurn_Execute(t *testing.T) {
	//TODO: Add tests for invalid player trying to submit turn
	var tests = []struct {
		name              string
		turn              MoveTurn
		cardEndsInDeck1   bool
		cardEndsInDeck2   bool
		shouldReturnError bool
	}{
		{
			//Valid turns should execute
			name: "Valid Turn",
			turn: MoveTurn{
				ActionId: "move",
				PieceId:  "card",
				TargetId: "place",
			},
			cardEndsInDeck1:   false,
			cardEndsInDeck2:   true,
			shouldReturnError: false,
		},
		// {
		// 	//Invalid turns should make no change
		// 	name: "Invalid Turn",
		// 	turn: MoveTurn{
		// 		ActionId: "invalid",
		// 		PieceId:  "card",
		// 		TargetId: "place",
		// 	},
		// 	cardEndsInDeck1:   true,
		// 	cardEndsInDeck2:   false,
		// 	shouldReturnError: true,
		// },
		{
			//Invalid piece should make no change
			name: "Invalid Piece",
			turn: MoveTurn{
				ActionId: "move",
				PieceId:  "invalid",
				TargetId: "place",
			},
			cardEndsInDeck1:   true,
			cardEndsInDeck2:   false,
			shouldReturnError: true,
		},
		{
			//Invalid target should make no change
			name: "Invalid Target",
			turn: MoveTurn{
				ActionId: "move",
				PieceId:  "card",
				TargetId: "invalid",
			},
			cardEndsInDeck1:   true,
			cardEndsInDeck2:   false,
			shouldReturnError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Setup: Create gamestate, player, and card to move from player to game
			cardToPlay := Pieces.Card{
				GamePiece: Pieces.GamePiece{
					Id: "card",
				},
			}

			me := Player.Player{
				Id: "ryan",
			}

			myState := PlayerState{
				Player: me,
				AllowedActions: Actions.ActionSet{
					Moves: []Actions.Move{
						{
							Action: Actions.Action{Id: "move"},
						},
					},
				},
			}

			gameState := GameState{
				Views: []Game.View{
					{
						Pieces: Pieces.PieceSet{
							Decks: []Pieces.Deck{
								{
									GamePiece: Pieces.GamePiece{
										Id: "deck",
									},
									Cards: []Pieces.Card{cardToPlay},
								},
								{
									GamePiece: Pieces.GamePiece{
										Id: "place",
									},
								},
							},
						},
					},
				},
				PlayerStates: []PlayerState{
					myState,
				},
			}

			//Execute the turn
			changelog, err := tt.turn.Execute(&gameState, me)

			if err != nil && !tt.shouldReturnError {
				t.Errorf("%s -- Error Not nil: Got {%s}", tt.name, err)
			}

			deck1 := gameState.Views[0].Pieces.Decks[0]
			deck2 := gameState.Views[0].Pieces.Decks[1]

			cardInDeck1 := cardInCollection(cardToPlay, deck1.Cards)
			cardInDeck2 := cardInCollection(cardToPlay, deck2.Cards)

			if tt.cardEndsInDeck1 != cardInDeck1 {
				t.Errorf("%s -- CardInPlayerHand: Expected {%t}, Got {%t}", tt.name, tt.cardEndsInDeck1, cardInDeck1)
			}
			if tt.cardEndsInDeck2 != cardInDeck2 {
				t.Errorf("%s -- CardInGame: Expected {%t}, Got {%t}", tt.name, tt.cardEndsInDeck2, cardInDeck2)
			}

			//For a valid move Changelog should contain the 2 Decks, since they were affected
			if tt.shouldReturnError {
				if len(changelog.Decks) != 0 || len(changelog.CardPlaces) != 0 {
					t.Errorf("Changelog recorded changes! len(Decks) == %d, len(CardPlaces) == %d", len(changelog.Decks), len(changelog.CardPlaces))
				}
			} else {
				if !slices.ContainsFunc(changelog.Decks, func(d Pieces.Deck) bool { return d.Id == deck1.Id }) {
					t.Error("Couldn't find Deck1 in changelog!")
				}
				if !slices.ContainsFunc(changelog.Decks, func(d Pieces.Deck) bool { return d.Id == deck2.Id }) {
					t.Error("Couldn't find Deck2 in changelog!")
				}
			}
		})
	}
}

func TestPieceUpdateTurn_Execute(t *testing.T) {
	//TODO: Add tests for invalid player trying to submit turn
	var tests = []struct {
		name              string
		turn              PieceUpdateTurn
		tagsShouldChange  bool
		shouldReturnError bool
	}{
		{
			//Valid turns should execute
			name: "Valid Turn",
			turn: PieceUpdateTurn{
				ActionId:      "update",
				TargetPieceId: "piece",
				NewTags:       map[string]string{"newTag": "newValue"},
			},
			tagsShouldChange:  true,
			shouldReturnError: false,
		},
		// {
		// 	//Invalid turns should make no change
		// 	name: "Invalid Turn",
		// 	turn: PieceUpdateTurn{
		// 		ActionId:      "invalid",
		// 		TargetPieceId: "piece",
		// 		NewTags:       map[string]string{"newTag": "newValue"},
		// 	},
		// 	tagsShouldChange:  false,
		// 	shouldReturnError: true,
		// },
		{
			//Invalid piece should make no change
			name: "Invalid Piece",
			turn: PieceUpdateTurn{
				ActionId:      "update",
				TargetPieceId: "invalid",
				NewTags:       map[string]string{"newTag": "newValue"},
			},
			tagsShouldChange:  false,
			shouldReturnError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Setup: Create gamestate and player
			me := Player.Player{
				Id: "ryan",
			}

			myState := PlayerState{
				Player: me,
				AllowedActions: Actions.ActionSet{
					PieceUpdates: []Actions.PieceUpdate{
						{
							Action: Actions.Action{Id: "update"},
						},
					},
				},
			}

			gameState := GameState{
				Views: []Game.View{
					{
						Pieces: Pieces.PieceSet{
							Decks: []Pieces.Deck{
								{
									GamePiece: Pieces.GamePiece{
										Id:   "piece",
										Tags: map[string]string{"old": "deleteMe"},
									},
								},
							},
						},
					},
				},
				PlayerStates: []PlayerState{
					myState,
				},
			}

			//Execute the turn
			changelog, err := tt.turn.Execute(&gameState, me)

			if (err != nil) != tt.shouldReturnError {
				t.Errorf("%s -- Error value: Expected {%t}, Got {%s}", tt.name, tt.shouldReturnError, err)
			}

			deck := gameState.Views[0].Pieces.Decks[0]
			if !tt.shouldReturnError && !maps.Equal(deck.GamePiece.Tags, tt.turn.NewTags) {
				fmt.Println(tt.turn.NewTags)
				fmt.Println(deck.GamePiece.Tags)
				t.Errorf("%s -- Tags do not line up", tt.name)
			}

			if tt.shouldReturnError {
				if len(changelog.Decks) != 0 {
					t.Error("Found Decks in changelog! Should not have found any!")
				}
			} else {
				if !slices.ContainsFunc(changelog.Decks, func(d Pieces.Deck) bool { return d.Id == deck.Id }) {
					t.Error("Changed deck not found in changelog!")
				}
			}
		})
	}
}

func TestTakeTurn_Execute(t *testing.T) {
	//TODO: Add tests for invalid player trying to submit turn
	var tests = []struct {
		name                 string
		turn                 TakeTurn
		cardEndsInPlayerHand bool
		cardEndsInGame       bool
		shouldReturnError    bool
	}{
		{
			//Valid turns should execute
			name: "Valid Turn Specified Piece",
			turn: TakeTurn{
				ActionId:     "take",
				PieceId:      "card",
				TakingFromId: "deck",
			},
			cardEndsInPlayerHand: true,
			cardEndsInGame:       false,
			shouldReturnError:    false,
		},
		{
			//Valid turns should execute
			name: "Valid Turn Unspecified Piece",
			turn: TakeTurn{
				ActionId:     "take",
				PieceId:      "",
				TakingFromId: "deck",
			},
			cardEndsInPlayerHand: true,
			cardEndsInGame:       true,
			shouldReturnError:    false,
		},
		// {
		// 	//Invalid turns should make no change
		// 	name: "Invalid Turn",
		// 	turn: PlacementTurn{
		// 		ActionId: "invalid",
		// 		PieceId:  "card",
		// 		TargetId: "place",
		// 	},
		// 	cardEndsInPlayerHand: true,
		// 	cardEndsInGame:       false,
		// 	shouldReturnError:    true,
		// },
		{
			//Invalid piece should make no change
			name: "Invalid PieceId",
			turn: TakeTurn{
				ActionId:     "take",
				PieceId:      "invalid",
				TakingFromId: "deck",
			},
			cardEndsInPlayerHand: false,
			cardEndsInGame:       true,
			shouldReturnError:    true,
		},
		{
			//Invalid target should make no change
			name: "Invalid TargetId",
			turn: TakeTurn{
				ActionId:     "take",
				PieceId:      "card",
				TakingFromId: "invalid",
			},
			cardEndsInPlayerHand: false,
			cardEndsInGame:       true,
			shouldReturnError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Setup: Create gamestate, player, and card to move from player to game
			cardToTake := Pieces.Card{
				GamePiece: Pieces.GamePiece{
					Id: "card",
				},
			}

			me := Player.Player{
				Id: "ryan",
				Hand: []Game.View{
					{
						Pieces: Pieces.PieceSet{
							Decks: []Pieces.Deck{
								{},
							},
						},
					},
				},
			}

			myState := PlayerState{
				Player: me,
				AllowedActions: Actions.ActionSet{
					Placements: []Actions.Placement{
						{
							Action: Actions.Action{Id: "take"},
						},
					},
				},
			}

			gameState := GameState{
				Views: []Game.View{
					{
						Pieces: Pieces.PieceSet{
							Decks: []Pieces.Deck{
								{
									GamePiece: Pieces.GamePiece{
										Id: "deck",
									},
									Cards: []Pieces.Card{
										cardToTake,
									},
								},
							},
						},
					},
				},
				PlayerStates: []PlayerState{
					myState,
				},
			}

			//Execute the turn
			changelog, err := tt.turn.Execute(&gameState, me)

			if (err != nil) != tt.shouldReturnError {
				t.Errorf("%s -- Error value: Expected {%t}, Got {%s}", tt.name, tt.shouldReturnError, err)
			}

			playerOrphansFromState := gameState.PlayerStates[0].Player.Hand[0].Pieces.Orphans.Cards

			if tt.turn.PieceId == "" { //Special case for unspecified piece
				if len(playerOrphansFromState) < 1 {
					t.Errorf("%s -- Length of Player's hand < 1", tt.name)
				}
			} else {
				playerHandOrphans := gameState.PlayerStates[0].Player.Hand[0].Pieces.Orphans.Cards
				gameDeck := gameState.Views[0].Pieces.Decks[0]

				cardInPlayerHand := cardInCollection(cardToTake, playerHandOrphans)
				cardInGame := cardInCollection(cardToTake, gameDeck.Cards)
				zeroedCardInGame := cardInCollection(Pieces.Card{GamePiece: Pieces.GamePiece{Id: ""}}, gameDeck.Cards)

				if zeroedCardInGame {
					t.Errorf("%s -- Zeroed card found in game!", tt.name)
				}

				if tt.cardEndsInPlayerHand != cardInPlayerHand {
					t.Errorf("%s -- CardInPlayerHand: Expected {%t}, Got {%t}", tt.name, tt.cardEndsInPlayerHand, cardInPlayerHand)
				}
				if tt.cardEndsInGame != cardInGame {
					t.Errorf("%s -- CardInGame: Expected {%t}, Got {%t}", tt.name, tt.cardEndsInGame, cardInGame)
				}

				//Changelog should only have entries if no error
				if tt.shouldReturnError {
					if len(changelog.CardPlaces) != 0 || len(changelog.Decks) != 0 {
						t.Errorf("Changelog recorded changes! len(Decks) == %d, len(CardPlaces) == %d", len(changelog.Decks), len(changelog.CardPlaces))
					}
				} else {
					if !slices.ContainsFunc(changelog.OrphanDecks[0].Cards, func(c Pieces.Card) bool {
						return c.Id == gameState.PlayerStates[0].Player.Hand[0].Pieces.Orphans.Cards[0].Id
					}) {
						t.Error("Card player took not found in changelog!")
					}
					if !slices.ContainsFunc(changelog.Decks, func(d Pieces.Deck) bool { return d.Id == gameDeck.Id }) {
						t.Error("Deck where card was taken from not found in changelog!")
					}
				}
			}
		})
	}
}

// Checks if a card with the same id as [card] exists in [collection]
func cardInCollection(card Pieces.Card, collection []Pieces.Card) bool {
	for _, c := range collection {
		if c.Id == card.Id {
			return true
		}
	}

	return false
}
