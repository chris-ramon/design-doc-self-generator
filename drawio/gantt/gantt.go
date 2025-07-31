package gantt

import "encoding/xml"

// MxFile represents the root element of a DrawIO file
type MxFile struct {
	XMLName  xml.Name  `xml:"mxfile"`
	Host     string    `xml:"host,attr"`
	Agent    string    `xml:"agent,attr"`
	Version  string    `xml:"version,attr"`
	Diagrams []Diagram `xml:"diagram"`
}

// Diagram represents a diagram within the DrawIO file
type Diagram struct {
	XMLName     xml.Name     `xml:"diagram"`
	Name        string       `xml:"name,attr"`
	ID          string       `xml:"id,attr"`
	MxGraphModel MxGraphModel `xml:"mxGraphModel"`
}

// MxGraphModel represents the graph model containing all the visual elements
type MxGraphModel struct {
	XMLName    xml.Name `xml:"mxGraphModel"`
	Dx         string   `xml:"dx,attr"`
	Dy         string   `xml:"dy,attr"`
	Grid       string   `xml:"grid,attr"`
	GridSize   string   `xml:"gridSize,attr"`
	Guides     string   `xml:"guides,attr"`
	Tooltips   string   `xml:"tooltips,attr"`
	Connect    string   `xml:"connect,attr"`
	Arrows     string   `xml:"arrows,attr"`
	Fold       string   `xml:"fold,attr"`
	Page       string   `xml:"page,attr"`
	PageScale  string   `xml:"pageScale,attr"`
	PageWidth  string   `xml:"pageWidth,attr"`
	PageHeight string   `xml:"pageHeight,attr"`
	Background string   `xml:"background,attr"`
	Math       string   `xml:"math,attr"`
	Shadow     string   `xml:"shadow,attr"`
	Root       Root     `xml:"root"`
}

// Root represents the root container for all cells
type Root struct {
	XMLName xml.Name `xml:"root"`
	Cells   []MxCell `xml:"mxCell"`
}

// MxCell represents a cell in the DrawIO diagram
type MxCell struct {
	XMLName    xml.Name    `xml:"mxCell"`
	ID         string      `xml:"id,attr"`
	Value      string      `xml:"value,attr,omitempty"`
	Style      string      `xml:"style,attr,omitempty"`
	Parent     string      `xml:"parent,attr,omitempty"`
	Vertex     string      `xml:"vertex,attr,omitempty"`
	MxGeometry *MxGeometry `xml:"mxGeometry,omitempty"`
}

// MxGeometry represents the geometry (position and size) of a cell
type MxGeometry struct {
	XMLName xml.Name `xml:"mxGeometry"`
	X       string   `xml:"x,attr,omitempty"`
	Y       string   `xml:"y,attr,omitempty"`
	Width   string   `xml:"width,attr,omitempty"`
	Height  string   `xml:"height,attr,omitempty"`
	As      string   `xml:"as,attr,omitempty"`
}
