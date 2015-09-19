package obj
import (
    "support"
    "math"
    "math/rand"
    "kdtree"
    "sort"
)

type Unit struct {
    id int
    X float32
    Y float32
    R float32
    SX float32
    SY float32
    T string
    H int
}

var currentId = 0
func NewUnit(x float32, y float32, radius float32) *Unit {
    currentId++
    return &Unit{currentId, x, y, radius, 0, 0, "en", 1}
}

func NewRandomUnit(steedRange float32, type_ string, radius float32) *Unit{
    unit := NewUnit(0, 0, radius)
    unit.T = type_

    swap := float32(0.0)
    t := 2 * math.Pi * rand.Float64()
    u := rand.Float32() + rand.Float32()
    if u > 1 { swap = 2 - u } else { swap = u }
    unit.SX = support.Round2(steedRange*swap*float32(math.Cos(t)))
    unit.SY = support.Round2(steedRange*swap*float32(math.Sin(t)))

    return unit
}

func (u *Unit) move() {
    u.X = support.Round2(u.X + u.SX)
    u.Y = support.Round2(u.Y + u.SY)
}

func (a *Unit) timeToHit(b *Unit) (bool, float32) {
    if a == b { return false, 0 }
    dx, dy := b.X - a.X, b.Y - a.Y
    dvx, dvy := b.SX - a.SX, b.SY - a.SY
    dvdr := dx*dvx + dy*dvy
    if dvdr > 0 { return false, 0 }
    dvdv := dvx*dvx + dvy*dvy
    drdr := dx*dx + dy*dy
    sigma := a.R + b.R
    d := dvdr*dvdr - dvdv*(drdr - sigma*sigma)
    if d < 0 { return false, 0 }
    return true, -(dvdr + float32(math.Sqrt(float64(d)))) / dvdv
}

func (player *Unit) SetPlayerMoveSpeed(pressedKeys map[string]interface {}) {
    player.SX, player.SY = 0, 0
    if pressedKeys["W"] != nil && pressedKeys["W"].(bool) { player.SY -= 5 }
    if pressedKeys["A"] != nil && pressedKeys["A"].(bool) { player.SX -= 5 }
    if pressedKeys["S"] != nil && pressedKeys["S"].(bool) { player.SY += 5 }
    if pressedKeys["D"] != nil && pressedKeys["D"].(bool) { player.SX += 5 }
    if player.SX != 0 && player.SY != 0 {
        player.SX *= 1.41/2
        player.SY *= 1.41/2
    }
}

type UnitCollision struct {
    Unit *Unit
    d float32
}
type UnitsCollisions []UnitCollision
func (a UnitsCollisions) Len() int           { return len(a) }
func (a UnitsCollisions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UnitsCollisions) Less(i, j int) bool { return a[i].d < a[j].d }

func (u *Unit) CollideWithShell(nearestNodes []*kdtree.T) {
    unitsCollisions := make([]UnitCollision, 0)
    for _, node := range nearestNodes {
        nodeUnit := node.Data.(*Unit)
        isCollision, d := u.timeToHit(nodeUnit)
        if isCollision && nodeUnit.H > 0 {
            unitsCollisions = append(unitsCollisions, UnitCollision{nodeUnit, d})
        }
    }
    sort.Sort(UnitsCollisions(unitsCollisions))
    for _, collision := range unitsCollisions {
        if collision.Unit.T == "fr" {
            u.H -= 1
            collision.Unit.H -= 1
        }
        if u.H <= 0 {
            break
        }
    }
}