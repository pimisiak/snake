package game

type snake struct {
	body      []coordinate
	direction direction
}

func newSnake(head coordinate, dir direction) *snake {
	var body []coordinate
	body = append(body, head)
	body = append([]coordinate{head.move(dir)}, body[:]...)
	body = append([]coordinate{head.move(dir).move(dir)}, body[:]...)
	return &snake{
		body:      body,
		direction: dir,
	}
}

func (s *snake) head() coordinate {
	return s.body[0]
}

func (s *snake) move() {
	head := []coordinate{s.head().move(s.direction)}
	s.body = append(head, s.body[:len(s.body)-1]...)
}

func (s *snake) eat(apple apple) {
	head := []coordinate{{apple.x, apple.y}}
	s.body = append(head, s.body[:]...)
}

func (s *snake) redirect(dir direction) {
	switch dir {
	case up:
		if s.direction != down {
			s.direction = dir
		}
	case down:
		if s.direction != up {
			s.direction = dir
		}
	case left:
		if s.direction != right {
			s.direction = dir
		}
	case right:
		if s.direction != left {
			s.direction = dir
		}
	}
}
