package Views

import (
	"encoding/json"
	"fmt"
)

// Tries to quietly Marshal the given [element], taking the given error and discarding it, unless it's not nil. If nil, will panic
func ToJsonRawMessage(element any) json.RawMessage {
	asJson, err := json.Marshal(element)
	if err != nil {
		panic(err)
	}
	return asJson
}

// Will try to convert the given [element]'s Element field to a struct of the given type T. The given type must line up
// with [element]'s Type field, or this function will panic
func FromJsonRawMessage[T View | Navigation | Deck | Space | CardZone | Card | Meeple | Die](element UI_Element) T {
	var toReturn T

	//Make sure the given T type matches element's Type
	typeMismatch := false
	switch any(toReturn).(type) {
	case View:
		if element.Type != Type_View {
			typeMismatch = true
		}
	case Navigation:
		if element.Type != Type_Navigation {
			typeMismatch = true
		}
	case Deck:
		if element.Type != Type_Deck {
			typeMismatch = true
		}
	case Space:
		if element.Type != Type_Space {
			typeMismatch = true
		}
	case CardZone:
		if element.Type != Type_CardZone {
			typeMismatch = true
		}
	case Card:
		if element.Type != Type_Card {
			typeMismatch = true
		}
	case Meeple:
		if element.Type != Type_Meeple {
			typeMismatch = true
		}
	case Die:
		if element.Type != Type_Die {
			typeMismatch = true
		}
	}

	if typeMismatch {
		panic(fmt.Sprintf("Given type does not match element.Type of {%s}!", element.Type))
	}

	err := json.Unmarshal(element.Element, &toReturn)
	if err != nil {
		panic(err)
	}

	return toReturn
}

// My solution for wanting an Element's Id to be visible from the UI_Element without having to unmarshall the whole thing. Since I know every single one will
// have the Id field, I just partially unmarshall it into a special struct to see the Id
func (el UI_Element) GetId() string {
	unmarshaled := struct {
		Id string `json:"id"`
	}{}
	err := json.Unmarshal(el.Element, &unmarshaled)
	if err != nil {
		panic(err)
	}
	return unmarshaled.Id
}

func (el UI_Element) GetCards() []Card {
	cardHolder := struct {
		Cards []Card
	}{}
	err := json.Unmarshal(el.Element, &cardHolder)
	if err != nil {
		panic(err)
	}
	return cardHolder.Cards
}
