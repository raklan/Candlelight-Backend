package Session

import (
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-models/Player"
	"fmt"
	"slices"
)

func (ins Insertion) Execute(gameState *GameState, playerId string) (Changelog, error) {
	changelog := Changelog{
		CurrentPlayer: gameState.CurrentPlayer,
	}
	var err error = nil

	playerToUse := findPlayerInGameState(playerId, gameState)
	if playerToUse == nil {
		return changelog, fmt.Errorf("could not find player in gamestate")
	}

	//IMPORTANT: DO ALL ERROR-CHECKING BEFORE CHANGING THE GAMESTATE

	//Check for Views first so we can get them in the Changelog if applicable
	takingFromView := findView(gameState, playerToUse, ins.FromView)
	if takingFromView == nil {
		err = fmt.Errorf("could not find View to take from with Id == {%s}", ins.FromView)
	} else {
		changelog.Views = append(changelog.Views, takingFromView)
	}

	intoView := findView(gameState, playerToUse, ins.InView)
	if intoView == nil {
		err = fmt.Errorf("could not find View to insert into with Id == {%s}", ins.InView)
	} else {
		changelog.Views = append(changelog.Views, intoView)
	}

	if err != nil {
		return changelog, err
	}

	cardToInsert := findCardInOrphans(ins.InsertCard, takingFromView)
	if cardToInsert == nil {
		return changelog, fmt.Errorf("could not find card to insert with Id == {%s} in given View", ins.InsertCard)
	}

	intoCollection := findCollectionInView(ins.ToCollection, intoView)
	if intoCollection == nil {
		return changelog, fmt.Errorf("could not find Collection to insert with with Id == {%s} in given View", ins.ToCollection)
	}

	//Also important. Copy the card since slices.DeleteFunc will 0 out that location in memory, and update the copy with its new ParentViewId
	cardCopy := *cardToInsert
	cardCopy.ParentView = intoView.Id
	//Remove the card from the View it's being removed from
	takingFromView.Pieces.Orphans = slices.DeleteFunc(takingFromView.Pieces.Orphans, func(c Pieces.Card) bool {
		return c.Id == cardToInsert.Id
	})
	//Insert that card into its new collection. Because this is a pointer, it should match up to the right place
	intoCollection.AddCardToCollection(cardCopy)

	return changelog, nil
}

func (with Withdrawal) Execute(gameState *GameState, playerId string) (Changelog, error) {
	changelog := Changelog{
		CurrentPlayer: gameState.CurrentPlayer,
	}
	var err error = nil

	//IMPORTANT: DO ALL ERROR-CHECKING BEFORE CHANGING THE GAMESTATE

	playerToUse := findPlayerInGameState(playerId, gameState)
	if playerToUse == nil {
		return changelog, fmt.Errorf("could not find player in gamestate")
	}

	//Check for Views first so we can get them in the Changelog if applicable
	takingFromView := findView(gameState, playerToUse, with.InView)
	if takingFromView == nil {
		err = fmt.Errorf("could not find View to take from with Id == {%s}", with.InView)
	} else {
		changelog.Views = append(changelog.Views, takingFromView)
	}

	intoView := findView(gameState, playerToUse, with.ToView)
	if intoView == nil {
		err = fmt.Errorf("could not find View to insert into with Id == {%s}", with.ToView)
	} else {
		changelog.Views = append(changelog.Views, intoView)
	}

	if err != nil {
		return changelog, err
	}

	fromCollection := findCollectionInView(with.FromCollection, takingFromView)
	if fromCollection == nil {
		return changelog, fmt.Errorf("could not find Collection to withdraw from with with Id == {%s} in given View", with.FromCollection)
	}

	var cardToWithdraw *Pieces.Card = nil
	if with.WithdrawCard == "" {
		cardToWithdraw = fromCollection.PickRandomCardFromCollection()
	} else {
		cardToWithdraw = fromCollection.FindCardInCollection(with.WithdrawCard)
	}
	if cardToWithdraw == nil {
		return changelog, fmt.Errorf("could not find Card to withdraw with Id == {%s} in given Collection", with.WithdrawCard)
	}

	//Also important: Make a copy of the card and update the copy's ParentViewId, as well as setting the Position to the Collection's X/Y to make it appear on top
	cardCopy := *cardToWithdraw
	cardCopy.ParentView = intoView.Id
	x, y := fromCollection.GetXY()
	cardCopy.X = x
	cardCopy.Y = y

	//Remove card from the collection it's being withdrawn from and add to the Orphans of the appropriate View
	fromCollection.RemoveCardFromCollection(*cardToWithdraw)
	intoView.Pieces.Orphans = append(intoView.Pieces.Orphans, cardCopy)

	return changelog, nil
}

func (move Movement) Execute(gameState *GameState, playerId string) (Changelog, error) {
	changelog := Changelog{
		CurrentPlayer: gameState.CurrentPlayer,
	}
	var err error = nil

	playerToUse := findPlayerInGameState(playerId, gameState)
	if playerToUse == nil {
		return changelog, fmt.Errorf("could not find player in gamestate")
	}

	takingFromView := findView(gameState, playerToUse, move.FromView)
	if takingFromView == nil {
		err = fmt.Errorf("could not find View to take from with Id == {%s}", move.FromView)
	} else {
		changelog.Views = append(changelog.Views, takingFromView)
	}

	intoView := findView(gameState, playerToUse, move.ToView)
	if intoView == nil {
		err = fmt.Errorf("could not find View to insert into with Id == {%s}", move.ToView)
	} else {
		changelog.Views = append(changelog.Views, intoView)
	}

	if err != nil {
		return changelog, err
	}

	pieceToMove := findCardInOrphans(move.CardId, takingFromView)
	if pieceToMove == nil {
		return changelog, fmt.Errorf("could not find Card to move with Id == {%s} in given View", move.CardId)
	}

	//Copy card and update ParentViewId and Position data
	copy := *pieceToMove
	copy.ParentView = intoView.Id
	copy.X = move.AtX
	copy.Y = move.AtY

	//Update appropriate Views, set changelog and return
	takingFromView.Pieces.Orphans = slices.DeleteFunc(takingFromView.Pieces.Orphans, func(c Pieces.Card) bool { return c.Id == pieceToMove.Id })
	intoView.Pieces.Orphans = append(intoView.Pieces.Orphans, copy)

	return changelog, nil
}

func (et EndTurn) Execute(gameState *GameState, playerId string) (Changelog, error) {
	changelog := Changelog{
		CurrentPlayer: gameState.CurrentPlayer,
	}

	//Get index of player whose turn it is
	currentPlayerIndex := slices.IndexFunc(gameState.Players, func(p Player.Player) bool { return p.Id == gameState.CurrentPlayer })

	if currentPlayerIndex < 0 {
		return changelog, fmt.Errorf("could not find current player in gameState.Players")
	}

	nextPlayerIndex := currentPlayerIndex //Default to not changing anything if something goes wrong
	if et.NextPlayer != "" {
		nextPlayerIndex = slices.IndexFunc(gameState.Players, func(p Player.Player) bool { return p.Id == et.NextPlayer })

		if nextPlayerIndex < 0 {
			return changelog, fmt.Errorf("could not find player to give turn to")
		}
	} else {
		//Get ID of next player (wrap if necessary)
		nextPlayerIndex = currentPlayerIndex + 1
		if nextPlayerIndex >= len(gameState.Players) {
			nextPlayerIndex -= len(gameState.Players)
		}
	}

	nextPlayerId := gameState.Players[nextPlayerIndex].Id

	//Update gameState and changelog
	gameState.CurrentPlayer = nextPlayerId
	changelog.CurrentPlayer = nextPlayerId

	//return
	return changelog, nil
}

func (cf Cardflip) Execute(gameState *GameState, playerId string) (Changelog, error) {
	changelog := Changelog{
		CurrentPlayer: gameState.CurrentPlayer,
	}

	player := findPlayerInGameState(playerId, gameState)

	if player == nil {
		return changelog, fmt.Errorf("could not find player in GameState")
	}

	parentView := findView(gameState, player, cf.InView)

	if parentView == nil {
		return changelog, fmt.Errorf("could not find view with Id == {%s} in GameState", cf.InView)
	}

	changelog.Views = append(changelog.Views, parentView)

	cardToFlip := findCardInOrphans(cf.FlipCard, parentView)

	if cardToFlip == nil {
		return changelog, fmt.Errorf("could not find card with Id == {%s} in View", cf.FlipCard)
	}

	cardToFlip.Facedown = !cardToFlip.Facedown

	return changelog, nil
}

func (reshuffle Reshuffle) Execute(gameState *GameState, playerId string) (Changelog, error) {
	changelog := Changelog{
		CurrentPlayer: gameState.CurrentPlayer,
	}
	var err error = nil

	playerToUse := findPlayerInGameState(playerId, gameState)
	if playerToUse == nil {
		return changelog, fmt.Errorf("could not find player in gamestate")
	}

	//IMPORTANT: DO ALL ERROR-CHECKING BEFORE CHANGING THE GAMESTATE

	//Check for Views first so we can get them in the Changelog if applicable
	takingFromView := findView(gameState, playerToUse, reshuffle.InView)
	if takingFromView == nil {
		err = fmt.Errorf("could not find View to take from with Id == {%s}", reshuffle.InView)
	} else {
		changelog.Views = append(changelog.Views, takingFromView)
	}

	toView := findView(gameState, playerToUse, reshuffle.ToView)
	if toView == nil {
		err = fmt.Errorf("could not find View to insert into with Id == {%s}", reshuffle.ToView)
	} else {
		changelog.Views = append(changelog.Views, toView)
	}

	if err != nil {
		return changelog, err
	}

	reshuffleCollection := findCollectionInView(reshuffle.ShuffleCardPlace, takingFromView)

	reshuffleCardPlace, ok := reshuffleCollection.(*Pieces.CardPlace)
	if !ok {
		return changelog, fmt.Errorf("could not find CardPlace to reshuffle in given View")
	}

	intoCollection := findCollectionInView(reshuffle.IntoDeck, toView)
	reshuffleDeck, ok := intoCollection.(*Pieces.Deck)
	if !ok {
		return changelog, fmt.Errorf("could not find Deck to reshuffle into with Id == {%s} in given View", reshuffle.IntoDeck)
	}

	transferAllCards(reshuffleCardPlace, reshuffleDeck)

	return changelog, nil
}

func findView(gameState *GameState, player *Player.Player, viewId string) *Game.View {
	//Check public views, then the given player's views
	for index, view := range gameState.Views {
		if view.Id == viewId {
			return &gameState.Views[index]
		}
	}

	for index, view := range player.Hand {
		if view.Id == viewId {
			return &player.Hand[index]
		}
	}

	return nil
}

func findCollectionInView(collectionId string, view *Game.View) Pieces.Card_Container {
	for _, cc := range view.Pieces.GetCollections() {
		if cc.GetId() == collectionId {
			return cc
		}
	}
	return nil
}

func findCardInOrphans(pieceId string, view *Game.View) *Pieces.Card {
	for index, card := range view.Pieces.Orphans {
		if card.Id == pieceId {
			return &view.Pieces.Orphans[index]
		}
	}
	return nil
}

func findPlayerInGameState(playerId string, gameState *GameState) *Player.Player {
	for index, player := range gameState.Players {
		if player.Id == playerId {
			return &gameState.Players[index]
		}
	}
	return nil
}

func transferAllCards(cardPlace *Pieces.CardPlace, deck *Pieces.Deck) {
	//Copy all cards into the deck
	cardCopy := Pieces.Card{}
	for i := range cardPlace.Cards {
		card := &cardPlace.Cards[i]
		cardCopy = *card
		deck.Cards = append(deck.Cards, cardCopy)
	}
	//Remove all cards from the CardPlace
	cardPlace.Cards = slices.DeleteFunc(cardPlace.Cards, func(cp Pieces.Card) bool { return true })
}
