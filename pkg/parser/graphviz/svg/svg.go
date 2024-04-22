package drawioxml

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

const ResourceInfoSeparator = "$$"

type SVGDiagram struct {
	g *graphviz.Graphviz
}

func NewSVGDiagram() *SVGDiagram {
	return &SVGDiagram{g: graphviz.New()}
}

func (d *SVGDiagram) Build(resCollection *resources.ResourceCollection) *SVG {
	graph, _ := d.g.Graph(graphviz.Directed)

	nodesByResourceID := map[string]*cgraph.Node{}
	for _, res := range resCollection.Resources {
		node, _ := graph.CreateNode(fmt.Sprintf("%s%s%s%s%s",
			res.ID(), ResourceInfoSeparator, res.Value(), ResourceInfoSeparator, res.ResourceType()))
		nodesByResourceID[res.ID()] = node
	}

	for _, rel := range resCollection.Relationships {
		_, _ = graph.CreateEdge("", nodesByResourceID[rel.Source.ID()], nodesByResourceID[rel.Target.ID()])
	}

	var buf bytes.Buffer
	_ = d.g.Render(graph, "svg", &buf)

	var svgData SVG
	_ = xml.Unmarshal(buf.Bytes(), &svgData)

	return &svgData
}
