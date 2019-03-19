package zpdf

import "encoding/xml"

type XMLElement struct {
	InnerXML   []byte         `xml:",innerxml"`
	Attributes *XMLAttributes `xml:",any,attr"`
	Children   []*XMLElement  `xml:",any"`
	XMLName    xml.Name
	Parent     *XMLElement
}

func (x *XMLElement) Init() {
	for _, c := range x.Children {
		c.Parent = x
		c.Init()
	}
}

type XMLAttributes struct {
	attr map[string]string
}

func (x *XMLAttributes) UnmarshalXMLAttr(attr xml.Attr) error {
	if x.attr == nil {
		x.attr = map[string]string{}
	}
	x.attr[attr.Name.Local] = attr.Value
	return nil
}
