package obj

type World struct {
	Height int
	Width  int
	Units  []*Unit
}

func NewWorld() *World {
	return &World{Height: 800, Width: 1000, Units: make([]*Unit, 0)}
}

func (w *World) AddUnit(u *Unit) {
	w.Units = append(w.Units, u)
}

func (w *World) AddUnits(units []*Unit) {
	w.Units = append(w.Units, units...)
}

func (w *World) RemoveUnit(i int) {
	w.Units = append(w.Units[:i], w.Units[i+1:]...)
}
