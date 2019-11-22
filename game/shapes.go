package game

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

var (
	SShape Shape
	ZShape Shape
	OShape Shape
	LShape Shape
	JShape Shape
	IShape Shape
	TShape Shape
	Shapes []Shape

	lastID = 0
)

type Shape struct {
	mask  [][]int
	color termbox.Attribute
}

type Block struct {
	x     int
	y     int
	mask  [][]int
	shape Shape
	id    int
}

func (t *Block) getCenter() (int, int) {
	dy := len(t.mask) / 2
	if len(t.mask)%2 == 0 {
		dy--
	}
	dx := len(t.mask[0]) / 2
	if len(t.mask[0])%2 == 0 {
		dx--
	}
	return dx, dy
}

func (t *Block) Draw(field *Field, fixed bool) error {
	dx, dy := t.getCenter()
	for i := 0; i < len(t.mask); i++ {
		for j := 0; j < len(t.mask[i]); j++ {
			if t.mask[i][j] > 0 {
				x := t.x + j - dx
				y := t.y + i - dy

				val := t.mask[i][j]
				if fixed {
					val = FixedCellValue
				}

				err := field.Set(x, y, val, t.shape.color)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (t *Block) MustDraw(field *Field, fixed bool) {
	err := t.Draw(field, fixed)
	if err != nil {
		panic(err)
	}
}

func (t *Block) Check(field *Field) error {
	dx, dy := t.getCenter()
	for i := 0; i < len(t.mask); i++ {
		for j := 0; j < len(t.mask[i]); j++ {
			val, _, _, err := field.Get(t.x+j-dx, t.y+i-dy)
			if err != nil {
				return err
			}
			if val != EmptyCellValue && t.mask[i][j] > 0 && t.mask[i][j] != val {
				return OverlappingError
			}
		}
	}
	return nil
}

func (t *Block) TryMove(dx, dy int, field *Field) bool {
	t.x += dx
	t.y += dy
	err := t.Check(field)
	if err != nil {
		t.x -= dx
		t.y -= dy
		return false
	}
	return true
}

func (t *Block) TryRotate(field *Field) bool {
	n, m := len(t.mask), len(t.mask[0])

	mask := make([][]int, m, m)
	for i := 0; i < m; i++ {
		mask[i] = make([]int, n, n)
	}

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			mask[x][n-y-1] = t.mask[y][x]
		}
	}

	old := t.mask
	t.mask = mask

	err := t.Check(field)
	if err != nil {
		t.mask = old
		return false
	}
	return true
}

func (t *Block) Copy(x, y int) *Block {
	return NewBlock(x, y, t.shape)
}

func NewBlock(x, y int, shape Shape) *Block {
	lastID++
	n, m := len(shape.mask), len(shape.mask[0])

	mask := make([][]int, n, n)
	for y := 0; y < n; y++ {
		mask[y] = make([]int, m, m)
		for x := 0; x < m; x++ {
			if shape.mask[y][x] > 0 {
				mask[y][x] = lastID
			} else {
				mask[y][x] = 0
			}
		}
	}

	return &Block{
		x:     x,
		y:     y,
		mask:  mask,
		shape: shape,
	}
}

func NewRandomBlock(x, y int) *Block {
	i := rand.Intn(len(Shapes))
	return NewBlock(x, y, Shapes[i])
}

func init() {
	ZShape = Shape{
		mask: [][]int{
			{1, 1, 0},
			{0, 1, 1},
		},
		color: termbox.ColorRed,
	}
	SShape = Shape{
		mask: [][]int{
			{0, 1, 1},
			{1, 1, 0},
		},
		color: termbox.ColorGreen,
	}
	OShape = Shape{
		mask: [][]int{
			{1, 1},
			{1, 1},
		},
		color: termbox.ColorYellow,
	}
	IShape = Shape{
		mask: [][]int{
			{1, 1, 1, 1},
		},
		color: termbox.ColorCyan,
	}
	TShape = Shape{
		mask: [][]int{
			{0, 1, 0},
			{1, 1, 1},
		},
		color: termbox.ColorMagenta,
	}
	JShape = Shape{
		mask: [][]int{
			{1, 0, 0},
			{1, 1, 1},
		},
		color: termbox.ColorBlue,
	}
	LShape = Shape{
		mask: [][]int{
			{1, 1, 1},
			{1, 0, 0},
		},
		color: termbox.ColorWhite,
	}
	Shapes = []Shape{
		ZShape, SShape,
		OShape, IShape,
		TShape, JShape,
		LShape,
	}
}
