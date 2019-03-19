package zpdf

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
)

type XPDFThickness struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

func (box *XPDFThickness) UnmarshalXMLAttr(attr xml.Attr) error {
	spl := strings.Split(strings.Trim(attr.Value, " \t"), " ")
	spl = deleteEmpty(spl)
	if len(spl) != 0 {
		if len(spl) == 1 {
			var v, e = strconv.ParseFloat(spl[0], 64)
			if e == nil {
				box.Left = v
				box.Top = v
				box.Right = v
				box.Bottom = v
				return nil
			}
		}
		if len(spl) == 2 {
			var v1, e1 = strconv.ParseFloat(spl[0], 64)
			var v2, e2 = strconv.ParseFloat(spl[1], 64)
			if e1 == nil && e2 == nil {
				box.Left = v1
				box.Top = v2
				box.Right = v1
				box.Bottom = v2
				return nil
			}
		}
		if len(spl) == 4 {
			var v1, e1 = strconv.ParseFloat(spl[0], 64)
			var v2, e2 = strconv.ParseFloat(spl[1], 64)
			var v3, e3 = strconv.ParseFloat(spl[2], 64)
			var v4, e4 = strconv.ParseFloat(spl[3], 64)
			if e1 == nil && e2 == nil && e3 == nil && e4 == nil {
				box.Left = v1
				box.Top = v2
				box.Right = v3
				box.Bottom = v4
				return nil
			}
		}
	}
	return errors.New("Invalid value for " + attr.Name.Local + ": " + attr.Value)
}

type XPDFBounds struct {
	Left      float64
	Top       float64
	Right     float64
	Bottom    float64
	StackingH bool
	StackingV bool
	GridCell  bool
}

func (box *XPDFBounds) GetHeight() float64 {
	return box.Bottom - box.Top
}
func (box *XPDFBounds) GetWidth() float64 {
	return box.Right - box.Left
}

type XPDFElementAlignment struct {
	HAlign string
	VAlign string
}

func (box *XPDFElementAlignment) GetHorizontalAlign() string {
	alignH := "L"
	if box.HAlign == "right" {
		alignH = "R"
	}
	if box.HAlign == "center" {
		alignH = "C"
	}
	return alignH
}
func (box *XPDFElementAlignment) GetVerticalAlign() string {
	alignV := "T"
	if box.HAlign == "bottom" {
		alignV = "B"
	}
	if box.HAlign == "middle" {
		alignV = "M"
	}
	return alignV
}

func (box *XPDFElementAlignment) UnmarshalXMLAttr(attr xml.Attr) error {
	spl := strings.Split(strings.Trim(attr.Value, " \t"), " ")
	box.HAlign = ""
	box.VAlign = ""
	for _, s := range spl {
		if s == "right" || s == "left" || s == "center" {
			box.HAlign = s
		} else if s == "bottom" || s == "top" || s == "middle" {
			box.VAlign = s
		} else if s == "fillh" {
			box.HAlign = ""
		} else if s == "fillv" {
			box.VAlign = ""
		} else if s == "fill" {
			box.HAlign = ""
			box.VAlign = ""
		} else {
			return errors.New("Invalid value for " + attr.Name.Local + ": " + attr.Value)
		}
	}
	return nil
}

type XPDFRect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type XPDFHorizontalAnchor int
type XPDFVerticalAnchor int
type XPDFElementPositioning int

const (
	HAnchorLeft   XPDFHorizontalAnchor   = 0
	HAnchorCenter XPDFHorizontalAnchor   = 1
	HAnchorRight  XPDFHorizontalAnchor   = 2
	HAnchorFill   XPDFHorizontalAnchor   = -1
	VAnchorTop    XPDFVerticalAnchor     = 0
	VAnchorMiddle XPDFVerticalAnchor     = 1
	VAnchorBottom XPDFVerticalAnchor     = 2
	VAnchorFill   XPDFVerticalAnchor     = -1
	EPosRelative  XPDFElementPositioning = 0
	EPosAbsolute  XPDFElementPositioning = 1
)

type XPDFDrawingBounds struct {
	Left, Top, Right, Bottom float64
	HAnchor                  XPDFHorizontalAnchor
	VAnchor                  XPDFVerticalAnchor
	Positioning              XPDFElementPositioning
	CurrentRow               int
	CurrentCol               int
}

func (box *XPDFDrawingBounds) GetHeight() float64 {
	return box.Bottom - box.Top
}
func (box *XPDFDrawingBounds) GetWidth() float64 {
	return box.Right - box.Left
}
