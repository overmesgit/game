package obj

import (
	"encoding/json"
	"math/rand"
)

type World struct {
	Height  int
	Width   int
	Units   map[int]*Unit
	Enemies map[int]*EnemyUnit
	Players map[int]*PlayerUnit
}

func NewWorld() *World {
	return &World{Height: 800, Width: 1000, Units: make(map[int]*Unit), Enemies: make(map[int]*Enemy), Players: make(map[int]*Player)}
}

func (w *World) RemovePlayer(player *Unit) {
	delete(w.Players, player.Id)
}

func (w *World) AddPlayer(player *Unit) {
	w.Players[player.Id] = player
}

func (w *World) addRandomEnemy() {
	u := NewRandomUnit(80, Enemy, 10)
	u.X = rand.Float32()*float32(w.Width-100) + 100
	u.Y = rand.Float32()*float32(w.Height-100) + 100
	w.Enemies[u.Id] = EnemyUnit{u, 1}
}

func (w *World) removeOutBoundUnits(unitsMap map[int]hasUnit) {
	for _, current := range unitsMap {
		unit := current.Unit()
		if !(0 < unit.X && unit.X < float32(w.Width) && 0 < unit.Y && unit.Y < float32(w.Height)) {
			delete(unitsMap, unit.Id)
		}
	}
}

func (w *World) UnitsToJSON() []byte {
	resp := map[string]interface{}{
		"get":     "units",
		"units":   w.Units,
		"enemies": w.Enemies,
		"players": w.Players,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil
	}
	return b
}
