package obj

import (
	"kdtree"
	"math"
	"math/rand"
	"sort"
	"support"
)

const (
    Enemy = "en"
    Player = "pl"
    Bullet = "bu"
)

type Unit struct {
	id int
	X  float32
	Y  float32
	R  float32 // radius
	SX float32 //speed x in 1 second
	SY float32 //speed y in 1 second
	T  string  //type
	H  int     //health
	DX float32 `json:"-"` //direction x
	DY float32 `json:"-"` //direction y
	F  bool    `json:"-"` //fire
}

var currentId = 0

func NewUnit(x float32, y float32, radius float32) *Unit {
	currentId++
	return &Unit{currentId, x, y, radius, 0, 0, Enemy, 1, 0, 0, false}
}

func NewRandomUnit(speed float32, type_ string, radius float32) *Unit {
	unit := NewUnit(0, 0, radius)
	unit.T = type_

	swap := float32(0.0)
	t := 2 * math.Pi * rand.Float64()
	u := rand.Float32() + rand.Float32()
	if u > 1 {
		swap = 2 - u
	} else {
		swap = u
	}
	unit.SX = support.Round2(speed * swap * float32(math.Cos(t)))
	unit.SY = support.Round2(speed * swap * float32(math.Sin(t)))

	return unit
}

func (u *Unit) move(gameStep int64) {
	u.X = support.Round2(u.X + u.SX*float32(gameStep)/1000)
	u.Y = support.Round2(u.Y + u.SY*float32(gameStep)/1000)
}

func (a *Unit) timeToHit(b *Unit) (bool, float32) {
	if a == b {
		return false, 0
	}
	dx, dy := b.X-a.X, b.Y-a.Y
	dvx, dvy := b.SX-a.SX, b.SY-a.SY
	dvdr := dx*dvx + dy*dvy
	if dvdr > 0 {
		return false, 0
	}
	dvdv := dvx*dvx + dvy*dvy
	drdr := dx*dx + dy*dy
	sigma := a.R + b.R
	d := dvdr*dvdr - dvdv*(drdr-sigma*sigma)
	if d < 0 {
		return false, 0
	}
	return true, -(dvdr + float32(math.Sqrt(float64(d)))) / dvdv
}

func (player *Unit) SetPlayerMoveSpeed(pressedKeys map[string]interface{}) {
	player.SX, player.SY = 0, 0
	if pressedKeys["W"] != nil && pressedKeys["W"].(bool) {
		player.SY -= 100
	}
	if pressedKeys["A"] != nil && pressedKeys["A"].(bool) {
		player.SX -= 100
	}
	if pressedKeys["S"] != nil && pressedKeys["S"].(bool) {
		player.SY += 100
	}
	if pressedKeys["D"] != nil && pressedKeys["D"].(bool) {
		player.SX += 100
	}
	if player.SX != 0 && player.SY != 0 {
		player.SX *= 1.41 / 2
		player.SY *= 1.41 / 2
	}
}

type UnitCollision struct {
	Unit *Unit
	d    float32
}
type UnitsCollisions []UnitCollision

func (a UnitsCollisions) Len() int           { return len(a) }
func (a UnitsCollisions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UnitsCollisions) Less(i, j int) bool { return a[i].d < a[j].d }

func (u *Unit) CollideWithShell(nearestNodes []*kdtree.T, maxTimeToHit float32) {
	unitsCollisions := make([]UnitCollision, 0)
	for _, node := range nearestNodes {
		nodeUnit := node.Data.(*Unit)
		isCollision, d := u.timeToHit(nodeUnit)
		if isCollision && nodeUnit.H > 0 && d < maxTimeToHit {
			unitsCollisions = append(unitsCollisions, UnitCollision{nodeUnit, d})
		}
	}
	sort.Sort(UnitsCollisions(unitsCollisions))
	for _, collision := range unitsCollisions {
		if collision.Unit.T == Bullet {
			u.H -= 1
			collision.Unit.H -= 1
		}
		if u.H <= 0 {
			break
		}
	}
}
