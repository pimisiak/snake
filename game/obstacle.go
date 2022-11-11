package game

type edge struct {
	start coordinate
	end   coordinate
}

type obstacle []edge

func generateObstacles() []obstacle {
	return []obstacle{}
}
