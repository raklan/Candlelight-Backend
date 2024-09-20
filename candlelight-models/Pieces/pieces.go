package Pieces

// A collection of GamePieces
type PieceSet struct {
	Decks      []Deck      `json:"decks"`
	CardPlaces []CardPlace `json:"cardPlaces"`
	Orphans    Deck        `json:"orphans"`
}

//Copies all Decks/Cardplaces/Cards in the Orphans deck from [second] into the caller
func (ps *PieceSet) Combine(second PieceSet) {
	ps.Decks = append(ps.Decks, second.Decks...)
	ps.CardPlaces = append(ps.CardPlaces, second.CardPlaces...)
	ps.Orphans.Cards = append(ps.Orphans.Cards, second.Orphans.Cards...)
}

// An outline for any piece the Game might use.
type GamePiece struct {
	//Id for book keeping
	Id string `json:"id"`
	//Name of this piece
	Name string `json:"name"`
	//A set of Player-defined properties for this Piece. Should have the form
	// property: value in JSON
	Tags map[string]string `json:"tags"`
	//Position data for this GamePiece. See Position struct
	Position Position `json:"position"`
	//Style data for this GamePiece. See Style struct
	Style Style `json:"style"`
	//Id of the View this GamePiece belongs to. Used primarily to match up this piece to the correct view when it appears in a changelog during gameplay
	ParentView string `json:"parentView"`
}

// Position data for a Piece
type Position struct {
	//The X position of this Piece. If this Piece is a child
	//of some other Piece, this is relative to the parent
	X float32 `json:"x"`
	//The Y position of this Piece. If this Piece is a child
	//of some other Piece, this is relative to the parent
	Y float32 `json:"y"`
}

// A list of CSS rules to apply to a Piece. Super rough, but I figured I'd throw it in and see what everyone thinks
type Style struct {
	//A map of CSS rules, where the Key is the name of the rule and the Value is the rule.
	//For example: Rules["color"] = "red" is equivalent to the CSS {color: red}
	Rules map[string]string `json:"rules"`
	//The [R,G,B,A] values for this card. As such, this array should be exactly 4 entries, in the order described
	Color []int `json:"color"`
}

// An outline for pieces that can contain other GamePieces. Whitelist is a dictionary
// of string -> array(string), where the key is the tag name and the value is a list of
// values such that any piece with that tag key-value pair is allowed to be placed within collections of this
// PieceContainer
type PieceContainer struct {
	//a dictionary of string -> array(string), where the key is the tag name and the value is a list of
	// values such that any piece with that tag key-value pair is allowed to be placed within collections of this
	// PieceContainer
	TagsWhitelist map[string][]string `json:"tagsWhitelist"`
}

// A deck simply serves to keep a collection of cards in one place.
type Deck struct {
	GamePiece
	PieceContainer
	//The cards in the deck
	Cards []Card `json:"cards"`
}

// A card. Hopefully if you're reading this code, you know what a card might
// be used for in a tabletop game.
type Card struct {
	GamePiece
	//Optional description
	Description string `json:"description"`
	//Optional value, if having/playing a certain card might be good/bad
	Value int `json:"value"`
}

/*
A place where a player can play their
cards. This might be shared between all players  (e.g. Uno)
or might be owned by one specific player (e.g. Cover Your Assets)
*/
type CardPlace struct {
	GamePiece
	PieceContainer
	//Cards currently in this CardPlace
	PlacedCards []Card `json:"placedCards"`
}

//An interface for any Piece that contains cards. Currently *Deck and *CardPlace implement this.
//(Pointers used to ensure the methods called actually change the object, instead of a copy of it)
type Card_Container interface {
	AddCardToCollection(cardToAdd Card)
	CardIsAllowed(card *Card) bool
	CollectionLength() int
	FindCardInCollection(cardId string) *Card
	PickRandomCardFromCollection() *Card
	RemoveCardFromCollection(cardToRemove Card)
}
