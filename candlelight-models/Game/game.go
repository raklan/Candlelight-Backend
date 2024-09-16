package Game

import (
	"candlelight-models/Actions"
	"candlelight-models/Pieces"
	"candlelight-models/Rules"
)

// The over-arching definition of a Game. Should contain everything needed for the
// Rule Engine to refer to when running the game
type Game struct {
	//Id for book-keeping purposes
	Id string `json:"id"`
	//User-defined name for this Game
	Name string `json:"name"`
	//User-defined Genre for this Game. Maybe remove?
	Genre string `json:"genre"`
	//The username of the user that created this game
	Author string `json:"author"`
	//Max number of allowed players in this Game
	MaxPlayers int `json:"maxPlayers"`
	//Resources this Game will use
	Resources []GameResource `json:"resources"`
	//Views this Game will use
	Views []View `json:"views"`
	//Actions this Game will use
	Actions Actions.ActionSet `json:"actions"`
	//Phases this Game will use
	Phases []Rules.GamePhase `json:"phases"`
	//Which GamePhase the game should start in
	BeginningPhase Rules.GamePhase `json:"beginningPhase"`
}

// A Resource that the Game will use/keep track of for every player
type GameResource struct {
	//Id for book-keeping
	Id string `json:"id"`
	//Name of the Resource
	Name string `json:"name"`
	//Optional description
	Description string `json:"description"`
	//Value that all Players should start with
	InitialValue int `json:"initialValue"`
	//Maximum allowed value for a Player to have
	MaxValue int `json:"maxValue"`
	//Minimum allowed value for a Player to have
	MinValue int `json:"minValue"`
}

// A collection of Pieces to display to a player.
type View struct {
	Id                string          `json:"id"`
	OwnerPlayerNumber int             `json:"ownerPlayerNumber"`
	Pieces            Pieces.PieceSet `json:"pieces"`
}

func (game Game) ViewsForPlayer(playerNum int) []View {
	toReturn := []View{}
	for _, view := range game.Views {
		if view.OwnerPlayerNumber == playerNum {
			toReturn = append(toReturn, view)
		}
	}
	return toReturn
}

func (game Game) PiecesForPlayer(playerNum int) Pieces.PieceSet {
	toReturn := Pieces.PieceSet{}
	for _, view := range game.Views {
		if view.OwnerPlayerNumber == playerNum {
			toReturn.Combine(view.Pieces)
		}
	}
	return toReturn
}
