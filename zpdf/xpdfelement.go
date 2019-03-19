package zpdf

import (
	"encoding/xml"
	"errors"
)

type XPDFElement struct {
	Width      float64              `xml:"width,attr"`
	Height     float64              `xml:"height,attr"`
	Padding    XPDFThickness        `xml:"padding,attr"`
	Margin     XPDFThickness        `xml:"margin,attr"`
	Location   XPDFThickness        `xml:"location,attr"`
	Border     XPDFThickness        `xml:"border,attr"`
	Row        int                  `xml:"row,attr"`
	Column     int                  `xml:"column,attr"`
	RowSpan    int                  `xml:"rowspan,attr"`
	ColSpan    int                  `xml:"colspan,attr"`
	Alignment  XPDFElementAlignment `xml:"align,attr"`
	Parent     interface{}
	Children   []*IXPDFRenderableElement
	validChild []string
}

func (ele *XPDFElement) Parse(e *XMLElement, parent interface{}) error {
	ele.Parent = parent
	err := xml.Unmarshal(e.Parent.InnerXML, ele)
	if err != nil {
		return err
	}
	ele.Children = []*IXPDFRenderableElement{}
	for _, c := range e.Children {
		if ele.validChild != nil {
			if !contains(ele.validChild, c.XMLName.Local) {
				return errors.New("A " + e.XMLName.Local + " element cannot contains " + c.XMLName.Local)
			}
		}
		iele := CreateElementByName(c.XMLName.Local)
		if iele != nil {
			iele.Parse(c, ele)
			ele.Children = append(ele.Children, &iele)
		}
	}
	return nil
}

func (ele *XPDFElement) Page() *XPDFPage {
	pgele, ok := ele.Parent.(*XPDFPage)
	if ok && pgele != nil {
		return pgele
	}
	oele, ok := ele.Parent.(*XPDFElement)
	if ok && oele != nil {
		return (*oele).Page()
	}
	return nil
}

func (ele *XPDFElement) GetBox(parentBound XPDFBounds, overflow bool) XPDFRect {
	res := XPDFRect{
		X:      parentBound.Left + ele.Margin.Left,
		Y:      parentBound.Top + ele.Margin.Top,
		Width:  parentBound.GetWidth() - ele.Margin.Left - ele.Margin.Right,
		Height: parentBound.GetHeight() - ele.Margin.Top - ele.Margin.Bottom,
	}

	if parentBound.StackingH {
		if ele.Height > 0 {
			res.Height = parentBound.GetHeight()
		}
	}

	if overflow {

	} else {

	}

	return res
}

func (ele *XPDFElement) AddOuterBox(size XPDFRect) XPDFRect {
	return XPDFRect{
		X:      size.X - ele.Padding.Left - ele.Margin.Left - ele.Border.Left,
		Y:      size.Y - ele.Padding.Top - ele.Margin.Top - ele.Border.Top,
		Width:  size.Width + ele.Padding.Left + ele.Padding.Right + ele.Margin.Right + ele.Margin.Left + ele.Border.Right + ele.Border.Left,
		Height: size.Height + ele.Padding.Top + ele.Padding.Bottom + ele.Margin.Top + ele.Margin.Bottom + ele.Border.Top + ele.Border.Bottom,
	}
}

func (ele *XPDFElement) AddPadding(size XPDFRect) XPDFRect {
	return XPDFRect{
		X:      size.X - ele.Padding.Left - ele.Padding.Right,
		Y:      size.Y - ele.Padding.Top - ele.Padding.Bottom,
		Width:  size.Width + ele.Padding.Left + ele.Padding.Right,
		Height: size.Height + ele.Padding.Top + ele.Padding.Bottom,
	}
}
