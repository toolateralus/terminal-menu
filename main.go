package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

type State struct {
	items    []string
	selected int
	running  bool
	isDir    bool
	path     string
	previous *State
}

func NewState(previous *State, path string) State {
	stats, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	if stats.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		items := []string{}
		for _, file := range files {
			items = append(items, file.Name())
		}
		
		return State{
			items:    items,
			selected: 0,
			running:  true,
			path:     path,
			previous: previous,
			isDir: 		true,
		}
	} else {
		contents, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		
		lines := strings.Split(string(contents), "\n")
		
		return State{
			items:    lines,
			selected: 0,
			running:  true,
			path:     path,
			previous: previous,
			isDir: 		false,
		}
			
	}
	
	
}

type InputHandler struct {
	state     *State
	inputDone chan bool
}

type Renderer struct {
	state *State
}

func termbox_print(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (renderer *Renderer) Draw() {
	for i, item := range renderer.state.items {
		if i == renderer.state.selected {
			termbox_print(0, i, termbox.ColorBlack, termbox.ColorWhite, "> "+item)
		} else {
			termbox_print(0, i, termbox.ColorWhite, termbox.ColorDefault, fmt.Sprintf("%d %s", i+1, item))
		}
	}
	termbox.Flush()
}

func (renderer *Renderer) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (handler *InputHandler) OnEnter() {
	if !handler.state.isDir {
		// cant go deeper into a file. only escpae.
		return
	}
	
	selected := handler.state.items[handler.state.selected]
	path := fmt.Sprintf("%s/%s", handler.state.path, selected)
	
	
	
	oldState := *handler.state
	*handler.state = NewState(&oldState, path)
}

func (handler *InputHandler) OnEscape() {
	if handler.state.previous != nil {
		*handler.state = *handler.state.previous
	}
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
			handler.OnEscape()
		case termbox.KeyCtrlC:
			handler.state.running = false
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
		time.Sleep(time.Millisecond * 16)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	state := NewState(nil, "/")

	inputDone := make(chan bool)

	handler := &InputHandler{
		state:     &state,
		inputDone: inputDone,
	}

	renderer := &Renderer{
		state: &state,
	}

	go handler.PollInput()

	Run(&state, renderer, handler)
}
