package Session

import (
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-models/Player"
	"fmt"
	"slices"
)

func (ins Insertion) Execute(gameState *GameState, player *Player.Player) (Changelog, error) {

	changelog := Changelog{}

	//IMPORTANT: DO ALL ERROR-CHECKING BEFORE CHANGING THE GAMESTATE
	takingFromView := findView(gameState, player, ins.FromView)
	if takingFromView == nil {
		return changelog, fmt.Errorf("could not find View to take from with Id == {%s}", ins.FromView)
	}

	cardToInsert := findCardInOrphans(ins.InsertCard, takingFromView)
	if cardToInsert == nil {
		return changelog, fmt.Errorf("could not find card to insert with Id == {%s} in given View", ins.InsertCard)
	}

	intoView := findView(gameState, player, ins.InView)
	if intoView == nil {
		return changelog, fmt.Errorf("could not find View to insert into with Id == {%s}", ins.InView)
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

	//Add affected views to the changelog and return
	changelog.Views = append(changelog.Views, *takingFromView)
	changelog.Views = append(changelog.Views, *intoView)

	return changelog, nil
}

func (with Withdrawl) Execute(gameState *GameState, player *Player.Player) (Changelog, error) {
	changelog := Changelog{}

	//IMPORTANT: DO ALL ERROR-CHECKING BEFORE CHANGING THE GAMESTATE
	takingFromView := findView(gameState, player, with.InView)
	if takingFromView == nil {
		return changelog, fmt.Errorf("could not find View to take from with Id == {%s}", with.InView)
	}

	fromCollection := findCollectionInView(with.FromCollection, takingFromView)
	if fromCollection == nil {
		return changelog, fmt.Errorf("could not find Collection to withdraw from with with Id == {%s} in given View", with.FromCollection)
	}

	cardToWithdraw := fromCollection.FindCardInCollection(with.WithdrawCard)
	if cardToWithdraw == nil {
		return changelog, fmt.Errorf("could not find Card to withdraw with Id == {%s} in given Collection", with.WithdrawCard)
	}

	intoView := findView(gameState, player, with.ToView)
	if intoView == nil {
		return changelog, fmt.Errorf("could not find View to insert into with Id == {%s}", with.ToView)
	}

	//Also important: Make a copy of the card and update the copy's ParentViewId
	cardCopy := *cardToWithdraw
	cardCopy.ParentView = intoView.Id

	//Remove card from the collection it's being withdrawn from and add to the Orphans of the appropriate View
	fromCollection.RemoveCardFromCollection(*cardToWithdraw)
	intoView.Pieces.Orphans = append(intoView.Pieces.Orphans, cardCopy)

	//Add the affected views to the changelog and return
	changelog.Views = append(changelog.Views, *takingFromView)
	changelog.Views = append(changelog.Views, *intoView)

	return changelog, nil
}

func (move Movement) Execute(gameState *GameState, player *Player.Player) (Changelog, error) {
	changelog := Changelog{}

	takingFromView := findView(gameState, player, move.FromView)
	if takingFromView == nil {
		return changelog, fmt.Errorf("could not find View to take from with Id == {%s}", move.FromView)
	}

	intoView := findView(gameState, player, move.ToView)
	if intoView == nil {
		return changelog, fmt.Errorf("could not find View to insert into with Id == {%s}", move.ToView)
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

	changelog.Views = append(changelog.Views, *takingFromView)
	changelog.Views = append(changelog.Views, *intoView)

	return changelog, nil
}

func findView(gameState *GameState, player *Player.Player, viewId string) *Game.View {
	//Check public views, then the given player's views
	for _, view := range gameState.Views {
		if view.Id == viewId {
			return &view
		}
	}

	for _, view := range player.Hand {
		if view.Id == viewId {
			return &view
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
	for _, card := range view.Pieces.Orphans {
		if card.Id == pieceId {
			return &card
		}
	}
	return nil
}
