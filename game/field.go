package game

import (
	"errors"

	"github.com/nsf/termbox-go"
)

const (
	FixedCellValue = -1
	EmptyCellValue = 0
)

var (
	OverlappingError      = errors.New("curBlock overlaps other one")
	IndexOutOfBoundsError = errors.New("out of field bounds")
)

type Cell struct {
	value int
	char  rune
	color termbox.Attribute
}

type Field struct {
	Width  int
	Height int
	cells  [][]Cell
}

func (f *Field) Set(x, y int, value int, color termbox.Attribute) error {
	if x < 0 || x >= f.Width || y < 0 || y >= f.Height {
		return IndexOutOfBoundsError
	}
	f.cells[y][x].color = color
	f.cells[y][x].value = value
	f.cells[y][x].char = BlockChar
	return nil
}

func (f *Field) Get(x, y int) (int, rune, termbox.Attribute, error) {
	if x < 0 || x >= f.Width || y < 0 || y >= f.Height {
		return 0, 0, 0, IndexOutOfBoundsError
	}
	cell := f.cells[y][x]
	return cell.value, cell.char, cell.color, nil
}

func (f *Field) Clear(full bool) {
	for i := 0; i < f.Height; i++ {
		for j := 0; j < f.Width; j++ {
			if full || f.cells[i][j].value != FixedCellValue {
				f.cells[i][j].color = BackgroundColor
				f.cells[i][j].char = BackgroundChar
				f.cells[i][j].value = EmptyCellValue
			}
		}
	}
}

func (f *Field) RemoveLines() int {
	removed := 0
	for i := 1; i < f.Height; i++ {
		full := true
		for j := 0; j < f.Width; j++ {
			if f.cells[i][j].value == EmptyCellValue {
				full = false
			}
		}
		if full {
			for k := i; k > 0; k-- {
				for j := 0; j < f.Width; j++ {
					f.cells[k][j] = f.cells[k-1][j]
				}
			}
			removed++
		}
	}
	return removed
}

func NewField(height, width int) *Field {
	f := Field{
		Width:  width,
		Height: height,
	}

	f.cells = make([][]Cell, height, height)
	for i := 0; i < height; i++ {
		f.cells[i] = make([]Cell, width, width)
	}
	f.Clear(true)

	return &f
}
