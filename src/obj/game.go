package obj
import (
    "support"
    "time"
    "math/rand"
    "kdtree"
)

type Game struct {
    World *World
    Step int64
    ToDel []*Unit
}

func NewGame() *Game {
    world := NewWorld()
    game := Game{world, 50, make([]*Unit, 0)}
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
        time.Sleep(time.Duration(diff)*time.Millisecond)
    }
}

func (g *Game) addRandomEnemy() {
    u := NewRandomUnit(4, "en", 3)
    u.X = rand.Float32()*float32(g.World.Height - 100) + 100
    u.Y = rand.Float32()*float32(g.World.Width - 100) + 100
    g.World.AddUnit(u)
}

func (g *Game) AddPlayer() *Unit{
    player := NewUnit(float32(g.World.Width/2), float32(g.World.Height/2), 10)
    player.T = "pl"
    g.World.AddUnit(player)
    return player
}

func (g *Game) DeleteUnit(unit *Unit) {
    g.ToDel = append(g.ToDel, unit)
}

func (g *Game) MakeBoom(x float32, y float32) {
    size := 200;
    newUnits := make([]*Unit, size)
    for i := 0; i < size; i++ {
        unit := NewRandomUnit(10, "fr", 1)
        unit.X, unit.Y = x, y
        newUnits[i] = unit
    }
    g.World.AddUnits(newUnits)
}

func (g *Game) removeUnits() {
    newUnits := make([]*Unit, 0)

    for _, unit := range g.World.Units {
        if 0 < unit.X && unit.X < float32(g.World.Width) && 0 < unit.Y && unit.Y < float32(g.World.Height) {
            newUnits = append(newUnits, unit)
        }
    }

    g.World.Units = newUnits
}

func (g *Game) deleteToDelUnits(units map[int]*Unit) {
    for _, unit := range g.ToDel {
        delete(units, unit.id)
    }
}

func (g *Game) makeTurn() {
    if len(g.World.Units) < 1000 {
        g.addRandomEnemy()
        g.addRandomEnemy()
        g.addRandomEnemy()
    }

    unitsTree := kdtree.New(nil)
    unitsMap := make(map[int]*Unit)

    for _, unit := range g.World.Units {
        unitT := new(kdtree.T)
        unitT.Point = kdtree.Point{float64(unit.X), float64(unit.Y)}
        unitT.Data = unit
        unitsTree = unitsTree.Insert(unitT)
        unitsMap[unit.id] = unit
    }
    g.deleteToDelUnits(unitsMap)
    g.enemyCollisionWithShell(unitsMap, unitsTree)

    newUnits := make([]*Unit, 0)
    for _, unit := range unitsMap {
        if unit.H >0 {
            g.collisionWithWall(unit)
            unit.move()
            newUnits = append(newUnits, unit)
        }
    }
    g.World.Units = newUnits
}