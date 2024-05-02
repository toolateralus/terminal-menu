package renderer

import (
	"fmt"

	"github.com/nsf/termbox-go"
	"github.com/toolateralus/menu/state"
)


type Renderer struct {
	State *state.State
}



func termbox_print(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (renderer *Renderer) Draw() {
	_, h := termbox.Size()
	page := renderer.State.Selected / h
	start := page * h
	end := min(len(renderer.State.Items), start+h)
	
	for i := start; i < end; i++ {
		item := renderer.State.Items[i]
		if i == renderer.State.Selected {
			termbox_print(0, i-start, termbox.ColorBlack, termbox.ColorWhite, "> "+item)
		} else {
			termbox_print(0, i-start, termbox.ColorWhite, termbox.ColorDefault, fmt.Sprintf("%d %s", i+1, item))
		}
	}
	
	extra := 1
	
	if renderer.State.ShowHelp {
		termbox_print(0, end + extra, termbox.ColorWhite, termbox.ColorBlack, "[esc] to go up, [enter] to go in, [up/down] to navigate\n")
		extra++
		termbox_print(0, end + extra, termbox.ColorWhite, termbox.ColorBlack, "[ctrl + S/W] to jump to bottom/top, [PageDown/PageUp] to jump pages\n")
		extra++
		} else {
			termbox_print(0, end + extra, termbox.ColorWhite, termbox.ColorBlack, "[Home] to toggle help")
			extra++
		}
		termbox_print(0, end + extra, termbox.ColorLightYellow, termbox.ColorBlack, fmt.Sprintf("%s\n", renderer.State.Path))
	
	termbox.Flush()
}

func (renderer *Renderer) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}