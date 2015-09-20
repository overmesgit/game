package obj

type World struct {
	Height int
	Width  int
	Units  []*Unit
    Players []*Unit
}

func NewWorld() *World {
	return &World{Height: 800, Width: 1000, Units: make([]*Unit, 0)}
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
