package obj
import "support"

type Unit struct {
    id int
    X float32
    Y float32
    R float32
    SX float32
    SY float32
    T string
}

var currentId = 0
func NewUnit(x float32, y float32, radius float32) *Unit {
    currentId++
    return &Unit{currentId, x, y, radius, 0, 0, "ba"}
}

func (u *Unit) move() {
    u.X = support.Round2(u.X + u.SX)
    u.Y = support.Round2(u.Y + u.SY)
}