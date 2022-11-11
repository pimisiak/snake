package game

type direction string

const (
	up    direction = "UP"
	down  direction = "DOWN"
	left  direction = "LEFT"
	right direction = "RIGHT"
)

func opposite(dir direction) direction {
	switch dir {
	case up:
		return down
	case down:
		return up
	case left:
		return right
	case right:
		return left
	default:
		return dir
	}
}
