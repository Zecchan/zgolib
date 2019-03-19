package zpdf

import "github.com/jung-kurt/gofpdf"

type IXPDFRenderableElement interface {
	Parse(ele *XMLElement, parent interface{}) error
	// Render should render the element within specified bounds
	// Occupied means occupied space within bounds
	// It should return the drawing rectangle that it uses
	Render(containerBounds XPDFDrawingBounds, pdf *gofpdf.Fpdf) (*XPDFRect, error)
}
