package Views

import (
	"testing"
)

func Test_Deck_AddPiece(t *testing.T) {
	deck := Deck{
		Zone: Zone{
			Id:   "deck",
			Name: "test deck",
		},
		Cards: []Card{},
	}

	cardToAdd := Card{
		GamePiece: GamePiece{
			Id: "card",
		},
	}

	//Pre-flight checks
	if len(deck.Cards) != 0 {
		t.Errorf("Error: Cards not empty to start")
	}

	if cardInCollection(cardToAdd, deck.Cards) {
		t.Error("Card already in Deck")
	}

	//Execution
	deck.AddPiece(cardToAdd)

	//Landing checks
	if len(deck.Cards) != 1 {
		t.Errorf("Incorrect # Cards in Deck. Got %d, Expected: %d", len(deck.Cards), 1)
	}

	if !cardInCollection(cardToAdd, deck.Cards) {
		t.Error("Could not find Card in Deck")
	}
}

func Test_Deck_RemovePiece(t *testing.T) {
	//Test valid piece
	cardToRemove := Card{
		GamePiece: GamePiece{
			Id: "ExistingCard",
		},
	}

	deck := Deck{
		Zone: Zone{
			Id:   "deck",
			Name: "Test Deck",
		},
		Cards: []Card{
			cardToRemove,
		},
	}

	//Pre-flight checks
	if len(deck.Cards) != 1 {
		t.Error("Deck doesn't have 1 card before removing")
	}

	if !cardInCollection(cardToRemove, deck.Cards) {
		t.Error("Was not able to find cardToRemove in Deck. It should be there")
	}

	//Execution
	deck.RemovePiece(cardToRemove)

	//Landing Checks
	if len(deck.Cards) != 0 {
		t.Error("Deck still has cards in it")
	}

	if cardInCollection(cardToRemove, deck.Cards) {
		t.Error("Was able to find cardToRemove in Deck. It should not be there")
	}

	//Test invalid piece
	deck = Deck{
		Zone: Zone{
			Id:   "deck",
			Name: "Test Deck",
		},
		Cards: []Card{
			{
				GamePiece: GamePiece{
					Id: "randomId",
				},
			},
		},
	}

	//Pre-flight checks
	if len(deck.Cards) != 1 {
		t.Error("Deck doesn't have 1 card before removing")
	}

	if cardInCollection(cardToRemove, deck.Cards) {
		t.Error("Was able to find cardToRemove in Deck. It should not be there")
	}

	//Execution
	deck.RemovePiece(cardToRemove)

	//Landing Checks
	if len(deck.Cards) != 1 {
		t.Errorf("Deck has wrong number of cards. Expected %d, Got %d", 1, len(deck.Cards))
	}
}

func Test_Deck_FindPiece(t *testing.T) {
	cardToFind := Card{
		GamePiece: GamePiece{
			Id: "valid",
		},
	}

	deck := Deck{
		Cards: []Card{
			cardToFind,
		},
	}

	//Find card that exists
	foundCard, _ := deck.FindPiece("valid")

	if foundCard == nil {
		t.Error("Could not find card")
	}

	if foundCard.Id() != cardToFind.Id() {
		t.Errorf("Found the wrong id. Expected %s, Got %s", cardToFind.Id(), foundCard.Id())
	}

	//Find card that doesn't exist. Should be nil
	nextFoundCard, _ := deck.FindPiece("randomId")
	if nextFoundCard != nil {
		t.Error("Found a card, should not have found one")
	}
}

func Test_Space_AddPiece(t *testing.T) {
	space := Space{
		Zone: Zone{
			Id:   "space",
			Name: "test deck",
		},
		Meeples: []Meeple{},
	}

	meepleToAdd := Meeple{
		GamePiece: GamePiece{
			Id: "meeple",
		},
	}

	//Pre-flight checks
	if len(space.Meeples) != 0 {
		t.Errorf("Error: Pieces not empty to start")
	}

	if meepleInCollection(meepleToAdd, space.Meeples) {
		t.Error("Meeple already in Space")
	}

	//Execution
	space.AddPiece(meepleToAdd)

	//Landing checks
	if len(space.Meeples) != 1 {
		t.Errorf("Incorrect # Meeples in Space. Got %d, Expected: %d", len(space.Meeples), 1)
	}

	if !meepleInCollection(meepleToAdd, space.Meeples) {
		t.Error("Could not find Meeple in Space")
	}
}

func Test_Space_RemovePiece(t *testing.T) {
	//Test valid piece
	meepleToRemove := Meeple{
		GamePiece: GamePiece{
			Id: "ExistingCard",
		},
	}

	space := Space{
		Zone: Zone{
			Id:   "deck",
			Name: "Test Deck",
		},
		Meeples: []Meeple{
			meepleToRemove,
		},
	}

	//Pre-flight checks
	if len(space.Meeples) != 1 {
		t.Error("Space doesn't have 1 piece before removing")
	}

	if !meepleInCollection(meepleToRemove, space.Meeples) {
		t.Error("Was not able to find meepleToRemove in Space. It should be there")
	}

	//Execution
	space.RemovePiece(meepleToRemove)

	//Landing Checks
	if len(space.Meeples) != 0 {
		t.Error("Space still has pieces in it")
	}

	if meepleInCollection(meepleToRemove, space.Meeples) {
		t.Error("Was able to find meepleToRemove in Space. It should not be there")
	}

	//Test invalid piece
	space = Space{
		Zone: Zone{
			Id:   "deck",
			Name: "Test Deck",
		},
		Meeples: []Meeple{
			{
				GamePiece: GamePiece{
					Id: "randomId",
				},
			},
		},
	}

	//Pre-flight checks
	if len(space.Meeples) != 1 {
		t.Error("Space doesn't have 1 piece before removing")
	}

	if meepleInCollection(meepleToRemove, space.Meeples) {
		t.Error("Was able to find meepleToRemove in Space. It should not be there")
	}

	//Execution
	space.RemovePiece(meepleToRemove)

	//Landing Checks
	if len(space.Meeples) != 1 {
		t.Errorf("Space has wrong number of cards. Expected %d, Got %d", 1, len(space.Meeples))
	}
}

func Test_Space_FindPiece(t *testing.T) {
	meepleToFind := Meeple{
		GamePiece: GamePiece{
			Id: "valid",
		},
	}

	space := Space{
		Meeples: []Meeple{
			meepleToFind,
		},
	}

	//Find card that exists
	foundPiece, _ := space.FindPiece("valid")

	if foundPiece == nil {
		t.Error("Could not find piece")
	}

	if foundPiece.Id() != meepleToFind.Id() {
		t.Errorf("Found the wrong id. Expected %s, Got %s", meepleToFind.Id(), foundPiece.Id())
	}

	//Find card that doesn't exist. Should be nil
	nextFoundPiece, _ := space.FindPiece("randomId")
	if nextFoundPiece != nil {
		t.Error("Found a Piece, should not have found one")
	}
}

func Test_CardZone_AddPiece(t *testing.T) {
	cardZone := CardZone{
		Zone: Zone{
			Id:   "deck",
			Name: "test deck",
		},
		Cards: []Card{},
	}

	cardToAdd := Card{
		GamePiece: GamePiece{
			Id: "card",
		},
	}

	//Pre-flight checks
	if len(cardZone.Cards) != 0 {
		t.Errorf("Error: Cards not empty to start")
	}

	if cardInCollection(cardToAdd, cardZone.Cards) {
		t.Error("Card already in CardZone")
	}

	//Execution
	cardZone.AddPiece(cardToAdd)

	//Landing checks
	if len(cardZone.Cards) != 1 {
		t.Errorf("Incorrect # Cards in CardZone. Got %d, Expected: %d", len(cardZone.Cards), 1)
	}

	if !cardInCollection(cardToAdd, cardZone.Cards) {
		t.Error("Could not find Card in CardZone")
	}
}

func Test_CardZone_RemovePiece(t *testing.T) {
	//Test valid piece
	cardToRemove := Card{
		GamePiece: GamePiece{
			Id: "ExistingCard",
		},
	}

	cardZone := CardZone{
		Zone: Zone{
			Id:   "deck",
			Name: "Test Deck",
		},
		Cards: []Card{
			cardToRemove,
		},
	}

	//Pre-flight checks
	if len(cardZone.Cards) != 1 {
		t.Error("CardZone doesn't have 1 card before removing")
	}

	if !cardInCollection(cardToRemove, cardZone.Cards) {
		t.Error("Was not able to find cardToRemove in CardZone. It should be there")
	}

	//Execution
	cardZone.RemovePiece(cardToRemove)

	//Landing Checks
	if len(cardZone.Cards) != 0 {
		t.Error("CardZone still has cards in it")
	}

	if cardInCollection(cardToRemove, cardZone.Cards) {
		t.Error("Was able to find cardToRemove in CardZone. It should not be there")
	}

	//Test invalid piece
	cardZone = CardZone{
		Zone: Zone{
			Id:   "deck",
			Name: "Test Deck",
		},
		Cards: []Card{
			{
				GamePiece: GamePiece{
					Id: "randomId",
				},
			},
		},
	}

	//Pre-flight checks
	if len(cardZone.Cards) != 1 {
		t.Error("CardZone doesn't have 1 card before removing")
	}

	if cardInCollection(cardToRemove, cardZone.Cards) {
		t.Error("Was able to find cardToRemove in CardZone. It should not be there")
	}

	//Execution
	cardZone.RemovePiece(cardToRemove)

	//Landing Checks
	if len(cardZone.Cards) != 1 {
		t.Errorf("CardZone has wrong number of cards. Expected %d, Got %d", 1, len(cardZone.Cards))
	}
}

func Test_CardZone_FindPiece(t *testing.T) {
	cardToFind := Card{
		GamePiece: GamePiece{
			Id: "valid",
		},
	}

	cardZone := CardZone{
		Cards: []Card{
			cardToFind,
		},
	}

	//Find card that exists
	foundCard, _ := cardZone.FindPiece("valid")

	if foundCard == nil {
		t.Error("Could not find card")
	}

	if foundCard.Id() != cardToFind.Id() {
		t.Errorf("Found the wrong id. Expected %s, Got %s", cardToFind.Id(), foundCard.Id())
	}

	//Find card that doesn't exist. Should be nil
	nextFoundCard, _ := cardZone.FindPiece("randomId")
	if nextFoundCard != nil {
		t.Error("Found a card, should not have found one")
	}
}

// =================HELPER FUNCTIONS===================
func cardInCollection(card Card, collection []Card) bool {
	for _, c := range collection {
		if c.Id() == card.Id() {
			return true
		}
	}

	return false
}

func meepleInCollection(meeple Meeple, collection []Meeple) bool {
	for _, c := range collection {
		if c.Id() == meeple.Id() {
			return true
		}
	}

	return false
}
