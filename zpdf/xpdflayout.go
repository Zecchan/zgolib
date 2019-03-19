package zpdf

import (
	"encoding/xml"

	"github.com/jung-kurt/gofpdf"
)

type XPDFStackLayout struct {
	*XPDFElement
	Direction      string `xml:"direction,attr"`
	AllowMultipage bool   `xml:"multipage,attr"`
}

func (ele *XPDFStackLayout) Parse(e *XMLElement, parent interface{}) error {
	ele.XPDFElement = &XPDFElement{}
	err := ele.XPDFElement.Parse(e, parent)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(e.Parent.InnerXML, ele)
	if err != nil {
		return err
	}
	if ele.Direction != "horizontal" {
		ele.Direction = "vertical"
	}
	return nil
}

func (ele *XPDFStackLayout) Render(bounds XPDFDrawingBounds, pdf *gofpdf.Fpdf) (*XPDFRect, error) {
	contentBounds := CalculateContentBounds(bounds, ele.XPDFElement)
	if ele.Direction == "vertical" {
		contentBounds.HAnchor = HAnchorFill
		contentBounds.VAnchor = VAnchorTop
	}
	if ele.Direction == "horizontal" {
		contentBounds.HAnchor = HAnchorLeft
		contentBounds.VAnchor = VAnchorFill
	}

	// Define border
	// border := GetBorderString(ele.Border)
	var size = XPDFRect{}

	// render elements
	for _, e := range ele.Children {
		if e != nil {
			rendSize, err := (*e).Render(contentBounds, pdf)
			if err != nil {
				return nil, err
			}
			if contentBounds.HAnchor != HAnchorFill {
				size.Height += rendSize.Height
				size.Width = contentBounds.GetWidth()
			}
			if contentBounds.VAnchor != VAnchorFill {
				size.Height = contentBounds.GetHeight()
				size.Width += rendSize.Width
			}
		}
	}

	size = ele.AddOuterBox(size)
	return &size, nil
}
