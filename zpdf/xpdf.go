package zpdf

import (
	"errors"
)

type XPDF struct {
	Page []*XPDFPage
	Data interface{}
}

func (x *XPDF) Parse(pdfEle *XMLElement, parent interface{}) error {
	x.Page = []*XPDFPage{}
	if pdfEle.XMLName.Local != "PDF" {
		return errors.New("This is not a zpdf template file")
	}

	for _, ce := range pdfEle.Children {
		if ce.XMLName.Local != "Page" {
			return errors.New("PDF can only contains Page")
		}
		pg := CreateElementByName("Page").(*XPDFPage)
		err := pg.Parse(ce, x)
		if err != nil {
			return err
		}
		x.Page = append(x.Page, pg)
	}

	return nil
}
