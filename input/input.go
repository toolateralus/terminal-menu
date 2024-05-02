package input

import (
	"fmt"

	"github.com/nsf/termbox-go"
	"github.com/toolateralus/menu/state"
)


type InputHandler struct {
	State     *state.State
	InputDone chan bool
}

func (handler *InputHandler) OnCtrlW() {
	handler.State.Selected = 0
}

func (handler *InputHandler) PollInput() {
	defer func() {
		handler.InputDone <- true
	}()
	
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyEsc:
			handler.OnEscape()
		case termbox.KeyCtrlC:
			handler.State.Running = false
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

func (handler *InputHandler) OnEnter() {
	if !handler.State.IsDir {
		// cant go deeper into a file. only escpae.
		return
	}

	selected := handler.State.Items[handler.State.Selected]
	path := fmt.Sprintf("%s/%s", handler.State.Path, selected)

	oldState := *handler.State
	*handler.State = state.NewState(&oldState, path)
}

func (handler *InputHandler) OnEscape() {
	if handler.State.Previous != nil {
		*handler.State = *handler.State.Previous
	}
}

func (handler *InputHandler) OnArrowUp() {
	handler.State.Selected--
	if handler.State.Selected < 0 {
		handler.State.Selected = 0
	}
}

func (handler *InputHandler) OnArrowDown() {
	handler.State.Selected++
	length := len(handler.State.Items) - 1
	if handler.State.Selected > length {
		handler.State.Selected = length
	}
}

func (handler *InputHandler) OnHome() {
	handler.State.ShowHelp = !handler.State.ShowHelp
}

func (handler *InputHandler) OnPageUp() {
	_, h := termbox.Size()
	page := handler.State.Selected / h
	start := page * h
	end := max(0, start-h)
	handler.State.Selected = end	
}

func (handler *InputHandler) OnPageDown() {
	_, h := termbox.Size()
	page := handler.State.Selected / h
	start := page * h
	end := min(len(handler.State.Items) - 1, start+h)
	handler.State.Selected = end
}

func (handler *InputHandler) OnCtrlS() {
	handler.State.Selected = len(handler.State.Items) -1
}