package Engine

import (
	"candlelight-models/Game"
	"candlelight-models/Player"
	"candlelight-models/Rules"
	"candlelight-models/Session"
	"testing"
)

const DUMMY_ID = "dummy"

// ==========================TESTS===============================
func TestSaveGameDefToDB(t *testing.T) {
	var tests = []struct {
		name                string
		game                Game.Game
		shouldReturnError   bool
		idShouldBeGenerated bool
	}{
		{
			name: "Predefined Id",
			game: Game.Game{
				Id:         "game123",
				Name:       "Epic Adventure",
				Genre:      "Strategy",
				MaxPlayers: 4,
				Resources: []Game.GameResource{
					{
						Id:           "resource1",
						Name:         "Gold",
						Description:  "Used to buy items and upgrades",
						InitialValue: 100,
						MaxValue:     1000,
						MinValue:     0,
					},
					{
						Id:           "resource2",
						Name:         "Mana",
						Description:  "Required for casting spells",
						InitialValue: 50,
						MaxValue:     500,
						MinValue:     0,
					},
				},
			},
			shouldReturnError:   false,
			idShouldBeGenerated: false,
		},
		{
			name: "No Predefined Id",
			game: Game.Game{
				Name:       "Epic Adventure",
				Genre:      "Strategy",
				MaxPlayers: 4,
				Resources: []Game.GameResource{
					{
						Id:           "resource1",
						Name:         "Gold",
						Description:  "Used to buy items and upgrades",
						InitialValue: 100,
						MaxValue:     1000,
						MinValue:     0,
					},
					{
						Id:           "resource2",
						Name:         "Mana",
						Description:  "Required for casting spells",
						InitialValue: 50,
						MaxValue:     500,
						MinValue:     0,
					},
				},
			},
			shouldReturnError:   false,
			idShouldBeGenerated: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saved, err := SaveGameDefToDB(tt.game)
			if saved.Id != "" {
				defer RDB.Del(RDB.Context(), "game:"+saved.Id)
			}

			//Error should be received IFF shouldReturnError == true
			if (err != nil) != tt.shouldReturnError {
				t.Errorf("%s -- Error value: Expected {%t}, Got {%s}", tt.name, tt.shouldReturnError, err)
			}

			if tt.shouldReturnError {
				//An Error should be accompanied by the game you tried to save
				if saved.Id != tt.game.Id {
					t.Errorf("%s -- Expected empty game, Got game with ID == {%s}", tt.name, saved.Id)
				}
			} else {
				//If we gave an ID, the saved Game's ID should be what we gave it
				if !tt.idShouldBeGenerated && tt.game.Id != saved.Id {
					t.Errorf("%s -- Game's ID was overwritten! Expected {%s}, Got {%s}", tt.name, tt.game.Id, saved.Id)
					//If we didn't give an ID, the Game should be given an ID
				} else if tt.idShouldBeGenerated && saved.Id == "" {
					t.Errorf("%s -- Game was not given an ID!", tt.name)
				}
			}
		})
	}
}

func TestGetGameDefFromDB(t *testing.T) {
	var tests = []struct {
		name              string
		gameId            string
		shouldReturnError bool
	}{
		{
			name:              "Valid GameId",
			gameId:            DUMMY_ID,
			shouldReturnError: false,
		},
		{
			name:              "Invalid GameId",
			gameId:            "invalid",
			shouldReturnError: true,
		},
	}

	//Test Setup
	saveDummyGameDef()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, err := GetGameDefFromDB(tt.gameId)

			//Error should be received IFF shouldReturnError == true
			if (err != nil) != tt.shouldReturnError {
				t.Errorf("%s -- Error value: Expected {%t}, Got {%s}", tt.name, tt.shouldReturnError, err)
			}

			if tt.shouldReturnError {
				//An Error should be accompanied by an empty game, just in case
				if game.Id != "" {
					t.Errorf("%s -- Expected empty game, Got game with ID == {%s}", tt.name, game.Id)
				}
			} else if game.Id != tt.gameId {
				t.Errorf("%s -- Expected GameId {%s}, Got {%s}", tt.name, tt.gameId, game.Id)
			}
		})
	}
}

func TestCreateRoom(t *testing.T) {
	var tests = []struct {
		name                   string
		gameDefId              string
		playerName             string
		expectedRoomCodeLength int
		shouldReturnError      bool
		shouldReturnRoomCode   bool
	}{
		{
			name:                   "Valid GameDefId",
			gameDefId:              DUMMY_ID,
			expectedRoomCodeLength: 4,
			shouldReturnError:      false,
			shouldReturnRoomCode:   true,
		},
		{
			name:                   "Invalid GameDefId",
			gameDefId:              "invalid",
			shouldReturnError:      true,
			shouldReturnRoomCode:   false,
			expectedRoomCodeLength: 0,
		},
	}

	//Test Setup
	saveDummyGameDef()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomCode, err := CreateRoom(tt.gameDefId)

			if (err != nil) != tt.shouldReturnError {
				t.Errorf("%s -- Error value: Expected {%t}, Got {%s}", tt.name, tt.shouldReturnError, err)
			}

			if tt.shouldReturnError {
				//If we get an error, we should not get a room code
				if roomCode != "" {
					t.Errorf("%s -- Should not have received RoomCode, got {%s}", tt.name, roomCode)
				}
			} else {
				//Should receive a 4-character room code
				if len(roomCode) != 4 {
					t.Errorf("%s -- Expected RoomCode of length {%d}, Got RoomCode of length {%d}", tt.name, tt.expectedRoomCodeLength, len(roomCode))
				}

				//The lobby should be usable to get a lobby from redis
				_, err := LoadLobbyFromRedis(roomCode)
				if err != nil {
					t.Errorf("%s -- Given room code could not load lobby. Got Error: {%s}", tt.name, err)
				}
			}
		})
	}
}

func TestJoinRoom(t *testing.T) {
	var tests = []struct {
		name              string
		lobbyRoomCode     string
		playerName        string
		shouldReturnError bool
	}{
		{
			name:              "Valid RoomCode",
			lobbyRoomCode:     DUMMY_ID,
			playerName:        "ryan",
			shouldReturnError: false,
		},
		{
			name:              "Invalid RoomCode",
			lobbyRoomCode:     "invalid",
			playerName:        "ryan",
			shouldReturnError: true,
		},
	}

	//Test Setup
	saveDummyLobby()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lobby, playerId, err := JoinRoom(tt.lobbyRoomCode, tt.playerName)

			//Error should be received IFF shouldReturnError == true
			if (err != nil) != tt.shouldReturnError {
				t.Errorf("%s -- Error value: Expected {%t}, Got {%s}", tt.name, tt.shouldReturnError, err)
			}

			if tt.shouldReturnError {
				//An Error should be accompanied by an empty lobby
				if lobby.RoomCode != "" {
					t.Errorf("%s -- Expected empty lobby, Got lobby with RoomCode == {%s}", tt.name, tt.lobbyRoomCode)
				}
			} else {
				if lobby.RoomCode != tt.lobbyRoomCode {
					t.Errorf("%s -- Received Lobby RoomCode mismatch! Expected {%s}, Got {%s}", tt.name, tt.lobbyRoomCode, lobby.RoomCode)
				}

				//Should receive a PlayerId for the frontend to save
				if playerId == "" {
					t.Errorf("%s -- Did not receive a PlayerId!", tt.name)
				}

				//Player should be in the player list
				if !playerExistsInLobby(lobby, tt.playerName, playerId) {
					t.Errorf("%s -- Player list mismatch! Player {%s} not found in Lobby!", tt.name, tt.playerName)
				}

				//Since there are no players in the dummy lobby, this player will be first to join and should be set as the host
				if !playerIsHost(lobby, tt.playerName, playerId) {
					t.Errorf("%s -- Host mismatch! Player {%s} was not set as host!", tt.name, tt.playerName)
				}
			}
		})
	}
}

// ==================HELPER FUNCTIONS=============================
func playerExistsInLobby(lobby Session.Lobby, playerName string, playerId string) bool {
	for _, player := range lobby.Players {
		if playerId != "" {
			if player.Id == playerId && player.Name == playerName {
				return true
			}
		} else {
			if player.Name == playerName {
				return true
			}
		}
	}
	return false
}

func playerIsHost(lobby Session.Lobby, playerName string, playerId string) bool {
	return lobby.Host.Name == playerName && lobby.Host.Id == playerId
}

func saveDummyGameDef() {
	SaveGameDefToDB(Game.Game{
		Id:         DUMMY_ID,
		Name:       "dummy",
		MaxPlayers: 4,
		Genre:      "test",
		Resources: []Game.GameResource{
			{
				Id:           "1",
				Name:         "magic",
				MaxValue:     60,
				MinValue:     0,
				InitialValue: 15,
			},
			{
				Id:           "2",
				Name:         "gold",
				MaxValue:     999,
				MinValue:     0,
				InitialValue: 0,
			},
		},
		BeginningPhase: Rules.GamePhase{
			Id:   "1",
			Name: "Beginning Phase",
		},
	})
}

func saveDummyLobby() {
	//Make sure the dummy game exists since we use it as the game def id
	saveDummyGameDef()
	SaveLobbyInRedis(Session.Lobby{
		RoomCode:         DUMMY_ID,
		GameDefinitionId: DUMMY_ID,
		Status:           Session.LobbyStatus_AwaitingStart,
		NumPlayers:       0,
		MaxPlayers:       4,
		Players:          []Player.Player{},
		Host:             Player.Player{},
	})
}
