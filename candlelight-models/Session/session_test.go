package Session

import (
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-models/Player"
	"slices"
	"testing"
)

func TestInsertion_Execute(t *testing.T) {
	var tests = []struct {
		Name                      string
		Insertion                 Insertion
		ExpectedChangelogLength   int
		ShouldReturnError         bool
		CardEndsInFirstView       bool
		CardEndsInDestinationView bool
	}{
		{
			Name: "Valid Insertion",
			Insertion: Insertion{
				InsertCard:   "cardToInsert",
				FromView:     "fromView",
				ToCollection: "toCollection",
				InView:       "inView",
			},
			ExpectedChangelogLength:   2,
			ShouldReturnError:         false,
			CardEndsInFirstView:       false,
			CardEndsInDestinationView: true,
		},
		{
			Name: "Invalid First View",
			Insertion: Insertion{
				InsertCard:   "cardToInsert",
				FromView:     "invalid",
				ToCollection: "toCollection",
				InView:       "inView",
			},
			ExpectedChangelogLength:   1,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
		{
			Name: "Invalid Card",
			Insertion: Insertion{
				InsertCard:   "invalid",
				FromView:     "fromView",
				ToCollection: "toCollection",
				InView:       "inView",
			},
			ExpectedChangelogLength:   2,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
		{
			Name: "Invalid Destination",
			Insertion: Insertion{
				InsertCard:   "cardToInsert",
				FromView:     "fromView",
				ToCollection: "toCollection",
				InView:       "invalid",
			},
			ExpectedChangelogLength:   1,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
		{
			Name: "Invalid Collection",
			Insertion: Insertion{
				InsertCard:   "cardToInsert",
				FromView:     "fromView",
				ToCollection: "invalid",
				InView:       "inView",
			},
			ExpectedChangelogLength:   2,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			cardInQuestion := Pieces.Card{
				GamePiece: Pieces.GamePiece{
					Id: "cardToInsert",
				},
			}

			fromView := Game.View{
				Id: "fromView",
				Pieces: Pieces.PieceSet{
					Orphans: []Pieces.Card{cardInQuestion},
				},
			}

			toView := Game.View{
				Id: "inView",
				Pieces: Pieces.PieceSet{
					Decks: []Pieces.Deck{
						{
							GamePiece: Pieces.GamePiece{
								Id: "toCollection",
							},
						},
					},
				},
			}

			me := Player.Player{
				Id: "me",
				Hand: []Game.View{
					fromView,
				},
			}

			gameState := GameState{
				Views:   []Game.View{toView},
				Players: []Player.Player{me},
			}

			changelog, err := tt.Insertion.Execute(&gameState, me.Id)

			//Check if we got an error when we shouldn't have, and vice versa
			if tt.ShouldReturnError != (err != nil) {
				t.Fatalf("ERROR in returned error value. Expected error: %t, err == %s", tt.ShouldReturnError, err)
			}

			cardIsInFirstView := cardInView("cardToInsert", &gameState.Players[0].Hand[0])
			cardIsInDestinationView := cardInView(tt.Insertion.InsertCard, &gameState.Views[0])
			//Check the gamestate that we passed in
			if tt.CardEndsInFirstView != cardIsInFirstView {
				t.Fatalf("Mismatch in First View checking. Expected: %t, Got %t", tt.CardEndsInFirstView, cardIsInFirstView)
			}

			if tt.CardEndsInDestinationView != cardIsInDestinationView {
				t.Fatalf("Mismatch in Destination View Checking. Expected %t, Got %t", tt.CardEndsInDestinationView, cardIsInDestinationView)
			}

			if len(changelog.Views) != tt.ExpectedChangelogLength {
				t.Fatalf("Wrong number of Views in Changelog. Expected: %d, Got: %d", tt.ExpectedChangelogLength, len(changelog.Views))
			}

			for _, view := range changelog.Views {
				if view.Id == toView.Id {
					if tt.CardEndsInDestinationView && !cardInView("cardToInsert", view) {
						t.Fatalf("Destination View in Changelog should have card, but doesn't")
					}
				} else if view.Id == fromView.Id {
					if tt.CardEndsInFirstView && !cardInView("cardToInsert", view) {
						t.Fatalf("First View in Changelog should have card, but doesn't")
					}
				}
			}
		})
	}
}

func TestWithdrawl_Execute(t *testing.T) {
	var tests = []struct {
		Name                      string
		Withdrawl                 Withdrawal
		ExpectedChangelogLength   int
		ShouldReturnError         bool
		CardEndsInFirstView       bool
		CardEndsInDestinationView bool
	}{
		{
			Name: "Valid Withdrawl",
			Withdrawl: Withdrawal{
				WithdrawCard:   "withdrawCard",
				FromCollection: "fromCollection",
				InView:         "inView",
				ToView:         "toView",
			},
			ExpectedChangelogLength:   2,
			ShouldReturnError:         false,
			CardEndsInFirstView:       false,
			CardEndsInDestinationView: true,
		},
		{
			Name: "Invalid Card",
			Withdrawl: Withdrawal{
				WithdrawCard:   "invalid",
				FromCollection: "fromCollection",
				InView:         "inView",
				ToView:         "toView",
			},
			ExpectedChangelogLength:   2,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
		{
			Name: "Invalid Collection",
			Withdrawl: Withdrawal{
				WithdrawCard:   "withdrawCard",
				FromCollection: "invalid",
				InView:         "inView",
				ToView:         "toView",
			},
			ExpectedChangelogLength:   2,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
		{
			Name: "Invalid First View",
			Withdrawl: Withdrawal{
				WithdrawCard:   "withdrawCard",
				FromCollection: "fromCollection",
				InView:         "invalid",
				ToView:         "toView",
			},
			ExpectedChangelogLength:   1,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
		{
			Name: "Invalid Destination",
			Withdrawl: Withdrawal{
				WithdrawCard:   "withdrawCard",
				FromCollection: "fromCollection",
				InView:         "inView",
				ToView:         "invalid",
			},
			ExpectedChangelogLength:   1,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			cardInQuestion := Pieces.Card{
				GamePiece: Pieces.GamePiece{
					Id: "withdrawCard",
				},
			}

			toView := Game.View{
				Id: "toView",
				Pieces: Pieces.PieceSet{
					Orphans: []Pieces.Card{},
				},
			}

			fromView := Game.View{
				Id: "inView",
				Pieces: Pieces.PieceSet{
					Decks: []Pieces.Deck{
						{
							GamePiece: Pieces.GamePiece{
								Id: "fromCollection",
								X:  50,
								Y:  50,
							},
							Cards: []Pieces.Card{cardInQuestion},
						},
					},
				},
			}

			me := Player.Player{
				Id: "me",
				Hand: []Game.View{
					toView,
				},
			}

			gameState := GameState{
				Views:   []Game.View{fromView},
				Players: []Player.Player{me},
			}

			changelog, err := tt.Withdrawl.Execute(&gameState, me.Id)

			//Check if we got an error when we shouldn't have, and vice versa
			if tt.ShouldReturnError != (err != nil) {
				t.Fatalf("ERROR in returned error value. Expected error: %t, err == %s", tt.ShouldReturnError, err)
			}

			cardIsInFirstView := cardInView("withdrawCard", &gameState.Views[0])
			cardIsInDestinationView := cardInView(tt.Withdrawl.WithdrawCard, &gameState.Players[0].Hand[0])
			//Check the gamestate that we passed in
			if tt.CardEndsInFirstView != cardIsInFirstView {
				t.Fatalf("Mismatch in First View checking. Expected: %t, Got %t", tt.CardEndsInFirstView, cardIsInFirstView)
			}

			if tt.CardEndsInDestinationView != cardIsInDestinationView {
				t.Fatalf("Mismatch in Destination View Checking. Expected %t, Got %t", tt.CardEndsInDestinationView, cardIsInDestinationView)
			}

			if len(changelog.Views) != tt.ExpectedChangelogLength {
				t.Fatalf("Wrong number of Views in Changelog. Expected: %d, Got: %d", tt.ExpectedChangelogLength, len(changelog.Views))
			}

			for _, view := range changelog.Views {
				if view.Id == toView.Id {
					if tt.CardEndsInDestinationView && !cardInView("withdrawCard", view) {
						t.Fatalf("Destination View in Changelog should have card, but doesn't")
					}
				} else if view.Id == fromView.Id {
					if tt.CardEndsInFirstView && !cardInView("withdrawCard", view) {
						t.Fatalf("First View in Changelog should have card, but doesn't")
					}
				}
			}
		})
	}
}

func TestMovement_Execute(t *testing.T) {
	var tests = []struct {
		Name                      string
		Movement                  Movement
		ExpectedChangelogLength   int
		ShouldReturnError         bool
		CardEndsInFirstView       bool
		CardEndsInDestinationView bool
	}{
		{
			Name: "Valid Movement",
			Movement: Movement{
				CardId:   "card",
				FromView: "fromView",
				ToView:   "toView",
				AtX:      100,
				AtY:      100,
			},
			ExpectedChangelogLength:   2,
			ShouldReturnError:         false,
			CardEndsInFirstView:       false,
			CardEndsInDestinationView: true,
		},
		{
			Name: "Invalid Card",
			Movement: Movement{
				CardId:   "invalid",
				FromView: "fromView",
				ToView:   "toView",
				AtX:      100,
				AtY:      100,
			},
			ExpectedChangelogLength:   2,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
		{
			Name: "Invalid From View",
			Movement: Movement{
				CardId:   "card",
				FromView: "invalid",
				ToView:   "toView",
				AtX:      100,
				AtY:      100,
			},
			ExpectedChangelogLength:   1,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
		{
			Name: "Invalid Destination",
			Movement: Movement{
				CardId:   "card",
				FromView: "fromView",
				ToView:   "invalid",
				AtX:      100,
				AtY:      100,
			},
			ExpectedChangelogLength:   1,
			ShouldReturnError:         true,
			CardEndsInFirstView:       true,
			CardEndsInDestinationView: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			cardInQuestion := Pieces.Card{
				GamePiece: Pieces.GamePiece{
					Id: "card",
				},
			}

			me := Player.Player{
				Id: "me",
			}

			fromView := Game.View{
				Id: "fromView",
				Pieces: Pieces.PieceSet{
					Orphans: []Pieces.Card{cardInQuestion},
				},
			}

			toView := Game.View{
				Id: "toView",
				Pieces: Pieces.PieceSet{
					Orphans: []Pieces.Card{},
				},
			}

			gameState := GameState{
				Players: []Player.Player{me},
				Views: []Game.View{
					fromView,
					toView,
				},
			}

			changelog, err := tt.Movement.Execute(&gameState, me.Id)

			//Check if we got an error when we shouldn't have, and vice versa
			if tt.ShouldReturnError != (err != nil) {
				t.Fatalf("ERROR in returned error value. Expected error: %t, err == %s", tt.ShouldReturnError, err)
			}

			cardIsInFirstView := cardInView("card", &gameState.Views[0])
			cardIsInDestinationView := cardInView(tt.Movement.CardId, &gameState.Views[1])
			//Check the gamestate that we passed in
			if tt.CardEndsInFirstView != cardIsInFirstView {
				t.Fatalf("Mismatch in First View checking. Expected: %t, Got %t", tt.CardEndsInFirstView, cardIsInFirstView)
			}

			if tt.CardEndsInDestinationView != cardIsInDestinationView {
				t.Fatalf("Mismatch in Destination View Checking. Expected %t, Got %t", tt.CardEndsInDestinationView, cardIsInDestinationView)
			}

			if len(changelog.Views) != tt.ExpectedChangelogLength {
				t.Fatalf("Wrong number of Views in Changelog. Expected: %d, Got: %d", tt.ExpectedChangelogLength, len(changelog.Views))
			}

			for _, view := range changelog.Views {
				if view.Id == toView.Id {
					if tt.CardEndsInDestinationView && !cardInView("card", view) {
						t.Fatalf("Destination View in Changelog should have card, but doesn't")
					}
				} else if view.Id == fromView.Id {
					if tt.CardEndsInFirstView && !cardInView("card", view) {
						t.Fatalf("First View in Changelog should have card, but doesn't")
					}
				}
			}
		})
	}
}

func TestEndTurn_Execute(t *testing.T) {

	var tests = []struct {
		name              string
		currentPlayer     string
		nextPlayer        string
		ShouldReturnError bool
	}{
		{
			name:              "Player 1, Unspecified Next Player",
			currentPlayer:     "player1",
			nextPlayer:        "",
			ShouldReturnError: false,
		},
		{
			name:              "Player 1, Next Player Player 3",
			currentPlayer:     "player1",
			nextPlayer:        "player3",
			ShouldReturnError: false,
		},
		{
			name:              "Player 4, Unspecified Next Player",
			currentPlayer:     "player4",
			nextPlayer:        "",
			ShouldReturnError: false,
		},
		{
			name:              "Player 4, Next Player Player 3",
			currentPlayer:     "player4",
			nextPlayer:        "player3",
			ShouldReturnError: false,
		},
		{
			name:              "Invalid Current Player",
			currentPlayer:     "invalid",
			nextPlayer:        "",
			ShouldReturnError: true,
		},
		{
			name:              "Invalid Next Player",
			currentPlayer:     "player1",
			nextPlayer:        "invalid",
			ShouldReturnError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player1 := Player.Player{
				Id: "player1",
			}

			player2 := Player.Player{
				Id: "player2",
			}

			player3 := Player.Player{
				Id: "player3",
			}

			player4 := Player.Player{
				Id: "player4",
			}

			gameState := GameState{
				Players:       []Player.Player{player1, player2, player3, player4},
				CurrentPlayer: tt.currentPlayer,
			}

			endTurn := EndTurn{
				NextPlayer: tt.nextPlayer,
			}

			changelog, err := endTurn.Execute(&gameState, tt.currentPlayer)

			//Check if we got an error when we shouldn't have, and vice versa
			if tt.ShouldReturnError != (err != nil) {
				t.Fatalf("ERROR in returned error value. Expected error: %t, err == %s", tt.ShouldReturnError, err)
			}

			if tt.ShouldReturnError {
				return
			}

			if tt.nextPlayer != "" {
				if changelog.CurrentPlayer != tt.nextPlayer {
					t.Fatalf("Error in changelog.CurrentPlayer. Expected: %s, Got %s", tt.nextPlayer, changelog.CurrentPlayer)
				}
			} else {
				indexOfCurrentPlayer := slices.IndexFunc(gameState.Players, func(p Player.Player) bool { return p.Id == tt.currentPlayer })
				indexOfNextPlayer := slices.IndexFunc(gameState.Players, func(p Player.Player) bool { return p.Id == changelog.CurrentPlayer })

				expectedNextIndex := indexOfCurrentPlayer + 1
				if expectedNextIndex >= len(gameState.Players) {
					expectedNextIndex -= len(gameState.Players)
				}

				if indexOfNextPlayer != expectedNextIndex {
					t.Fatalf("Error in expected next player. Expected to get index %d, got index %d", expectedNextIndex, indexOfNextPlayer)
				}
			}
		})
	}
}

func TestCardflip_Execute(t *testing.T) {
	var tests = []struct {
		name              string
		cardToFlip        string
		inView            string
		ShouldReturnError bool
	}{
		{
			name:              "Valid Flip",
			cardToFlip:        "card",
			inView:            "view",
			ShouldReturnError: false,
		},
		{
			name:              "Invalid Card",
			cardToFlip:        "invalid",
			inView:            "view",
			ShouldReturnError: true,
		},
		{
			name:              "Invalid View",
			cardToFlip:        "card",
			inView:            "invalid",
			ShouldReturnError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardInQuestion := Pieces.Card{
				GamePiece: Pieces.GamePiece{
					Id: "card",
				},
			}

			me := Player.Player{
				Id: "me",
			}

			inView := Game.View{
				Id: "view",
				Pieces: Pieces.PieceSet{
					Orphans: []Pieces.Card{cardInQuestion},
				},
			}

			gameState := GameState{
				Players: []Player.Player{me},
				Views: []Game.View{
					inView,
				},
			}

			turn := Cardflip{
				FlipCard: tt.cardToFlip,
				InView:   tt.inView,
			}

			changelog, err := turn.Execute(&gameState, me.Id)

			//Check if we got an error when we shouldn't have, and vice versa
			if tt.ShouldReturnError != (err != nil) {
				t.Fatalf("ERROR in returned error value. Expected error: %t, err == %s", tt.ShouldReturnError, err)
			}

			if tt.ShouldReturnError {
				return
			}

			//Check if the card is flipped in the Changelog
			changelogView := changelog.Views[0]

			if changelogView.Pieces.Orphans[0].Facedown != true {
				t.Fatalf("Card is not Facedown in Changelog!")
			}

			//Check if the card is flipped in the GameState
			gameStateView := gameState.Views[0]

			if gameStateView.Pieces.Orphans[0].Facedown != true {
				t.Fatalf("Card is not Facedown in GameState!")
			}
		})
	}
}

func cardInView(cardId string, view *Game.View) bool {
	for _, collection := range view.Pieces.GetCollections() {
		if collection.FindCardInCollection(cardId) != nil {
			return true
		}
	}
	for _, card := range view.Pieces.Orphans {
		if card.Id == cardId {
			return true
		}
	}
	return false
}
