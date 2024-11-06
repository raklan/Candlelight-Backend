package Sparks

//A collection of all Sparks available to a user
type Sparks struct {
	Dealer  Dealer  `json:"dealer"`
	Flipper Flipper `json:"flipper"`
}

//A Dealer will move [NumToDeal] random cards from [DeckToUse] to each player in the game
type Dealer struct {
	//Whether to perform this Spark
	Enabled bool `json:"enabled"`
	//Number of random cards to deal to each player
	NumToDeal int `json:"numToDeal"`
	//Id of the Deck from which each card should come from
	DeckToUse string `json:"deckToUse"`
}

type Flipper struct {
	//Whether to perform this Spark
	Enabled bool `json:"enabled"`
	//How many cards to move from [DeckToUse] to [CardPlaceToUse]
	NumToFlip int `json:"numToFlip"`
	//The Deck to take cards from
	DeckToUse string `json:"deckToUse"`
	//The CardPlace to put the cards in
	CardPlaceToUse string `json:"cardPlaceToUse"`
}
