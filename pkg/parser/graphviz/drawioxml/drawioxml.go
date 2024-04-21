package drawioxml

import (
	"bytes"
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/diagram-code-generator/resources/pkg/resources"
	pdrawioxml "github.com/joselitofilho/drawio-parser-go/pkg/parser/xml"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

// Help with tests.
var RandRead = rand.Read

func Build(resCollection *resources.ResourceCollection, config *Config, g *graphviz.Graphviz) *pdrawioxml.MxFile {
	svgData := buildSVGGraph(resCollection, g)

	mxCells := buildMxCells(resCollection, svgData, config)

	return &pdrawioxml.MxFile{
		Diagram: pdrawioxml.Diagram{
			MxGraphModel: pdrawioxml.MxGraphModel{
				Root: pdrawioxml.Root{
					MxCells: mxCells,
				},
			},
		},
	}
}

func buildSVGGraph(resCollection *resources.ResourceCollection, g *graphviz.Graphviz) *SVG {
	graph, _ := g.Graph(graphviz.Directed)

	nodesByResourceID := map[string]*cgraph.Node{}
	for _, res := range resCollection.Resources {
		node, _ := graph.CreateNode(fmt.Sprintf("%s$$%s$$%s", res.ID(), res.Value(), res.ResourceType()))
		nodesByResourceID[res.ID()] = node
	}

	for _, rel := range resCollection.Relationships {
		_, _ = graph.CreateEdge("", nodesByResourceID[rel.Source.ID()], nodesByResourceID[rel.Target.ID()])
	}

	var buf bytes.Buffer
	_ = g.Render(graph, "svg", &buf)

	var svgData SVG
	_ = xml.Unmarshal(buf.Bytes(), &svgData)

	return &svgData
}

func buildMxCells(resCollection *resources.ResourceCollection, svgData *SVG, config *Config) []pdrawioxml.MxCell {
	edgeID := len(resCollection.Resources) + 1
	baseID := generateBaseID(20)

	styles := config.NodeStyles
	if styles == nil {
		styles = map[string]string{}
	}

	mxCells := make([]pdrawioxml.MxCell, 0, len(resCollection.Resources)+len(resCollection.Relationships)+2)

	mxCells = append(mxCells, pdrawioxml.MxCell{ID: "0"}, pdrawioxml.MxCell{ID: "1", Parent: "0"})

	for _, node := range svgData.G.Nodes {
		content := node.Text.Content
		if content == "" {
			continue
		}

		parts := strings.Split(content, "$$")
		id, value, resType := parts[0], parts[1], parts[2]
		x, y := node.Text.X, node.Text.Y

		sourceID := fmt.Sprintf("%s-%s", baseID, id)

		mxCells = append(mxCells, pdrawioxml.MxCell{
			ID:       sourceID,
			Value:    value,
			Style:    styles[resType],
			Vertex:   "1",
			Parent:   "1",
			Geometry: &pdrawioxml.Geometry{X: x, Y: y, Width: 40, Height: 40, As: "geometry"},
		})
	}

	for _, rel := range resCollection.Relationships {
		sourceID := fmt.Sprintf("%s-%s", baseID, rel.Source.ID())
		targetID := fmt.Sprintf("%s-%s", baseID, rel.Target.ID())

		mxCells = append(mxCells, pdrawioxml.MxCell{
			ID:       fmt.Sprintf("%s-%d", baseID, edgeID),
			Source:   sourceID,
			Target:   targetID,
			Edge:     "1",
			Parent:   "1",
			Geometry: &pdrawioxml.Geometry{As: "geometry"},
		})

		edgeID++
	}

	return mxCells
}

func generateBaseID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Calculate the number of bytes needed to represent the string.
	byteLength := length * 3 / 4

	bytes := make([]byte, byteLength)
	_, err := RandRead(bytes)
	if err != nil {
		return ""
	}

	// Convert bytes to string using the custom character set.
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes)
}
