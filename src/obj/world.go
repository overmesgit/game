package obj
import (
    "math/rand"
    "encoding/json"
)

type World struct {
	Height int
	Width  int
	Units  []*Unit
    Players []*Unit
    ToDel []*Unit
}

func NewWorld() *World {
	return &World{Height: 800, Width: 1000, Units: make([]*Unit, 0),ToDel: make([]*Unit, 0)}
}

func (w *World) AddUnit(u *Unit) {
	w.Units = append(w.Units, u)
}

func (w *World) AddPlayer(u *Unit) {
    w.Players = append(w.Players, u)
    w.AddUnit(u)
}

func (w *World) AddUnits(units []*Unit) {
	w.Units = append(w.Units, units...)
}

func (w *World) RemovePlayer(player *Unit) {
    newPlayers := make([]*Unit, 0)
    for _, p := range w.Players {
        if p != player {
            newPlayers = append(newPlayers, p)
        }
    }
    w.Players = newPlayers
}

func (w *World) addRandomEnemy() {
	u := NewRandomUnit(80, Enemy, 10)
	u.X = rand.Float32()*float32(w.Width-100) + 100
	u.Y = rand.Float32()*float32(w.Height-100) + 100
	w.AddUnit(u)
}

func (w *World) DeleteUnit(unit *Unit) {
    if unit.T == Player {
        w.RemovePlayer(unit)
    }
	w.ToDel = append(w.ToDel, unit)
}

func (w *World) removeOutBoundUnits(unitsMap map[int]*Unit) {
	for _, unit := range unitsMap {
		if !(0 < unit.X && unit.X < float32(w.Width) && 0 < unit.Y && unit.Y < float32(w.Height)) {
			w.DeleteUnit(unit)
		}
	}
}

func (w *World) deleteToDelUnits(units map[int]*Unit) {
	for _, unit := range w.ToDel {
		delete(units, unit.id)
	}
}

func (w *World) UnitsToJSON() []byte {
	resp := map[string]interface{}{
		"get":   "units",
		"units": w.Units,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil
	}
	return b
}
