package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	term "github.com/nsf/termbox-go"
)

type cell struct {
	char rune
	fg   term.Attribute
	bg   term.Attribute
}

var (
	difficulty      = 3
	difficultyCount = 0
	currentScore    = 0
)

func exit() {
	term.Close()
	os.Exit(0)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	err := term.Init()
	if err != nil {
		panic(err)
	}

	_ = term.Clear(term.ColorWhite, term.ColorBlack)

	width, higth := term.Size()
	playerCol := width / 2

	screenBuffer := make([][]cell, width)
	for col := 0; col < width; col++ {
		screenBuffer[col] = make([]cell, higth)
		for row := 0; row < higth; row++ {
			screenBuffer[col][row] = cell{
				fg:   term.ColorWhite,
				bg:   term.ColorBlack,
				char: ' ',
			}
		}
	}

	eventQueue := make(chan term.Event)
	go func() {
		for {
			eventQueue <- term.PollEvent()
		}
	}()
	drawTick := time.NewTicker(150 * time.Millisecond)

	for {
		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case term.EventKey:
				switch ev.Key {
				case term.KeyEsc:
					exit()
				case term.KeyArrowUp:
					fmt.Println("Arrow Up pressed")
				case term.KeyArrowDown:
					fmt.Println("Arrow Down pressed")
				case term.KeyArrowLeft:
					playerCol--
				case term.KeyArrowRight:
					playerCol++
				default:
					fmt.Println("ASCII : ", ev.Ch)
				}
			case term.EventError:
				panic(ev.Err)
			}
		case <-drawTick.C:
			copyUp(screenBuffer)
			for i := 0; i < difficulty; i++ {
				addObstacle(screenBuffer)
			}

			currentScore++
			difficultyCount++
			if difficultyCount == 300 {
				difficulty++
				difficultyCount = 0
			}

			// collision detection
			if screenBuffer[playerCol][higth/2].char == '*' {
				fmt.Println("colision!")
				exit()
			}

			screenBuffer[playerCol][higth/2].char = 'V'
			printScreen(
				screenBuffer, 0, 0,
				term.ColorWhite,
				term.ColorBlack,
				fmt.Sprintf("Current score %v", currentScore))
			draw(screenBuffer)
		}
	}
}

func draw(screen [][]cell) {
	width, higth := term.Size()
	for col := 0; col < width; col++ {
		for row := 0; row < higth; row++ {
			term.SetCell(
				col,
				row,
				screen[col][row].char,
				screen[col][row].fg,
				screen[col][row].bg)
		}
	}
	_ = term.Flush()
}

func copyUp(screen [][]cell) {
	width, higth := term.Size()
	for col := 0; col < width; col++ {
		for row := 0; row < higth; row++ {
			if screen[col][row].char == 'V' {
				if row-1 >= 0 {
					screen[col][row-1] = screen[col][row]
					screen[col][row-1].fg = screen[col][row].fg + 1
					if screen[col][row-1].fg >= 8 {
						screen[col][row-1].fg = 2
					}
				}
				screen[col][row].char = ' '
			}
			if screen[col][row].char == '*' {
				if row-1 >= 0 {
					screen[col][row-1] = screen[col][row]
				}
				screen[col][row].char = ' '
			}
		}
	}
}

func addObstacle(screen [][]cell) {
	w, h := term.Size()
	w = rand.Intn(w - 1)
	screen[w][h-1].char = '*'
}

func printScreen(
	screen [][]cell,
	col, row int,
	fg, bg term.Attribute,
	msg string) {
	for c := 0; c < len(msg); c++ {
		screen[col+c][row].char = rune(msg[c])
		screen[col+c][row].fg = fg
		screen[col+c][row].bg = bg
	}
}
