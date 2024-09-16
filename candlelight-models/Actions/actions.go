package Actions

import "candlelight-models/Pieces"

// A set of Actions. Can be used to define what actions are available in a game, or what actions should be executed as a group, or more
type ActionSet struct {
	Moves        []Move        `json:"moves"`
	Trades       []Trade       `json:"trades"`
	Transitions  []Transition  `json:"transitions"`
	Takes        []Take        `json:"takes"`
	Placements   []Placement   `json:"placements"`
	PieceUpdates []PieceUpdate `json:"pieceUpdates"`
}

// Any Action that the Game or a Player may take. Serves as a base for actual Actions
// to build off of.
type Action struct {
	//Id for book keeping
	Id string `json:"id"`
	//Name for the Action
	Name string `json:"name"`
}

//This is used to Transition the Game to a new "phase" which is useful
//for advancing the game through what's allowed and what's not.
//To avoid the circular import, I will assume these only ever appear
//when you're working with an EventActionPair, and thus already have scope
//into the Rules package. Therefore, the NewPhase field will have to be
//an id to a GamePhase
type Transition struct {
	Action
	//The id of the Phase the game should transition to as a result of this Action
	NewPhase string `json:"newPhase"`
}

//A Trade is defined simply as a list of Transactions. The Transaction objects
//will contain the affected players/resources
type Trade struct {
	Action
	//All Resources allowed to be traded as part of this one Trade. Is an Id to avoid the cirular import
	AllowedResources []string `json:"allowedResources"`
	//The maximum amount this Trade is allowed to change any of the AllowedResources by
	MaximumAllowedChange int `json:"maximumAllowedChange"`
	//If TRUE, MaximumAllowedChange will actually be enforced. Useful if you don't actually want to define a maximum - you can just set this to false and keep the maximum at 0
	EnforceMaximumAllowedChange bool `json:"enforceMaximumAllowedChange"`
	//The minimum amount this Trade is allowed to change any AllowedResources by
	MinimumAllowedChange int `json:"minimumAllowedChange"`
	//If TRUE, MinimumAllowedChange will actually be enfored. Useful if you don't want to define an actual minimum, just set this to false and keep the minimum at 0
	EnforceMinimumAllowedChange bool `json:"enforceMinimumAllowedChange"`
}

//Defines a Player taking a GamePiece from somewhere to add it into their hand
type Take struct {
	Action
	Num int `json:"num"`
	//Allowed places a Player can take from
	TakeFrom Pieces.PieceSet `json:"takeFrom"`
	//If TRUE, the Player should not be allowed to choose where they draw from, instead taking each piece from a random target in TakeFrom
	TakeRandomly bool `json:"drawRandomly"`
}

//Defines a Player placing/playing something from their hand. Note that a Player is only allowed
// to Place something that's in their hand. To move something that's not currently in their hand, a Move should be used
type Placement struct {
	Action
	//Number of pieces to be placed in this Action
	Num int `json:"num"`
	//Allowed places a Player may Place something
	AllowedTargets Pieces.PieceSet `json:"allowedTargets"`
	//If TRUE, treat [AllowedTargets] as a Blacklist instead of a Whitelist
	FlipAllowedTargets bool `json:"flipAllowedTargets"`
	//The Pieces allowed to be used in this Action. (Note the piece must still be in the Player's hand)
	AllowedPieces Pieces.PieceSet `json:"allowedPieces"`
	//If TRUE, treat [AllowedPieces] as a blacklist instead of a whitelist
	FlipAllowedPieces bool `json:"flipAllowedPieces"`
	//If TRUE, the Player should not be allowed to pick where each Piece is played, instead choosing randomly from Targets
	PlaceRandomly bool `json:"placeRandomly"`
}

// This Action is used to update the Tags field found on all GamePieces. This is useful in
// scenarios where a player might have a choice to alter the properties of something that's
// in their hand or has been played, such as in Uno when you play a Wild card and can choose its color
// NOTE: The TagUpdates field will OVERWRITE all current Tags on the GamePiece
type PieceUpdate struct {
	Action
	TargetPieces Pieces.PieceSet   `json:"targetPieces"`
	TagUpdates   map[string]string `json:"tagNames"`
}

//This Action is used for Moving things to other things. This might be moving
//a card between decks or pieces on a game board, for example
type Move struct {
	Action
	//Move up to NumPieces GamePieces
	NumPieces int `json:"numPieces"`
	//Pieces this Action is allowed to move
	AllowedPieces Pieces.PieceSet `json:"allowedPieces"`
	//Pieces this Action is allowed to move pieces TO
	AllowedTargets Pieces.PieceSet `json:"allowedTargets"`
	//Treat AllowedPieces as a blacklist instead of a whitelist
	FlipAllowedPieces bool `json:"flipAllowedPieces"`
	//Treat AllowedTargets as a blacklist instead of a whitelist
	FlipAllowedTargets bool `json:"flipAllowedTargets"`
}
