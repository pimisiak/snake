package game

type snake struct {
	body      []coordinate
	direction direction
}

// Creates snake moving in specified direction with head at provided coordinate.
// Inital body length is 3.
func newSnake(head coordinate, dir direction) *snake {
	var body []coordinate
	body = append(body, head)
	body = append(body, body[0].move(opposite(dir)))
	body = append(body, body[1].move(opposite(dir)))
	return &snake{
		body:      body,
		direction: dir,
	}
}

func (s *snake) head() coordinate {
	return s.body[0]
}

func (s *snake) move() {
	head := s.head().move(s.direction)
	s.body = append([]coordinate{head}, s.body[:len(s.body)-1]...)
}

func (s *snake) eat(apple apple) {
	head := []coordinate{{apple.x, apple.y}}
	s.body = append(head, s.body[:]...)
}

func (s *snake) redirect(dir direction) {
	if opposite(s.direction) != dir {
		s.direction = dir
	}
}
