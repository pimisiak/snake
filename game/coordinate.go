package game

type coordinate struct {
	x, y int
}

func (c coordinate) move(direction direction) coordinate {
	switch direction {
	case up:
		return coordinate{c.x, c.y - 1}
	case down:
		return coordinate{c.x, c.y + 1}
	case left:
		return coordinate{c.x - 1, c.y}
	case right:
		return coordinate{c.x + 1, c.y}
	default:
		return c
	}
}

func (c coordinate) equals(other coordinate) bool {
	return c.x == other.x && c.y == other.y
}
