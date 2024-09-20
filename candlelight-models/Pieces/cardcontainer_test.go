package Pieces

import "testing"

func TestAddCards(t *testing.T) {
	var tests = []struct {
		name      string
		container Card_Container
		cards     []Card
	}{
		{
			name:      "Deck Valid Card To Add",
			container: getEmptyDeck(),
			cards:     []Card{getDummyCard()},
		},
		{
			name:      "Deck Duplicate Card",
			container: getEmptyDeck(),
			cards:     []Card{getDummyCard(), getDummyCard()},
		},
		{
			name:      "CardPlace Valid Card To Add",
			container: getEmptyCardPlace(),
			cards:     []Card{getDummyCard()},
		},
		{
			name:      "CardPlace Duplicate Card",
			container: getEmptyCardPlace(),
			cards:     []Card{getDummyCard(), getDummyCard()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, card := range tt.cards {
				tt.container.AddCardToCollection(card)
			}

			//Coerce it into something specific so we can manually check its cards
			if c, ok := tt.container.(*Deck); ok {
				//Starting with empty container, so we can just compare the length of the
				//card container with the length of the cards we were supposed to add
				if len(c.Cards) != len(tt.cards) {
					t.Errorf("%s -- Wrong # of Cards. Expected {%d}, Got {%d}", tt.name, len(c.Cards), len(tt.cards))
				}

				for _, card := range c.Cards {
					if !cardIsInDeck(c, card) {
						t.Errorf("%s -- Couldn't find Card with ID == {%s} in Deck's Cards", tt.name, card.Id)
					}
				}
			} else if c, ok := tt.container.(*CardPlace); ok {
				//Starting with empty container, so we can just compare the length of the
				//card container with the length of the cards we were supposed to add
				if len(c.PlacedCards) != len(tt.cards) {
					t.Errorf("%s -- Wrong # of Cards. Expected {%d}, Got {%d}", tt.name, len(c.PlacedCards), len(tt.cards))
				}

				for _, card := range c.PlacedCards {
					if !cardIsInCardPlace(c, card) {
						t.Errorf("%s -- Couldn't find Card with ID == {%s} in CardPlace's Cards", tt.name, card.Id)
					}
				}
			} else {
				t.Errorf("%s -- Couldn't coerce container into known type. Something is very wrong", tt.name)
			}

		})
	}
}

func TestFindCard(t *testing.T) {
	var tests = []struct {
		name           string
		container      Card_Container
		card           Card
		shouldFindCard bool
	}{
		{
			name:           "Empty Deck",
			container:      getEmptyDeck(),
			card:           getDummyCard(),
			shouldFindCard: false,
		},
		{
			name:           "Full Deck Invalid Card",
			container:      getDeckWithCards(),
			card:           getDummyCard(),
			shouldFindCard: false,
		},
		{
			name:      "Full Deck Valid Card",
			container: getDeckWithCards(),
			card: Card{
				GamePiece: GamePiece{
					Id: "1",
				},
			},
			shouldFindCard: true,
		},
		{
			name:           "Empty CardPlace",
			container:      getEmptyCardPlace(),
			card:           getDummyCard(),
			shouldFindCard: false,
		},
		{
			name:           "Full CardPlace Invalid Card",
			container:      getCardPlaceWithCards(),
			card:           getDummyCard(),
			shouldFindCard: false,
		},
		{
			name:      "Full CardPlace Valid Card",
			container: getCardPlaceWithCards(),
			card: Card{
				GamePiece: GamePiece{
					Id: "1",
				},
			},
			shouldFindCard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := tt.container.FindCardInCollection(tt.card.Id)

			//Found != nil if card WAS found. Therefore, found != nil should == tt.shouldFindCard
			if (found != nil) != tt.shouldFindCard {
				t.Errorf("%s -- Unexpected value for the found card. Expected to find? {%t}, Did Find? {%t}", tt.name, tt.shouldFindCard, found != nil)
			}
		})
	}
}

func TestRemoveCard(t *testing.T) {
	var tests = []struct {
		name                string
		container           Card_Container
		card                Card
		cardShouldBeRemoved bool
	}{
		{
			name:                "Empty Deck",
			container:           getEmptyDeck(),
			card:                getDummyCard(),
			cardShouldBeRemoved: false,
		},
		{
			name:                "Full Deck Invalid Card",
			container:           getDeckWithCards(),
			card:                getDummyCard(),
			cardShouldBeRemoved: false,
		},
		{
			name:                "Full Deck Valid Card",
			container:           getDeckWithCards(),
			card:                Card{GamePiece: GamePiece{Id: "1"}},
			cardShouldBeRemoved: true,
		},
		{
			name:                "Empty CardPlace",
			container:           getEmptyCardPlace(),
			card:                getDummyCard(),
			cardShouldBeRemoved: false,
		},
		{
			name:                "Full CardPlace Invalid Card",
			container:           getCardPlaceWithCards(),
			card:                getDummyCard(),
			cardShouldBeRemoved: false,
		},
		{
			name:                "Full CardPlace Valid Card",
			container:           getCardPlaceWithCards(),
			card:                Card{GamePiece: GamePiece{Id: "1"}},
			cardShouldBeRemoved: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalCollectionLength := tt.container.CollectionLength()
			tt.container.RemoveCardFromCollection(tt.card)

			if c, ok := tt.container.(*Deck); ok {
				if tt.cardShouldBeRemoved {
					if len(c.Cards) != originalCollectionLength-1 {
						t.Errorf("%s -- Length of Deck's cards not right! Expected {%d}, Got {%d}", tt.name, originalCollectionLength-1, len(c.Cards))
					}

					if cardIsInDeck(c, tt.card) {
						t.Errorf("%s -- Card still present in Deck's Cards!", tt.name)
					}
				} else {
					if len(c.Cards) != originalCollectionLength {
						t.Errorf("%s -- Length of Deck's cards changed! Original Len {%d}, New Len {%d}", tt.name, originalCollectionLength, len(c.Cards))
					}
				}
			} else if c, ok := tt.container.(*CardPlace); ok {
				if tt.cardShouldBeRemoved {
					if len(c.PlacedCards) != originalCollectionLength-1 {
						t.Errorf("%s -- Length of Deck's cards not right! Expected {%d}, Got {%d}", tt.name, originalCollectionLength-1, len(c.PlacedCards))
					}

					if cardIsInCardPlace(c, tt.card) {
						t.Errorf("%s -- Card still present in CardPlace's Cards!", tt.name)
					}
				} else {
					if len(c.PlacedCards) != originalCollectionLength {
						t.Errorf("%s -- Length of CardPlace's cards changed! Original Len {%d}, New Len {%d}", tt.name, originalCollectionLength, len(c.PlacedCards))
					}
				}
			} else {
				t.Errorf("%s -- Couldn't coerce container into known type. Something is very wrong", tt.name)
			}
		})
	}
}

func TestRandomCard(t *testing.T) {
	var tests = []struct {
		name      string
		container Card_Container
	}{
		{
			name:      "Deck Test",
			container: getDeckWithCards(),
		},
		{
			name:      "CardPlace Test",
			container: getCardPlaceWithCards(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := tt.container.PickRandomCardFromCollection()

			if card == nil {
				t.Error("Should have gotten a card but didn't!")
			}
		})
	}
}

func TestCardIsAllowed(t *testing.T) {
	var tests = []struct {
		name            string
		container       Card_Container
		card            Card
		shouldBeAllowed bool
	}{
		{
			name: "Deck Valid Card",
			container: &Deck{
				GamePiece: GamePiece{
					Id:   "emptydeck",
					Name: "Empty Deck",
					Tags: map[string]string{},
				},
				PieceContainer: PieceContainer{
					TagsWhitelist: map[string][]string{
						"color": []string{"red"},
					},
				},
				Cards: []Card{},
			},
			card: Card{
				GamePiece: GamePiece{
					Id:   "dummycard",
					Name: "Dummy Card",
					Tags: map[string]string{
						"color": "red",
					},
				},
				Description: "A dummy card",
				Value:       0,
			},
			shouldBeAllowed: true,
		},
		{
			name: "Deck Invalid Card",
			container: &Deck{
				GamePiece: GamePiece{
					Id:   "emptydeck",
					Name: "Empty Deck",
					Tags: map[string]string{},
				},
				PieceContainer: PieceContainer{
					TagsWhitelist: map[string][]string{
						"color": []string{"red"},
					},
				},
				Cards: []Card{},
			},
			card: Card{
				GamePiece: GamePiece{
					Id:   "dummycard",
					Name: "Dummy Card",
					Tags: map[string]string{
						"color": "blue",
					},
				},
				Description: "A dummy card",
				Value:       0,
			},
			shouldBeAllowed: false,
		},
		{
			name: "CardPlace Valid Card",
			container: &CardPlace{
				GamePiece: GamePiece{
					Id:   "emptydeck",
					Name: "Empty Deck",
					Tags: map[string]string{},
				},
				PieceContainer: PieceContainer{
					TagsWhitelist: map[string][]string{
						"color": []string{"red"},
					},
				},
				PlacedCards: []Card{},
			},
			card: Card{
				GamePiece: GamePiece{
					Id:   "dummycard",
					Name: "Dummy Card",
					Tags: map[string]string{
						"color": "red",
					},
				},
				Description: "A dummy card",
				Value:       0,
			},
			shouldBeAllowed: true,
		},
		{
			name: "CardPlace Invalid Card",
			container: &CardPlace{
				GamePiece: GamePiece{
					Id:   "emptydeck",
					Name: "Empty Deck",
					Tags: map[string]string{},
				},
				PieceContainer: PieceContainer{
					TagsWhitelist: map[string][]string{
						"color": []string{"red"},
					},
				},
				PlacedCards: []Card{},
			},
			card: Card{
				GamePiece: GamePiece{
					Id:   "dummycard",
					Name: "Dummy Card",
					Tags: map[string]string{
						"color": "blue",
					},
				},
				Description: "A dummy card",
				Value:       0,
			},
			shouldBeAllowed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isAllowed := tt.container.CardIsAllowed(&tt.card)

			if tt.shouldBeAllowed != isAllowed {
				t.Fatalf("Mismatch in whether the card is allowed! Expected: %t, Got %t", tt.shouldBeAllowed, isAllowed)
			}

		})
	}
}

func TestCollectionLength(t *testing.T) {
	var tests = []struct {
		name           string
		container      Card_Container
		expectedLength int
	}{
		{
			name: "Deck With 1 Card",
			container: &Deck{
				Cards: []Card{getDummyCard()},
			},
			expectedLength: 1,
		},
		{
			name: "Deck With 3 Cards",
			container: &Deck{
				Cards: []Card{getDummyCard(), getDummyCard(), getDummyCard()},
			},
			expectedLength: 3,
		},
		{
			name: "CardPlace With 1 Card",
			container: &CardPlace{
				PlacedCards: []Card{getDummyCard()},
			},
			expectedLength: 1,
		},
		{
			name: "CardPlace With 3 Cards",
			container: &CardPlace{
				PlacedCards: []Card{getDummyCard(), getDummyCard(), getDummyCard()},
			},
			expectedLength: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.container.CollectionLength() != tt.expectedLength {
				t.Fatalf("Mismatch in collection length! Expected %d, Got %d", tt.expectedLength, tt.container.CollectionLength())
			}
		})
	}
}

func TestPieceSetCombine(t *testing.T) {
	ps1 := PieceSet{
		Decks:      []Deck{*getEmptyDeck()},
		CardPlaces: []CardPlace{*getEmptyCardPlace()},
	}
	ps2 := PieceSet{
		Decks:      []Deck{*getDeckWithCards()},
		CardPlaces: []CardPlace{*getCardPlaceWithCards()},
	}

	ps1.Combine(ps2)
}

//==================HELPER FUNCTIONS=================
func cardIsInDeck(deck *Deck, card Card) bool { //Have to manually implement these
	//Can't compare the structs directly right now, so just compare IDs
	for _, c := range deck.Cards {
		if card.Id == c.Id {
			return true
		}
	}
	return false
}

func cardIsInCardPlace(cardPlace *CardPlace, card Card) bool { //Have to manually implement these
	//Can't compare the structs directly right now, so just compare IDs
	for _, c := range cardPlace.PlacedCards {
		if card.Id == c.Id {
			return true
		}
	}
	return false
}

func getEmptyDeck() *Deck {
	return &Deck{
		GamePiece: GamePiece{
			Id:   "emptydeck",
			Name: "Empty Deck",
			Tags: map[string]string{},
		},
		PieceContainer: PieceContainer{
			TagsWhitelist: map[string][]string{},
		},
		Cards: []Card{},
	}
}

func getDeckWithCards() *Deck {
	return &Deck{
		GamePiece: GamePiece{
			Id:   "emptydeck",
			Name: "Empty Deck",
			Tags: map[string]string{},
		},
		PieceContainer: PieceContainer{
			TagsWhitelist: map[string][]string{},
		},
		Cards: []Card{
			{
				GamePiece: GamePiece{
					Id: "1",
				},
			},
			{
				GamePiece: GamePiece{
					Id: "2",
				},
			},
			{
				GamePiece: GamePiece{
					Id: "3",
				},
			},
			{
				GamePiece: GamePiece{
					Id: "4",
				},
			},
		},
	}
}

func getEmptyCardPlace() *CardPlace {
	return &CardPlace{
		GamePiece: GamePiece{
			Id:   "emptycardplace",
			Name: "Empty CardPlace",
			Tags: map[string]string{},
		},
		PieceContainer: PieceContainer{
			TagsWhitelist: map[string][]string{},
		},
		PlacedCards: []Card{},
	}
}

func getCardPlaceWithCards() *CardPlace {
	return &CardPlace{
		GamePiece: GamePiece{
			Id:   "emptycardplace",
			Name: "Empty CardPlace",
			Tags: map[string]string{},
		},
		PieceContainer: PieceContainer{
			TagsWhitelist: map[string][]string{},
		},
		PlacedCards: []Card{
			{
				GamePiece: GamePiece{
					Id: "1",
				},
			},
			{
				GamePiece: GamePiece{
					Id: "2",
				},
			},
			{
				GamePiece: GamePiece{
					Id: "3",
				},
			},
			{
				GamePiece: GamePiece{
					Id: "4",
				},
			},
		},
	}
}

func getDummyCard() Card {
	return Card{
		GamePiece: GamePiece{
			Id:   "dummycard",
			Name: "Dummy Card",
			Tags: map[string]string{},
		},
		Description: "A dummy card",
		Value:       0,
	}
}
