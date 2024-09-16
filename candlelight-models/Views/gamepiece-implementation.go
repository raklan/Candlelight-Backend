package Views

//====================Card Implementation===================
func (c Card) Id() string {
	return c.GamePiece.Id
}

func (c Card) Type() string {
	return Type_Card
}

func (c Card) Tags() map[string]string {
	return c.GamePiece.Tags
}

//====================Meeple Implementation==================

func (m Meeple) Id() string {
	return m.GamePiece.Id
}

func (m Meeple) Type() string {
	return Type_Meeple
}

func (m Meeple) Tags() map[string]string {
	return m.GamePiece.Tags
}

//======================Die Implementation=====================

func (d Die) Id() string {
	return d.GamePiece.Id
}

func (d Die) Type() string {
	return Type_Die
}

func (d Die) Tags() map[string]string {
	return d.GamePiece.Tags
}
