package Pieces

// A collection of GamePieces
type PieceSet struct {
	Decks      []Deck      `json:"decks"`
	CardPlaces []CardPlace `json:"cardPlaces"`
	Orphans    []Card      `json:"orphans"`
}

//Copies all Decks/Cardplaces/Cards in the Orphans deck from [second] into the caller
func (ps *PieceSet) Combine(second PieceSet) {
	ps.Decks = append(ps.Decks, second.Decks...)
	ps.CardPlaces = append(ps.CardPlaces, second.CardPlaces...)
	ps.Orphans = append(ps.Orphans, second.Orphans...)
}

func (ps *PieceSet) GetCollections() []Card_Container {
	toReturn := []Card_Container{}
	for index := range ps.Decks {
		toReturn = append(toReturn, &ps.Decks[index])
	}
	for index := range ps.CardPlaces {
		toReturn = append(toReturn, &ps.CardPlaces[index])
	}
	return toReturn
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

	//The [R,G,B,A] values for this card. As such, this array should be exactly 4 entries, in the order described
	Color []float32 `json:"color"`
	//The [R,G,B,A] values for this card. As such, this array should be exactly 4 entries, in the order described
	PickColor []float32 `json:"pickColor"`
	//Text that appears on this piece
	Text string `json:"text"`

	//The X position of this Piece. This is relative to the parent view
	X float32 `json:"x"`
	//The Y position of this Piece. This is relative to the parent view
	Y float32 `json:"y"`

	//Id of the View this GamePiece belongs to. Used primarily to match up this piece to the correct view when it appears in a changelog during gameplay
	ParentView string `json:"parentView"`
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
	Cards []Card `json:"cards"`
}

//An interface for any Piece that contains cards. Currently *Deck and *CardPlace implement this.
//(Pointers used to ensure the methods called actually change the object, instead of a copy of it)
type Card_Container interface {
	GetId() string
	GetXY() (float32, float32)
	AddCardToCollection(cardToAdd Card)
	CardIsAllowed(card *Card) bool
	CollectionLength() int
	FindCardInCollection(cardId string) *Card
	PickRandomCardFromCollection() *Card
	RemoveCardFromCollection(cardToRemove Card)
}
