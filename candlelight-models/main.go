package main

import (
	"candlelight-models/Game"
	"candlelight-models/Pieces"
	"candlelight-models/Session"
	"fmt"
)

func main() {
	// cards := []Pieces.Card{
	// 	{
	// 		GamePiece: Pieces.GamePiece{
	// 			Id: "1",
	// 		},
	// 	},
	// 	{
	// 		GamePiece: Pieces.GamePiece{
	// 			Id: "2",
	// 		},
	// 	},
	// }

	// cardCopy := cards[0]

	// newCards := slices.DeleteFunc(cards, func(c Pieces.Card) bool { return c.Id == "1" })

	// fmt.Println(cards)
	// fmt.Println(cardCopy)
	// fmt.Println(newCards)
	gameState := Session.GameState{
		Views: []Game.View{
			{
				Id: "1",
				Pieces: Pieces.PieceSet{
					Orphans: []Pieces.Card{
						{
							GamePiece: Pieces.GamePiece{
								Id: "1",
							},
						},
					},
				},
			},
			{
				Id: "2",
			},
		},
	}

	copy := makeCopy(gameState)

	fmt.Println(copy)
	gameState.Views[0].Pieces.Orphans[0].Id = "changed"
	fmt.Println(copy)

}

func makeCopy(toCopy Session.GameState) Session.GameState {
	return toCopy
}
