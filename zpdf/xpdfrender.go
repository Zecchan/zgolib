package zpdf

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

func (x *XPDF) Render(path string) error {
	if len(x.Page) == 0 {
		return errors.New("There is no page to render")
	}
	fp := x.Page[0]
	ori := "P"
	if fp.Orientation == "landscape" {
		ori = "L"
	}
	pdf := gofpdf.New(ori, fp.Measurement, fp.Size, "C:\\Users\\ASUS\\go\\src\\eaciit\\proactive-inv\\library\\fontpdf\\")
	pdf.SetMargins(fp.Margin.Left, fp.Margin.Top, fp.Margin.Right)

	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")

	for i, page := range x.Page {
		pdf.SetFont("Century_Gothic", "", 12)
		_, err := page.Render(XPDFDrawingBounds{
			Left:   0,
			Top:    0,
			Right:  page.Width,
			Bottom: page.Height,
		}, pdf)
		if err != nil {
			fmt.Println("Error while rendering Page element #" + strconv.Itoa(i+1) + ": " + err.Error())
		}
	}

	return pdf.OutputFileAndClose(path)
}
