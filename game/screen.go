package game

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

const (
	BackgroundChar  = ' '
	BlockChar       = ' '
	BackgroundColor = termbox.ColorDefault

	FieldXScale   = 2
	FieldYScale   = 1
	ScreenMinLeft = 0
	ScreenTop     = 2

	FieldBoxLeftWidth  = 2
	FieldBoxLeftChars  = "<|"
	FieldBoxRightWidth = 2
	FieldBoxRightChars = "|>"
	FieldBoxBottomChar = '='
	FieldBoxColor      = termbox.ColorDefault

	LeftPromptWidth  = 21
	LeftPromptLeft   = 6
	RightPromptWidth = 21
	RightPromptLeft  = 3

	LogWidth  = 50
	LogHeight = FieldHeight + 1

	StatePausedPrompt   = "Paused"
	StateRunningPrompt  = "Running"
	StateFinishedPrompt = "Game Over"
	StatePausedColor    = termbox.ColorYellow
	StateRunningColor   = termbox.ColorGreen
	StateFinishedColor  = termbox.ColorRed

	StatsPromptHeight = 4
	StatsPrompt       = "" +
		"Level:  %4s\n" +
		"Time:   %4d\n" +
		"Blocks: %4d\n" +
		"Lines:  %4d\n" +
		"Score:  %4d"

	NextBlockPrompt = "Next block:"
	NextBlockLeft   = 2
	NextBlockTop    = 1

	CopyrightPromptHeight = 2
	CopyrightPromptColor  = termbox.ColorYellow
	CopyrightPrompt       = "" +
		"Gogi's Tetris\n" +
		"   is 300$   "

	HelpPrompt = "" +
		"Move left:     ←\n" +
		"Move right:    →\n" +
		"Speed up:      ↓\n" +
		"Rotate:        ↑\n" +
		"Close game:    esc\n" +
		"Pause/resume:  p\n" +
		"Restart:       n"

	Header = "" +
		"\n" +
		"  _                      _    \n" +
		">(.)__  Gogi's Tetris  >(.)__ \n" +
		" (___/     is 300$      (___/ \n"
)

type Screen struct {
	debug   bool
	field   *Field
	preview *Field

	Top  int
	Left int
}

func (s *Screen) width() int {
	w := LeftPromptWidth +
		FieldBoxLeftWidth +
		s.field.Width*FieldXScale +
		FieldBoxRightWidth +
		RightPromptWidth
	if s.debug {
		w += LogWidth
	}
	return w
}

func (s *Screen) resize() {
	w, _ := termbox.Size()
	left := (w - s.width()) / 2
	if ScreenMinLeft < left {
		s.Left = left
	} else {
		s.Left = ScreenMinLeft
	}
}

func (s *Screen) drawString(left, top int, str string, color termbox.Attribute) {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		str := []rune(line)
		for j, char := range str {
			termbox.SetCell(left+j, top+i, char, color, BackgroundColor)
		}
	}
}

func (s *Screen) drawDebugInfo() {
	bottom := s.Top + s.field.Height*FieldYScale + 2
	header := strings.Repeat("0", LeftPromptWidth)
	header += strings.Repeat("1", FieldBoxLeftWidth)
	header += strings.Repeat("2", s.field.Width*FieldXScale)
	header += strings.Repeat("3", FieldBoxRightWidth)
	header += strings.Repeat("4", RightPromptWidth)
	header += strings.Repeat("5", LogWidth)
	s.drawString(s.Left, bottom, header, termbox.ColorDefault)

	left := s.Left +
		LeftPromptWidth +
		FieldBoxLeftWidth +
		s.field.Width*FieldXScale +
		FieldBoxRightWidth +
		RightPromptWidth
	top := s.Top
	log := Log.String(LogHeight, LogWidth-4)
	s.drawString(left+4, top, log, termbox.ColorDefault)
}

func (s *Screen) drawRightPrompt() {
	left := s.Left +
		LeftPromptWidth +
		FieldBoxLeftWidth +
		s.field.Width*FieldXScale +
		FieldBoxRightWidth +
		RightPromptLeft
	s.drawString(left, s.Top, CopyrightPrompt, CopyrightPromptColor)

	top := s.Top + CopyrightPromptHeight + 1
	s.drawString(left, top, HelpPrompt, termbox.ColorDefault)
}

func (s *Screen) drawLeftPrompt(state int, stats *Stats) {
	left := s.Left + LeftPromptLeft
	top := s.Top
	switch state {
	case StateRunning:
		s.drawString(left, top, StateRunningPrompt, StateRunningColor)
	case StatePaused:
		s.drawString(left, top, StatePausedPrompt, StatePausedColor)
	case StateFinished:
		s.drawString(left, top, StateFinishedPrompt, StateFinishedColor)
	}

	str := fmt.Sprintf(
		StatsPrompt,
		stats.Level, int(stats.Elapsed),
		stats.Blocks, stats.Lines, stats.Score,
	)
	s.drawString(left, top+2, str, termbox.ColorDefault)

	top = top + 2 + StatsPromptHeight + 2
	s.drawString(left, top, NextBlockPrompt, termbox.ColorDefault)
	s.drawField(left+NextBlockLeft, top+NextBlockTop+1, s.preview)
}

func (s *Screen) drawFrame() {
	height := s.field.Height * FieldYScale
	width := s.field.Width * FieldXScale

	top := s.Top
	bottom := top + height
	left := s.Left + LeftPromptWidth
	right := left + width + FieldBoxRightWidth

	for i := 0; i < height+1; i++ {
		str := []rune(FieldBoxLeftChars)
		for j, char := range str {
			termbox.SetCell(
				left+j, top+i, char,
				FieldBoxColor, BackgroundColor,
			)
		}

		str = []rune(FieldBoxRightChars)
		for j, char := range str {
			termbox.SetCell(
				right+j, top+i, char,
				FieldBoxColor, BackgroundColor,
			)
		}
	}

	for j := 0; j < width; j++ {
		termbox.SetCell(
			left+j+FieldBoxLeftWidth, bottom, FieldBoxBottomChar,
			FieldBoxColor, BackgroundColor,
		)
	}
}

func (s *Screen) drawField(left, top int, field *Field) {
	for i := 0; i < field.Height; i++ {
		for dj := 0; dj < FieldXScale; dj++ {
			for j := 0; j < field.Width; j++ {
				for di := 0; di < FieldYScale; di++ {
					x := left + j*FieldXScale + dj
					y := top + i*FieldYScale + di

					_, char, color, err := field.Get(j, i)
					if err != nil {
						panic(err)
					}
					termbox.SetCell(x, y, char, BackgroundColor, color)
				}
			}
		}
	}
}

func (s *Screen) Draw(state int, stats *Stats) {
	err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if err != nil {
		panic(err)
	}

	s.resize()
	s.drawFrame()
	s.drawRightPrompt()
	s.drawLeftPrompt(state, stats)

	if s.debug {
		s.drawDebugInfo()
	}

	top := s.Top
	left := s.Left + LeftPromptWidth + FieldBoxLeftWidth
	s.drawField(left, top, s.field)

	err = termbox.Flush()
	if err != nil {
		panic(err)
	}
	termbox.SetCursor(0, 0)
}

func NewScreen(field *Field, preview *Field, debug bool) *Screen {
	s := Screen{
		Top:     ScreenTop,
		Left:    ScreenMinLeft,
		field:   field,
		preview: preview,
		debug:   debug,
	}
	return &s
}
