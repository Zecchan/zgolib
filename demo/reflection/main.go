package main

import (
	"fmt"

	"github.com/zecchan/zgolib/reflection"
)

type UnitLocation struct {
	X, Y int
	Name string
}

func main() {
	objA := UnitLocation{8, 8, "Colony"}
	objB := struct {
		X, Y int
	}{18, 98}

	typA := reflection.GetGoType(objA)
	val, err := typA.Create(objB)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		objC := typA.Instantiate(val).(UnitLocation)
		objC.Name = "Alice"
		fmt.Println(objC)
	}
}
