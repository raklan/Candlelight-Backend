# Overview
Actions that are sent through the websocket with the intent of the backend applying some action on behalf of a client should be wrapped in a [SubmitAction](https://capstone-cs.eng.utah.edu/candlelight/candlelight-backend/-/wikis/Websocket-Communication#submitaction) message. Within the `data` field of that object, clients should put a [SubmittedAction](#submittedaction)-type object.

# SubmittedAction
The SubmittedAction is designed to be place in the `data` field of a [SubmitAction](https://capstone-cs.eng.utah.edu/candlelight/candlelight-backend/-/wikis/Websocket-Communication#submitaction) message. They have the following structure: 
```json
{
  "type": "the string type. See below for options",
  "turn": {"the turn object": "goes here"}
}
```

# Actions
There are currently 6 supported actions that clients can take which will be sent out to other clients as well. One of the following strings should be placed in the `type` field of the SubmittedAction, with a matching object placed in the `turn` field
- ["Insertion"](#insertion)
- ["Withdrawal"](#withdrawal)
- ["Movement"](#movement)
- ["EndTurn"](#endturn)
- ["Cardflip"](#cardflip)
- ["Reshuffle"](#reshuffle)

## Insertion
An Insertion is defined as a Player inserting an Orphan into some Card Collection (Currently Decks or CardPlaces). They have the following structure:
```json
{
  "insertCard": "id of the card being inserted",
  "fromView": "the id of the View which [insertCard] is an Orphan of _before_ the insertion",
  "toCollection": "the id of the Collection into which [insertCard] should be inserted into",
  "inView": "the id of the View which [toCollection] belongs to"
}
```

## Withdrawal
A Withdrawal is defined as a Player moving a Card out of a Card Collection into the Orphans of a given View.
```json
{
  "withdrawCard": "the id of the card to withdraw. If left blank, a random card is chosen from the given collection instead",
  "fromCollection": "the id of the collection that [withdrawCard] is to be taken from",
  "inView": "the id of the View to which [fromCollection] belongs",
  "toView": "the id of the View that [withdrawCard] should be moved into as an Orphan. Can be the same as [inView]"
}
```

## Movement
A Movement is defined as a Player moving an Orphan from one (x,y) position to another (x,y) position, optionally between Views as well.
```json
{
  "cardId": "the id of the Card to move",
  "fromView": "the id of the View that [cardId] belongs to _before_ moving",
  "toView": "the id of the View that [cardId] should be moved into. Can be the same as [fromView], if desired",
  "atX": "the new x coordinate that should be assigned to [cardId]",
  "atY": "the new y coordinate that should be assigned to [cardId]"
}
```

## EndTurn
An EndTurn is only used if the Game's rules have been marked with `EnforceTurnOrder` as true. EndTurn will end the "Turn" of the current player, updating the GameState and putting the id of the next player whose turn it is into the Changelog within the `currentPlayer` field.
```json
{
  "nextPlayer": "An optional string specifying the ID of the player who should be given the next turn. If left blank, the turn will pass to whichever player is next in the GameState's player list, wrapping around in the event of the last playing submitting an EndTurn"
}
```

## Cardflip
A Cardflip will reverse the `flipped` boolean of a given card, specifying that whichever side is currently not showing on the card should now be face-up
```json
{
  "flipCard": "the id of the card to flip. This card must be an Orphan",
  "inView": "the id of the View in which [flipCard] can be found"
}
```

## Reshuffle
A reshuffle will move all cards from a given CardPlace into a given Deck.
```json
{
  "shuffleCardPlace": "the id of the CardPlace to reshuffle the cards from",
  "inView": "the id of the View in which [shuffleCardPlace] is found",
  "toView": "the id of the View in which [intoDeck] is found",
  "intoDeck": "the id of the Deck to reshuffle into"
}
```
