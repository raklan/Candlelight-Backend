package Session

import (
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-models/Player"
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

			changelog, err := tt.Insertion.Execute(&gameState, &me)

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

			changelog, err := tt.Withdrawl.Execute(&gameState, &me)

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

			changelog, err := tt.Movement.Execute(&gameState, &me)

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
