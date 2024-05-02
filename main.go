package main

import (
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/toolateralus/menu/input"
	"github.com/toolateralus/menu/renderer"
	"github.com/toolateralus/menu/state"
)


func Run(state *state.State, renderer *renderer.Renderer, handler *input.InputHandler) {
	go handler.PollInput()
	renderer.Clear()
	renderer.Draw()
	for state.Running {
		select {
		case <-handler.InputDone:
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
	state := state.NewState(nil, homeDir)
	
	handler := &input.InputHandler{
		State:     &state,
		InputDone: make(chan bool),
	}
	
	renderer := &renderer.Renderer{
		State: &state,
	}
	
	// blocking call to run the app's lifetime.
	Run(&state, renderer, handler)
}
