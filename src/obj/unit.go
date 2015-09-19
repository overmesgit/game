package obj
import "support"

type Unit struct {
    id int
    X float32
    Y float32
    SpeedX float32
    SpeedY float32
    Type string
}

var currentId = 0
func NewUnit(x float32, y float32) *Unit {
    currentId++
    return &Unit{currentId, x, y, 0, 0, "ba"}
}

func (u *Unit) move() {
    u.X = support.Round2(u.X + u.SpeedX)
    u.Y = support.Round2(u.Y + u.SpeedY)
}