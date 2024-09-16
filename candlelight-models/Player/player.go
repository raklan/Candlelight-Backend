package Player

import "candlelight-models/Game"

// A Player within your game.
type Player struct {
	//Id just for book keeping
	Id string `json:"id"`
	//Display name of this Player. Should be shown to other players in-game
	Name string `json:"name"`
	//All the Pieces this Player currently has
	Hand []Game.View `json:"hand"`
	//All the Resources this Player currently has
	Resources []PlayerResource `json:"resources"`
}

// A resource that a Player currently possesses. Should have a name identical to
// a name in the list of GameResources to ensure that min/max values stay correct
type PlayerResource struct {
	//Should line up EXACTLY with the name of one of the GameResources in the Game object
	Name string `json:"name"`
	//How many of this resource the player currently has
	CurrentValue int `json:"currentValue"`
	//Maximum amount of this resource. Possibly redundant because it lines up with
	//a GameResource. Remove maybe?
	MaxValue int `json:"maxValue"`
}
