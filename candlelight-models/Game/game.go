package Game

import (
	"candlelight-models/Pieces"
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
	//Whether players should be able to see details about other players such as how many cards are in their hands
	ShowOtherPlayerDetails bool `json:"showOtherPlayerDetails"`
	//Resources this Game will use
	//Resources []GameResource `json:"resources"`
	//Views this Game will use
	Views []View `json:"views"`
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
	//An Id for this View. Used to fill in ParentView on all GamePieces belonging to this View
	Id string `json:"id"`
	//The PlayerNumber of the Owner of this view. 0 is a special, reserved number for the Game itself. Any
	//view with OwnerPlayerNumber == 0 is public and accessible by all Players
	OwnerPlayerNumber int `json:"ownerPlayerNumber"`
	//The PieceSet belonging to (and rendered within) this View
	Pieces Pieces.PieceSet `json:"pieces"`
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
