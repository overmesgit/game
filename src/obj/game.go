package obj

import (
	"kdtree"
	"support"
	"time"
)

type Game struct {
	World *World
	Step  int64
}

const (
	FrameStep    int64   = 50
	MaximumUnits int     = 100
	BulletSpeed  float32 = 800
	PlayerSpeed  float32 = 100
	EnemySpeed   float32 = 80
)

func NewGame() *Game {
	world := NewWorld()
	game := Game{world, FrameStep}
	return &game
}

func (g *Game) Start() {
	for {
		g.turn()
	}
}

func (g *Game) turn() {
	start := support.MakeTimestamp()

	g.makeTurn()

	end := support.MakeTimestamp()
	if diff := g.Step - (end - start); diff > 0 {
		time.Sleep(time.Duration(diff) * time.Millisecond)
	}
}

func (g *Game) makeTurn() {
	if len(g.World.Units) < MaximumUnits {
		g.World.addRandomEnemy()
	}

	unitsTree := kdtree.New(nil)
	unitsMap := make(map[int]*Unit)

	for _, unit := range g.World.Units {
		if unit.Health > 0 {
			unitsTree = insertUnitToKdTree(unitsTree, unit)
			unitsMap[unit.Id] = unit

			if unit.Fire {
				bullet := unit.unitBullet(BulletSpeed)
				unitsTree = insertUnitToKdTree(unitsTree, bullet)
				unitsMap[bullet.Id] = bullet
			}

			if unit.Type == Enemy {
				unit.moveToNearestPlayer(g.World.Players, EnemySpeed)
			}

		}

	}

	g.World.deleteToDelUnits(unitsMap)
	g.World.removeOutBoundUnits(unitsMap)
	g.enemyCollisionWithShell(unitsMap, unitsTree, BulletSpeed)

	newUnits := make([]*Unit, 0)
	for _, unit := range unitsMap {
		unit.move(g.Step)
		newUnits = append(newUnits, unit)
	}
	g.World.Units = newUnits

}

func insertUnitToKdTree(tree *kdtree.T, unit *Unit) *kdtree.T {
	unitT := new(kdtree.T)
	unitT.Point = kdtree.Point{float64(unit.X), float64(unit.Y)}
	unitT.Data = unit
	return tree.Insert(unitT)
}

func (g *Game) AddPlayer() *Unit {
	player := NewPlayer(float32(g.World.Width/2), float32(g.World.Height/2), 10, 100)
	g.World.AddPlayer(player)
	return player
}

func (g *Game) RemovePlayer(player *Unit) {
	delete(g.World.Players, player.Id)
}

func (g *Game) MakeBoom(x float32, y float32) {
	size := 200
	newUnits := make([]*Unit, size)
	for i := 0; i < size; i++ {
		unit := NewRandomUnit(200, Bullet, 1)
		unit.X, unit.Y = x, y
		newUnits[i] = unit
	}
	g.World.AddUnits(newUnits)
}
