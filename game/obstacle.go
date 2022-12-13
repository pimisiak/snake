package game

import (
	"math"
	"math/rand"
	"time"
)

func generateObstacles(width, height int) [][]int {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	obstacles := make([][]int, height)
	for i := range obstacles {
		obstacles[i] = make([]int, width)
	}

	// base count on screen size
	count := int(math.Floor(0.0085 * float64(width) * float64(height)))

	// generate the specified number of obstacles
	for i := 0; i < count; i++ {

		// choose a random starting coordinate
		current := coordinate{random.Intn(width), random.Intn(height)}
		obstacles[current.y][current.x] = 1

		// generate the obstacle cluster
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
		smoothed := false

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if obstacles[y][x] == 0 {
					continue
				}

				current := coordinate{x, y}
				leftUp := current.move(left).move(up)
				rightUp := current.move(right).move(up)
				leftDown := current.move(left).move(down)
				rightDown := current.move(right).move(down)

				lu := 0
				if leftUp.x >= 0 && leftUp.y >= 0 {
					lu = obstacles[leftUp.y][leftUp.x]
				}

				ru := 0
				if rightUp.x < width && rightUp.y >= 0 {
					ru = obstacles[rightUp.y][rightUp.x]
				}

				ld := 0
				if leftDown.x >= 0 && leftDown.y < height {
					ld = obstacles[leftDown.y][leftDown.x]
				}

				rd := 0
				if rightDown.x < width && rightDown.y < height {
					rd = obstacles[rightDown.y][rightDown.x]
				}

				if lu == 1 && ru == 1 && obstacles[y-1][x] != 1 {
					obstacles[y-1][x] = 1
					smoothed = true
				}

				if ld == 1 && rd == 1 && obstacles[y+1][x] != 1 {
					obstacles[y+1][x] = 1
					smoothed = true
				}

				if ld == 1 && lu == 1 && obstacles[y][x-1] != 1 {
					obstacles[y][x-1] = 1
					smoothed = true
				}

				if rd == 1 && ru == 1 && obstacles[y][x+1] != 1 {
					obstacles[y][x+1] = 1
					smoothed = true
				}
			}
		}

		if !smoothed {
			break
		}
	}

	return obstacles
}
