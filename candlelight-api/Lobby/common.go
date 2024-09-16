package Lobby

import "candlelight-models/Session"

type ClientMessage struct {
	PlayerId string
	Message  []byte
}

type LobbyMessage struct {
	PlayerID  string        `json:"playerID"`
	LobbyInfo Session.Lobby `json:"lobbyInfo"`
}

type SocketError struct {
	Message string `json:"message"`
}
