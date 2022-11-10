package game

type Snake []Coordinate

func NewSnake(coor Coordinate) *Snake {
	var snake Snake
	snake = append(snake, coor)
	return &snake
}
