package Lobby

import (
	"candlelight-models/Player"
	"candlelight-models/Session"
	"candlelight-ruleengine/Engine"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

// Mutex to control access to the clients map
// do PlayerId instead
// Outer key is the room code
// in key is the playId
// value is the websocket
var gamesClients = make(map[string]map[string]*websocket.Conn)
var gamesClientsMutex = sync.Mutex{}

// Map of each lobby to its respective message buffer
// Key is the lobby code, value is its ClientMessage channel
var messageBuffers = make(map[string]chan ClientMessage)
var messageBufferMutex = sync.Mutex{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  32768, // Setting read buffer size to 32 KB
	WriteBufferSize: 32768, // Setting write buffer size to 32 KB
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HostLobby creates the waiting lobby, joins on behalf of the given player, and upgrades the host into a websocket.
// The lobby (which contains the Room Code used for other people to join) is then passed back into the websocket.
func HostLobby(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting hostLobby")
	gameDefId := r.URL.Query().Get("gameId")
	playerName := r.URL.Query().Get("playerName")

	if gameDefId == "" || playerName == "" {
		log.Println("Missing gameId or playerName in request")
		http.Error(w, "Missing gameId or playerName in request", http.StatusBadRequest)
		return
	}

	lobbyCode, err := Engine.CreateRoom(gameDefId) // Assuming Engine.CreateRoom initializes room in DB
	if err != nil {
		http.Error(w, "Unable to create room", http.StatusInternalServerError)
		return
	}
	lobbyInfo, playerID, err := Engine.JoinRoom(lobbyCode, playerName)
	if err != nil {
		log.Printf("Error joining room: %v\n", err)
		http.Error(w, "Unable to join room", http.StatusInternalServerError)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	msg := LobbyMessage{
		PlayerID:  playerID,
		LobbyInfo: lobbyInfo,
	}
	conn.WriteJSON(msg)
	gamesClientsMutex.Lock()
	if _, exists := gamesClients[lobbyCode]; !exists {
		gamesClients[lobbyCode] = make(map[string]*websocket.Conn)
	}
	gamesClients[lobbyCode][playerID] = conn
	gamesClientsMutex.Unlock()

	go manageLobby(lobbyCode)
}

// Given a player and Room Code, tries to join that room for the given player. If joining was successul,
// the client's connection is upgraded to a websocket. Once complete, the client receives the lobby info
func HandleJoinLobby(w http.ResponseWriter, r *http.Request) {
	lobbyCode := r.URL.Query().Get("roomCode")
	playerName := r.URL.Query().Get("playerName")

	if lobbyCode == "" || playerName == "" {
		http.Error(w, "Please provide roomCode and playerName", http.StatusBadRequest)
		return
	}

	lobbyInfo, playerID, err := Engine.JoinRoom(lobbyCode, playerName)
	if err != nil {
		log.Printf("Error joining room: %v\n", err)
		http.Error(w, "Unable to join room", http.StatusNotFound)
		return
	}

	gamesClientsMutex.Lock()
	if _, exists := gamesClients[lobbyCode]; !exists {
		http.Error(w, "Lobby is not being tracked by server.", http.StatusInternalServerError)
		gamesClientsMutex.Unlock()
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	msg := LobbyMessage{
		PlayerID:  playerID,
		LobbyInfo: lobbyInfo,
	}
	conn.WriteJSON(msg)
	gamesClients[lobbyCode][playerID] = conn
	gamesClientsMutex.Unlock()
	go handShake(lobbyCode, playerID)
}

func HandleRejoinLobby(w http.ResponseWriter, r *http.Request) {
	lobbyCode := r.URL.Query().Get("roomCode")
	playerId := r.URL.Query().Get("playerId")

	if lobbyCode == "" || playerId == "" {
		http.Error(w, "Please send both roomCode and playerId", http.StatusBadRequest)
		return
	}

	lobbyInfo, err := Engine.LoadLobbyFromRedis(lobbyCode)
	if err != nil {
		http.Error(w, "Could not find requested lobby", http.StatusNotFound)
		return
	}

	//Make sure this player has joined the game before
	if !slices.ContainsFunc(lobbyInfo.Players, func(p Player.Player) bool { return p.Id == playerId }) {
		http.Error(w, "No player with given ID found in lobby", http.StatusNotFound)
		return
	}

	//Make sure that player does not already have an open connection
	gamesClientsMutex.Lock()
	if _, exists := gamesClients[lobbyCode][playerId]; exists {
		http.Error(w, "Found already open connection for player", http.StatusBadRequest)
		return
	}

	//Now we should know the player is allowed to rejoin. Upgrade to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	gamesClients[lobbyCode][playerId] = conn
	lobby := gamesClients[lobbyCode]
	gamesClientsMutex.Unlock()

	messageBufferMutex.Lock()
	buffer := messageBuffers[lobbyCode]
	messageBufferMutex.Unlock()

	//Check if game is started. If so, send the GameState instead of the lobby
	if lobbyInfo.GameStateId != "" {
		gameState, err := Engine.GetCachedGameStateFromRedis(lobbyInfo.GameStateId)
		if err != nil {
			conn.WriteJSON(SocketError{Message: err.Error()})
		}
		conn.WriteJSON(gameState)
	} else {
		msg := LobbyMessage{
			PlayerID:  playerId,
			LobbyInfo: lobbyInfo,
		}
		conn.WriteJSON(msg)
	}
	go awaitClientMessage(lobby, playerId, conn, buffer)
}

// handShake sends out the lobby info to everyone currently in the room, along with the
// name of the freshly joined player
func handShake(lobbyCode string, newPlayerId string) {
	gamesClientsMutex.Lock()
	lobby := gamesClients[lobbyCode]
	gamesClientsMutex.Unlock()
	jsonLobby, err := Engine.LoadLobbyFromRedis(lobbyCode)

	msg := LobbyMessage{
		PlayerID:  newPlayerId,
		LobbyInfo: jsonLobby,
	}
	for playerId, conn := range lobby {
		if err != nil {
			log.Printf("error connecting: {%s}", err)
		}
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Println("Error sending handshake, aborting connection ", playerId)
			gamesClientsMutex.Lock()
			lobby[playerId].Close()
			delete(lobby, playerId)
			gamesClientsMutex.Unlock()
			continue // I don't handle disconnection rn
		}
	}

	//Last thing we need to do is start listening for messages from this player
	messageBufferMutex.Lock()
	buffer := messageBuffers[lobbyCode]
	messageBufferMutex.Unlock()
	go awaitClientMessage(lobby, newPlayerId, lobby[newPlayerId], buffer)
}

// Manage all clients assosciated with one lobbyCode
func manageLobby(lobbyCode string) {
	gamesClientsMutex.Lock()
	lobby := gamesClients[lobbyCode]
	gamesClientsMutex.Unlock()

	//A buffer for incoming messages to be placed into
	messages := make(chan ClientMessage, len(lobby))

	messageBufferMutex.Lock()
	messageBuffers[lobbyCode] = messages
	messageBufferMutex.Unlock()

	//Set up routines to wait for client messages
	for playerId, conn := range lobby {
		go awaitClientMessage(lobby, playerId, conn, messages)
	}

	for {
		select {
		case data := <-messages:
			log.Printf("Received message: %s", data.Message)
			processMessage(lobbyCode, data.PlayerId, data.Message)
			//Since the only way we got here is by one of the clients sending something (and therefore completing awaitClientMessage for that client)
			//we need to spin off another goroutine to wait for this client's next message. However, if the connection was closed as a result of the last message,
			//we don't want to start listening again
			if lobby[data.PlayerId] != nil {
				go awaitClientMessage(lobby, data.PlayerId, lobby[data.PlayerId], messages)
			}
		default:
			continue
		}
		//TODO: Add appropriate delay or termination condition
	}
}

func awaitClientMessage(lobby map[string]*websocket.Conn, playerId string, conn *websocket.Conn, messageBuffer chan ClientMessage) {
	//Panic handling for when a client closes the connection on their end.
	defer socketRecovery(lobby, playerId)

	log.Printf("Starting listener for playerId %s", playerId)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		//Disconnect client by closing connection and removing player from lobby
		gamesClientsMutex.Lock()
		lobby[playerId].Close()
		delete(lobby, playerId)
		gamesClientsMutex.Unlock()
	}
	messageBuffer <- ClientMessage{PlayerId: playerId, Message: msg}
}

// This gets called on loop per lobby
func processMessage(roomCode string, playerId string, message []byte) {

	var msg struct {
		JsonType string          `json:"jsonType"`
		Data     json.RawMessage `json:"data"` //Raw message delays the parsing
	}
	json.Unmarshal(message, &msg)

	lobby := gamesClients[roomCode]

	switch msg.JsonType {
	case "startGame":
		game, err := Engine.GetInitialGameState(roomCode)
		if err != nil {
			log.Fatal("GAME NOT STARTED, ABORTING", err, playerId)
			return
		}

		sendMessageToAllPlayers(lobby, game)
	case "submitAction":
		var action struct {
			GameId string                  `json:"gameId"`
			Action Session.SubmittedAction `json:"action"`
		}
		if err := json.Unmarshal(msg.Data, &action); err != nil {
			log.Fatal("error decoding submitAction: {}", err, playerId)
		}

		changelog, err := Engine.SubmitAction(action.GameId, action.Action)
		if err != nil {
			log.Fatalf("error with submitAction: {%s}", err)
		}

		sendMessageToAllPlayers(lobby, changelog)
	case "leaveLobby":
		updatedLobby, err := endPlayerConnection(roomCode, playerId, lobby)

		if err != nil {
			socketError := SocketError{
				Message: err.Error(),
			}
			lobby[playerId].WriteJSON(socketError)
			break
		}

		sendMessageToAllPlayers(lobby, LobbyMessage{PlayerID: "", LobbyInfo: updatedLobby})
	case "kickPlayer":
		var action struct {
			PlayerToKick string `json:"playerToKick"`
		}
		if err := json.Unmarshal(msg.Data, &action); err != nil {
			log.Printf("Error trying to unmarshal kick request into struct with field 'playerToKick' ... Please ensure field exists in Data")
			socketError := SocketError{
				Message: "Message is malformed. Please ensure field 'playerToKick' is found in message object's 'Data' field!",
			}
			lobby[playerId].WriteJSON(socketError)
			break
		}

		dbLobby, err := Engine.LoadLobbyFromRedis(roomCode)

		if err != nil {
			log.Printf("Error trying to find lobby")
			socketError := SocketError{
				Message: "Could not find lobby. Something has gone terribly wrong",
			}
			lobby[playerId].WriteJSON(socketError)
			break
		}

		if dbLobby.Host.Id != playerId {
			socketError := SocketError{
				Message: "Player submitting kick request is not the host of the lobby!",
			}
			lobby[playerId].WriteJSON(socketError)
			break
		}

		updatedLobby, err := endPlayerConnection(roomCode, action.PlayerToKick, lobby)

		if err != nil {
			socketError := SocketError{
				Message: err.Error(),
			}
			lobby[playerId].WriteJSON(socketError)
			break
		}

		sendMessageToAllPlayers(lobby, LobbyMessage{PlayerID: "", LobbyInfo: updatedLobby})
	default:
		log.Println("Unknown type sent, ignoring message recieved", msg)
	}
}

func endPlayerConnection(roomCode string, playerId string, lobby map[string]*websocket.Conn) (Session.Lobby, error) {
	//Tell the engine to remove the player from the DB copy of the lobby
	updatedLobby, err := Engine.LeaveRoom(roomCode, playerId)
	if err != nil {
		return Session.Lobby{}, err
	}

	//Tell the client that the connection is closing, then close connection
	conn := lobby[playerId]
	msg := struct {
		Message string
	}{
		Message: "Player has been removed from Lobby. Closing connection",
	}
	conn.WriteJSON(msg)
	conn.Close()

	//Remove connection from lobby map so we don't try to send them any more messages
	delete(lobby, playerId)

	return updatedLobby, nil
}

func sendMessageToAllPlayers(lobby map[string]*websocket.Conn, message interface{}) {
	log.Printf("Sending message to every player: %s", message)
	for playerId, conn := range lobby {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Println("Error sending message, skipping meesage to ", playerId)
			continue // I don't handle disconnection rn
		}
	}
}

// Defer this function whenever you try to read from a socket. If ReadMessage panics, this will kick in. Note: This must be set up (deferred) **BEFORE** calling ReadMessage
func socketRecovery(lobby map[string]*websocket.Conn, playerId string) {
	if r := recover(); r != nil {
		log.Printf("Something went wrong trying to read from the connection of Player: {%s} -- %s", playerId, r)
		//Disconnect client by closing connection and removing player from lobby
		gamesClientsMutex.Lock()
		//lobby[playerId].Close() Assume the socket is already closed if we're here
		delete(lobby, playerId)
		gamesClientsMutex.Unlock()
	}
}
