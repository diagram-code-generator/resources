package resourcestodrawio

import (
	"crypto/rand"
	"fmt"
	"strings"

	pdrawioxml "github.com/joselitofilho/drawio-parser-go/pkg/parser/xml"

	svg "github.com/diagram-code-generator/resources/pkg/parser/graphviz/svg"
	"github.com/diagram-code-generator/resources/pkg/resources"
)

// Help tests.
var randRead = rand.Read

type Transformer struct {
	config        *Config
	resCollection *resources.ResourceCollection
}

func NewTransformer(resCollection *resources.ResourceCollection, config *Config) *Transformer {
	return &Transformer{
		config:        config,
		resCollection: resCollection,
	}
}

func (t *Transformer) Transform() *pdrawioxml.MxFile {
	if t.config == nil {
		t.config = &Config{}
	}

	svgData := svg.NewSVGDiagram().Build(t.resCollection)

	return &pdrawioxml.MxFile{
		Diagram: pdrawioxml.Diagram{
			MxGraphModel: pdrawioxml.MxGraphModel{
				Root: pdrawioxml.Root{
					MxCells: t.buildMxCells(t.resCollection, svgData),
				},
			},
		},
	}
}

func (t *Transformer) buildMxCells(resCollection *resources.ResourceCollection, svgData *svg.SVG) []pdrawioxml.MxCell {
	edgeID := len(resCollection.Resources) + 1
	baseID := generateBaseID(20)

	mxCells := make([]pdrawioxml.MxCell, 0, len(resCollection.Resources)+len(resCollection.Relationships)+2)

	mxCells = append(mxCells, pdrawioxml.MxCell{ID: "0"}, pdrawioxml.MxCell{ID: "1", Parent: "0"})

	for _, node := range svgData.G.Nodes {
		content := node.Text.Content
		if content == "" {
			continue
		}

		parts := strings.Split(content, svg.ResourceInfoSeparator)
		id, value, resType := parts[0], parts[1], parts[2]
		x, y := node.Text.X, node.Text.Y

		sourceID := fmt.Sprintf("%s-%s", baseID, id)

		mxCells = append(mxCells, pdrawioxml.MxCell{
			ID:       sourceID,
			Value:    value,
			Style:    t.config.NodeStyles[resType],
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
