package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	showHelp bool
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
			// try hide the hidden files.
			if !strings.HasPrefix(file.Name(), ".") {
				fileInfo, err := os.Stat(filepath.Join(path, file.Name()))
				if err != nil {
					log.Fatal(err)
				}
				// Check if the file has read permission
				if fileInfo.Mode().Perm()&0400 != 0 {
					items = append(items, file.Name())
				}
			}
		}
		
		return State{
			items:    items,
			selected: 0,
			running:  true,
			path:     path,
			previous: previous,
			isDir:    true,
			showHelp: true,
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
			isDir:    false,
			showHelp: false,
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
	_, h := termbox.Size()
	page := renderer.state.selected / h
	start := page * h
	end := min(len(renderer.state.items), start+h)
	
	for i := start; i < end; i++ {
		item := renderer.state.items[i]
		if i == renderer.state.selected {
			termbox_print(0, i-start, termbox.ColorBlack, termbox.ColorWhite, "> "+item)
		} else {
			termbox_print(0, i-start, termbox.ColorWhite, termbox.ColorDefault, fmt.Sprintf("%d %s", i+1, item))
		}
	}
	
	extra := 1
	
	if renderer.state.showHelp {
		termbox_print(0, end + extra, termbox.ColorWhite, termbox.ColorBlack, "[esc] to go up, [enter] to go in, [up/down] to navigate\n")
		extra++
		termbox_print(0, end + extra, termbox.ColorWhite, termbox.ColorBlack, "[ctrl + S/W] to jump to bottom/top, [PageDown/PageUp] to jump pages\n")
		extra++
		} else {
			termbox_print(0, end + extra, termbox.ColorWhite, termbox.ColorBlack, "[Home] to toggle help")
			extra++
		}
		termbox_print(0, end + extra, termbox.ColorLightYellow, termbox.ColorBlack, fmt.Sprintf("%s\n", renderer.state.path))
	
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

func (handler *InputHandler) OnHome() {
	handler.state.showHelp = !handler.state.showHelp
}

func (handler *InputHandler) OnPageUp() {
	_, h := termbox.Size()
	page := handler.state.selected / h
	start := page * h
	end := max(0, start-h)
	handler.state.selected = end	
}

func (handler *InputHandler) OnPageDown() {
	_, h := termbox.Size()
	page := handler.state.selected / h
	start := page * h
	end := min(len(handler.state.items) - 1, start+h)
	handler.state.selected = end
}

func (handler *InputHandler) OnCtrlS() {
	handler.state.selected = len(handler.state.items) -1
}

func (handler *InputHandler) OnCtrlW() {
	handler.state.selected = 0
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
		case termbox.KeyPgdn: 
			handler.OnPageDown()
		case termbox.KeyPgup:
			handler.OnPageUp()
		case termbox.KeyCtrlS:
			handler.OnCtrlS()
		case termbox.KeyCtrlW:
			handler.OnCtrlW()
		case termbox.KeyHome:
			handler.OnHome()
		}
		
	default:
		break
	}
}

func Run(state *State, renderer *Renderer, handler *InputHandler) {
	go handler.PollInput()
	renderer.Clear()
	renderer.Draw()
	for state.running {
		select {
		case <-handler.inputDone:
			go handler.PollInput()
			renderer.Clear()
			renderer.Draw()
		default:
		}

		time.Sleep(time.Millisecond * 16)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	homeDir := os.Getenv("HOME")
	state := NewState(nil, homeDir)

	handler := &InputHandler{
		state:     &state,
		inputDone: make(chan bool),
	}

	renderer := &Renderer{
		state: &state,
	}

	Run(&state, renderer, handler)
}
