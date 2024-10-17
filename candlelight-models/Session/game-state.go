package Session

import (
	"candlelight-models/Game"
	"candlelight-models/Player"
	"encoding/json"
)

const (
	ActionType_Insertion = "Insertion"
	ActionType_Withdrawl = "Withdrawl"
	ActionType_Movement  = "Movement"
)

/*
This is intended to be the actual data the backend sends to the frontend to have it render things for the players. This is separate from the Game
definition found throughout the other packages
*/
type GameState struct {
	//This is solely for book-keeping. The front end should submit this Id along with SubmittedActions to update the GameState
	Id string `json:"id"`
	//The ID of the GameDefinition that this game state tracks
	GameDefinitionId string `json:"gameDefinitionId"`
	//The name of the GameDefinition that this game tracks. Added for rejoining players to be able to see the game's name
	GameName string `json:"gameName"`
	//A list of the states of each Player in the game.
	Players []Player.Player `json:"playerStates"`
	//The player whose turn it currently is
	//CurrentPlayer Player.Player `json:"currentPlayer"`
	//The pieces (and their locations) as they are currently
	Views []Game.View `json:"views"`
}

// A struct containing any and all Views that were affected by a SubmittedAction, as well
// as those objects' new states. One of these is generated and returned any time a Client submits an action,
// regardless of whether the action was successful.
type Changelog struct {
	Views []*Game.View `json:"views"`
}

// This is the way the frontend will send data to the backend during gameplay. They will
// send one of these objects, then the Rule Engine will take it, perform any updates to the
// internal model of the Game, then respond with a GameState
type SubmittedAction struct {
	//The type of Turn you want to take. Should match exactly with the name of one of the below structs (i.e. "Movement", "Insertion", etc)
	Type string `json:"type"`
	//The actual turn object. Should have all the fields within the struct that you're wanting
	Turn json.RawMessage `json:"turn"`
	//The player who is trying to submit this action
	Player Player.Player `json:"player"`
}

// An Insertion is defined as a Player inserting an Orphan into a Card Collection, whether that's a Deck or CardPlace
type Insertion struct {
	//The Id of the Card being inserted
	InsertCard string `json:"insertCard"`
	//The Id of the View which [InsertCard] is an Orphan of before the Insertion
	FromView string `json:"fromView"`
	//The Id of the Collection which [InsertCard] is to be inserted into
	ToCollection string `json:"toCollection"`
	//The Id of the View to which [ToCollection] belongs (and which [InsertCart] will belong to after the Insertion)
	InView string `json:"inView"`
}

// A Movement is defined as a Player moving an Orphan from one position to another, optionally between Views
type Movement struct {
	//The Id of the card being moved
	CardId string `json:"cardId"`
	//The View that [CardId] belongs to before moving
	FromView string `json:"fromView"`
	//The View that [CardId] is moving into. Can be the same as [FromView] if desired
	ToView string `json:"toView"`
	//The new X position that should be assigned to [CardId]
	AtX float32 `json:"toX"`
	//The new Y position that should be assigned to [CardId]
	AtY float32 `json:"toY"`
}

// A Withdrawl is defined as a Player moving a Card out of a Card Collection into the Orphans of a certain View
type Withdrawl struct {
	//The Id of the Card to Withdraw. Leave blank to be given an random card from [FromCollection]
	WithdrawCard string `json:"withdrawCard"`
	//The Collection a [WithdrawCard ]is to be withdrawn from.
	FromCollection string `json:"fromCollection"`
	//The View to which [FromCollection] belongs
	InView string `json:"inView"`
	//The View to which [WithdrawCard] should be moved into as an Orphan
	ToView string `json:"toView"`
}

// One of the possible Turn objects. This is solely for backend reference, and you should not have
// to ever think about this on the frontend
type Turn interface {
	Execute(gameState *GameState, player *Player.Player) (Changelog, error)
}
