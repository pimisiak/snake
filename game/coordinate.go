package game

type Coordinate struct {
	X, Y int
}

func NewCoordinate(x, y int) Coordinate {
	return Coordinate{x, y}
}
