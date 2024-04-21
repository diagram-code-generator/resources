package drawioxml

import "encoding/xml"

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	G       G        `xml:"g"`
}

type G struct {
	ID    string `xml:"id,attr"`
	Class string `xml:"class,attr"`
	Nodes []Node `xml:"g"`
}

type Node struct {
	Title string `xml:"title"`
	Text  Text   `xml:"text"`
}

type Text struct {
	Content string `xml:",chardata"`
	X       string `xml:"x,attr"`
	Y       string `xml:"y,attr"`
}
