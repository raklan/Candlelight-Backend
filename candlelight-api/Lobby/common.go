package Lobby

import "candlelight-models/Session"

// A message that is awaiting processing after being sent from a client. The client's raw message is put into [Message], while [PlayerId] is provided
// by the receiving function, detailing which player this message came from
type ClientMessage struct {
	PlayerId string
	Message  []byte
}

// The different types of messages the server might send to a client connected via websocket.
const (
	WebsocketMessage_Changelog = "Changelog"
	WebsocketMessage_Close     = "Close"
	WebsocketMessage_Error     = "Error"
	WebsocketMessage_GameState = "GameState"
	WebsocketMessage_LobbyInfo = "LobbyInfo"
)

// A message sent from the server to a client. The frontend can check [Type] to determine how to parse the object in [Data]
type WebsocketMessage struct {
	//One of the above constants. That constant will tell you which of the below structs is found in the [Data] field
	Type string `json:"type"`
	//One of the below structs, a Changelog, or a GameState. Its exact type is recorded in [Type]
	Data any `json:"data"`
}

// A message containing a Player's assigned ID and the details of the lobby after they've joined it, whether by hosting it or joining a pre-existing lobby.
// The frontend should store this PlayerID.
type LobbyInfo struct {
	PlayerID  string        `json:"playerID"`
	LobbyInfo Session.Lobby `json:"lobbyInfo"`
}

// I know SocketError and SocketClose have the same exact structure, but I've separated them for both clarity in the code and in case we end up wanting to put additional (unique) data in one or both

// If some message from a client causes any error, one of these is sent back to the client
type SocketError struct {
	Message string `json:"message"`
}

// If a connection is about to be closed by the server, it will send a SocketClose, followed by immediately closing the connection
type SocketClose struct {
	Message string `json:"message"`
}
