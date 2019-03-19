package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"time"

	"github.com/zecchan/zgolib/zgraph"
)

func main() {
	fmt.Println(time.Now())
	var g = zgraph.ZGraph{}
	g.Init(image.Rect(0, 0, 3000, 3000), image.Rect(0, 0, 100, 100))
	g.Func2D = func(x, y float64, cc zgraph.ColorChannel) float64 {
		if cc == zgraph.ChannelA {
			return 1
		}
		if cc == zgraph.ChannelR {
			return 1 - math.Min(100, math.Sqrt(math.Pow(x-30, 2)+math.Pow(y-30, 2)))/100
		}
		if cc == zgraph.ChannelG {
			return 1 - (x+y)/200
		}
		if cc == zgraph.ChannelB {
			return math.Abs(x-y) / 100
		}
		return 1
	}
	gr, e := g.Draw2D(nil)
	if e != nil {
		fmt.Println("Failed to draw")
		return
	}

	os.Remove("D:\\testarea\\test.png")
	fo, er := os.Create("D:\\testarea\\test.png")
	if er != nil {
		fmt.Println("Cannot create file")
		return
	}
	defer fo.Close()
	png.Encode(fo, gr)
	fmt.Println("Picture generated")
	fmt.Println(time.Now())
}
