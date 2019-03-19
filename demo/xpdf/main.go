package main

import (
	"fmt"
	"os"

	"github.com/zecchan/zgolib/zpdf"
)

func main() {
	o, er := zpdf.LoadXMLFromFile("test.xml")
	if er != nil && er.Error() != "EOF" {
		fmt.Println(er.Error())
	} else {
		_, er := os.Stat("D:\\test.pdf")
		if er == nil {
			os.Remove("D:\\test.pdf")
		}
		er = o.Render("D:\\test.pdf")
		if er != nil {
			fmt.Println(er.Error())
		}
		fmt.Println("Done~")
	}
}
