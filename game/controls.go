package game

import (
    "github.com/nsf/termbox-go"
)

const (
    EventExit = iota
    EventDown
    EventUp
    EventLeft
    EventRight
    EventNewGame
    EventPauseResume
    EventResize
)

func NewController() *Controller {
    return &Controller{
        handlers: make(map[int]func()),
    }
}

type Controller struct {
    handlers map[int]func()
}

func (c *Controller) handle(event int) {
    handler, ok := c.handlers[event]
    if ok {
        handler()
    }
}

func (c *Controller) RegisterHandler(event int, handler func()) {
    c.handlers[event] = handler
}

func (c *Controller) Run() {
    for {
        e := termbox.PollEvent()
        if e.Type == termbox.EventResize {
            c.handle(EventResize)
        }

        if e.Type == termbox.EventKey {
            switch e.Ch {
            case 'p':
                c.handle(EventPauseResume)
            case 'n':
                c.handle(EventNewGame)
            }
        }

        switch e.Key {
        case termbox.KeyArrowUp:
            c.handle(EventUp)
        case termbox.KeyArrowLeft:
            c.handle(EventLeft)
        case termbox.KeyArrowRight:
            c.handle(EventRight)
        case termbox.KeyArrowDown:
            c.handle(EventDown)
        case termbox.KeyCtrlC, termbox.KeyEsc, termbox.KeyCtrlD:
            c.handle(EventExit)
        }
    }
}
