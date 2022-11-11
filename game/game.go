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
	start     bool
	over      bool
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
		apple:     apple{w/2 - 4, h/2 + 4},
		obstacles: generateObstacles(),
		speed:     time.Duration(50) * time.Millisecond,
		pause:     false,
		start:     false,
		over:      false,
	}

	return game, nil
}

func (g *game) spawnApple() {
	w, h := g.screen.Size()
	g.apple = apple{g.random.Intn(w), g.random.Intn(h)}
}

func (g *game) update() {
	g.mutex.Lock()

	// snakey bit me, and that really hurt snakey
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

func (g *game) render() {
	// draw snake
	for _, co := range g.snake.body {
		g.screen.SetContent(co.x, co.y, '#', nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
	}

	// draw apple
	g.screen.SetContent(g.apple.x, g.apple.y, '*', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))

	// draw obstacles

	// draw game over
	if g.over {
		g.screen.Beep()
		g.screen.Fill('.', tcell.StyleDefault)
	}
}

func (g *game) registerKeys(quit chan struct{}) {
	for {
		switch event := g.screen.PollEvent().(type) {
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				close(quit)
				g.screen.Fini()
				os.Exit(0)
			case tcell.KeyDown:
				g.mutex.Lock()
				g.snake.redirect(down)
				g.mutex.Unlock()
			case tcell.KeyUp:
				g.mutex.Lock()
				g.snake.redirect(up)
				g.mutex.Unlock()
			case tcell.KeyLeft:
				g.mutex.Lock()
				g.snake.redirect(left)
				g.mutex.Unlock()
			case tcell.KeyRight:
				g.mutex.Lock()
				g.snake.redirect(right)
				g.mutex.Unlock()
			case tcell.KeyCtrlP:
				g.mutex.Lock()
				g.pause = !g.pause
				g.mutex.Unlock()
			}
		default:
			g.screen.Sync()
		}
	}
}

func (g *game) loop(quit chan struct{}) {
	for {
		select {
		case <-quit:
			return
		case <-time.After(g.speed):
			if !g.pause && !g.over {
				g.screen.Clear()
				g.update()
				g.render()
				g.screen.Show()
			}
		}
	}
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
