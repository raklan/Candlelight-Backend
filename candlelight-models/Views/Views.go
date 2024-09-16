package Views

import (
	"encoding/json"
	"math/rand"
)

// A list of constants to be used in the Type field of UI_Element.
// This tells the JSON parsing what struct to deserialize the Element field into
// just like for SubmittedActions
const (
	//General Types
	Type_View       = "View"
	Type_Navigation = "Navigation"

	//Types for Zones
	Type_Deck     = "Deck"
	Type_Space    = "Space"
	Type_CardZone = "CardZone"

	//Types for Pieces
	Type_Card   = "Card"
	Type_Meeple = "Meeple"
	Type_Die    = "Die"
)

// A parent struct for all common fields between the types of elements
type UI_Element struct {
	//One of the above constants. Tells the JSON parser what type of struct
	//is in the Element field
	Type string `json:"type"`
	//Whether the user should be allowed to interact with the data of this
	//UI_Element. For example, on a Card that they move around, this will be true. For
	//a View (or something else that's background-only) it will likely be false
	Interactable bool `json:"interactable"`
	//The Element itself. This will be one of the below structs
	Element json.RawMessage `json:"element"`
	//Position/Location data for this UI_Element. See Position struct
	Position Position `json:"position"`
	//A style object to be directly applied to this UI_Element. See Style struct
	Styling Style `json:"styling"`
}

// Position data for a UI_Element
type Position struct {
	//The X position of this UI_Element. If this UI_Element is a child
	//of some other UI_Element, this is relative to the parent
	X int `json:"x"`
	//The Y position of this UI_Element. If this UI_Element is a child
	//of some other UI_Element, this is relative to the parent
	Y int `json:"y"`
	//Width of this UI_Element
	Width int `json:"width"`
	//Height of this UI_Element
	Height int `json:"height"`
}

// A list of CSS rules to apply to a UI_Element. Super rough, but I figured I'd throw it in and see what everyone thinks
type Style struct {
	//A map of CSS rules, where the Key is the name of the rule and the Value is the rule.
	//For example: Rules["color"] = "red" is equivalent to the CSS {color: red}
	Rules map[string]string `json:"rules"`
}

//===========================================================================
/*All the possible Elements. One of these is placed within the Element field
of UI_Element*/
//===========================================================================

// A box that contains other UI_Elements, which acts as a "stage" for other UI_Elements. This is the only
// type allowed to not have a "parent"
type View struct {
	//A list of all children within this View. The X + Y positions
	//of all children are relative to this view. (i.e. If a child is at
	//[100,150] then it's absolute X + Y are going to be at [this View's X
	// + 100, this View's Y + 150] )
	Children []UI_Element `json:"children"`
	//Id. Self-explanatory and non-user-generated
	Id string `json:"id"`
}

// Returns "View" for determining what type of UI_Element you're looking at
func (v View) Type() string {
	return Type_View
}

// A link to a View
type Navigation struct {
	//The ID of the View this Navigation links to. Any interactions (clicking, dropping a piece, etc.)
	//will automatically fall through to this linked View
	TargetView string `json:"targetView"`
	//Id. Self-explanatory and non-user-generated
	Id string `json:"id"`
}

// Returns "Navigation" for determining what type of UI_Element you're looking at
func (n Navigation) Type() string {
	return Type_Navigation
}

//===========Zones============

// A spot that pieces can be contained in. This serves as the "parent class"
// for all Zone-type things, such as Decks, Spaces, and CardZones. Embed this
// struct in all such structs
type Zone struct {
	//Id. Self-explanatory and non-user-generated
	Id   string `json:"id"`
	Name string `json:"name"`
	//a dictionary of string -> array(string), where the key is the tag name and the value is a list of
	// values such that any piece with that tag key-value pair is allowed to be placed within collections of this
	// Zone. If this is empty, no requirements are enforced. If there's > 0 entries, only Pieces
	// with an approved value are allowed
	TagsWhitelist map[string][]string `json:"tagsWhitelist"`
}

// A container for cards only. Cards within a deck are NOT flipped until drawn
type Deck struct {
	Zone
	Cards []Card `json:"cards"`
}

// A container for Meeples. Unused, but will likely be useful if/when we get to board games
type Space struct {
	Zone
	Meeples []Meeple `json:"meeples"`
}

// A place for Cards to be placed, such as a play area or discard pile. Cards can be flipped or unflipped here
type CardZone struct {
	Zone
	Cards []Card `json:"cards"`
}

// All Zone types should implement these methods in piececontainer-implementation.go
// for their given Piece. (i.e. Decks are GamePiece[Card])
type PieceContainer[T Piece] interface {
	//Adds [pieceToAdd] to this Zone's piece collection.
	AddPiece(pieceToAdd T)
	//Removes any pieces from this Zone's collection with an ID matching [pieceToRemove] using DeleteFunc
	RemovePiece(pieceToRemove T)
	//Finds and returns the address of the first GamePiece in this Zone's collections with an ID matching [id]
	FindPiece(id string) (*T, error)
	//Selects and returns the address of a random GamePiece in this Zone's collection
	PickRandomPiece() *T
	//Returns the Type constant (see above) matching this Zone. Maybe useful?
	Type() string
}

//==================Pieces=================

// A piece that can (typically) be moved around and interacted with, such as a card or meeple. Serves
// as the "parent class" to all types of Pieces. Embed this struct in all such structs
type GamePiece struct {
	//Id. Self-explanatory and non-user-generated
	Id string `json:"id"`
	//A list of key-value pairs to serve as the tags of this piece. Conceptual use is to
	//filter what pieces can be put in a zone using TagsWhitelist of that Zone and each piece's Tags
	Tags map[string]string `json:"tags"`
}

// A game card. Cards have 2 sides, with the Flip() method serving as a way
// to transition between which side is showing. By default, Back is showing
type Card struct {
	GamePiece
	//Not sure what type to make these yet
	Back  string `json:"back"`
	Front string `json:"front"`
	//Whether [Front] is showing
	Flipped bool `json:"flipped"`
}

// Flips this Card.
func (c *Card) Flip() {
	c.Flipped = !c.Flipped
}

// Skeleton for the any future Meeple-type piece
type Meeple struct {
	GamePiece
	Shape string `json:"shape"`
}

// A Die with any number of sides. Use Roll() to roll this die
type Die struct {
	GamePiece
	NumSides int `json:"numSides"`
}

// Rolls the die, picking a random number between 1 and [NumSides] inclusive
func (d Die) Roll() int {
	return rand.Intn(d.NumSides) + 1
}

// A Piece for games. This interface exists mostly for the PieceContainer interface to be able
// to be generic. Each Piece type should implement these methods
type Piece interface {
	//Returns the ID of this Piece
	Id() string
	//Returns the Type constant matching this Piece type. Not sure if useful
	Type() string
	//Returns the Tags map of this Piece
	Tags() map[string]string
}
