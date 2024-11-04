package Sparks

//A collection of all Sparks available to a user
type Sparks struct {
	Dealer Dealer `json:"dealer"`
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
