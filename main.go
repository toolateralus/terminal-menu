package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

type Context struct {
	items     []string
	selected  int
	running   bool
	inputDone chan bool
}


type Renderer struct {}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (renderer *Renderer) Draw(c *Context) {
	for i, item := range c.items {
		if i == c.selected {
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
func (ctx *Context) OnEnter() {}

func (ctx *Context) OnEscape() {
	ctx.running = false
}
func (ctx *Context) OnArrowUp() {
	ctx.selected--
	if ctx.selected < 0 {
		ctx.selected = 0
	}
}
func (ctx *Context) OnArrowDown() {
	ctx.selected++
	length := len(ctx.items) - 1
	if ctx.selected > length {
		ctx.selected = length
	}
}

func (ctx *Context) PollInput() {
	defer func() {
		ctx.inputDone <- true
	}()

	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyEsc:
		case termbox.KeyCtrlC:
			ctx.OnEscape()
		case termbox.KeyEnter:
			ctx.OnEnter()
		case termbox.KeyArrowUp:
			ctx.OnArrowUp()
		case termbox.KeyArrowDown:
			ctx.OnArrowDown()
		}

	default:
		break
	}
}

func (ctx Context) Run(r *Renderer) {
	for ctx.running {
		select {
		case <-ctx.inputDone:
			go ctx.PollInput()
		default:
		}
		r.Clear()
		r.Draw(&ctx)
		time.Sleep(time.Millisecond * 8)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	ctx := Context{
		items:     []string{"One", "Two", "Three"},
		selected:  0,
		running:   true,
		inputDone: make(chan bool),
	}
	go ctx.PollInput()
	renderer := Renderer{}
	ctx.Run(&renderer)
}


