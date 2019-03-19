package main

import (
	"fmt"

	"github.com/zecchan/zgolib/strformat"
)

func main() {
	id := strformat.NumeralCreateEnglish()
	tbl := id.Convert(98237221597389.37182, 3)
	fmt.Println(strformat.Capitalize(tbl))
}
