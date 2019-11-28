package game

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/CatWantsMeow/gtetris/log"
)

const (
	FastGameMultiplier = 10

	FieldWidth    = 14
	FieldHeight   = 22
	PreviewWidth  = 4
	PreviewHeight = 2
	PreviewTop    = 0
	PreviewLeft   = 1

	StateInit = iota
	StateRunning
	StatePaused
	StateFinished
	StateClosed
	StateExiting
)

type Level struct {
	Name        string
	Delay       int
	LinePoints  int
	BlockPoints int
	TickPoints  int
	StartsAfter int
}

var (
	level1 = Level{
		Name:        "A",
		Delay:       250,
		LinePoints:  100,
		BlockPoints: 10,
		TickPoints:  0,
		StartsAfter: 0,
	}
	level2 = Level{
		Name:        "B",
		Delay:       200,
		LinePoints:  200,
		BlockPoints: 20,
		TickPoints:  1,
		StartsAfter: 180,
	}
	level3 = Level{
		Name:        "C",
		Delay:       150,
		LinePoints:  300,
		BlockPoints: 30,
		TickPoints:  2,
		StartsAfter: 360,
	}
	level4 = Level{
		Name:        "D",
		Delay:       100,
		LinePoints:  500,
		BlockPoints: 50,
		TickPoints:  3,
		StartsAfter: 720,
	}
	level5 = Level{
		Name:        "E",
		Delay:       50,
		LinePoints:  1000,
		BlockPoints: 100,
		TickPoints:  4,
		StartsAfter: 1500,
	}
	levels = []*Level{&level1, &level2, &level3, &level4, &level5}
)

type Stats struct {
	Level   string
	Score   int
	Lines   int
	Blocks  int
	Elapsed float64
}

type Game struct {
	screen  *Screen
	ctrl 	*Controller
	stats   *Stats
	field   *Field
	preview *Field

	state     int
	level     *Level
	curBlock  *Block
	nextBlock *Block

	fast  bool
	debug bool

	mu sync.Mutex
}

func (g *Game) generateBlock() {
	left := g.field.Width/2 - 1
	if g.nextBlock == nil {
		g.curBlock = NewRandomBlock(left, 0)
	} else {
		g.curBlock = g.nextBlock.Copy(left, 0)
	}
	g.nextBlock = NewRandomBlock(PreviewLeft, PreviewTop)
	log.Debug("Generated new block.")

	g.tryChangeLevel()
	g.preview.Clear(false)
	g.nextBlock.MustDraw(g.preview, false)
	g.stats.Score += g.level.BlockPoints
	g.stats.Blocks++

	g.field.Clear(false)
	if g.curBlock.Overlaps(g.field) {
		g.state = StateFinished
		log.Info("Changed state to finished.")
	}
}

func (g *Game) moveDown() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.curBlock != nil {
		ok := g.curBlock.TryMove(0, 1, g.field)
		if !ok {
			log.Debug("Failed to move down.")
			g.curBlock.MustDraw(g.field, true)
			g.removeLines()
			g.generateBlock()
		} else {
			log.Debug("Moved down.")
		}
		g.redraw()
	}
}

func (g *Game) removeLines() {
	removed := g.field.RemoveFilledLines()
	if removed > 0 {
		g.stats.Score += g.level.LinePoints * int(math.Pow(2, float64(removed)))
		g.stats.Lines += removed
	}
}

func (g *Game) tryChangeLevel() {
	for _, level := range levels {
		elapsed := int(g.stats.Elapsed)
		if g.fast {
			elapsed *= FastGameMultiplier
		}
		if level.StartsAfter <= elapsed {
			g.level = level
		}
	}
	log.Info("Changed level to %s.", g.level.Name)
	g.stats.Level = g.level.Name
}

func (g *Game) redraw() {
	g.field.Clear(false)
	g.curBlock.MustDraw(g.field, false)
	g.screen.Draw(g.state, g.stats)
}

func (g *Game) init() {
	g.ctrl.RegisterHandler(EventExit, func() {
		g.state = StateExiting
	})

	g.ctrl.RegisterHandler(EventResize, func() {
		g.screen.Resize()
		g.redraw()
	})

	g.ctrl.RegisterHandler(EventNewGame, func() {
		log.Info("Restarting game.")
		g.start()
	})

	g.ctrl.RegisterHandler(EventPauseResume, func() {
		switch g.state {
		case StatePaused:
			g.state = StateRunning
			log.Info("Changed state to running.")
		case StateRunning:
			g.state = StatePaused
			log.Info("Changed state to paused.")
		}
		g.redraw()
	})

	g.ctrl.RegisterHandler(EventUp, func() {
		if g.state == StateRunning {
			g.mu.Lock()
			defer g.mu.Unlock()

			if g.curBlock != nil {
				ok := g.curBlock.TryRotate(g.field)
				if ok {
					log.Debug("Rotated.")
				} else {
					log.Debug("Failed to rotate.")
				}
				g.redraw()
			}
		}
	})

	g.ctrl.RegisterHandler(EventLeft, func() {
		if g.state == StateRunning {
			g.mu.Lock()
			defer g.mu.Unlock()

			if g.curBlock != nil {
				ok := g.curBlock.TryMove(-1, 0, g.field)
				if ok {
					log.Debug("Moved left.")
				} else {
					log.Debug("Failed to move left.")
				}
				g.redraw()
			}
		}
	})

	g.ctrl.RegisterHandler(EventRight, func() {
		if g.state == StateRunning {
			g.mu.Lock()
			defer g.mu.Unlock()

			if g.curBlock != nil {
				ok := g.curBlock.TryMove(1, 0, g.field)
				if ok {
					log.Debug("Moved right.")
				} else {
					log.Debug("Failed to move right.")
				}
				g.redraw()
			}
		}
	})

	g.ctrl.RegisterHandler(EventDown, func() {
		if g.state == StateRunning {
			g.moveDown()
		}
	})
}

func (g *Game) start() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.stats.Elapsed = 0
	g.stats.Score = 0
	g.stats.Blocks = 0
	g.stats.Lines = 0

	g.level = levels[0]
	g.state = StateRunning
	g.field.Clear(true)
	g.generateBlock()
}

func (g *Game) Run() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
    defer termbox.Close()

	rand.Seed(time.Now().UnixNano())
	g.init()
	g.start()

	go g.ctrl.Run()
	for {
		switch g.state {
		case StateRunning:
			g.moveDown()
			g.stats.Elapsed += float64(g.level.Delay) / 1000
			g.stats.Score += g.level.TickPoints
		case StateExiting:
			return
		}
		time.Sleep(time.Millisecond * time.Duration(g.level.Delay))
	}
}

func Run(debug bool, fast bool) {
	field := NewField(FieldHeight, FieldWidth)
	preview := NewField(PreviewHeight, PreviewWidth)
	g := Game{
		stats:   &Stats{},
		field:   field,
		preview: preview,
		screen:  NewScreen(field, preview, debug),
		ctrl:    NewController(),
		state:   StateInit,
		fast:    fast,
		debug:   debug,
	}
	g.Run()
}
