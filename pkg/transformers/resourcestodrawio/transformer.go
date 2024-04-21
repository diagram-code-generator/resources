package resourcestodrawio

import (
	"bytes"
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/diagram-code-generator/resources/pkg/resources"

	drawioxml "github.com/joselitofilho/drawio-parser-go/pkg/parser/xml"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

// Help with tests.
var randRead = rand.Read

type Transformer struct {
	config        *Config
	resCollection *resources.ResourceCollection
	g             *graphviz.Graphviz
}

func NewTransformer(config *Config, resCollection *resources.ResourceCollection, g *graphviz.Graphviz) *Transformer {
	return &Transformer{
		config:        config,
		resCollection: resCollection,
		g:             g,
	}
}

func (t *Transformer) Transform() *drawioxml.MxFile {
	svgData := t.buildGraph()

	mxCells := t.buildMxCells(svgData)

	return &drawioxml.MxFile{
		Diagram: drawioxml.Diagram{
			MxGraphModel: drawioxml.MxGraphModel{
				Root: drawioxml.Root{
					MxCells: mxCells,
				},
			},
		},
	}
}

func (t *Transformer) buildGraph() *SVG {
	graph := t.newGraph()

	nodesByResourceID := map[string]*cgraph.Node{}
	for _, res := range t.resCollection.Resources {
		node, _ := graph.CreateNode(fmt.Sprintf("%s$$%s$$%s", res.ID(), res.Value(), res.ResourceType()))
		nodesByResourceID[res.ID()] = node
	}

	for _, rel := range t.resCollection.Relationships {
		_, _ = graph.CreateEdge("", nodesByResourceID[rel.Source.ID()], nodesByResourceID[rel.Target.ID()])
	}

	var buf bytes.Buffer
	_ = t.g.Render(graph, "svg", &buf)

	var svgData SVG
	_ = xml.Unmarshal(buf.Bytes(), &svgData)

	return &svgData
}

func (t *Transformer) buildMxCells(svgData *SVG) []drawioxml.MxCell {
	edgeID := len(t.resCollection.Resources) + 1
	baseID := generateBaseID(20)

	styles := t.config.Styles
	if styles == nil {
		styles = map[string]string{}
	}

	mxCells := make([]drawioxml.MxCell, 0, len(t.resCollection.Resources)+len(t.resCollection.Relationships)+2)

	mxCells = append(mxCells, drawioxml.MxCell{ID: "0"}, drawioxml.MxCell{ID: "1", Parent: "0"})

	for _, node := range svgData.G.Nodes {
		content := node.Text.Content
		if content == "" {
			continue
		}

		parts := strings.Split(content, "$$")
		id, value, resType := parts[0], parts[1], parts[2]
		x, y := node.Text.X, node.Text.Y

		sourceID := fmt.Sprintf("%s-%s", baseID, id)

		mxCells = append(mxCells, drawioxml.MxCell{
			ID:       sourceID,
			Value:    value,
			Style:    styles[resType],
			Vertex:   "1",
			Parent:   "1",
			Geometry: &drawioxml.Geometry{X: x, Y: y, Width: 40, Height: 40, As: "geometry"},
		})
	}

	for _, rel := range t.resCollection.Relationships {
		sourceID := fmt.Sprintf("%s-%s", baseID, rel.Source.ID())
		targetID := fmt.Sprintf("%s-%s", baseID, rel.Target.ID())

		mxCells = append(mxCells, drawioxml.MxCell{
			ID:       fmt.Sprintf("%s-%d", baseID, edgeID),
			Source:   sourceID,
			Target:   targetID,
			Edge:     "1",
			Parent:   "1",
			Geometry: &drawioxml.Geometry{As: "geometry"},
		})

		edgeID++
	}

	return mxCells
}

func (t *Transformer) newGraph() (graph *cgraph.Graph) {
	graph, _ = t.g.Graph(graphviz.UnDirected)
	return graph
}

func generateBaseID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Calculate the number of bytes needed to represent the string.
	byteLength := length * 3 / 4

	bytes := make([]byte, byteLength)
	_, err := randRead(bytes)
	if err != nil {
		return ""
	}

	// Convert bytes to string using the custom character set.
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes)
}
