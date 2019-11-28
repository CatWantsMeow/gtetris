package game

import (
	"errors"
)

const (
	EmptyCellValue byte = iota
	MovingCellValue
	FixedCellValue
)

var (
	IndexOutOfBoundsError = errors.New("out of field bounds")
)

type Cell struct {
	value byte
	color uint16
}

type Field struct {
	Width  int
	Height int
	cells  [][]Cell
}

func (f *Field) Set(x, y int, value byte, color uint16) error {
	if x < 0 || x >= f.Width || y < 0 || y >= f.Height {
		return IndexOutOfBoundsError
	}
	f.cells[y][x].color = color
	f.cells[y][x].value = value
	return nil
}

func (f *Field) Get(x, y int) (value byte, color uint16, err error) {
	if x < 0 || x >= f.Width || y < 0 || y >= f.Height {
		return 0, 0, IndexOutOfBoundsError
	}
	cell := f.cells[y][x]
	return cell.value, cell.color, nil
}

func (f *Field) Clear(full bool) {
	for i := 0; i < f.Height; i++ {
		for j := 0; j < f.Width; j++ {
			if full || f.cells[i][j].value != FixedCellValue {
				f.cells[i][j].value = EmptyCellValue
				f.cells[i][j].color = 0
			}
		}
	}
}

func (f *Field) RemoveFilledLines() (removed int) {
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
