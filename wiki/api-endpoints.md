# Accounts
## /createAccount
- Method: POST
  - Body: JSON serialization of the following:
  ```json
  {
    "username": "myUsername",
    "password": "myPassword"
  }
  ```
  - On Success: 
    - Status Code: 200
    - Body: JSON serialization of the following:
    ```json
    {
      "username": "myUsername"
    }
    ```
  - On Failure:
    - Status Code: 400 (If any error occurs)
    - Body: Error message

## /login
- Method: POST
  - Body: JSON serialization of the following:
  ```json
  {
    "username": "myUsername",
    "password": "myPassword"
  }
  ```
  - On Success: 
    - Status Code: 200
    - Body: JSON serialization of the following:
    ```json
    {
      "username": "myUsername"
    }
    ```
  - On Failure:
    - Status Codes:
      - 400 (If the body is missing `username` and/or `password`, or if an error occurs trying to deserialize the body)
      - 401 (If the given username and password combination do not represent a user)
    - Body: Error message

# Creator Studio
## /allGames
- Method: GET
  - Query Params:
    - slimmed: boolean _optional_
      - `true`: Sets the objects in the returned array to only have the game name and its ID
      - Any other value or unset: The objects in the returned array are the full game definition objects
  - On Success:
    - Status Code: 200
    - Body: JSON serialization of an array containing game definition objects, slimmed if requested

## /studio
- Method: GET
  - Query Params:
    - id: string **required**
      - The id of the game definition you're requesting
  - On Success:
    - Status Code: 200
    - Body: JSON serialization of the requested Game Definition object
  - On Failure:
    - Status Codes: 
      - 400 (If `id` is missing from the query string)
      - 404 (If a game definition with the given `id` is not found)
    - Body: Error message
- Method: POST
  - Body: JSON serialization of the Game Definition to save. If the object's Id field is empty, one will be assigned and a new entry is added to the DB. If the `Id` field is filled in, the entry of the given Id is overwritten with the given Game Definition
  - On Success:
    - Status Code: 200
    - Body: JSON serialization of the saved Game Definition
  - On Failure:
    - Status Code: 400
    - Body: Error Message
- Method: DELETE
  - Query Params:
    - id: string **required**
      - The id of the game definition to delete
  - On Success:
    - Status Code: 200
    - Body: JSON serialization of the following:
    ```go
    {
      DeletedId: string
    }
    ```

# Gameplay
## /hostLobby
- Method: GET
  - Query Params: 
    - gameId: string **required**
      - The ID of the game definition to create a lobby from. Must be a valid Id
    - playerName: string **required**
      - The display name of the player hosting the lobby. Can be anything.
  - On Success:
    - Connection upgraded to websocket. Connection remains open.
    - Websocket receives JSON serialization of the new Lobby object.
  - On Failure:
    - Status Codes:
      - 400 (If `gameId` and/or `playerName` is missing from the query string)
      - 500 (If anything else goes wrong)
    - Body: Error Message

## /joinLobby
- Method: GET
  - Query Params:
    - roomCode: string **required**
      - The Room Code of a previously created Lobby.
    - playerName: string **required**
      - The display name of the player joining the lobby
  - On Success:
    - Connection upgraded to websocket. Connection remains open.
    - Websocket immediately receives JSON serialization of the joined Lobby object. All other players in Lobby receive an updated Lobby object through their respective websockets
  - On Failure:
    - Status Codes: 
      - 400 (If `roomCode` and/or `playerName` is missing from the query string)
      - 404 (If a lobby with the given room code does not exist)
      - 500 (If websocket upgrade fails for any other reason)
    - Body: Error Message

## /rejoinLobby
- Method: GET
  - Query Params:
    - roomCode: string **required**
      - The Room Code of the Lobby the player is trying to rejoin
    - playerId: string
      - The PlayerId found with the lobby info returned by the /joinLobby handshake
  - On Success:
    - Connection upgraded to websocket. Connection remains open
    - Websocket immediately receives Lobby Info if the server detects that the game has not started yet. If the server detects that the game has started, first message will be the current GameState
  - On Failure:
    - Status Codes:
      - 400 (If the query string is missing `roomCode` and/or `playerId`, or if there is already an open connection for the given `playerId`)
      - 404 (If a lobby with the given `roomCode` or a Player with the given `playerId` cannot be found)
      - 500 (If websocket upgrade fails for any other reason)
    - Body: Error Message

# Misc
## /heartbeat
- Method: GET
  - Query Params: None
  - On Success:
    - Status Code: 200
    - Body: A string reading "Buh-dump, buh-dump"

## /version
- Method: GET
  -  Query Params: None
  - On Success:
    - Status Code: 200
    - Body: A string reading v0.x.y - z
      - X == 1 for Alpha general release, 2 for Beta general release, and final release will be v1.0.y - z
      - Y == The most recent issue number that was implemented
      - Z == The date of the implementation being checked in
