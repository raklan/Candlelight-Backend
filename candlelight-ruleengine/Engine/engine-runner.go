package Engine

import (
	"candlelight-api/LogUtil"
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-models/Player"
	"candlelight-models/Session"
	"candlelight-models/Sparks"
	"slices"
	"time"

	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Saves the given [game] in the database, which is currently Redis. If the save is successful, [error] will be nil
func SaveGameDefToDB(game Game.Game) (Game.Game, error) {
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	funcLogPrefix := "==SaveGameDefToDB==:"
	log.Printf("%s Saving Game with id=={%s}", funcLogPrefix, game.Id)

	// If the Game doesn't have an ID yet, generate one
	id := game.Id
	if id == "" {
		log.Printf("%s Game doesn't have an Id. Generating one...", funcLogPrefix)
		id = GenerateId()
		log.Printf("%s Id successfully generated. Assigning Id {%s} to Game", funcLogPrefix, id)
		game.Id = id
	}

	asJson, err := json.Marshal(game)
	if err != nil {
		LogError(funcLogPrefix, err)
		return game, err
	}

	key := "game:" + id
	err = RDB.Set(ctx, key, asJson, 0).Err()
	if err != nil {
		LogError(funcLogPrefix, err)
		return game, err
	}

	log.Printf("%s GameDefinition saved with key == {%s}", funcLogPrefix, key)

	return game, nil
}

// Grabs a game from Redis for the given [id]. Returns nil for [error] if the returned Game is an actual Game that can be used
func GetGameDefFromDB(id string) (Game.Game, error) {
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	funcLogPrefix := "==GetGameDefFromDB==:"

	log.Printf("%s Retrieving Game with id=={%s} from DB", funcLogPrefix, id)

	game := Game.Game{}

	//Catch empty ID
	if id == "" {
		log.Printf("%s ERROR! Id cannot be empty. Returning empty Game Definition", funcLogPrefix)
		return game, fmt.Errorf("%s Id cannot be empty", funcLogPrefix)
	}

	//Try to get the Game from Redis. If it doesn't exist, give a specific error for that
	def, err := RDB.Get(ctx, "game:"+id).Result()
	if err == redis.Nil {
		log.Printf("%s Could not find cached Game for id \"%s\"...Returning Empty Game", funcLogPrefix, id)
		return game, fmt.Errorf("%s No game for Id=={%s} found", funcLogPrefix, id)
	} else if err != nil {
		LogError(funcLogPrefix, err)
		return game, err
	}

	//Result is just a JSON string, so we still need to deserialize/unmarshal it
	err = json.Unmarshal([]byte(def), &game)
	if err != nil {
		LogError(funcLogPrefix, err)
		return game, err
	}

	log.Printf("%s Found a Game, returning result", funcLogPrefix)
	return game, nil
}

// Deletes a game from Redis for the given [id]. Returns the id of the deleted game and nil for [error] if the game is successfully deleted
func DeleteGameDefFromDB(id string) (string, error) {
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	funcLogPrefix := "==DeleteGameDefFromDB==:"

	log.Printf("%s Retrieving Game with id=={%s} from DB", funcLogPrefix, id)

	//Catch empty ID
	if id == "" {
		log.Printf("%s ERROR! Id cannot be empty. Returning empty string", funcLogPrefix)
		return "", fmt.Errorf("%s Id cannot be empty", funcLogPrefix)
	}

	numDeleted, err := RDB.Del(ctx, "game:"+id).Result()
	if numDeleted == 0 { //No game was deleted, which means it doesn't exist
		log.Printf("%s Couldn't find game with ID == {%s}", funcLogPrefix, id)
		return "", fmt.Errorf("could not find game with id == {%s}", id)
	} else if err != nil {
		log.Printf("%s Error trying to delete Game", funcLogPrefix)
		return "", fmt.Errorf("unknown error trying to delete game")
	}

	log.Printf("%s Deleted a Game, returning id of deleted game", funcLogPrefix)
	return id, nil
}

// Retrieves all games in the DB and returns them in a list
func GetAllGamesFromDB(criteria Criteria) ([]Game.Game, error) {
	funcLogPrefix := "==GetAllGamesFromDB_Slimmed=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Gettings all GAMES from DB...", funcLogPrefix)

	var cursor uint64
	toReturn := []Game.Game{}
	for { //Iterate through all keys beginning with "game:" and break when cursor is 0.
		var keys []string
		var err error
		keys, cursor, err = RDB.Scan(ctx, cursor, "game:*", 1000).Result()
		if err != nil {
			LogError(funcLogPrefix, err)
			return toReturn, err
		}

		//SCAN returns subsets of matching keys, so we need to iterate through each subset
		//as it comes back, get each game for the keys returned, and add it to the results
		for _, k := range keys {
			//Get Key from DB
			gameJSON, err := RDB.Get(ctx, k).Result()
			if err != nil {
				LogError(funcLogPrefix, err)
				return toReturn, err
			}

			//Result is a JSON string, so deserialize it and append to results
			game := Game.Game{}
			err = json.Unmarshal([]byte(gameJSON), &game)
			if err != nil {
				LogError(funcLogPrefix, err)
				return toReturn, err
			}

			//Check if the Game matches the given criteria, if any
			if criteria.Check(game) {
				toReturn = append(toReturn, game)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return toReturn, nil
}

// Caches the given [gameState] in redis. Returns nil for [error] if everything goes well
func CacheGameStateInRedis(gameState Session.GameState) (Session.GameState, error) {
	funcLogPrefix := "==CacheGameStateInRedis==:"
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Received GameState to cache", funcLogPrefix)

	//If the gameState doesn't have an ID yet,
	//Generate one for it by simply using the Current UNIX time in milliseconds
	id := gameState.Id
	if id == "" {
		log.Printf("%s GameState does not yet have an ID. Generating new one.", funcLogPrefix)
		id = GenerateId()
		log.Printf("%s ID successfully generated. Assigning ID {%s} to GameState", funcLogPrefix, id)
		gameState.Id = id
	}

	//Convert to string and save to Redis
	asJson, err := json.Marshal(gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	key := "gameState:" + id
	expiry, _ := time.ParseDuration("168h")
	err = RDB.Set(ctx, key, asJson, expiry).Err()
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	log.Printf("%s GameState cached with key=={%s}", funcLogPrefix, key)
	return gameState, nil
}

// Retrieves a gameState with an id == [id] from Redis. If everything goes well, then [error] is nil
func GetCachedGameStateFromRedis(id string) (Session.GameState, error) {
	funcLogPrefix := "==GetCachedGameStateFromRedis==:"
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Received request to get cached GameState from Redis", funcLogPrefix)

	gameState := Session.GameState{}

	//Catch empty id string early
	if id == "" {
		log.Printf("%s ERROR! Id cannot be empty. Returning empty GameState", funcLogPrefix)
		return gameState, fmt.Errorf("%s Id cannot be empty", funcLogPrefix)
	}

	//Try to get the game from Redis. If it doesn't exist, fail gracefully
	game, err := RDB.Get(ctx, "gameState:"+id).Result()
	if err == redis.Nil {
		log.Printf("%s Could not find cached GameState for key \"%s\"...Returning Empty GameState", funcLogPrefix, id)
		return gameState, fmt.Errorf("%s No game for Id=={%s} found", funcLogPrefix, id)
	} else if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	//game is a JSON string of a GameState, so unmarshal it
	err = json.Unmarshal([]byte(game), &gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	log.Printf("%s Found a GameState, returning result", funcLogPrefix)
	return gameState, nil
}

// Given an id to a Game defition, constructs and returns an initial GameState for it. This is essentially
// how to start the game
func GetInitialGameState(roomCode string) (Session.GameState, error) {
	funcLogPrefix := "==GetInitialGameState=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	gameState := Session.GameState{}

	lobby, err := LoadLobbyFromRedis(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	//Check if the lobby is already started.
	if lobby.Status == Session.LobbyStatus_InProgress {
		err := fmt.Errorf("tried to start game, but Lobby {%s} has been marked as In Progress and has a GameStateId == {%s}", lobby.RoomCode, lobby.GameStateId)
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	gameDef, err := GetGameDefFromDB(lobby.GameDefinitionId)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	gameState.GameDefinitionId = gameDef.Id
	gameState.GameName = gameDef.Name
	gameState.Rules = gameDef.Rules
	gameState.SplashText = gameDef.SplashText
	gameState.Views = gameDef.ViewsForPlayer(0) //Player 0 == public/table-owned

	//startingResources := make([]Player.PlayerResource, len(gameDef.Resources))

	//Construct starting resources for each player
	// for _, element := range gameDef.Resources {
	// 	startingResources = append(startingResources, Player.PlayerResource{
	// 		Name:         element.Name,
	// 		CurrentValue: element.InitialValue,
	// 		MaxValue:     element.MaxValue,
	// 	})
	// }

	gameState.Players = []Player.Player{}

	for index, element := range lobby.Players {
		gameState.Players = append(gameState.Players, Player.Player{
			Id:   element.Id,
			Name: element.Name,
			Hand: gameDef.ViewsForPlayer(index + 1), //TODO: Need a more in-depth discussion about what to do in terms of determining starting pieces
			//Resources: slices.Clone(startingResources),
		})
	}

	gameState.CurrentPlayer = gameState.Players[0].Id //TODO: Make a better way to determine a starting player maybe?

	applySparks(&gameState, gameDef.Sparks)

	gameState, err = CacheGameStateInRedis(gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	//Mark the lobby as started and fill in GameStateId
	lobby.GameStateId = gameState.Id
	lobby.Status = Session.LobbyStatus_InProgress
	_, err = SaveLobbyInRedis(lobby)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	return gameState, nil
}

func applySparks(gameState *Session.GameState, sparks Sparks.Sparks) {
	if sparks.Dealer.Enabled {
		applyDealer(gameState, sparks.Dealer)
	}
	if sparks.Flipper.Enabled {
		applyFlipper(gameState, sparks.Flipper)
	}
}

func applyDealer(gameState *Session.GameState, dealer Sparks.Dealer) {
	var deckToUse (*Pieces.Deck) = nil
	for _, view := range gameState.Views {
		indexOfDeck := slices.IndexFunc(view.Pieces.Decks, func(d Pieces.Deck) bool { return d.Id == dealer.DeckToUse })
		if indexOfDeck != -1 {
			deckToUse = &view.Pieces.Decks[indexOfDeck]
			break
		}
	}

	if deckToUse == nil {
		LogError("==applyDealer==", fmt.Errorf("could not find deck to deal from. Ignoring Dealer"))
		return
	}

	for _, player := range gameState.Players {
		for x := range dealer.NumToDeal {
			cardWithdraw := deckToUse.PickRandomCardFromCollection()
			cardCopy := *cardWithdraw
			cardCopy.ParentView = player.Hand[0].Id
			//Put X as 0, 20, 40, etc
			cardCopy.X = float32(x * 20)
			cardCopy.Y = 0

			deckToUse.RemoveCardFromCollection(*cardWithdraw)
			player.Hand[0].Pieces.Orphans = append(player.Hand[0].Pieces.Orphans, cardCopy)
		}
	}
}

func applyFlipper(gameState *Session.GameState, flipper Sparks.Flipper) {
	var deckToUse (*Pieces.Deck) = nil
	var cardPlaceToUse (*Pieces.CardPlace) = nil
	foundDeck, foundCardPlace := false, false
	indexOfDeck, indexOfCardPlace := -1, -1
	for _, view := range gameState.Views {
		if !foundDeck {
			indexOfDeck = slices.IndexFunc(view.Pieces.Decks, func(d Pieces.Deck) bool { return d.Id == flipper.DeckToUse })
		}
		if !foundCardPlace {
			indexOfCardPlace = slices.IndexFunc(view.Pieces.CardPlaces, func(cp Pieces.CardPlace) bool { return cp.Id == flipper.CardPlaceToUse })
		}

		if indexOfDeck != -1 {
			foundDeck = true
			deckToUse = &view.Pieces.Decks[indexOfDeck]
		}
		if indexOfCardPlace != -1 {
			foundCardPlace = true
			cardPlaceToUse = &view.Pieces.CardPlaces[indexOfCardPlace]
		}

		if foundDeck && foundCardPlace {
			break
		}
	}

	if deckToUse == nil {
		LogError("==applyFlipper==", fmt.Errorf("could not find deck to deal from. Ignoring Flipper"))
		return
	}
	if cardPlaceToUse == nil {
		LogError("==applyFlipper==", fmt.Errorf("could not find cardplace to put cards in. Ignoring Flipper"))
		return
	}

	for range flipper.NumToFlip {
		cardWithdraw := deckToUse.PickRandomCardFromCollection()
		cardCopy := *cardWithdraw
		cardCopy.ParentView = cardPlaceToUse.ParentView

		deckToUse.RemoveCardFromCollection(*cardWithdraw)
		cardPlaceToUse.AddCardToCollection(cardCopy)
	}

}

func EndGame(roomCode string, playerId string) error {
	funcLogPrefix := "==EndGame=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	lobby, err := LoadLobbyFromRedis(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	//Make sure that A) this player is the host and therefore allowed to end the game, and B) this game isn't already ended

	if lobby.Host.Id != playerId {
		return fmt.Errorf("player trying to end game is not host of lobby")
	}

	if lobby.Status == Session.LobbyStatus_Ended {
		return fmt.Errorf("game has already been marked as ended")
	}

	//Mark Game as ended and resave
	lobby.Status = Session.LobbyStatus_Ended

	_, err = SaveLobbyInRedis(lobby)

	//Return any error that occurred during saving, if any
	return err
}

// Submits an Action to the GameState with id == [gameId]. Will always return some GameState, even if something goes wrong, in which case [error] will not be nil.
// If the action is not allowed, [error] will indicate so, and it will simply return the GameState without any changes
func SubmitAction(gameId string, action Session.SubmittedAction) (Session.Changelog, error) {
	funcLogPrefix := "==SubmitAction=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	changelog := Session.Changelog{}

	//Grab last cached gameState
	gameState, err := GetCachedGameStateFromRedis(gameId)
	if err != nil {
		LogError(funcLogPrefix, err)
		return changelog, err
	}

	changelog = Session.Changelog{
		Views:         []*Game.View{},
		CurrentPlayer: gameState.CurrentPlayer,
	}

	//Only allow the player whose turn it is to take an action
	if gameState.Rules.EnforceTurnOrder {
		if gameState.CurrentPlayer != action.PlayerId {
			log.Printf("Player %s has tried to submit an action when it's not their turn. (CurrentPlayer == %s) Ignoring action", action.PlayerId, gameState.CurrentPlayer)
			//Add every view to the changelog as unchanged so any rejected action from the client gets overwritten
			for _, view := range gameState.Views {
				changelog.Views = append(changelog.Views, &view)
			}
			for _, player := range gameState.Players {
				for _, view := range player.Hand {
					changelog.Views = append(changelog.Views, &view)
				}
			}
			return changelog, fmt.Errorf("Your action was rejected because it's not your turn!") //Error message formatted for display directly to user as requested by Brian
		}
	}

	switch action.Type {
	case Session.ActionType_Insertion:
		turn := Session.Insertion{}
		err = json.Unmarshal(action.Turn, &turn)
		if err != nil {
			LogError(funcLogPrefix, err)
			return changelog, fmt.Errorf("%s Error trying to unmarshal turn into Insertion: %s", funcLogPrefix, err)
		}
		changelog, _ = turn.Execute(&gameState, action.PlayerId)
	case Session.ActionType_Withdrawal:
		turn := Session.Withdrawal{}
		err = json.Unmarshal(action.Turn, &turn)
		if err != nil {
			LogError(funcLogPrefix, err)
			return changelog, fmt.Errorf("%s Error trying to unmarshal turn into Withdrawl: %s", funcLogPrefix, err)
		}
		changelog, _ = turn.Execute(&gameState, action.PlayerId)
	case Session.ActionType_Movement:
		turn := Session.Movement{}
		err = json.Unmarshal(action.Turn, &turn)
		if err != nil {
			LogError(funcLogPrefix, err)
			return changelog, fmt.Errorf("%s Error trying to unmarshal turn into Movement: %s", funcLogPrefix, err)
		}
		changelog, _ = turn.Execute(&gameState, action.PlayerId)
	case Session.ActionType_EndTurn:
		turn := Session.EndTurn{}
		err = json.Unmarshal(action.Turn, &turn)
		if err != nil {
			LogError(funcLogPrefix, err)
			return changelog, fmt.Errorf("%s Error trying to unmarshal turn into EndTurn: %s", funcLogPrefix, err)
		}
		changelog, _ = turn.Execute(&gameState, action.PlayerId)
	case Session.ActionType_CardFlip:
		turn := Session.Cardflip{}
		err = json.Unmarshal(action.Turn, &turn)
		if err != nil {
			LogError(funcLogPrefix, err)
			return changelog, fmt.Errorf("%s Error trying to unmarshal turn into Cardflip: %s", funcLogPrefix, err)
		}
		changelog, _ = turn.Execute(&gameState, action.PlayerId)
	case Session.ActionType_Reshuffle:
		turn := Session.Reshuffle{}
		err = json.Unmarshal(action.Turn, &turn)
		if err != nil {
			LogError(funcLogPrefix, err)
			return changelog, fmt.Errorf("%s Error trying to unmarshal turn into Reshuffle: %s", funcLogPrefix, err)
		}
		changelog, _ = turn.Execute(&gameState, action.PlayerId)
	default:
		return changelog, fmt.Errorf("%s Error - Submitted Action's type {%s} not recognized", funcLogPrefix, action.Type)
	}

	//Cache the updated gameState in Redis
	_, err = CacheGameStateInRedis(gameState)
	if err != nil {
		return changelog, fmt.Errorf("%s Error trying to cache updated gameState. Action may not properly persist! %s", funcLogPrefix, err)
	}

	return changelog, nil
}

func SaveLobbyInRedis(lobby Session.Lobby) (Session.Lobby, error) {
	funcLogPrefix := "==SaveLobbyInRedis=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Recieved request to save lobby in Redis", funcLogPrefix)

	asJson, err := json.Marshal(lobby)
	if err != nil {
		LogError(funcLogPrefix, err)
		return Session.Lobby{}, err
	}

	key := "lobby:" + lobby.RoomCode
	expiry, _ := time.ParseDuration("168h")
	err = RDB.Set(ctx, key, asJson, expiry).Err()
	if err != nil {
		LogError(funcLogPrefix, err)
		return Session.Lobby{}, err
	}

	log.Printf("%s Lobby saved in Redis with key == {%s}", funcLogPrefix, key)
	return lobby, nil
}

func LoadLobbyFromRedis(roomCode string) (Session.Lobby, error) {
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	funcLogPrefix := "==LoadLobbyFromRedis==:"

	log.Printf("%s Retrieving Lobby with RoomCode=={%s} from DB", funcLogPrefix, roomCode)

	lobby := Session.Lobby{}

	//Catch empty ID
	if roomCode == "" {
		log.Printf("%s ERROR! RoomCode cannot be empty. Returning empty Lobby", funcLogPrefix)
		return lobby, fmt.Errorf("%s Id cannot be empty", funcLogPrefix)
	}

	//Try to get the Game from Redis. If it doesn't exist, give a specific error for that
	def, err := RDB.Get(ctx, "lobby:"+roomCode).Result()
	if err == redis.Nil {
		log.Printf("%s Could not find cached lobby for roomCode \"%s\"...Returning Empty Lobby", funcLogPrefix, roomCode)
		return lobby, fmt.Errorf("%s No game for Id=={%s} found", funcLogPrefix, roomCode)
	} else if err != nil {
		LogError(funcLogPrefix, err)
		return lobby, err
	}

	//Result is just a JSON string, so we still need to deserialize/unmarshal it
	err = json.Unmarshal([]byte(def), &lobby)
	if err != nil {
		LogError(funcLogPrefix, err)
		return lobby, err
	}

	log.Printf("%s Found a lobby, returning result", funcLogPrefix)
	return lobby, nil
}

// Given a GameDefinition's ID and a player name, creates and saves a new lobby for that player's game, returning the Lobby's room code.
func CreateRoom(gameDefId string) (string, error) {
	funcLogPrefix := "==CreateRoom=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting GameDef for Id == {%s}", funcLogPrefix, gameDefId)
	requestedGame, err := GetGameDefFromDB(gameDefId)
	if err != nil {
		LogError(funcLogPrefix, err)
		return "", err
	}

	log.Printf("%s Creating lobby object", funcLogPrefix)
	lobby := Session.Lobby{
		GameDefinitionId: requestedGame.Id,
		Status:           Session.LobbyStatus_AwaitingStart,
		GameName:         requestedGame.Name,
		MaxPlayers:       requestedGame.MaxPlayers,
		NumPlayers:       0,
		Players:          []Player.Player{},
		Host:             Player.Player{},
	}

	log.Printf("%s Generating Room Code", funcLogPrefix)
	roomCode := generateRoomCode()

	log.Printf("%s Room Code successfully generated. Assigning RoomCode {%s} to Lobby", funcLogPrefix, roomCode)
	lobby.RoomCode = roomCode

	log.Printf("%s Saving Lobby to Redis", funcLogPrefix)
	lobby, err = SaveLobbyInRedis(lobby)
	if err != nil {
		LogError(funcLogPrefix, err)
		return "", err
	}

	log.Printf("%s Lobby Created & Saved. Returning RoomCode", funcLogPrefix)
	return roomCode, nil
}

// Creates a Player object for the given PlayerName and attempts to add them to the lobby with the given RoomCode. On Success, returns the new state
// of the lobby, the Player's assigned Id, and any error that occurred
func JoinRoom(roomCode string, playerName string) (Session.Lobby, string, error) {
	funcLogPrefix := "==JoinRoom=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Recieved request from Player {%s} to join lobby with RoomCode == {%s}", funcLogPrefix, playerName, roomCode)

	lobby, err := LoadLobbyFromRedis(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return Session.Lobby{}, "", err
	}

	//Only allow player to join if there's room & the game hasn't started yet (i.e. Status == LobbyStatus_AwaitingStart)
	if lobby.NumPlayers >= lobby.MaxPlayers {
		log.Printf("%s ERROR: Lobby's max player count {%d} already reached. Player cannot join!", funcLogPrefix, lobby.MaxPlayers)
		return Session.Lobby{}, "", fmt.Errorf("Lobby's max player count {%d} already reached", lobby.MaxPlayers)
	}
	if lobby.Status != Session.LobbyStatus_AwaitingStart {
		log.Printf("%s Error: Game has already started. Player cannot join!", funcLogPrefix)
		return Session.Lobby{}, "", fmt.Errorf("Game has already started!")
	}
	if slices.ContainsFunc(lobby.Players, func(p Player.Player) bool { return p.Name == playerName }) {
		log.Printf("%s Error: Player name {%s} already taken. Player cannot join!", funcLogPrefix, playerName)
		return Session.Lobby{}, "", fmt.Errorf("Name already taken!")
	}

	thisPlayer := createPlayerObject(playerName)

	//Create a copy, in case anything goes wrong
	updatedLobby := lobby
	updatedLobby.Players = slices.Clone(lobby.Players)

	log.Printf("%s Adding player {%s} to lobby's Player List", funcLogPrefix, playerName)

	//If this player is the first to join, set them as the host
	if updatedLobby.NumPlayers == 0 {
		updatedLobby.Host = thisPlayer
	}

	updatedLobby.Players = append(lobby.Players, thisPlayer)
	updatedLobby.NumPlayers = len(updatedLobby.Players)

	log.Printf("%s Player added. Caching new Lobby", funcLogPrefix)
	saved, err := SaveLobbyInRedis(updatedLobby)
	if err != nil { //If something goes wrong, re-save and return the version without any changes
		SaveLobbyInRedis(lobby)
		return Session.Lobby{}, "", err
	}

	log.Printf("%s Lobby joined and saved. Returning Lobby", funcLogPrefix)
	return saved, thisPlayer.Id, nil
}

func LeaveRoom(roomCode string, playerId string) (Session.Lobby, error) {
	funcLogPrefix := "==LeaveRoom=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Recieved request to remove Player {%s} from lobby with RoomCode == {%s}", funcLogPrefix, playerId, roomCode)

	lobby, err := LoadLobbyFromRedis(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return Session.Lobby{}, err
	}

	//Create a copy, in case anything goes wrong
	updatedLobby := lobby
	updatedLobby.Players = slices.Clone(lobby.Players)

	log.Printf("%s Removing player {%s} from lobby's Player List", funcLogPrefix, playerId)

	newPlayers := []Player.Player{}

	for _, player := range updatedLobby.Players {
		if player.Id != playerId {
			newPlayers = append(newPlayers, player)
		}
	}

	updatedLobby.Players = newPlayers
	updatedLobby.NumPlayers = len(newPlayers)

	log.Printf("%s Player Removed. Caching new Lobby", funcLogPrefix)
	saved, err := SaveLobbyInRedis(updatedLobby)
	if err != nil { //If something goes wrong, re-save and return the version without any changes
		SaveLobbyInRedis(lobby)
		return Session.Lobby{}, err
	}

	//If the game has started, we need to remove them from the GameState too
	if saved.Status == Session.LobbyStatus_InProgress {
		log.Println("Player is being removed from an in-progress game. Removing player from GameState...")
		gameState, err := GetCachedGameStateFromRedis(saved.GameStateId)
		if err != nil {
			LogError(funcLogPrefix, err)
		}

		//If it's this player's turn, end their turn before removing them
		if gameState.CurrentPlayer == playerId {
			log.Println("GameState is listing Player as CurrentPlayer. Ending their turn before removal...")
			Session.EndTurn{}.Execute(&gameState, playerId)
		}

		//Remove from Player list
		currentPlayers := slices.Clone(gameState.Players)
		newPlayers = []Player.Player{}

		for _, player := range currentPlayers {
			if player.Id != playerId {
				newPlayers = append(newPlayers, player)
			}
		}

		gameState.Players = newPlayers

		log.Println("Player has been removed from GameState. (NOTE: THIS HAS ALSO REMOVED ALL PIECES IN THEIR HAND FROM THE GAME. WILL FIX LATER) Caching new GameState now...")
		_, err = CacheGameStateInRedis(gameState)
		if err != nil {
			LogError(funcLogPrefix, err)
		}

	}

	log.Printf("%s Left Lobby. Returning Lobby", funcLogPrefix)
	return saved, nil
}

func createPlayerObject(name string) Player.Player {
	funcLogPrefix := "==CreatePlayerObject=="
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	defer LogUtil.EnsureLogPrefixIsReset()

	log.Printf("%s Creating Player object for Player name {%s}", funcLogPrefix, name)

	return Player.Player{
		Id:   GenerateId(),
		Name: name,
		Hand: []Game.View{},
		//Resources: []Player.PlayerResource{},
	}
}
