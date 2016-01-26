package obj

import (
	"kdtree"
	"sort"
)

func (g *Game) collisionWithWall(unit *Unit) {
	if unit.X-unit.Radius < 0 {
		unit.SpeedX = -unit.SpeedX
	}

	if float32(g.World.Width) < unit.X+unit.Radius {
		unit.SpeedX = -unit.SpeedX
	}

	if unit.Y-unit.Radius < 0 {
		unit.SpeedY = -unit.SpeedY
	}
	if float32(g.World.Height) < unit.Y+unit.Radius {
		unit.SpeedY = -unit.SpeedY
	}
}

func (g *Game) enemyCollisionWithShell(unitsMap map[int]*Unit, unitsTree *kdtree.T, radius float32) {
	for _, unit := range unitsMap {
		if unit.Type == Enemy {
			nearestNodes := unitsTree.InRange(kdtree.Point{float64(unit.X), float64(unit.Y)}, float64(radius), nil)
			if len(nearestNodes) > 1 {
				collideWithShell(nearestNodes, unit, float32(g.Step)/1000)
			}
		}
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

func collideWithShell(nearestNodes []*kdtree.T, unit *Unit, maxTimeToHit float32) {
	unitsCollisions := make([]UnitCollision, 0)
	for _, node := range nearestNodes {
		nodeUnit := node.Data.(*Unit)
		isCollision, d := unit.timeToHit(nodeUnit)
		if isCollision && nodeUnit.Health > 0 && d < maxTimeToHit {
			unitsCollisions = append(unitsCollisions, UnitCollision{nodeUnit, d})
		}
	}
	sort.Sort(UnitsCollisions(unitsCollisions))
	for _, collision := range unitsCollisions {
		if collision.Unit.Type == Bullet {
			unit.Health -= 1
			collision.Unit.Health -= 1
		}
		if unit.Health <= 0 {
			break
		}
	}
}
