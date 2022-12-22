package game

import (
	"math"
	"math/rand"
)

const (
	obstacleRate = 0.0195
)

func generateObstacles(width, height int, random rand.Rand) [][]int {
	obstacles := make([][]int, height)
	for i := range obstacles {
		obstacles[i] = make([]int, width)
	}

	// base count on screen size
	count := int(math.Floor(obstacleRate * float64(width) * float64(height)))

	// generate the specified number of obstacles
	for i := 0; i < count; i++ {

		// choose a random starting coordinate
		x := random.Intn(width)
		y := random.Intn(height)

		// continue until it is out of snake safe starting zone
		for isInSafeZone(x, y, width, height) {
			x = random.Intn(width)
			y = random.Intn(height)
		}

		// generate the obstacle cluster
		obstacles[y][x] = 1
		current := coordinate{x, y}

		for {
			end := random.Intn(1000)
			if end >= 975 {
				break
			}

			// choose a random direction to move in
			dir := directions[rand.Intn(len(directions))]
			current = current.move(dir)

			// if the new location is out of bounds, continue to the next iteration
			if current.x < 0 || current.x >= width || current.y < 0 || current.y >= height {
				continue
			}

			// if the new location is within snake safe starting zone, continue to the next iteration
			if isInSafeZone(current.x, current.y, width, height) {
				continue
			}

			// if the new location is already part of the cluster, continue to the next iteration
			if obstacles[current.y][current.x] == 1 {
				continue
			}

			// set the new location to 1 to indicate it is part of the cluster
			obstacles[current.y][current.x] = 1
		}
	}

	// smooth edges
	for i := 1; i <= 100_000; i++ {
		done := true
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if obstacles[y][x] == 0 {
					continue
				}

				var (
					current = coordinate{x, y}
					lu      = current.move(left).move(up)
					ru      = current.move(right).move(up)
					ld      = current.move(left).move(down)
					rd      = current.move(right).move(down)
					luval   = 0
					ruval   = 0
					ldval   = 0
					rdval   = 0
				)

				if lu.x >= 0 && lu.y >= 0 {
					luval = obstacles[lu.y][lu.x]
				}

				if ru.x < width && ru.y >= 0 {
					ruval = obstacles[ru.y][ru.x]
				}

				if ld.x >= 0 && ld.y < height {
					ldval = obstacles[ld.y][ld.x]
				}

				if rd.x < width && rd.y < height {
					rdval = obstacles[rd.y][rd.x]
				}

				if luval == 1 && ruval == 1 && obstacles[y-1][x] != 1 {
					obstacles[y-1][x] = 1
					done = false
				}

				if ldval == 1 && rdval == 1 && obstacles[y+1][x] != 1 {
					obstacles[y+1][x] = 1
					done = false
				}

				if ldval == 1 && luval == 1 && obstacles[y][x-1] != 1 {
					obstacles[y][x-1] = 1
					done = false
				}

				if rdval == 1 && ruval == 1 && obstacles[y][x+1] != 1 {
					obstacles[y][x+1] = 1
					done = false
				}
			}
		}
		if done {
			break
		}
	}

	return obstacles
}

func isInSafeZone(x, y, w, h int) bool {
	return x >= w/2-10 && x <= w/2+10 && y >= h/2-5 && y <= h/2+5
}
