package game

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type apple = coordinate

type game struct {
	mutex     sync.Mutex
	random    *rand.Rand
	screen    tcell.Screen
	snake     *snake
	apple     apple
	obstacles []obstacle
	speed     time.Duration
	pause     bool
	over      bool
}

func Run() {
	game, err := newGame()
	if err != nil {
		log.Fatalln(err)
		return
	}
	quit := make(chan struct{})
	go game.registerKeys(quit)
	game.loop(quit)
}

func newGame() (*game, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("Couldn not create new screen %x", err)
	}

	err = screen.Init()
	if err != nil {
		return nil, fmt.Errorf("Could not initialize screen: %x", err)
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	w, h := screen.Size()

	game := &game{
		random:    random,
		screen:    screen,
		snake:     newSnake(coordinate{w / 2, h / 2}, left),
		apple:     apple{random.Intn(w), random.Intn(h)},
		obstacles: generateObstacles(),
		speed:     time.Duration(50) * time.Millisecond,
		pause:     false,
		over:      false,
	}

	return game, nil
}

func (g *game) loop(quit chan struct{}) {
	for {
		select {
		case <-quit:
			return
		case <-time.After(g.speed):
			if !g.over {
				g.screen.Clear()
				if !g.pause {
					g.update()
				}
				g.render()
				g.screen.Show()
			}
		}
	}
}

func (g *game) registerKeys(quit chan struct{}) {
	for {
		switch event := g.screen.PollEvent().(type) {
		case *tcell.EventKey:
			g.mutex.Lock()
			switch event.Key() {

			// game control
			case tcell.KeyEscape, tcell.KeyCtrlC:
				close(quit)
				g.screen.Fini()
				os.Exit(0)
			case tcell.KeyCtrlP:
				if !g.over {
					g.pause = !g.pause
				}
			case tcell.KeyCtrlR:
				g.restartGame()

			// movement
			case tcell.KeyDown:
				g.snake.redirect(down)
			case tcell.KeyUp:
				g.snake.redirect(up)
			case tcell.KeyLeft:
				g.snake.redirect(left)
			case tcell.KeyRight:
				g.snake.redirect(right)

			// vim fun
			case tcell.KeyRune:
				switch event.Rune() {
				case 'j':
					g.snake.redirect(down)
				case 'k':
					g.snake.redirect(up)
				case 'h':
					g.snake.redirect(left)
				case 'l':
					g.snake.redirect(right)
				}

			}
			g.mutex.Unlock()
		default:
			g.screen.Sync()
		}
	}
}

func (g *game) restartGame() {
	w, h := g.screen.Size()
	g.snake = newSnake(coordinate{w / 2, h / 2}, left)
	g.apple = apple{g.random.Intn(w), g.random.Intn(h)}
	g.obstacles = generateObstacles()
	g.pause = false
	g.over = false
}

func (g *game) update() {
	g.mutex.Lock()

	// snakey bit me, and that really hurtz snakey
	for i := 1; i < len(g.snake.body); i++ {
		if g.snake.head().equals(g.snake.body[i]) {
			g.over = true
		}
	}

	// check if snake ate apple
	if g.snake.head().equals(g.apple) {
		g.snake.eat(g.apple)
		g.spawnApple()
	}

	// check if snake hit obstacle
	//

	// move snake
	g.snake.move()

	g.mutex.Unlock()
}

func (g *game) spawnApple() {
	w, h := g.screen.Size()
	g.apple = apple{g.random.Intn(w), g.random.Intn(h)}
}

func (g *game) render() {
	g.drawSnake()
	g.drawApple()

	// draw obstacles

	if g.pause {
		g.drawPause()
	}

	if g.over {
		g.drawGameOver()
	}
}

func (g *game) drawSnake() {
	for _, co := range g.snake.body {
		g.screen.SetContent(co.x, co.y, '#', nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
	}
}

func (g *game) drawApple() {
	g.screen.SetContent(g.apple.x, g.apple.y, '*', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
}

func (g *game) drawGameOver() {
	g.screen.Beep()
	g.screen.Fill('.', tcell.StyleDefault.Foreground(tcell.ColorSlateGray))
	w, h := g.screen.Size()
	g.drawStr(w/2-7, h/2, tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorDarkRed), "Game Over!")
	g.drawStr(w/2-11, h/2+1, tcell.StyleDefault.Foreground(tcell.ColorGreen), "Press ESC to exit.")
	g.drawStr(w/2-14, h/2+2, tcell.StyleDefault.Foreground(tcell.ColorGreen), "Press Ctr-r to restart.")
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
