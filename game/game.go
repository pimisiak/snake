package game

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type apple = coordinate

type game struct {
	random    *rand.Rand
	screen    tcell.Screen
	snake     *snake
	apple     apple
	obstacles [][]int
	speed     time.Duration
	pause     bool
	over      bool
	score     int
}

func Run() {
	game, err := newGame()
	if err != nil {
		log.Fatalln(err)
		return
	}
	updates := game.registerKeys()
	game.loop(updates)
}

func newGame() (*game, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("could not create new screen %x", err)
	}

	err = screen.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize screen: %x", err)
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	w, h := screen.Size()

	game := &game{
		random:    random,
		screen:    screen,
		snake:     newSnake(coordinate{w / 2, h / 2}, left),
		obstacles: generateObstacles(w, h),
		speed:     time.Duration(90) * time.Millisecond,
		pause:     false,
		over:      false,
		score:     0,
	}
	game.spawnApple()

	return game, nil
}

func (g *game) loop(updates <-chan string) {
	updateAndRender := func() {
		if !g.over {
			g.screen.Clear()
			g.update()
			g.render()
			g.screen.Show()
		}
	}
	for {
		select {
		case update := <-updates:
			switch update {
			case "UP":
				g.snake.turn(up)
			case "DOWN":
				g.snake.turn(down)
			case "RIGHT":
				g.snake.turn(right)
			case "LEFT":
				g.snake.turn(left)
			case "PAUSE":
				g.pause = !g.pause
			case "RESTART":
				g.restartGame()
			case "SYNC":
				g.screen.Sync()
			case "QUIT":
				g.screen.Fini()
				return
			}
			updateAndRender()
		case <-time.After(g.speed):
			updateAndRender()
		}
	}
}

func (g *game) restartGame() {
	w, h := g.screen.Size()
	g.snake = newSnake(coordinate{w / 2, h / 2}, left)
	g.obstacles = generateObstacles(w, h)
	g.score = 0
	g.pause = false
	g.over = false
	g.spawnApple()
}

func (g *game) update() {
	if g.pause || g.over {
		return
	}

	// snakey bit me, and that really hurtz snakey
	for i := 1; i < len(g.snake.body); i++ {
		if g.snake.head().equals(g.snake.body[i]) {
			g.over = true
		}
	}

	// check if snake hit obstacle
	if g.obstacles[g.snake.head().y][g.snake.head().x] == 1 {
		if len(g.snake.body) == 1 {
			g.over = true
		} else {
			g.obstacles[g.snake.head().y][g.snake.head().x] = 0
			g.snake.puke()
			g.score++
		}
	}

	// check if snake ate apple
	if g.snake.head().equals(g.apple) {
		g.snake.eat(g.apple)
		g.spawnApple()
	}

	// move snake
	g.moveSnake(g.screen.Size())
}

func (g *game) spawnApple() {
	w, h := g.screen.Size()
	x, y := g.random.Intn(w), g.random.Intn(h)
	for g.obstacles[y][x] == 1 {
		x, y = g.random.Intn(w), g.random.Intn(h)
	}
	g.apple = apple{x, y}
}

func (g *game) moveSnake(width, height int) {
	g.snake.move()

	// check if head out of bounds
	if g.snake.body[0].x >= width {
		g.snake.body[0].x = 0
	}
	if g.snake.body[0].x < 0 {
		g.snake.body[0].x = width - 1
	}
	if g.snake.body[0].y >= height {
		g.snake.body[0].y = 0
	}
	if g.snake.body[0].y < 0 {
		g.snake.body[0].y = height - 1
	}
}

func (g *game) render() {
	if g.over {
		g.drawGameOver()
		return
	}
	g.drawSnake()
	g.drawApple()
	g.drawObstacles()
	if g.pause {
		g.drawPause()
	}
}

func (g *game) drawSnake() {
	for _, co := range g.snake.body {
		g.screen.SetContent(co.x, co.y, '#', nil, tcell.StyleDefault.Foreground(tcell.ColorGreenYellow))
	}
}

func (g *game) drawApple() {
	g.screen.SetContent(g.apple.x, g.apple.y, '*', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
}

func (g *game) drawObstacles() {
	for y := range g.obstacles {
		for x := range g.obstacles[y] {
			if g.obstacles[y][x] == 1 {
				g.screen.SetContent(x, y, '@', nil, tcell.StyleDefault.Foreground(tcell.ColorSlateGray))
			}
		}
	}
}

func (g *game) drawGameOver() {
	g.screen.Beep()
	g.screen.Fill('.', tcell.StyleDefault.Foreground(tcell.ColorSlateGray))
	w, h := g.screen.Size()
	g.drawStr(w/2-6, h/2-5, tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorDarkRed), "Game Over!")
	g.drawStr(w/2-5, h/2-3, tcell.StyleDefault.Foreground(tcell.ColorGreenYellow), fmt.Sprintf("Score: %d", g.score))
	g.drawStr(w/2-9, h/2-1, tcell.StyleDefault.Foreground(tcell.ColorGreen), "Press ESC to exit")
	g.drawStr(w/2-11, h/2, tcell.StyleDefault.Foreground(tcell.ColorGreen), "Press Ctr-r to restart")
}

func (g *game) drawPause() {
	g.drawStr(1, 1, tcell.StyleDefault.Foreground(tcell.ColorSlateGray), "Pause...")
}

func (g *game) drawStr(x, y int, style tcell.Style, str string) {
	for _, c := range str {
		g.screen.SetContent(x, y, c, nil, style)
		x++
	}
}

func (g *game) registerKeys() <-chan string {
	updates := make(chan string)
	go func() {
		for {
			switch event := g.screen.PollEvent().(type) {
			case *tcell.EventKey:
				switch event.Key() {

				// game control
				case tcell.KeyEscape, tcell.KeyCtrlC:
					updates <- "QUIT"
					os.Exit(0)
				case tcell.KeyCtrlP:
					updates <- "PAUSE"
				case tcell.KeyCtrlR:
					updates <- "RESTART"

				// movement
				case tcell.KeyDown:
					updates <- "DOWN"
				case tcell.KeyUp:
					updates <- "UP"
				case tcell.KeyLeft:
					updates <- "LEFT"
				case tcell.KeyRight:
					updates <- "RIGHT"

				// vim fun
				case tcell.KeyRune:
					switch event.Rune() {
					case 'j':
						updates <- "DOWN"
					case 'k':
						updates <- "UP"
					case 'h':
						updates <- "LEFT"
					case 'l':
						updates <- "RIGHT"
					case 'r':
						updates <- "RESTART"
					case 'p':
						updates <- "PAUSE"
					case 'q':
						updates <- "QUIT"
					}
				}
			default:
				updates <- "SYNC"
			}
		}
	}()
	return updates
}
