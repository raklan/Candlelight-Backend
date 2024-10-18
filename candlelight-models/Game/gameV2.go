package Game

import (
	"candlelight-models/Views"
)

//If you're looking at this file looking for how GameDefs are saved, you're in the wrong spot. Ignore GameV2

type GameV2 struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Genre      string `json:"genre"`
	MaxPlayers int    `json:"maxPlayers"`
	//Resources      []GameResource     `json:"resources"`
	UI_Elements    []Views.UI_Element `json:"uiElements"`
	BeginningPhase string             `json:"beginningPhase"`
}
