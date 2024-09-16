package Session

import (
	"candlelight-models/Actions"
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-models/Player"
	"candlelight-models/Rules"
	"encoding/json"
)

const (
	MoveTurnType        = "MoveTurn"
	PieceUpdateTurnType = "PieceUpdateTurn"
	PlacementTurnType   = "PlacementTurn"
	TakeTurnType        = "TakeTurn"
	TradeTurnType       = "TradeTurn"
	TransitionTurnType  = "TransitionTurn"
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
	//A list of the states of each Player in the game.
	PlayerStates []PlayerState `json:"playerStates"`
	//The player whose turn it currently is
	CurrentPlayer PlayerState `json:"currentPlayer"`
	//The current Phase the Game is in
	CurrentPhase Rules.GamePhase `json:"currentPhase"`
	//The pieces (and their locations) as they are currently
	Views []Game.View `json:"views"`
}

// The State of a certain Player. Right now, it just pairs a Player object with
// a list of actions that Player is currently allowed to take
type PlayerState struct {
	//The list of Actions the given Player make take at this time.
	//Any time you go through this, you should
	//filter it by the Constraints within the Game's Rules list
	AllowedActions Actions.ActionSet `json:"allowedActions"`
	//The Player for this PlayerState
	Player Player.Player `json:"player"`
}

// This is the way the frontend will send data to the backend during gameplay. They will
// send one of these objects, then the Rule Engine will take it, perform any updates to the
// internal model of the Game, then respond with a GameState
type SubmittedAction struct {
	//The type of Turn you want to take. Should match exactly with the name of one of the below structs (i.e. "MoveTurn", "PlacementTurn", etc)
	Type string `json:"type"`
	//The actual turn object. Should have all the fields within the struct that you're wanting
	Turn json.RawMessage `json:"turn"`
	//The player who is trying to submit this action
	Player Player.Player `json:"player"`
}

// A Move is defined as some piece moving from one place (in play) to another. For example, in Chess, you might
// take a Move turn to move one of your pieces.
type MoveTurn struct {
	//The Id of the Action you're trying to submit. This should be directly copied from the list of allowed actions for this player given in the GameState's PlayerState
	ActionId string `json:"actionId"`
	//The Id of the Piece to Move
	PieceId string `json:"pieceId"`
	//The place to which the Piece identified by [PieceId] should move TO
	TargetId string `json:"targetId"`
}

// A PieceUpdate writes a new set of Tags to the target piece. This might be used for something like Uno to change the
// suit of the deck after a wild card is played
type PieceUpdateTurn struct {
	//The Id of the Action you're trying to submit. This should be directly copied from the list of allowed actions for this player given in the GameState's PlayerState
	ActionId string `json:"actionId"`
	//The target GamePiece to which the updates should be applied
	TargetPieceId string `json:"pieceId"`
	//The new set of Tags to give to the GamePiece identified by [TargetPieceId]. ALL EXISTING TAGS WILL BE OVERWRITTEN BY THE GIVEN TAGS
	NewTags map[string]string `json:"newTags"`
}

// A placement is defined as a player moving a Piece from their hand to a Target, somewhere on the board. For example, in Uno, you might
// take a PlacementTurn to play a card from your hand.
type PlacementTurn struct {
	//The Id of the Action you're trying to submit. This should be directly copied from the list of allowed actions for this player given in the GameState's PlayerState
	ActionId string `json:"actionId"`
	//The Id of the Piece to play from a Player's hand. Obviously, it must be the Id of a piece they actually have in their hand
	PieceId string `json:"pieceId"`
	//The place the the Piece identified by [PieceId] should move to.
	TargetId string `json:"targetId"`
}

// A Take is defined as a player taking something from the board and putting it in their hand. For example, in Uno, you might
// take a TakeTurn to draw a card
type TakeTurn struct {
	//The Id of the Action you're trying to submit. This should be directly copied from the list of allowed actions for this player given in the GameState's PlayerState
	ActionId string `json:"actionId"`
	//The piece a Player wishes to take from the board. Leave empty if the player does not choose which piece specifically to take
	PieceId string `json:"pieceId"`
	//The place a Piece should be taken from. If [PieceId] is empty, a random Piece from the GamePiece identified by [TakingFromId] will be chosen.
	// If [PieceId] is not empty, it must identify a Piece within or equal to [TakingFromId]
	TakingFromId string `json:"takingFromId"`
}

// A Trade is defined as a change of Resources for one or more players. See the comments for a Transaction to see how this works
type TradeTurn struct {
	//The Id of the Action you're trying to submit. This should be directly copied from the list of allowed actions for this player given in the GameState's PlayerState
	ActionId string `json:"actionId"`
	//The list of all transactions for this Trade
	Transactions []Transaction `json:"transactions"`
}

// The definition for what players get what out of a Trade. Every thing that
// is included in this Trade should have its own Transaction. For example, if player 1 is giving
// something to player 2, there should be one transaction for player 1 (with a negative ChangeAmount) and
// one transaction for player 2 (with a positive ChangeAmount)
type Transaction struct {
	//Id of the player affected by this transaction
	PlayerId string `json:"playerId"`
	//Id of the Resource affected by this transaction
	ResourceId string `json:"resourceId"`
	//Amount to change the given resource by
	ChangeAmount int `json:"changeAmount"`
}

// A transition is defined as the Game moving to a new GamePhase. The applications here are fairly broad and depend on the
// defined GamePhases for the current Game
type TransitionTurn struct {
	//The Id of the Action you're trying to submit. This should be directly copied from the list of allowed actions for this player given in the GameState's PlayerState
	ActionId string `json:"actionId"`
	NewPhase string `json:"newPhase"`
}

type EndTurn struct {
}

// One of the possible Turn objects. This is solely for backend reference, and you should not have
// to ever think about this on the frontend
type Turn interface {
	Execute(gameState *GameState) (Pieces.PieceSet, error)
}
