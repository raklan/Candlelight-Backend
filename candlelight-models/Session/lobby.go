package Session

import "candlelight-models/Player"

const (
	LobbyStatus_AwaitingStart = "Awaiting Start"
	LobbyStatus_InProgress    = "In Progress"
	LobbyStatus_Ended         = "Game Ended"
)

//A lobby is a collection of players waiting for a game to start. This is created by the /createRoom endpoint
//which returns the room code. From there, you can pass the room code to /joinRoom which will put you in the room
//and return the state of the lobby
type Lobby struct {
	//The Room Code used for players to join this lobby. Passed to /joinRoom
	RoomCode string `json:"roomCode"`
	//The GameDefinition this Lobby is going to play
	GameDefinitionId string `json:"gameDefinitionId"`
	//The GameState created from this lobby. This field is empty until the game is started,
	//at which point the API server will fill it in.
	GameStateId string `json:"gameStateId"`
	//The status of the game, really only used for the backend to determine whether a game has started/ended. Will be one of the above constants
	Status string `json:"status"`
	//Name of the Game being played
	GameName string `json:"gameName"`
	//Current number of players in the lobby
	NumPlayers int `json:"numPlayers"`
	//Maximum allowed players in the lobby. Determined from the GameDef's MaxPlayers
	MaxPlayers int `json:"maxPlayers"`
	//Current list of joined players
	Players []Player.Player `json:"players"`
	//The player that created the Lobby
	Host Player.Player `json:"host"`
}
