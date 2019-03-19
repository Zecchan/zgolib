package main

import (
	"fmt"

	"github.com/zecchan/zgolib/syntax/grammar"
)

func main() {
	g := grammar.Grammar{}
	err := g.Parse("objmember -> <strlit> (colon) agaga\r\n// This is a comment*/\r\nagaga  -> <strlit> (colon) objmember\r\nnigga->nigga")
	if err == nil {
		fmt.Println(g)
	} else {
		fmt.Println(err.Error())
	}
}
