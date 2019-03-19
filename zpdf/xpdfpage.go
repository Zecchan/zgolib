package zpdf

import (
	"encoding/xml"

	"github.com/jung-kurt/gofpdf"
)

type XPDFPageMember struct {
}

type XPDFPage struct {
	Measurement string `xml:"measurement,attr"`
	Size        string `xml:"size,attr"`
	Orientation string `xml:"orientation,attr"`
	PPI         int    `xml:"ppi,attr"`
	*XPDFElement
}

func (ele *XPDFPage) Render(bounds XPDFDrawingBounds, pdf *gofpdf.Fpdf) (*XPDFRect, error) {
	initPageBound := CalculateContentBounds(bounds, ele.XPDFElement)
	curPageBound := initPageBound
	contentSize := XPDFRect{}
	pdf.AddPage()
	pdf.SetPageBox("bleedbox", initPageBound.Left, initPageBound.Top, initPageBound.GetWidth(), initPageBound.GetHeight())
	for _, layout := range ele.Children {
		if layout != nil {
			rendSize, err := (*layout).Render(curPageBound, pdf)
			if err != nil {
				return nil, err
			}
			// Layout will always occupies full width
			var bottom = curPageBound.Top + rendSize.Height
			for bottom > initPageBound.GetHeight() {
				bottom -= initPageBound.GetHeight()
			}
			contentSize.Height += rendSize.Height
			contentSize.Width = bounds.GetWidth()
			curPageBound.Top = bottom
		}
	}
	return &contentSize, nil
}

func (ele *XPDFPage) Parse(e *XMLElement, parent interface{}) error {
	ele.XPDFElement = &XPDFElement{}
	ele.XPDFElement.validChild = []string{"StackLayout"}
	err := ele.XPDFElement.Parse(e, parent)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(e.Parent.InnerXML, ele)
	if err != nil {
		return err
	}
	// parse str to float array
	if ele.Measurement != "mm" && ele.Measurement != "cm" && ele.Measurement != "in" && ele.Measurement != "px" {
		ele.Measurement = "mm"
	}
	if ele.PPI == 0 {
		ele.PPI = 300
	}
	var sizeW, sizeH float64
	switch ele.Size {
	case "A5":
		sizeW = ele.ConvertMeasurement(148, "mm", ele.Measurement)
		sizeH = ele.ConvertMeasurement(210, "mm", ele.Measurement)
	default:
		sizeW = ele.ConvertMeasurement(210, "mm", ele.Measurement)
		sizeH = ele.ConvertMeasurement(297, "mm", ele.Measurement)
		ele.Size = "A4"
	}

	if ele.Width == 0 || true {
		ele.XPDFElement.Width = sizeW
	}
	if ele.Height == 0 || true {
		ele.XPDFElement.Height = sizeH
	}

	if ele.Orientation != "landscape" {
		ele.Orientation = "portrait"
	}

	return nil
}

func (ele *XPDFPage) ConvertMeasurement(value float64, fromUnit, toUnit string) float64 {
	if fromUnit == toUnit {
		return value
	}
	if fromUnit == "mm" && toUnit == "cm" {
		return value / 10
	}
	if fromUnit == "mm" && toUnit == "in" {
		return value / 25.4
	}
	if fromUnit == "mm" && toUnit == "px" {
		return value / 25.4 * float64(ele.PPI)
	}
	if fromUnit == "cm" && toUnit == "mm" {
		return value * 10
	}
	if fromUnit == "cm" && toUnit == "in" {
		return value / 2.54
	}
	if fromUnit == "cm" && toUnit == "px" {
		return value / 2.54 * float64(ele.PPI)
	}
	if fromUnit == "in" && toUnit == "mm" {
		return value * 25.4
	}
	if fromUnit == "in" && toUnit == "cm" {
		return value * 2.54
	}
	if fromUnit == "in" && toUnit == "px" {
		return value * float64(ele.PPI)
	}
	if fromUnit == "px" && toUnit == "mm" {
		return value / float64(ele.PPI) * 25.4
	}
	if fromUnit == "px" && toUnit == "cm" {
		return value / float64(ele.PPI) * 2.54
	}
	if fromUnit == "px" && toUnit == "in" {
		return value / float64(ele.PPI)
	}
	return value
}
