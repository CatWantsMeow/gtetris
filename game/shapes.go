package game

import (
    "math/rand"
)

var (
    ZShape = Shape{
        mask: [][]byte{
            {1, 1, 0},
            {0, 1, 1},
        },
        color: 1,
    }
    SShape = Shape{
        mask: [][]byte{
            {0, 1, 1},
            {1, 1, 0},
        },
        color: 2,
    }
    OShape = Shape{
        mask: [][]byte{
            {1, 1},
            {1, 1},
        },
        color: 3,
    }
    IShape = Shape{
        mask: [][]byte{
            {1, 1, 1, 1},
        },
        color: 4,
    }
    TShape = Shape{
        mask: [][]byte{
            {0, 1, 0},
            {1, 1, 1},
        },
        color: 5,
    }
    JShape = Shape{
        mask: [][]byte{
            {1, 0, 0},
            {1, 1, 1},
        },
        color: 6,
    }
    LShape = Shape{
        mask: [][]byte{
            {1, 1, 1},
            {1, 0, 0},
        },
        color: 7,
    }
    Shapes = []Shape{
        ZShape, SShape,
        OShape, IShape,
        TShape, JShape,
        LShape,
    }
)

type Shape struct {
    mask  [][]byte
    color uint16
}

type Block struct {
    x     int
    y     int
    mask  [][]byte
    shape Shape
}

func (b *Block) pos() (x int, y int) {
    dy := len(b.mask) / 2
    if len(b.mask)%2 == 0 {
        dy--
    }
    dx := len(b.mask[0]) / 2
    if len(b.mask[0])%2 == 0 {
        dx--
    }
    return b.x - dx, b.y - dy
}

func (b *Block) Draw(field *Field, fixed bool) error {
    x, y := b.pos()
    for i := 0; i < len(b.mask); i++ {
        for j := 0; j < len(b.mask[i]); j++ {
            if b.mask[i][j] != 0 {
                x := x + j
                y := y + i

                val := MovingCellValue
                if fixed {
                    val = FixedCellValue
                }

                err := field.Set(x, y, val, b.shape.color)
                if err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func (b *Block) MustDraw(field *Field, fixed bool) {
    err := b.Draw(field, fixed)
    if err != nil {
        panic(err)
    }
}

func (b *Block) Overlaps(field *Field) bool {
    x, y := b.pos()
    for i := 0; i < len(b.mask); i++ {
        for j := 0; j < len(b.mask[i]); j++ {
            val, _, err := field.Get(x+j, y+i)
            if err != nil {
                return true
            }
            if b.mask[i][j] > 0 && val == FixedCellValue {
                return true
            }
        }
    }
    return false
}

func (b *Block) TryMove(dx, dy int, field *Field) bool {
    b.x += dx
    b.y += dy
    if b.Overlaps(field) {
        b.x -= dx
        b.y -= dy
        return false
    }
    return true
}

func (b *Block) TryRotate(field *Field) bool {
    n := len(b.mask)
    m := len(b.mask[0])

    rotated := make([][]byte, m, m)
    for i := 0; i < m; i++ {
        rotated[i] = make([]byte, n, n)
    }

    for y := 0; y < n; y++ {
        for x := 0; x < m; x++ {
            rotated[x][n-y-1] = b.mask[y][x]
        }
    }

    old := b.mask
    b.mask = rotated
    if b.Overlaps(field) {
        b.mask = old
        return false
    }
    return true
}

func (b *Block) Copy(x, y int) *Block {
    return NewBlock(x, y, b.shape)
}

func NewBlock(x, y int, shape Shape) *Block {
    n := len(shape.mask)
    m := len(shape.mask[0])

    mask := make([][]byte, n, n)
    for i := 0; i < n; i++ {
        mask[i] = make([]byte, m, m)
        copy(mask[i], shape.mask[i])
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
