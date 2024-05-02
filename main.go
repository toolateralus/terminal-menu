package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

type State struct {
	items    []string
	selected int
	running  bool
}

type InputHandler struct {
	state     *State
	inputDone chan bool
}

type Renderer struct {
	state *State
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (renderer *Renderer) Draw() {
	for i, item := range renderer.state.items {
		if i == renderer.state.selected {
			tbprint(0, i, termbox.ColorBlack, termbox.ColorWhite, "> "+item)
		} else {
			tbprint(0, i, termbox.ColorWhite, termbox.ColorDefault, fmt.Sprintf("%d %s", i+1, item))
		}
	}
	termbox.Flush()
}

func (renderer *Renderer) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (handler *InputHandler) OnEnter() {}

func (handler *InputHandler) OnEscape() {
	handler.state.running = false
}

func (handler *InputHandler) OnArrowUp() {
	handler.state.selected--
	if handler.state.selected < 0 {
		handler.state.selected = 0
	}
}

func (handler *InputHandler) OnArrowDown() {
	handler.state.selected++
	length := len(handler.state.items) - 1
	if handler.state.selected > length {
		handler.state.selected = length
	}
}

func (handler *InputHandler) PollInput() {
	defer func() {
		handler.inputDone <- true
	}()

	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyEsc:
		case termbox.KeyCtrlC:
			handler.OnEscape()
		case termbox.KeyEnter:
			handler.OnEnter()
		case termbox.KeyArrowUp:
			handler.OnArrowUp()
		case termbox.KeyArrowDown:
			handler.OnArrowDown()
		}

	default:
		break
	}
}

func Run(state *State, renderer *Renderer, handler *InputHandler) {
	for state.running {
		select {
		case <-handler.inputDone:
			go handler.PollInput()
		default:
		}
		renderer.Clear()
		renderer.Draw()
		time.Sleep(time.Millisecond * 8)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	state := &State{
		items:    []string{"One", "Two", "Three"},
		selected: 0,
		running:  true,
	}

	inputDone := make(chan bool)

	handler := &InputHandler{
		state:     state,
		inputDone: inputDone,
	}

	renderer := &Renderer{
		state: state,
	}

	go handler.PollInput()

	Run(state, renderer, handler)
}