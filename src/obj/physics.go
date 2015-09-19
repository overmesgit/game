package obj
import "kdtree"

func (g *Game) collisionWithWall(unit *Unit) {
    if unit.X - unit.R < 0 {
        unit.SX = -unit.SX
    }

    if float32(g.World.Width) < unit.X + unit.R {
        unit.SX = -unit.SX
    }

    if unit.Y - unit.R < 0 {
        unit.SY = -unit.SY
    }
    if float32(g.World.Height) < unit.Y + unit.R {
        unit.SY = -unit.SY
    }
}

func (g *Game) enemyCollisionWithShell(unitsMap map[int]*Unit, unitsTree *kdtree.T) {
    for _, unit := range g.World.Units {
        if unit.T == "en" {
            nearestNodes := unitsTree.InRange(kdtree.Point{float64(unit.X), float64(unit.Y)}, 10, nil)
            if len(nearestNodes) > 1 {
                unit.CollideWithShell(nearestNodes)
            }
        }
    }
}

