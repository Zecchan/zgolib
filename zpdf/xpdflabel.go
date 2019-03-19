package zpdf

import (
	"encoding/xml"
	"math"

	"github.com/jung-kurt/gofpdf"
)

type XPDFLabel struct {
	*XPDFElement
	Text      string               `xml:"text,attr"`
	TextAlign XPDFElementAlignment `xml:"textalign,attr"`
}

func (ele *XPDFLabel) Parse(e *XMLElement, parent interface{}) error {
	ele.XPDFElement = &XPDFElement{}
	err := ele.XPDFElement.Parse(e, parent)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(e.Parent.InnerXML, ele)
	if err != nil {
		return err
	}
	return nil
}

func (ele *XPDFLabel) Render(bounds XPDFDrawingBounds, pdf *gofpdf.Fpdf) (*XPDFRect, error) {
	// Define H align
	alignH := ele.TextAlign.GetHorizontalAlign()
	if bounds.VAnchor == VAnchorFill {
		alignH = "L"
	}

	// Define V align
	alignV := ele.TextAlign.GetVerticalAlign()
	if bounds.HAnchor == HAnchorFill {
		alignV = "T"
	}

	// Define border
	border := GetBorderString(ele.Border)

	// Define cell rect
	cellRect := ele.CalculateCellRect(bounds, pdf)
	pdf.SetXY(cellRect.X, cellRect.Y)

	if cellRect.Width > 0 && cellRect.Height > 0 {
		// Draw the cell
		pdf.CellFormat(cellRect.Width, cellRect.Height, ele.Text, "", 0, alignH+alignV, false, 0, "")

		if border != "" {
			borderRect := ele.AddPadding(cellRect)
			pdf.SetXY(borderRect.X, borderRect.Y)
			pdf.CellFormat(borderRect.Width, borderRect.Height, "", border, 0, alignH+alignV, false, 0, "")
		}
		cellRect = ele.AddOuterBox(cellRect)
		return &cellRect, nil
	} else {
		return &XPDFRect{
			Width:  math.Max(cellRect.Width, 0),
			Height: math.Max(cellRect.Height, 0),
		}, nil
	}
}

func (ele *XPDFLabel) CalculateCellRect(bounds XPDFDrawingBounds, pdf *gofpdf.Fpdf) XPDFRect {
	rect := CalculateContentBounds(bounds, ele.XPDFElement)

	_, ftH := pdf.GetFontSize()
	ftW := pdf.GetStringWidth(ele.Text)
	if bounds.HAnchor != HAnchorFill {
		if ele.Width > 0 {
			rect.Right = math.Min(rect.Right, rect.Left+ele.Width-ele.Padding.Right-ele.Padding.Left)
		} else {
			rect.Right = math.Min(rect.Right, rect.Left+ftW)
		}
	}
	if bounds.VAnchor == VAnchorFill {
		if ele.Height > 0 {
			rect.Bottom = math.Min(rect.Bottom, rect.Top+ele.Height-ele.Padding.Bottom-ele.Padding.Top)
		} else {
			rect.Right = math.Min(rect.Right, rect.Left+ftH)
		}
	}

	return XPDFRect{
		X:      rect.Left,
		Y:      rect.Top,
		Width:  rect.GetWidth(),
		Height: rect.GetHeight(),
	}
}
