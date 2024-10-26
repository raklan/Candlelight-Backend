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
	msg := WebsocketMessage{
		Type: WebsocketMessage_LobbyInfo,
		Data: LobbyInfo{
			PlayerID:  playerID,
			LobbyInfo: lobbyInfo,
		},
	}
	conn.WriteJSON(msg)
	gamesClientsMutex.Lock()
	if _, exists := gamesClients[lobbyCode]; !exists {
		gamesClients[lobbyCode] = make(map[string]*websocket.Conn)
	}
	gamesClients[lobbyCode][playerID] = conn
	gamesClientsMutex.Unlock()

	go manageClient(lobbyCode, gamesClients[lobbyCode], playerID, conn)
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

	_, playerID, err := Engine.JoinRoom(lobbyCode, playerName)
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

	// msg := WebsocketMessage{
	// 	Type: WebsocketMessage_LobbyInfo,
	// 	Data: LobbyInfo{
	// 		PlayerID:  playerID,
	// 		LobbyInfo: lobbyInfo,
	// 	},
	// }
	//conn.WriteJSON(msg)
	gamesClients[lobbyCode][playerID] = conn
	gamesClientsMutex.Unlock()
	go handShake(lobbyCode, playerID)
}

func HandleRejoinLobby(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to rejoin a lobby!")

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
	log.Printf("Making sure player {%s} has joined this game before", playerId)
	if !slices.ContainsFunc(lobbyInfo.Players, func(p Player.Player) bool { return p.Id == playerId }) {
		http.Error(w, "No player with given ID found in lobby", http.StatusNotFound)
		return
	}

	//Make sure that player does not already have an open connection
	log.Printf("Making sure player {%s} does not already have an open connection", playerId)
	gamesClientsMutex.Lock()
	if _, exists := gamesClients[lobbyCode][playerId]; exists {
		http.Error(w, "Found already open connection for player", http.StatusBadRequest)
		gamesClientsMutex.Unlock()
		return
	}
	gamesClientsMutex.Unlock()

	//Now we should know the player is allowed to rejoin. Upgrade to websocket
	log.Printf("Player {%s} is allowed to rejoin. Upgrading connection to websocket.", playerId)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	log.Println("Upgrade finished. Checking status of lobby to give accurate first message...")

	//Check if game is started. If so, send the GameState instead of the lobby
	if lobbyInfo.Status == Session.LobbyStatus_InProgress {
		log.Println("Game has started. Sending GameState")
		gameState, err := Engine.GetCachedGameStateFromRedis(lobbyInfo.GameStateId)
		if err != nil {
			conn.WriteJSON(WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{Message: err.Error()},
			})
		}
		conn.WriteJSON(WebsocketMessage{Type: WebsocketMessage_GameState, Data: gameState})
	} else if lobbyInfo.Status == Session.LobbyStatus_AwaitingStart {
		log.Println("Game has not started yet. Sending LobbyInfo")
		msg := WebsocketMessage{
			Type: WebsocketMessage_LobbyInfo,
			Data: LobbyInfo{
				PlayerID:  playerId,
				LobbyInfo: lobbyInfo,
			},
		}
		conn.WriteJSON(msg)
	} else { //I don't love how I handle this case, upgrading the socket only to immediately close it feels bad
		log.Println("Game has already ended! Player cannot join!")
		msg := WebsocketMessage{
			Type: WebsocketMessage_Error,
			Data: SocketError{
				Message: "Game has already ended. Cannot rejoin",
			},
		}
		conn.WriteJSON(msg)
		conn.Close()
		return
	}

	log.Printf("Player {%s} has been given first message. Beginning to track connection for further communication...", playerId)

	gamesClientsMutex.Lock()
	gamesClients[lobbyCode][playerId] = conn
	lobby := gamesClients[lobbyCode]
	gamesClientsMutex.Unlock()

	go manageClient(lobbyCode, lobby, playerId, conn)
}

// handShake sends out the lobby info to everyone currently in the room, along with the
// name of the freshly joined player
func handShake(lobbyCode string, newPlayerId string) {
	gamesClientsMutex.Lock()
	lobby := gamesClients[lobbyCode]
	gamesClientsMutex.Unlock()
	jsonLobby, err := Engine.LoadLobbyFromRedis(lobbyCode)

	msg := WebsocketMessage{
		Type: WebsocketMessage_LobbyInfo,
		Data: LobbyInfo{
			PlayerID:  newPlayerId,
			LobbyInfo: jsonLobby,
		},
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
	go manageClient(lobbyCode, lobby, newPlayerId, lobby[newPlayerId])
}

func manageClient(lobbyCode string, lobbyMap map[string]*websocket.Conn, playerId string, conn *websocket.Conn) {
	defer socketRecovery(lobbyMap, playerId)

	log.Printf("Managing Connection for playerId %s and waiting for message", playerId)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		//Disconnect client by closing connection and removing player from lobby
		gamesClientsMutex.Lock()
		lobbyMap[playerId].Close()
		delete(lobbyMap, playerId)
		gamesClientsMutex.Unlock()
	}

	processMessage(lobbyCode, playerId, msg)
}

// This gets called on loop per lobby
func processMessage(roomCode string, playerId string, message []byte) {

	var msg struct {
		JsonType string          `json:"jsonType"`
		Data     json.RawMessage `json:"data"` //Raw message delays the parsing
	}
	json.Unmarshal(message, &msg)

	gamesClientsMutex.Lock()
	lobby := gamesClients[roomCode]
	gamesClientsMutex.Unlock()

	switch msg.JsonType {
	case "startGame":
		game, err := Engine.GetInitialGameState(roomCode)
		if err != nil {
			log.Printf("ERROR: GAME NOT STARTED, ABORTING...%s", err)
			return
		}

		sendMessageToAllPlayers(lobby, WebsocketMessage{Type: WebsocketMessage_GameState, Data: game})
	case "endGame":
		err := Engine.EndGame(roomCode, playerId)
		if err != nil {
			log.Printf("ERROR: Trying to end game...%s", err)
			return
		}
		closeAllConnections(lobby)
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

		sendMessageToAllPlayers(lobby, WebsocketMessage{Type: WebsocketMessage_Changelog, Data: changelog})
	case "leaveLobby":
		updatedLobby, err := endPlayerConnection(roomCode, playerId, lobby)

		if err != nil {
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: err.Error(),
				},
			}
			lobby[playerId].WriteJSON(socketError)
			break
		}

		sendMessageToAllPlayers(lobby, WebsocketMessage{Type: WebsocketMessage_LobbyInfo, Data: LobbyInfo{PlayerID: "", LobbyInfo: updatedLobby}})
	case "kickPlayer":
		var action struct {
			PlayerToKick string `json:"playerToKick"`
		}
		if err := json.Unmarshal(msg.Data, &action); err != nil {
			log.Printf("Error trying to unmarshal kick request into struct with field 'playerToKick' ... Please ensure field exists in Data")
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: "Message is malformed. Please ensure field 'playerToKick' is found in message object's 'Data' field!",
				},
			}
			lobby[playerId].WriteJSON(socketError)
			break
		}

		dbLobby, err := Engine.LoadLobbyFromRedis(roomCode)

		if err != nil {
			log.Printf("Error trying to find lobby")
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: "Could not find lobby. Something has gone terribly wrong",
				},
			}
			lobby[playerId].WriteJSON(socketError)
			break
		}

		if dbLobby.Host.Id != playerId {
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: "Player submitting kick request is not the host of the lobby!",
				},
			}
			lobby[playerId].WriteJSON(socketError)
			break
		}

		updatedLobby, err := endPlayerConnection(roomCode, action.PlayerToKick, lobby)

		if err != nil {
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: err.Error(),
				},
			}
			lobby[playerId].WriteJSON(socketError)
			break
		}

		sendMessageToAllPlayers(lobby, WebsocketMessage{Type: WebsocketMessage_LobbyInfo, Data: LobbyInfo{PlayerID: "", LobbyInfo: updatedLobby}})
	default:
		log.Println("Unknown type sent, ignoring message recieved", msg)
	}

	//Listen for the next message from this client. Not using `go manageClient(...)` because this is already happening in a goroutine
	manageClient(roomCode, lobby, playerId, lobby[playerId])
}

func endPlayerConnection(roomCode string, playerId string, lobby map[string]*websocket.Conn) (Session.Lobby, error) {
	//Tell the engine to remove the player from the DB copy of the lobby
	updatedLobby, err := Engine.LeaveRoom(roomCode, playerId)
	if err != nil {
		return Session.Lobby{}, err
	}

	//Tell the client that the connection is closing, then close connection
	conn := lobby[playerId]
	msg := WebsocketMessage{
		Type: WebsocketMessage_Close,
		Data: SocketClose{
			Message: "Player has been removed from Lobby. Closing connection",
		},
	}
	conn.WriteJSON(msg)
	conn.Close()

	//Remove connection from lobby map so we don't try to send them any more messages
	delete(lobby, playerId)

	return updatedLobby, nil
}

func closeAllConnections(lobby map[string]*websocket.Conn) {
	message := WebsocketMessage{
		Type: WebsocketMessage_Close,
		Data: SocketClose{
			Message: "Game has ended. Closing connection",
		},
	}
	for playerId, conn := range lobby {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending message to %s. Aborting message, but closing connection anyways", playerId)
		}
		conn.Close()
		delete(lobby, playerId)
	}
}

func sendMessageToAllPlayers(lobby map[string]*websocket.Conn, message WebsocketMessage) {
	//log.Printf("Sending message to every player: %s", message)

	if message.Type == "" {
		log.Println("WARNING: Websocket message being sent has no Type set! Frontend will likely not know how to handle the message!")
	}

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
		log.Printf("Something went wrong trying to read from the connection of Player, likely due to an unexpected closing of the Websocket connection: {%s} -- %s", playerId, r)
		//Disconnect client by closing connection and removing player from lobby
		gamesClientsMutex.Lock()
		//lobby[playerId].Close() Assume the socket is already closed if we're here
		delete(lobby, playerId)
		gamesClientsMutex.Unlock()
	}
}
