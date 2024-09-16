package Game

import (
	"candlelight-models/Actions"
	"candlelight-models/Rules"
	"candlelight-models/Views"
)

/*Changelog
-Removed Rules in favor of piece-by-piece rules
*/

type GameV2 struct {
	Id             string             `json:"id"`
	Name           string             `json:"name"`
	Genre          string             `json:"genre"`
	MaxPlayers     int                `json:"maxPlayers"`
	Resources      []GameResource     `json:"resources"`
	UI_Elements    []Views.UI_Element `json:"uiElements"`
	Actions        Actions.ActionSet  `json:"actions"`
	Phases         []Rules.GamePhase  `json:"phases"`
	BeginningPhase string             `json:"beginningPhase"`
}
