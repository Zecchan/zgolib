package zgraph

import (
	"errors"
	"image"
	"image/color"
	"math"
)

// ColorChannel is used to distinguish color channel (RGBA) for ZGraph functions
type ColorChannel int

const (
	// ChannelR the Red channel of the pixel
	ChannelR ColorChannel = 0
	// ChannelG the Green channel of the pixel
	ChannelG ColorChannel = 1
	// ChannelB the Blue channel of the pixel
	ChannelB ColorChannel = 2
	// ChannelA the Alpha channel of the pixel
	ChannelA ColorChannel = 3
)

// ZGraph is used to create 2D graph
type ZGraph struct {
	// Image is the image object of this graph
	Image image.RGBA
	// DataBounds is the bound of the data represented on the image
	DataBounds image.Rectangle
	// ImageBounds is the bound of the image
	ImageBounds image.Rectangle

	// Func1D is the function to generate y data from point x in dataSeries
	Func1D func(x float64, dataSeries int) float64
	// Func2D is the function to generate color channel strength from point (x,y). Should only return 0.0~1.0
	Func2D func(x float64, y float64, channel ColorChannel) float64

	// Series1DColors is a map of color used by Func1D graphs
	Series1DColors []color.RGBA

	// ThreadCount specifies how many threads to create for drawing the graphs
	ThreadCount int
}

// Init will initialize the ZGraph
func (z *ZGraph) Init(bounds image.Rectangle, dataBounds image.Rectangle) {
	z.Image = *image.NewRGBA(bounds)
	z.ImageBounds = bounds
	z.DataBounds = dataBounds
	z.ThreadCount = 4
}

func (z *ZGraph) dataBoundToPixel(x, y float64) (int, int) {
	var imW = z.ImageBounds.Max.X - z.ImageBounds.Min.X
	var imH = z.ImageBounds.Max.Y - z.ImageBounds.Min.Y
	var dbW = z.DataBounds.Max.X - z.DataBounds.Min.X
	var dbH = z.DataBounds.Max.Y - z.DataBounds.Min.Y

	var dbX = x - float64(z.DataBounds.Min.X)
	var dbY = y - float64(z.DataBounds.Min.Y)

	var pxX = dbX*float64(imW)/float64(dbW) + float64(z.ImageBounds.Min.X)
	var pxY = float64(z.ImageBounds.Max.Y) - dbY*float64(imH)/float64(dbH)
	return int(math.Round(pxX)), int(math.Round(pxY))
}
func (z *ZGraph) pixelToDataBound(x, y int) (float64, float64) {
	var imW = z.ImageBounds.Max.X - z.ImageBounds.Min.X
	var imH = z.ImageBounds.Max.Y - z.ImageBounds.Min.Y
	var dbW = z.DataBounds.Max.X - z.DataBounds.Min.X
	var dbH = z.DataBounds.Max.Y - z.DataBounds.Min.Y

	var pxX = float64(x - z.ImageBounds.Min.X)
	var pxY = float64(y - z.ImageBounds.Min.Y)

	var dbX = pxX*float64(dbW)/float64(imW) + float64(z.DataBounds.Min.X)
	var dbY = float64(z.DataBounds.Max.Y) - pxY*float64(dbH)/float64(imH)
	return dbX, dbY
}

// Set1DSeriesColor will set the color of the series that is contained in 1D graph
// It is also defines the number of series it has
func (z *ZGraph) Set1DSeriesColor(colors ...color.RGBA) {
	z.Series1DColors = colors
}

// Draw1D draws a graph based on function, the func1D must return value between 0..1 inclusive
func (z *ZGraph) Draw1D(func1D func(x float64, dataSeries int) float64) (*image.RGBA, error) {
	var f = func1D
	if f == nil {
		f = z.Func1D
	}
	if f == nil {
		return nil, errors.New("Func1D is not set")
	}
	var series = len(z.Series1DColors)
	if series == 0 {
		return nil, errors.New("Series is not set")
	}

	var thCount = 0
	if z.ThreadCount <= 0 {
		z.ThreadCount = 1
	}
	var chans = make(chan interface{}, z.ThreadCount)
	for x := z.ImageBounds.Min.X; x <= z.ImageBounds.Max.X; x++ {
		dbX, _ := z.pixelToDataBound(x, 0)
		if thCount < z.ThreadCount {
			chans <- z.plot1D(f, dbX)
		} else {
			_ = <-chans
			chans <- z.plot1D(f, dbX)
		}
	}
	for thCount > 0 {
		_ = <-chans
		thCount--
	}

	return &z.Image, nil
}

func (z *ZGraph) plot1D(f func(x float64, dataSeries int) float64, dataX float64) interface{} {
	for s := 0; s < len(z.Series1DColors); s++ {
		dataY := f(dataX, s)
		pxX, pxY := z.dataBoundToPixel(dataX, dataY)
		serie := z.Series1DColors[s]
		z.Image.SetRGBA(pxX, pxY, serie)
	}
	return nil
}

// Draw2D draws a 2D image based on function, the func2D must return value between 0..1 inclusive
func (z *ZGraph) Draw2D(func2D func(x float64, y float64, channel ColorChannel) float64) (image.Image, error) {
	var f = func2D
	if f == nil {
		f = z.Func2D
	}
	if f == nil {
		return nil, errors.New("Func2D is not set")
	}

	var thCount = 0
	if z.ThreadCount <= 0 {
		z.ThreadCount = 1
	}
	var chans = make(chan interface{}, z.ThreadCount)
	for y := z.ImageBounds.Min.Y; y <= z.ImageBounds.Max.Y; y++ {
		for x := z.ImageBounds.Min.X; x <= z.ImageBounds.Max.X; x++ {
			dbX, dbY := z.pixelToDataBound(x, y)
			if thCount < z.ThreadCount {
				chans <- z.plot2D(f, dbX, dbY)
				thCount++
			} else {
				_ = <-chans
				chans <- z.plot2D(f, dbX, dbY)
			}
		}
	}
	for thCount > 0 {
		_ = <-chans
		thCount--
	}

	return &z.Image, nil
}

func (z *ZGraph) plot2D(f func(x float64, y float64, channel ColorChannel) float64, dataX float64, dataY float64) interface{} {
	pxX, pxY := z.dataBoundToPixel(dataX, dataY)
	var col = color.RGBA{
		A: uint8(f(dataX, dataY, ChannelA) * 255),
		R: uint8(f(dataX, dataY, ChannelR) * 255),
		G: uint8(f(dataX, dataY, ChannelG) * 255),
		B: uint8(f(dataX, dataY, ChannelB) * 255),
	}
	z.Image.SetRGBA(pxX, pxY, col)
	return nil
}
