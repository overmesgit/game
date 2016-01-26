package obj

import (
	"math"
	"math/rand"
)

const (
	Enemy  = 1
	Player = 2
	Bullet = 3
)

type Unit struct {
	Id     int
	X      float32
	Y      float32
	Radius float32
	SpeedX float32
	SpeedY float32
	Type   int
}

type hasUnit interface {
	Unit() *Unit
	Move()
}

func (unit *Unit) Unit() *Unit {
	return unit
}

var currentId = 0

func NewUnit(x float32, y float32, radius float32, type_ int) *Unit {
	currentId++
	return &Unit{currentId, x, y, radius, 0, 0, type_}
}

func NewRandomUnit(speed float32, type_ int, radius float32) *Unit {
	unit := NewUnit(0, 0, radius, type_)

	swap := float32(0.0)
	t := 2 * math.Pi * rand.Float64()
	u := rand.Float32() + rand.Float32()
	if u > 1 {
		swap = 2 - u
	} else {
		swap = u
	}
	unit.SpeedX = speed * swap * float32(math.Cos(t))
	unit.SpeedY = speed * swap * float32(math.Sin(t))

	return unit
}

//===================================== Enemy ======================================================================
type EnemyUnit struct {
	*Unit
	Health    int
	DeathTurn int
}

func (enemy *EnemyUnit) Unit() *Unit {
	return enemy.Unit
}

func NewEnemy(x float32, y float32, radius float32, health int) *EnemyUnit {
	return &EnemyUnit{NewUnit(x, y, radius, Enemy), health}
}

//===================================== Player =====================================================================
type PlayerUnit struct {
	*Unit
	Health  int
	TargetX float32
	TargetY float32
	Fire    bool
}

func (player *PlayerUnit) Unit() *Unit {
	return player.Unit
}

func NewPlayer(x float32, y float32, radius float32, health int) *PlayerUnit {
	return &EnemyUnit{NewUnit(x, y, radius, Player), health}
}

//===================================== Api ======================================================================
func (u *Unit) move(gameStep int64) {
	u.X = u.X + u.SpeedX*float32(gameStep)/1000
	u.Y = u.Y + u.SpeedY*float32(gameStep)/1000
}

func (a *Unit) timeToHit(b *Unit) (bool, float32) {
	if a == b {
		return false, 0
	}
	dx, dy := b.X-a.X, b.Y-a.Y
	dvx, dvy := b.SpeedX-a.SpeedX, b.SpeedY-a.SpeedY
	dvdr := dx*dvx + dy*dvy
	if dvdr > 0 {
		return false, 0
	}
	dvdv := dvx*dvx + dvy*dvy
	drdr := dx*dx + dy*dy
	sigma := a.Radius + b.Radius
	d := dvdr*dvdr - dvdv*(drdr-sigma*sigma)
	if d < 0 {
		return false, 0
	}
	return true, -(dvdr + float32(math.Sqrt(float64(d)))) / dvdv
}

func (player *PlayerUnit) SetPlayerMoveSpeed(pressedKeys map[string]interface{}) {
	player.SpeedX, player.SpeedY = 0, 0
	targetX, targetY := player.X, player.Y
	if pressedKeys["W"] != nil && pressedKeys["W"].(bool) {
		targetY -= 1
	}
	if pressedKeys["A"] != nil && pressedKeys["A"].(bool) {
		targetX -= 1
	}
	if pressedKeys["S"] != nil && pressedKeys["S"].(bool) {
		targetY += 1
	}
	if pressedKeys["D"] != nil && pressedKeys["D"].(bool) {
		targetX += 1
	}
	if player.X-targetX != 0 || player.Y-targetY != 0 {
		player.setSpeedToXY(targetX, targetY, PlayerSpeed)
	}
}

func (u *Unit) setSpeedToXY(targetX float32, targetY float32, speed float32) {
	c := math.Hypot(float64(targetX-u.X), float64(targetY-u.Y))
	alpha := math.Asin(float64(u.Y-targetY) / c)
	if targetX > u.X {
		u.SpeedX = speed * float32(math.Cos(alpha))
		u.SpeedY = -speed * float32(math.Sin(alpha))
	} else {
		u.SpeedX = -speed * float32(math.Cos(-alpha))
		u.SpeedY = speed * float32(math.Sin(-alpha))
	}
}

func (u *Unit) setSpeedToUnit(target *Unit, speed float32) {
	c := u.distance(target)
	alpha := math.Asin(float64(u.Y-target.Y) / c)
	if target.X > u.X {
		u.SpeedX = speed * float32(math.Cos(alpha))
		u.SpeedY = -speed * float32(math.Sin(alpha))
	} else {
		u.SpeedX = -speed * float32(math.Cos(-alpha))
		u.SpeedY = speed * float32(math.Sin(-alpha))
	}
}

func (u *Unit) distance(target *Unit) float64 {
	return math.Hypot(float64(target.X-u.X), float64(target.Y-u.Y))
}

func (player *PlayerUnit) playerBullet(speed float32) *Unit {
	bullet := NewUnit(player.X, player.Y, 1, Bullet)
	bullet.setSpeedToXY(player.TargetX, player.TargetY, speed)
	return bullet
}

func (u *Unit) moveToNearestPlayer(players []*Unit, speed float32) {
	if len(players) > 0 {
		currentMin := float64(99999)
		nearestPlayer := players[0]
		for _, p := range players {
			d := u.distance(p)
			if d < currentMin {
				currentMin = d
				nearestPlayer = p
			}
		}
		if currentMin > 100 {
			u.setSpeedToUnit(nearestPlayer, speed)
		} else {
			u.setSpeedToUnit(nearestPlayer, 0)
		}
	}

}
