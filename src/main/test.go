package main
import (
    "code.google.com/p/eaburns/kdtree"
    "fmt"
)

func main() {
    tree := kdtree.New(nil)

    unitT1 := new(kdtree.T)
    unitT1.Point = kdtree.Point{0, 0}

    tree = tree.Insert(unitT1)

    unitT2 := new(kdtree.T)
    unitT2.Point = kdtree.Point{10, 10}

    tree = tree.Insert(unitT2)

    nearest := tree.InRange(unitT1.Point, 5, nil)
    fmt.Println(len(nearest))
    for _, n := range nearest {
        fmt.Println(n)
    }
    fmt.Println("end")
}
