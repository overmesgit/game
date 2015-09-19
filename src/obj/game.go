package obj
import (
    "support"
    "time"
    "math"
    "math/rand"
    "kdtree"
)

type Game struct {
    World *World
    Step int64
}

func NewGame() *Game {
    world := NewWorld()
    game := Game{world, 50}
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

func (g *Game) addRandomUnits() {
    u := g.MakeRandomUnit(4, "ba")
    u.X = rand.Float32()*float32(g.World.Height)
    u.Y = rand.Float32()*float32(g.World.Width)
    g.World.AddUnit(u)
}

func (g *Game) MakeRandomUnit(steedRange float32, type_ string) *Unit{
    unit := NewUnit(0, 0)
    unit.Type = type_

    swap := float32(0.0)
    t := 2 * math.Pi * rand.Float64()
    u := rand.Float32() + rand.Float32()
    if u > 1 { swap = 2 - u } else { swap = u }
    unit.SpeedX = support.Round2(steedRange*swap*float32(math.Cos(t)))
    unit.SpeedY = support.Round2(steedRange*swap*float32(math.Sin(t)))

    return unit
}

func (g *Game) MakeBoom(x float32, y float32) {
    size := 200;
    newUnits := make([]*Unit, size)
    for i := 0; i < size; i++ {
        unit := g.MakeRandomUnit(10, "fr")
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

func (g *Game) makeTurn() {
    if len(g.World.Units) < 1000 {
        g.addRandomUnits()
        g.addRandomUnits()
        g.addRandomUnits()
    }

    unitsTree := kdtree.New(nil)
    unitsMap := make(map[int]*Unit)

    for _, unit := range g.World.Units {
        unitT := new(kdtree.T)
        unitT.Point = kdtree.Point{float64(unit.X), float64(unit.Y)}
        unitT.Data = unit.id
        unitsTree = unitsTree.Insert(unitT)
        unitsMap[unit.id] = unit
    }

    for _, unit := range g.World.Units {
        if unit.Type == "ba" {
            nearestNodes := unitsTree.InRange(kdtree.Point{float64(unit.X), float64(unit.Y)}, 10, nil)
            if len(nearestNodes) > 1 {
                deleteUnit := false
                for _, node := range nearestNodes {
                    nodeId := node.Data.(int)
                    if unitsMap[nodeId] != nil && unitsMap[nodeId].Type == "fr" {
                        deleteUnit = true
                        delete(unitsMap, nodeId)
                        break
                    }
                }
                if deleteUnit {
                    delete(unitsMap, unit.id)
                }
            }
        }
    }

    newUnits := make([]*Unit, len(unitsMap))
    i := 0
    for _, unit := range unitsMap {
        unit.move()
        newUnits[i] = unit
        i++
    }
    g.World.Units = newUnits

    g.removeUnits()
}