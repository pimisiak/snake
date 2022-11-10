package game

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Apple Coordinate

type Game struct {
	screen tcell.Screen
	snake  *Snake
	apple  Apple
	over   bool
	speed  time.Duration
}

func NewGame() (*Game, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("Couldn not create new screen %x", err)
	}

	err = screen.Init()
	if err != nil {
		return nil, fmt.Errorf("Could not initialize screen: %x", err)
	}

	x, y := screen.Size()
	game := &Game{
		screen: screen,
		snake:  NewSnake(NewCoordinate(x/2, y/2)),
		apple:  Apple{},
		over:   false,
	}
	game.spawnApple()

	return game, nil
}

func (g *Game) spawnApple() {
	source := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(source)
	w, h := g.screen.Size()
	g.apple = Apple{randomizer.Intn(w), randomizer.Intn(h)}
}

func (g *Game) drawSnake() {
	for _, coor := range *g.snake {
		g.screen.SetContent(coor.X, coor.Y, '#', nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
	}
}

func (g *Game) drawApple() {
	g.screen.SetContent(g.apple.X, g.apple.Y, '*', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
}

func Run() {
	game, err := NewGame()
	if err != nil {
		log.Fatalln(err)
		return
	}

	// game loop
	for !game.over {
		game.screen.Clear()

		switch event := game.screen.PollEvent().(type) {
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				game.screen.Fini()
				os.Exit(0)
			case tcell.KeyDown:
				(*game.snake)[0].Y += 1
			case tcell.KeyUp:
				(*game.snake)[0].Y -= 1
			case tcell.KeyLeft:
				(*game.snake)[0].X -= 1
			case tcell.KeyRight:
				(*game.snake)[0].X += 1
			}
		default:
			game.screen.Sync()
		}

		game.drawSnake()
		game.drawApple()

		game.screen.Show()
	}
}
