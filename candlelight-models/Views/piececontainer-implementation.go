package Views

import (
	"fmt"
	"math/rand"
	"slices"
)

// =========================Deck Implementation=======================
func (d *Deck) AddPiece(pieceToAdd Card) {
	//Cards in a deck should NOT be flipped, so ensure it's not flipped here
	pieceToAdd.Flipped = false
	d.Cards = append(d.Cards, pieceToAdd)
}

// Removes any pieces from this Deck's collection with an ID matching [pieceToRemove] using DeleteFunc
func (d *Deck) RemovePiece(pieceToRemove Card) {
	//Note: Because DeleteFunc is used, the data is 0'd out instead of removed. This might cause unintended side-effects
	d.Cards = slices.DeleteFunc(d.Cards, func(c Card) bool { return c.Id() == pieceToRemove.Id() })
}

func (d *Deck) FindPiece(id string) (*Card, error) {
	foundIndex := slices.IndexFunc(d.Cards, func(c Card) bool { return c.Id() == id })
	if foundIndex != -1 {
		return &d.Cards[foundIndex], nil
	}
	return nil, fmt.Errorf("could not find Card with ID == {%s}", id)
}

func (d *Deck) PickRandomPiece() *Card {
	index := rand.Intn(len(d.Cards))
	return &d.Cards[index]
}

func (d *Deck) Type() string {
	return Type_Deck
}

// ===============================Space Implementation========================
func (s *Space) AddPiece(pieceToAdd Meeple) {
	s.Meeples = append(s.Meeples, pieceToAdd)
}

// Removes any pieces from this Space's collection with an ID matching [pieceToRemove] using DeleteFunc
func (s *Space) RemovePiece(pieceToRemove Meeple) {
	//Note: Because DeleteFunc is used, the data is 0'd out instead of removed. This might cause unintended side-effects
	s.Meeples = slices.DeleteFunc(s.Meeples, func(m Meeple) bool { return m.Id() == pieceToRemove.Id() })
}

func (s *Space) FindPiece(id string) (*Meeple, error) {
	foundIndex := slices.IndexFunc(s.Meeples, func(c Meeple) bool { return c.Id() == id })
	if foundIndex != -1 {
		return &s.Meeples[foundIndex], nil
	}
	return nil, fmt.Errorf("could not find Meeple with ID == {%s}", id)
}

func (s *Space) PickRandomPiece() *Meeple {
	index := rand.Intn(len(s.Meeples))
	return &s.Meeples[index]
}

func (s *Space) Type() string {
	return Type_Space
}

//========================CardZone Implementation====================

func (cz *CardZone) AddPiece(pieceToAdd Card) {
	cz.Cards = append(cz.Cards, pieceToAdd)
}

// Removes any pieces from this CardZone's collection with an ID matching [pieceToRemove] using DeleteFunc
func (cz *CardZone) RemovePiece(pieceToRemove Card) {
	//Note: Because DeleteFunc is used, the data is 0'd out instead of removed. This might cause unintended side-effects
	cz.Cards = slices.DeleteFunc(cz.Cards, func(c Card) bool { return c.Id() == pieceToRemove.Id() })
}

func (cz *CardZone) FindPiece(id string) (*Card, error) {
	foundIndex := slices.IndexFunc(cz.Cards, func(c Card) bool { return c.Id() == id })
	if foundIndex != -1 {
		return &cz.Cards[foundIndex], nil
	}
	return nil, fmt.Errorf("could not find Card with ID == {%s}", id)
}

func (cz *CardZone) PickRandomPiece() *Card {
	index := rand.Intn(len(cz.Cards))
	return &cz.Cards[index]
}

func (cz *CardZone) Type() string {
	return Type_CardZone
}
