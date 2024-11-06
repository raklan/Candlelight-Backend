package Session

import (
	"candlelight-models/Game"
	"candlelight-models/Player"
	"encoding/json"
)

// Supported valued for SubmittedAction.Type. Make sure this matches up with the object you put
// in the Turn field
const (
	ActionType_Insertion  = "Insertion"
	ActionType_Withdrawal = "Withdrawal"
	ActionType_Movement   = "Movement"
	ActionType_EndTurn    = "EndTurn"
	ActionType_CardFlip   = "Cardflip"
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
	Players []Player.Player `json:"players"`
	//Id of the Player whose turn it currently is
	CurrentPlayer string `json:"currentPlayer"`
	//Set of rules Candlelight should use while running this game
	Rules Game.GameRules `json:"rules"`
	//The pieces (and their locations) as they are currently
	Views []Game.View `json:"views"`
}

// A struct containing any and all Views that could have been affected by a SubmittedAction, as well
// as those objects' new states. One of these is generated and returned any time a Client submits an action,
// regardless of whether the action was successful.
type Changelog struct {
	//Any views that were affected by the most recent SubmittedAction, represented as their current state after applying the Action
	Views []*Game.View `json:"views"`
	//Id of the Player whose turn it is after applying the most recent SubmittedAction
	CurrentPlayer string `json:"currentPlayer"`
	//A description of the action that just took place. Will be empty if the most recent SubmittedAction had no effect for any reason
	MostRecentAction string `json:"mostRecentAction"`
}

// This is the way the frontend will send data to the backend during gameplay. They will
// send one of these objects, then the Rule Engine will take it, perform any updates to the
// internal model of the Game, then respond with a Changelog
type SubmittedAction struct {
	//The type of Turn you want to take. Should match exactly with the name of one of the below structs (i.e. "Movement", "Insertion", etc)
	Type string `json:"type"`
	//The actual turn object. Should have all the fields within the struct that you're wanting
	Turn json.RawMessage `json:"turn"`
	//ID of the player who is trying to submit this action. This is now supplied by the backend
	PlayerId string `json:"playerId"`
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
	AtX float32 `json:"atX"`
	//The new Y position that should be assigned to [CardId]
	AtY float32 `json:"atY"`
}

// A Withdrawal is defined as a Player moving a Card out of a Card Collection into the Orphans of a certain View
type Withdrawal struct {
	//The Id of the Card to Withdraw. Leave blank to be given an random card from [FromCollection]
	WithdrawCard string `json:"withdrawCard"`
	//The Collection a [WithdrawCard ]is to be withdrawn from.
	FromCollection string `json:"fromCollection"`
	//The View to which [FromCollection] belongs
	InView string `json:"inView"`
	//The View to which [WithdrawCard] should be moved into as an Orphan
	ToView string `json:"toView"`
}

type EndTurn struct {
	//Optional string specifying the ID of the player who should be given the next turn. If not specified, the turn will pass to whichever player is next in
	//the gameState's player list, wrapping around in the event of the last player submitting an EndTurn
	NextPlayer string `json:"nextPlayer"`
}

type Cardflip struct {
	//Id of the card which should be flipped to its opposite side as that which is showing. This card must be an Orphan
	FlipCard string `json:"flipCard"`
	//Id of the View in which [FlipCard] can be found
	InView string `json:"inView"`
}

// One of the possible Turn objects. This is solely for backend reference, and you should not have
// to ever think about this on the frontend
type Turn interface {
	Execute(gameState *GameState, playerId string) (Changelog, error)
}
