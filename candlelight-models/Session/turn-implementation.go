package Session

import (
	"candlelight-models/Actions"
	"candlelight-models/Pieces"
	"candlelight-models/Player"
	"fmt"
	"slices"
)

func (mt MoveTurn) Execute(gameState *GameState, player Player.Player) (Pieces.PieceSet, error) {
	//Make sure the Player trying to play this is tracked in the gameState, and if so, grab his state off
	//the gameState to make updates there
	currentPlayerState := findPlayerInGameState(gameState, &player)
	changelog := Pieces.PieceSet{}

	if currentPlayerState == nil {
		return changelog, fmt.Errorf("ERROR: Could not find PlayerState within GameState for given player")
	}

	//Make sure the Action they're trying to take is within their allowed actions
	if !actionInAllowedActions(mt.ActionId, &currentPlayerState.AllowedActions, MoveTurnType) {
		return changelog, fmt.Errorf("ERROR: Attempted action not found within Player's allowed actions")
	}

	target := findCardContainerInGame(mt.TargetId, gameState)
	if target == nil {
		return changelog, fmt.Errorf("couldn't find target with targetId == {%s} in Game", mt.TargetId)
	}

	removedCard, removedFrom := attemptToRemoveCardFromGame(mt.PieceId, gameState)
	if removedCard == nil {
		return changelog, fmt.Errorf("couldn't find card with PieceId == {%s} in GameState", mt.PieceId)
	}

	target.AddCardToCollection(*removedCard)

	addCardContainerToChangelog(&changelog, removedFrom)
	addCardContainerToChangelog(&changelog, target)

	return changelog, nil
}

func (puT PieceUpdateTurn) Execute(gameState *GameState, player Player.Player) (Pieces.PieceSet, error) {
	currentPlayerState := findPlayerInGameState(gameState, &player)
	changelog := Pieces.PieceSet{}

	if currentPlayerState == nil {
		return changelog, fmt.Errorf("ERROR: Could not find PlayerState within GameState for given player")
	}

	if !actionInAllowedActions(puT.ActionId, &currentPlayerState.AllowedActions, PieceUpdateTurnType) {
		return changelog, fmt.Errorf("ERROR: Attempted action not found within Player's allowed actions")
	}

	for i := range gameState.Views {
		for index, cp := range gameState.Views[i].Pieces.CardPlaces {
			if puT.TargetPieceId == cp.Id {
				gameState.Views[i].Pieces.CardPlaces[index].Tags = puT.NewTags
				addCardContainerToChangelog(&changelog, &gameState.Views[i].Pieces.CardPlaces[index])
				return changelog, nil
			}
		}
		for index, deck := range gameState.Views[i].Pieces.Decks {
			if puT.TargetPieceId == deck.Id {
				gameState.Views[i].Pieces.Decks[index].Tags = puT.NewTags
				addCardContainerToChangelog(&changelog, &gameState.Views[i].Pieces.Decks[index])
				return changelog, nil
			}
		}

	}

	return changelog, fmt.Errorf("could not find target with TargetPieceId=={%s}", puT.TargetPieceId)
}

func (pt PlacementTurn) Execute(gameState *GameState, player Player.Player) (Pieces.PieceSet, error) {

	currentPlayerState := findPlayerInGameState(gameState, &player)
	changelog := Pieces.PieceSet{}

	if currentPlayerState == nil {
		return changelog, fmt.Errorf("ERROR: Could not find PlayerState within GameState for given player")
	}

	//Make sure the Action they're trying to take is within their allowed actions
	if !actionInAllowedActions(pt.ActionId, &currentPlayerState.AllowedActions, PlacementTurnType) {
		return changelog, fmt.Errorf("ERROR: Attempted action not found within Player's allowed actions")
	}

	target := findCardContainerInGame(pt.TargetId, gameState)
	if target == nil {
		return changelog, fmt.Errorf("couldn't find target with targetId == {%s} in Game", pt.TargetId)
	}

	removedCard, removedFrom := attemptToRemoveCardFromPlayer(pt.PieceId, &currentPlayerState.Player)
	if removedCard == nil {
		return changelog, fmt.Errorf("couldn't find card with PieceId == {%s} in player's hand", pt.PieceId)
	}

	target.AddCardToCollection(*removedCard)

	addCardContainerToChangelog(&changelog, removedFrom)
	addCardContainerToChangelog(&changelog, target)

	return changelog, nil
}

func (tt TakeTurn) Execute(gameState *GameState, player Player.Player) (Pieces.PieceSet, error) {
	currentPlayerState := findPlayerInGameState(gameState, &player)
	changelog := Pieces.PieceSet{}

	if currentPlayerState == nil {
		return changelog, fmt.Errorf("ERROR: Could not find PlayerState within GameState for given player")
	}

	if !actionInAllowedActions(tt.ActionId, &currentPlayerState.AllowedActions, TakeTurnType) {
		return changelog, fmt.Errorf("ERROR: Attempted action not found within Player's allowed actions")
	}

	target := findCardContainerInGame(tt.TakingFromId, gameState)

	if target == nil {
		return changelog, fmt.Errorf("could not find target with Id == {%s}", tt.TakingFromId)
	}

	var cardToTake *Pieces.Card = nil
	if tt.PieceId == "" {
		//PieceId is empty, so take a random card
		cardToTake = target.PickRandomCardFromCollection()
	} else {
		cardToTake = target.FindCardInCollection(tt.PieceId)
	}

	if cardToTake == nil { //If we haven't found it, ERROR
		return changelog, fmt.Errorf("could not find piece with ID == {%s} within target with ID == {%s}", tt.PieceId, tt.TakingFromId)
	}

	copy := *cardToTake
	target.RemoveCardFromCollection(*cardToTake)
	//Add card to player's Orphaned cards
	currentPlayerState.Player.Hand[0].Pieces.Orphans.Cards = append(currentPlayerState.Player.Hand[0].Pieces.Orphans.Cards, copy)
	addCardContainerToChangelog(&changelog, target)
	changelog.Orphans.Cards = append(changelog.Orphans.Cards, copy)

	return changelog, nil

}

func (mt TradeTurn) Execute(gameState *GameState, player Player.Player) (Pieces.PieceSet, error) {
	return Pieces.PieceSet{}, fmt.Errorf("TradeTurn is not implemented yet")
}

func (tt TransitionTurn) Execute(gameState *GameState, player Player.Player) (Pieces.PieceSet, error) {

	currentPlayerState := findPlayerInGameState(gameState, &player)
	changelog := Pieces.PieceSet{}

	if currentPlayerState == nil {
		return changelog, fmt.Errorf("ERROR: Could not find PlayerState within GameState for given player")
	}

	//Make sure the Action they're trying to take is within their allowed actions
	if !actionInAllowedActions(tt.ActionId, &currentPlayerState.AllowedActions, PlacementTurnType) {
		return changelog, fmt.Errorf("ERROR: Attempted action not found within Player's allowed actions")
	}

	return changelog, nil
}

// Checks if the given [actionId] is found on any of the [allowedActions] of type [actionType]
func actionInAllowedActions(actionId string, allowedActions *Actions.ActionSet, actionType string) bool {
	return true
	// switch actionType {
	// case MoveTurnType:
	// 	return slices.ContainsFunc(allowedActions.Moves, func(action Actions.Move) bool { return action.Id == actionId })
	// case PieceUpdateTurnType:
	// 	return slices.ContainsFunc(allowedActions.PieceUpdates, func(action Actions.PieceUpdate) bool { return action.Id == actionId })
	// case PlacementTurnType:
	// 	return slices.ContainsFunc(allowedActions.Placements, func(action Actions.Placement) bool { return action.Id == actionId })
	// case TakeTurnType:
	// 	return slices.ContainsFunc(allowedActions.Takes, func(action Actions.Take) bool { return action.Id == actionId })
	// case TradeTurnType:
	// 	return slices.ContainsFunc(allowedActions.Trades, func(action Actions.Trade) bool { return action.Id == actionId })
	// case TransitionTurnType:
	// 	return slices.ContainsFunc(allowedActions.Transitions, func(action Actions.Transition) bool { return action.Id == actionId })
	// default:
	// 	return false
	// }
}

func findPlayerInGameState(gameState *GameState, player *Player.Player) *PlayerState {
	for index, element := range gameState.PlayerStates {
		if element.Player.Id == player.Id {
			return &gameState.PlayerStates[index]
		}
	}
	return nil
}

// func attemptToRemoveCardFromPlayer(pieceId string, player *Player.Player) *Views.Card {
// 	for _, el := range player.Hand.Children {
// 		foundIndex := slices.IndexFunc(el.GetCards(), func(c Views.Card) bool { return c.Id() == pieceId })
// 		if foundIndex != -1 {
// 			switch el.Type {
// 			case Views.Type_Deck:
// 				deck := Views.FromJsonRawMessage[Views.Deck](el)
// 				card := &deck.Cards[foundIndex]
// 				copy := *card
// 				deck.RemovePiece(deck.Cards[foundIndex])
// 				return &copy
// 			case Views.Type_CardZone:
// 				cardZone := Views.FromJsonRawMessage[Views.CardZone](el)
// 				card := &cardZone.Cards[foundIndex]
// 				copy := *card
// 				cardZone.RemovePiece(cardZone.Cards[foundIndex])
// 				return &copy
// 			}
// 		}
// 	}
// 	return nil
// }

func attemptToRemoveCardFromPlayer(pieceId string, player *Player.Player) (*Pieces.Card, Pieces.Card_Container) {
	for i := range player.Hand {
		//Try to find the card in their CardPlaces
		for index := range player.Hand[i].Pieces.CardPlaces {
			cardPlace := &player.Hand[i].Pieces.CardPlaces[index]
			foundCard := cardPlace.FindCardInCollection(pieceId)
			if foundCard != nil { //[foundCard] != nil if the card was found
				//Need to make a copy because RemoveCardFromCollection uses slices.Delete, which will 0 out the
				//card at the address pointed to by foundCard
				copy := *foundCard
				cardPlace.RemoveCardFromCollection(*foundCard)
				return &copy, cardPlace
			}
		}

		for index := range player.Hand[i].Pieces.Decks {
			deck := &player.Hand[i].Pieces.Decks[index]
			foundCard := deck.FindCardInCollection(pieceId)
			if foundCard != nil { //[foundCard] != nil if the card was found
				//Need to make a copy because RemoveCardFromCollection uses slices.Delete, which will 0 out the
				//card at the address pointed to by foundCard
				copy := *foundCard
				deck.RemoveCardFromCollection(*foundCard)
				return &copy, deck
			}
		}

		foundCard := (&player.Hand[i].Pieces.Orphans).FindCardInCollection(pieceId)
		if foundCard != nil {
			copy := *foundCard
			player.Hand[i].Pieces.Orphans.RemoveCardFromCollection(*foundCard)
			return &copy, &player.Hand[i].Pieces.Orphans
		}
	}

	return nil, nil
}

// func attemptToRemoveCardFromGame(pieceId string, gameState *GameState) *Views.Card {
// 	// for _, view := range gameState.Table {
// 	// 	for _, el := range view.Children {
// 	// 		foundIndex := slices.IndexFunc(el.GetCards(), func(c Views.Card) bool { return c.Id() == pieceId })
// 	// 		if foundIndex != -1 {
// 	// 			switch el.Type {
// 	// 			case Views.Type_Deck:
// 	// 				deck := Views.FromJsonRawMessage[Views.Deck](el)
// 	// 				card := &deck.Cards[foundIndex]
// 	// 				copy := *card
// 	// 				deck.RemovePiece(deck.Cards[foundIndex])
// 	// 				return &copy
// 	// 			case Views.Type_CardZone:
// 	// 				cardZone := Views.FromJsonRawMessage[Views.CardZone](el)
// 	// 				card := &cardZone.Cards[foundIndex]
// 	// 				copy := *card
// 	// 				cardZone.RemovePiece(cardZone.Cards[foundIndex])
// 	// 				return &copy
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }
// 	return nil
// }

// Finds the card with the given [pieceId] in the entire gameState, removes it from the collection it's found in, and returns a copy of it
func attemptToRemoveCardFromGame(pieceId string, gameState *GameState) (*Pieces.Card, Pieces.Card_Container) {
	//Try to find the card in their CardPlaces
	for i := range gameState.Views {
		for index := range gameState.Views[i].Pieces.CardPlaces {
			cardPlace := &gameState.Views[i].Pieces.CardPlaces[index]
			foundCard := cardPlace.FindCardInCollection(pieceId)
			if foundCard != nil { //[foundCard] != nil if the card was found
				//Need to make a copy because RemoveCardFromCollection uses slices.Delete, which will 0 out the
				//card at the address pointed to by foundCard
				copy := *foundCard
				cardPlace.RemoveCardFromCollection(*foundCard)
				return &copy, cardPlace
			}
		}
	}

	for i := range gameState.Views {
		for index := range gameState.Views[i].Pieces.Decks {
			deck := &gameState.Views[i].Pieces.Decks[index]
			foundCard := deck.FindCardInCollection(pieceId)
			if foundCard != nil { //[foundCard] != nil if the card was found
				//Need to make a copy because RemoveCardFromCollection uses slices.Delete, which will 0 out the
				//card at the address pointed to by foundCard
				copy := *foundCard
				deck.RemoveCardFromCollection(*foundCard)
				return &copy, deck
			}
		}
	}

	return nil, nil
}

func findCardContainerInGame(targetId string, gameState *GameState) Pieces.Card_Container {
	foundIndex := -1
	for i := range gameState.Views {
		foundIndex = slices.IndexFunc(gameState.Views[i].Pieces.CardPlaces, func(c Pieces.CardPlace) bool { return c.Id == targetId })
		if foundIndex != -1 {
			return &gameState.Views[i].Pieces.CardPlaces[foundIndex]
		}
		foundIndex = slices.IndexFunc(gameState.Views[i].Pieces.Decks, func(c Pieces.Deck) bool { return c.Id == targetId })
		if foundIndex != -1 {
			return &gameState.Views[i].Pieces.Decks[foundIndex]
		}
	}
	return nil
}

func addCardContainerToChangelog(changelog *Pieces.PieceSet, container Pieces.Card_Container) {
	if deckPointer, ok := container.(*Pieces.Deck); ok {
		deck := *deckPointer
		changelog.Decks = append(changelog.Decks, deck)
	} else if cardPlacePointer, ok := container.(*Pieces.CardPlace); ok {
		cardPlace := *cardPlacePointer
		changelog.CardPlaces = append(changelog.CardPlaces, cardPlace)
	} else {
		panic("Couldn't coerce container into known instance of Card_Container")
	}
}

// func findCardContainerInGame(targetId string, gameState *GameState) Views.PieceContainer[Views.Card] {
// 	// for _, v := range gameState.Table {
// 	// 	for _, el := range v.Children {
// 	// 		if el.GetId() == targetId {
// 	// 			switch el.Type {
// 	// 			case Views.Type_Deck:
// 	// 				deck := Views.FromJsonRawMessage[Views.Deck](el)
// 	// 				return &deck
// 	// 			case Views.Type_CardZone:
// 	// 				zone := Views.FromJsonRawMessage[Views.CardZone](el)
// 	// 				return &zone
// 	// 			default:
// 	// 				return nil
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }
// 	return nil
// }

// func removeUiElementFromView(targetId string, gameState *GameState) *Views.UI_Element {
// 	return nil
// }

// func ExtractFromView_Recursion(targetPieceId string, element *Views.UI_Element) *Views.UI_Element { //TODO: Make private once tested
// 	if element.Type != Views.Type_View {
// 		return nil
// 	}
// 	//Examine all children of this View. If one has the targetId, extract/delete from children, update Serialization, return pointer
// 	//If none have targetId, call self on each View child. Effectively a DFS
// 	var copy Views.UI_Element
// 	found := false
// 	view := Views.FromJsonRawMessage[Views.View](*element)
// 	for _, child := range view.Children {
// 		if child.GetId() == targetPieceId {
// 			copyAddr := &child
// 			copy = *copyAddr
// 			found = true
// 			break
// 		} else {
// 			addr := ExtractFromView_Recursion(targetPieceId, &child)
// 			if addr != nil {
// 				found = true
// 				copy = *addr
// 			}
// 		}
// 	}
// 	if found {
// 		view.Children = slices.DeleteFunc(view.Children, func(e Views.UI_Element) bool { return e.GetId() == targetPieceId })
// 		element.Element = Views.ToJsonRawMessage(view)
// 		return &copy
// 	} else {
// 		return nil
// 	}
// }
