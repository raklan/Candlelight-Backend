package Lobby

import "candlelight-models/Session"

// A message that is awaiting processing after being sent from a client. The client's raw message is put into [Message], while [PlayerId] is provided
// by the receiving function, detailing which player this message came from
type ClientMessage struct {
	PlayerId string
	Message  []byte
}

// A message containing a Player's assigned ID and the details of the lobby after they've joined it, whether by hosting it or joining a pre-existing lobby.
// The frontend should store this PlayerID.
type LobbyMessage struct {
	PlayerID  string        `json:"playerID"`
	LobbyInfo Session.Lobby `json:"lobbyInfo"`
}

// If some message from a client causes any error, one of these is sent back to the client
type SocketError struct {
	Message string `json:"message"`
}

// The different types of messages the server might send a client connected via websocket.
const (
	WebsocketMessage_Changelog = "Changelog"
	WebsocketMessage_Close     = "Close"
	WebsocketMessage_Error     = "Error"
	WebsocketMessage_GameState = "GameState"
	WebsocketMessage_LobbyInfo = "LobbyInfo"
)

// A message sent from the server to a client. The frontend can check [Type] to determine how to parse the object in [Data]
type WebsocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
