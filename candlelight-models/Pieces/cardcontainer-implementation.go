package Pieces

import (
	"math/rand"
	"slices"
)

//============Deck Implementation==================

// Attempts to add the given card to Cards. Does no error checking
func (deck *Deck) AddCardToCollection(cardToAdd Card) {
	deck.Cards = append(deck.Cards, cardToAdd)
}

// Attempts to remove any Cards with an ID == [card].Id -- Does not do any error checking
func (deck *Deck) RemoveCardFromCollection(card Card) {
	deck.Cards = slices.DeleteFunc(deck.Cards, func(c Card) bool { return c.Id == card.Id })
}

// Attempts to find a card with the given id in Cards. Returns a pointer to the found card
// if found, or nil otherwise
func (deck *Deck) FindCardInCollection(cardId string) *Card {
	foundIndex := slices.IndexFunc(deck.Cards, func(c Card) bool { return c.Id == cardId })
	if foundIndex != -1 {
		return &deck.Cards[foundIndex]
	}
	return nil
}

func (deck *Deck) PickRandomCardFromCollection() *Card {
	index := rand.Intn(len(deck.Cards))
	return &(deck.Cards[index])
}

func (deck *Deck) CollectionLength() int {
	return len(deck.Cards)
}

func (deck *Deck) CardIsAllowed(card *Card) bool {
	/*
		Tags whitelist is a map[string]string where the Key is a possible tag key for the Card and the value is a
		list of all values that tag on the Card can have
	*/
	for key, values := range deck.TagsWhitelist {
		if card.Tags[key] != "" { //If the card does NOT have a tag with the given key, card.Tags[key] == ""
			//If it DOES exist, it will be some string value. So we check if that string value exists in the list
			//of approved values for this tag within the whitelist
			if slices.ContainsFunc(values, func(s string) bool { return s == card.Tags[key] }) {
				return true
			}
		}
	}
	//If we haven't found a match and the TagsWhitelist is non-empty, return false. If it's empty, of course
	//no match CAN be made so return true
	return len(deck.TagsWhitelist) == 0
}

//============CardPlace Implementation==================

// Attempts to add the given card to PlacedCards. Does no error checking
func (cp *CardPlace) AddCardToCollection(cardToAdd Card) {
	cp.PlacedCards = append(cp.PlacedCards, cardToAdd)
}

// Attempts to remove any Cards with an ID == [card].Id -- Does not do any error checking
func (cp *CardPlace) RemoveCardFromCollection(card Card) {
	cp.PlacedCards = slices.DeleteFunc(cp.PlacedCards, func(c Card) bool { return c.Id == card.Id })
}

// Attempts to find a card with the given id in PlacedCards. Returns a pointer to the found card
// if found, or nil otherwise
func (cp *CardPlace) FindCardInCollection(cardId string) *Card {
	foundIndex := slices.IndexFunc(cp.PlacedCards, func(c Card) bool { return c.Id == cardId })
	if foundIndex != -1 {
		return &cp.PlacedCards[foundIndex]
	}
	return nil
}

func (cp *CardPlace) PickRandomCardFromCollection() *Card {
	index := rand.Intn(len(cp.PlacedCards))
	return &(cp.PlacedCards[index])
}

func (cp *CardPlace) CardIsAllowed(card *Card) bool {
	/*
		Tags whitelist is a map[string]string where the Key is a possible tag key for the Card and the value is a
		list of all values that tag on the Card can have
	*/
	for key, values := range cp.TagsWhitelist {
		if card.Tags[key] != "" { //If the card does NOT have a tag with the given key, card.Tags[key] == ""
			//If it DOES exist, it will be some string value. So we check if that string value exists in the list
			//of approved values for this tag within the whitelist
			if slices.ContainsFunc(values, func(s string) bool { return s == card.Tags[key] }) {
				return true
			}
		}
	}
	//If we haven't found a match and the TagsWhitelist is non-empty, return false. If it's empty, of course
	//no match CAN be made so return true
	return len(cp.TagsWhitelist) == 0
}

func (cp *CardPlace) CollectionLength() int {
	return len(cp.PlacedCards)
}
