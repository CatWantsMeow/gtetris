package game

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
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
	Log.Debug("Generated new block.")

	g.tryChangeLevel()
	g.preview.Clear(false)
	g.nextBlock.MustDraw(g.preview, false)
	g.stats.Score += g.level.BlockPoints
	g.stats.Blocks++

	g.field.Clear(false)
	err := g.curBlock.Check(g.field)
	if err == OverlappingError {
		g.finish()
	}
}

func (g *Game) removeLines() {
	removed := g.field.RemoveLines()
	if removed > 0 {
		g.stats.Score += g.level.LinePoints * int(math.Pow(2, float64(removed)))
		g.stats.Lines += removed
	}
}

func (g *Game) rotate() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.curBlock != nil {
		ok := g.curBlock.TryRotate(g.field)
		if ok {
			Log.Debug("Rotated.")
		} else {
			Log.Debug("Failed to rotate.")
		}
		g.redraw()
	}
}

func (g *Game) moveLeft() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.curBlock != nil {
		ok := g.curBlock.TryMove(-1, 0, g.field)
		if ok {
			Log.Debug("Moved left.")
		} else {
			Log.Debug("Failed to move left.")
		}
		g.redraw()
	}
}

func (g *Game) moveRight() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.curBlock != nil {
		ok := g.curBlock.TryMove(1, 0, g.field)
		if ok {
			Log.Debug("Moved right.")
		} else {
			Log.Debug("Failed to move right.")
		}
		g.redraw()
	}
}

func (g *Game) moveDown() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.curBlock != nil {
		ok := g.curBlock.TryMove(0, 1, g.field)
		if !ok {
			Log.Debug("Failed to move down.")
			g.curBlock.MustDraw(g.field, true)
			g.removeLines()
			g.generateBlock()
		} else {
			Log.Debug("Moved down.")
		}
		g.redraw()
	}
}

func (g *Game) redraw() {
	g.field.Clear(false)
	g.curBlock.MustDraw(g.field, false)
	g.screen.Draw(g.state, g.stats)
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

func (g *Game) finish() {
	g.state = StateFinished
	Log.Info("Changed state to finished.")
}

func (g *Game) pauseOrResume() {
	switch g.state {
	case StatePaused:
		g.state = StateRunning
		Log.Info("Changed state to running.")
	case StateRunning:
		g.state = StatePaused
		Log.Info("Changed state to paused.")
	}
	g.redraw()
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
	Log.Info("Changed level to %s.", g.level.Name)
	g.stats.Level = g.level.Name
}

func (g *Game) keyboardLoop(done chan bool) {
	defer g.Close()
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventResize {
			g.screen.resize()
			g.redraw()
		}

		if ev.Type == termbox.EventKey {
			switch ev.Ch {
			case 'p':
				g.pauseOrResume()
				continue
			case 'n':
				Log.Info("Restarting game.")
				g.start()
				continue
			}

			switch ev.Key {
			case termbox.KeyCtrlC, termbox.KeyEsc, termbox.KeyCtrlD:
				Log.Info("Received exit key.")
				done <- true
			}

			if g.state == StateRunning {
				switch ev.Key {
				case termbox.KeyArrowUp:
					g.rotate()
				case termbox.KeyArrowLeft:
					g.moveLeft()
				case termbox.KeyArrowRight:
					g.moveRight()
				case termbox.KeyArrowDown:
					g.moveDown()
				}
			}
		}
	}
}

func (g *Game) mainLoop() {
	defer g.Close()
	for {
		if g.state == StateRunning {
			g.moveDown()
			g.stats.Elapsed += float64(g.level.Delay) / 1000
			g.stats.Score += g.level.TickPoints
		}
		time.Sleep(time.Millisecond * time.Duration(g.level.Delay))
	}
}

func (g *Game) Run() {
	rand.Seed(time.Now().UnixNano())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	done := make(chan bool)
	go g.keyboardLoop(done)
	go g.mainLoop()

	g.start()
	select {
	case <-done:
		Log.Info("Exiting...")
		g.Close()
		return
	}
}

func (g *Game) Close() {
	if g.state != StateClosed {
		g.state = StateClosed
		termbox.Close()
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
		state:   StateInit,
		fast:    fast,
		debug:   debug,
	}
	g.Run()
}
