# Messages To the Server
Any message sent to the server should follow this general outline:
```json
{
  "jsonType": "string",
  "data": {
    "an object containing your message": "goes here"
  }
}
```

There are currently 5 supported message types:
- [startGame](#startgame)
- [endGame](#endGame)
- [submitAction](#submitaction)
- [leaveLobby](#leavelobby)
- [kickPlayer](#kickplayer)
- [disconnect](#disconnect)

### startGame
If the host of a lobby sends this message, it will start the game, mark the lobby as "In Progress" and send out a [GameState](#gamestate) message to ever client in the lobby. A non-host player sending this message is ignored.
```json
{
  "jsonType": "startGame",
  "data": {
    "the data field": "is ignored"
  }
}
```

### endGame
If the host of a lobby sends this message, it will end the game, mark the Lobby as "Game Ended" and send a [GameOver](#gameover) message immediately followed by a [Close](#close) message to every open connection in the lobby. A non-host player sending this message is ignored.
```json
{
  "jsonType": "endGame",
  "data": {
    "the data field": "is ignored"
  }
}
```

### submitAction
A client will send this message any time they want to affect something within the gamestate. Regardless of whether this action changes anything, every client in the lobby will receive a [Changelog](#changelog). The object within the `data` field should be one of the accepted [SubmittedActions](https://capstone-cs.eng.utah.edu/candlelight/candlelight-backend/-/wikis/Submitted-Actions)
```json
{
  "jsonType": "submitAction",
  "data": {
    "the action to submit": "an action to submit. See Submitted-Actions wiki page for details"
  }
}
```

### leaveLobby
A client can submit this message to be removed from the Lobby and have their connection closed. Currently, this is somewhat bugged and any card in their hand will simply be removed from the game as well. Once the player has been removed from the lobby, they will receive a [Close](#close) message, and every other client in the lobby will receive a [LobbyInfo](#lobbyinfo) message with the updated Lobby object. 
```json
{
  "jsonType": "leaveLobby",
  "data": {
    "the data field": "is ignored"
  }
}
```

### kickPlayer
The host **(and only the host)** can submit this message to remove another player from the lobby and have their connection closed.  Currently, this is somewhat bugged and any card in their hand will simply be removed from the game as well. Once the player has been removed from the lobby, they will receive a [Close](#close) message, and every other client in the lobby will receive a [LobbyInfo](#lobbyinfo) message with the updated Lobby object. A non-host player sending this message is replied to with a [Error](#error) message
```json
{
  "jsonType": "kickPlayer"
  "data": {
    "playerToKick": "The ID of the player to kick from the lobby"
  }
}
```

### disconnect
Any client can submit this message type to request the Server to close their connection without altering the underlying Lobby, for example, to leave room to rejoin later. Upon receiving this message type, the Server will respond with a [Close](#close) Message acknowledging the request, then immediately close the connection. This differs from leaveLobby or kickPlayer in that those messages will remove the player from the underlying lobby and inform the rest of the Room of such, while a disconnect message simply closes the connection
```json
{
  "jsonType": "disconnect",
  "data": {
    "the data field": "is ignored"
  }
}
```

# Messages From the Server
Any websocket message sent from the server will follow this general outline:
```json
{
  "type": "string",
  "data": {
    "the underlying message": "goes here"
  }
}
```

There are currently 5 types of websocket messages the server might send, which are (in alphabetical order):
- [Changelog](#changelog)
- [Close](#close)
- [Error](#error)
- [GameOver](#gameover)
- [GameState](#gamestate)
- [LobbyInfo](#lobbyinfo)

### Changelog
Changelog messages are sent out after a client sends a "submitAction" message. They follow this structure:
```json
{
  "type": "Changelog",
  "data": {
    "views": ["an array containing any views that might have been affected by the most recent SubmittedAction"],
    "currentPlayer": "the id of the Player whose turn it is after applying the most recent SubmittedAction",
    "mostRecentAction": "A string describing the most recent action that just took place. Will be empty if the most recent SubmittedAction had no effect"
  }
}
```

### Close
A Close message is sent out any time the server is about to terminate a websocket connection. The server will immediately close a websocket connection after sending a Close message. Currently, there are 4 cases in which this might happen:
- The host sends a [endGame](#endgame) message, in which case, every connection will receive a Close message
- The host sends a [kickPlayer](#kickplayer) message, in which case the affected player will receive a Close message, and every other player will receive a [LobbyInfo](#lobbyinfo) message to reflect the new state of the lobby.
- A player sends a [leaveLobby](#leavelobby) message, in which case that player will receive a Close message, and every other player will receive a [LobbyInfo](#lobbyinfo) message to reflect the new state of the lobby
- A player sends a [disconnect](#disconnect) message, in which case that player will receive a Close message
```json
{
  "type": "Close",
  "data": {
    "message": "A message about why the connection is closing i.e. Player left the lobby or something"
  }
}
```

### Error
Error messages are returned anytime a client submits a message that is deemed invalid by the server. The data object will contain a message with more details
```json
{
  "type": "Error",
  "data": {
    "message": "The error message"
  }
}
```

### GameOver
A GameOver message is sent as the first of two messages in response to an [endGame](#endgame) message from the host. It currently has nothing in the data field, but the Type is set to "GameOver"
```json
{
  "type": "GameOver",
  "data": {}
}
```

### GameState
GameState messages are sent to everyone in a lobby in response to the host submitting a "startGame" message, or to a client that has just reconnected via the /rejoinLobby endpoint if the game has already started. A rejoining client is sent a [LobbyInfo](#lobbyinfo) if the game has not started yet
```json
{
  "type": "GameState",
  "data": {
    "a gamestate object": "you know what these look like"
  }
}
```

### LobbyInfo
LobbyInfo messages are sent out to a player who has just connected to a lobby by either hosting or joining. They are also sent out to every player in a lobby any time another client connects or is disconnected by the server. The "playerID" field will be empty if this is being sent out in response to a player being removed from the lobby
```json
{
  "type": "LobbyInfo",
  "data": {
    "playerID": "the player's newly assigned ID",
    "lobbyInfo": "a lobby object. You know what these look like"
  }
}
```
