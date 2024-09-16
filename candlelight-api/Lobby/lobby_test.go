package Lobby

import (
	"candlelight-models/Player"
	"candlelight-ruleengine/Engine"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestHostLobby(t *testing.T) {
	tests := []struct {
		name          string
		queryString   string
		shouldSucceed bool
	}{
		{
			name:          "Valid Call",
			queryString:   "?gameId=game123&playerName=testplayer",
			shouldSucceed: true,
		},
		{
			name:          "Invalid Game Id",
			queryString:   "?playerName=testplayer&gameId=" + Engine.GenerateId(),
			shouldSucceed: false,
		},
		{
			name:          "Missing Game Id",
			queryString:   "?playerName=testplayer",
			shouldSucceed: false,
		},
		{
			name:          "Missing Player Name",
			queryString:   "?gameId=game123",
			shouldSucceed: false,
		},
		{
			name:          "Extra Query Params",
			queryString:   "?gameId=game123&playerName=testplayer&extraparam=true&otherparam=hello",
			shouldSucceed: true,
		},
		{
			name:          "No Query Params",
			queryString:   "",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server using the hostLobby handler.
			// Dummy game must be created prior to running
			server := httptest.NewServer(http.HandlerFunc(HostLobby))
			defer server.Close()

			// Modify the URL for WebSocket usage
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + tt.queryString

			// Connect to the WebSocket server
			ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

			if tt.shouldSucceed {
				if err != nil {
					t.Fatalf("Failed to open ws connection on %s: %v", wsURL, err)
				}

				defer ws.Close()

				// Read the message from the WebSocket connection
				_, bytes, err := ws.ReadMessage()
				if err != nil {
					t.Errorf("Failed to read message from WebSocket connection: %v", err)
					return
				}

				// Check that the received data is not empty
				receivedData := LobbyMessage{}
				if len(bytes) < 1 {
					t.Errorf("Received no data, expected a non-empty value")
				}

				err = json.Unmarshal(bytes, &receivedData)
				if err != nil {
					t.Errorf("error trying to unmarshal received JSON: %s", err)
				}

				//Check that we were given a PlayerId
				if receivedData.PlayerID == "" {
					t.Errorf("Did not receive PlayerID")
				}

				lobby := receivedData.LobbyInfo
				//Set up the Cleanup
				defer Engine.RDB.Del(Engine.RDB.Context(), "lobby:"+lobby.RoomCode)
				defer testRecovery(t, "lobby:"+lobby.RoomCode)

				//Given lobby should contain one player whose name is "testplayer" as given by the query string
				if len(lobby.Players) != 1 {
					t.Errorf("Lobby does not have 1 player! Actual: %d", len(lobby.Players))
				}

				if !slices.ContainsFunc(lobby.Players, func(p Player.Player) bool { return p.Name == "testplayer" }) {
					t.Error("Could not find 'testplayer' in Player list!")
				}

				//Given lobby should have a room code that is 4 letters
				if len(lobby.RoomCode) != 4 {
					t.Errorf("Error with length of given room code! Expected 4, got %d", len(lobby.RoomCode))
				}

				alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
				for index, letter := range lobby.RoomCode {
					if !strings.ContainsRune(alphabet, letter) {
						t.Errorf("Could not find letter %c at index %d from Lobby Code {%s} in standard alphabet!", letter, index, lobby.RoomCode)
					}
				}

			} else {
				if err == nil {
					t.Errorf("Was able to connect with faulty params!")
				}
			}
		})
	}
}

func TestJoinLobby(t *testing.T) {

	roomCode, err := Engine.CreateRoom("game123")
	if err != nil {
		t.Fatal("Couldn't Create Lobby for dummy game! Ensure function createJSON has been called or a GET request has been sent to /dummy")
	}

	defer testRecovery(t, "lobby:"+roomCode)
	defer Engine.RDB.Del(Engine.RDB.Context(), "lobby:"+roomCode)

	//Hack our newly created lobby into the gamesClients tracker
	gamesClients[roomCode] = make(map[string]*websocket.Conn)

	tests := []struct {
		name          string
		playerName    string
		queryString   string
		shouldSucceed bool
	}{
		{
			name:          "Valid Join Request",
			playerName:    "validtestplayer",
			queryString:   "?playerName=validtestplayer&roomCode=" + roomCode,
			shouldSucceed: true,
		},
		{
			name:          "Invalid Room Code",
			playerName:    "invalidroomcode",
			queryString:   "?playerName=invalidroomcode&roomCode=invalid",
			shouldSucceed: false,
		},
		{
			name:          "Missing Room Code",
			playerName:    "missingroomcode",
			queryString:   "?playerName=missingroomcode",
			shouldSucceed: false,
		},
		{
			name:          "Missing Player Name",
			playerName:    "",
			queryString:   "?roomCode=" + roomCode,
			shouldSucceed: false,
		},
		{
			name:          "No Query Params",
			playerName:    "",
			queryString:   "",
			shouldSucceed: false,
		},
		{
			name:          "Extra Query Param",
			playerName:    "extraqueryparam",
			queryString:   "?playerName=extraqueryparam&roomCode=" + roomCode + "&extraparam=extravalue",
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer testRecovery(t, "lobby:"+roomCode)
			// Create a test server using the hostLobby handler.
			// Dummy game must be created prior to running
			server := httptest.NewServer(http.HandlerFunc(HandleJoinLobby))
			defer server.Close()

			// Modify the URL for WebSocket usage
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + tt.queryString

			ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if tt.shouldSucceed {
				if err != nil {
					t.Fatalf("Failed to open ws connection on %s: %v", wsURL, err)
				}
				defer ws.Close()

				// Read the message from the WebSocket connection
				_, bytes, err := ws.ReadMessage()
				if err != nil {
					t.Errorf("Failed to read message from WebSocket connection: %v", err)
					return
				}

				// Check that the received data is not empty
				receivedData := LobbyMessage{}
				if len(bytes) < 1 {
					t.Errorf("Received no data, expected a non-empty value")
				}

				err = json.Unmarshal(bytes, &receivedData)
				if err != nil {
					t.Errorf("error trying to unmarshal received JSON: %s", err)
				}

				//Check that we were given a PlayerId
				if receivedData.PlayerID == "" {
					t.Errorf("Did not receive PlayerID")
				}

				lobby := receivedData.LobbyInfo

				if lobby.RoomCode != roomCode {
					t.Errorf("Received lobby roomcode mismatch! Expected {%s}, Actual {%s}", roomCode, lobby.RoomCode)
				}

				//Make sure there's a player in the lobby with our name
				if !slices.ContainsFunc(lobby.Players, func(p Player.Player) bool { return p.Name == tt.playerName }) {
					t.Errorf("Couldn't find player with name == {%s} in lobby!", tt.playerName)
				}

			} else {
				if err == nil {
					t.Error("Did not receive expected error trying to connect websocket")
				} else {
					return //Succeed
				}
			}
			ws.ReadMessage()
		})
	}
}

func TestRejoinLobby(t *testing.T) { //TODO: Tests don't take their player out of the lobby so we're hitting the max player count

	roomCode, err := Engine.CreateRoom("game123")
	if err != nil {
		t.Fatal("Couldn't Create Lobby for dummy game! Ensure function createJSON has been called or a GET request has been sent to /dummy")
	}

	defer testRecovery(t, "lobby:"+roomCode)
	defer Engine.RDB.Del(Engine.RDB.Context(), "lobby:"+roomCode)

	//Hack our newly created lobby into the gamesClients tracker
	gamesClients[roomCode] = make(map[string]*websocket.Conn)

	tests := []struct {
		name                 string
		firstJoinQueryString string
		rejoinQueryString    string
		shouldSucceed        bool
		addPlayerId          bool
	}{
		{
			name:              "Valid Rejoin",
			rejoinQueryString: "?roomCode=" + roomCode + "&playerId=",
			shouldSucceed:     true,
			addPlayerId:       true, //Add the PlayerId from first joining to the end of the query string before trying to rejoin
		},
		{
			name:              "No RoomCode",
			rejoinQueryString: "?playerId=",
			shouldSucceed:     false,
			addPlayerId:       true,
		},
		{
			name:              "No PlayerId",
			rejoinQueryString: "?roomCode=" + roomCode,
			shouldSucceed:     false,
			addPlayerId:       false,
		},
		{
			name:              "Duplicate Connection",
			rejoinQueryString: "?roomCode=" + roomCode + "&playerId=",
			shouldSucceed:     false,
			addPlayerId:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer testRecovery(t, "lobby:"+roomCode)
			// Create a test server using the hostLobby handler.
			// Dummy game must be created prior to running
			server := httptest.NewServer(http.HandlerFunc(HandleJoinLobby))
			defer server.Close()

			// Modify the URL for WebSocket usage
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?playerName=rejointester" + Engine.GenerateId() + "&roomCode=" + roomCode

			//First join. Make sure everything looks good, then close connection
			ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				t.Fatalf("Couldn't open websocket on first join! %s", err)
			}

			// Read the message from the WebSocket connection
			_, bytes, err := ws.ReadMessage()
			if err != nil {
				t.Errorf("Failed to read message from WebSocket connection: %v", err)
				return
			}

			// Check that the received data is not empty
			receivedData := LobbyMessage{}
			if len(bytes) < 1 {
				t.Errorf("Received no data, expected a non-empty value")
			}

			err = json.Unmarshal(bytes, &receivedData)
			if err != nil {
				t.Errorf("error trying to unmarshal received JSON: %s", err)
			}

			//Check that we were given a PlayerId
			if receivedData.PlayerID == "" {
				t.Errorf("Did not receive PlayerID")
			}

			if tt.name != "Duplicate Connection" { //Leave the connection open for this specific test
				ws.Close()
			}

			//Now rejoin
			server.Config.Handler = http.HandlerFunc(HandleRejoinLobby)
			if tt.addPlayerId {
				wsURL = "ws" + strings.TrimPrefix(server.URL, "http") + tt.rejoinQueryString + receivedData.PlayerID
			} else {
				wsURL = "ws" + strings.TrimPrefix(server.URL, "http") + tt.rejoinQueryString
			}
			ws, _, err = websocket.DefaultDialer.Dial(wsURL, nil)

			if tt.shouldSucceed {
				if err != nil {
					t.Fatalf("Couldn't open websocket on rejoin! %s", err)
				}
				defer ws.Close()

				// Read the message from the WebSocket connection
				_, bytes, err := ws.ReadMessage()
				if err != nil {
					t.Errorf("Failed to read message from WebSocket connection: %v", err)
					return
				}

				// Check that the received data is not empty
				receivedData := LobbyMessage{}
				if len(bytes) < 1 {
					t.Errorf("Received no data, expected a non-empty value")
				}

				err = json.Unmarshal(bytes, &receivedData)
				if err != nil {
					t.Errorf("error trying to unmarshal received JSON: %s", err)
				}

				//Check that we were given a PlayerId
				if receivedData.PlayerID == "" {
					t.Errorf("Did not receive PlayerID")
				}
			} else {
				if err == nil {
					t.Error("Did not receive expected error trying to connect websocket")
				} else {
					return //Succeed
				}
			}
		})
	}
}

func testRecovery(t *testing.T, dbCleanupKey string) {
	t.Helper()

	if r := recover(); r != nil {
		if dbCleanupKey != "" {
			Engine.RDB.Del(Engine.RDB.Context(), dbCleanupKey)
		}
		t.Fatal("Go panicked")
	}
}
