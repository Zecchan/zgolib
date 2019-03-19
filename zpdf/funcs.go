package zpdf

import (
	"encoding/xml"
	"io/ioutil"
	"math"

	"github.com/jung-kurt/gofpdf"
)

func LoadXML(stream []byte) (*XPDF, error) {
	v := XMLElement{}
	err := xml.Unmarshal(stream, &v)
	v.Init()
	if err != nil {
		return nil, err
	}
	res := XPDF{}
	err = res.Parse(&v, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func LoadXMLFromFile(path string) (*XPDF, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadXML(b)
}

func ComputeContentBounds() {

}

func (rect *XPDFRect) ReadPdfPosition(pdf *gofpdf.Fpdf) {
	x, y := pdf.GetXY()
	rect.X = x
	rect.Y = y
}

// CalculateContentBounds returns renderable bounds for elements INSIDE of currentContainer
func CalculateContentBounds(parentContainerBounds XPDFDrawingBounds, currentContainer *XPDFElement) XPDFDrawingBounds {
	res := XPDFDrawingBounds{
		Left:        parentContainerBounds.Left + currentContainer.Padding.Left + currentContainer.Margin.Left + currentContainer.Border.Left,
		Top:         parentContainerBounds.Top + currentContainer.Padding.Top + currentContainer.Margin.Top + currentContainer.Border.Top,
		Right:       parentContainerBounds.Right - currentContainer.Padding.Right - currentContainer.Margin.Right - currentContainer.Border.Right,
		Bottom:      parentContainerBounds.Bottom - currentContainer.Padding.Bottom - currentContainer.Margin.Bottom - currentContainer.Border.Bottom,
		CurrentCol:  currentContainer.Column,
		CurrentRow:  currentContainer.Row,
		HAnchor:     parentContainerBounds.HAnchor,
		VAnchor:     parentContainerBounds.VAnchor,
		Positioning: parentContainerBounds.Positioning,
	}
	if currentContainer.Width > 0 {
		res.Right = math.Min(res.Right, parentContainerBounds.Left+currentContainer.Margin.Left+currentContainer.Border.Left+currentContainer.Width-currentContainer.Padding.Right)
	}
	if currentContainer.Height > 0 {
		res.Bottom = math.Min(res.Bottom, parentContainerBounds.Top+currentContainer.Margin.Top+currentContainer.Border.Top+currentContainer.Height-currentContainer.Padding.Bottom)
	}
	return res
}

func GetBorderString(border XPDFThickness) string {
	var b = ""
	if border.Left > 0 {
		b += "L"
	}
	if border.Top > 0 {
		b += "T"
	}
	if border.Right > 0 {
		b += "R"
	}
	if border.Bottom > 0 {
		b += "B"
	}
	return b
}

func CreateElementByName(name string) IXPDFRenderableElement {
	switch name {
	case "Page":
		return &(XPDFPage{})
	case "StackLayout":
		return &(XPDFStackLayout{})
	case "Label":
		return &(XPDFLabel{})
	default:
		return nil
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
