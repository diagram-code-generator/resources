package resourcestodrawio

import (
	pdrawioxml "github.com/joselitofilho/drawio-parser-go/pkg/parser/xml"

	gdrawioxml "github.com/diagram-code-generator/resources/pkg/parser/graphviz/drawioxml"
	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/goccy/go-graphviz"
)

type Transformer struct {
	config        *gdrawioxml.Config
	resCollection *resources.ResourceCollection
	g             *graphviz.Graphviz
}

func NewTransformer(resCollection *resources.ResourceCollection, config *gdrawioxml.Config, g *graphviz.Graphviz) *Transformer {
	return &Transformer{
		config:        config,
		resCollection: resCollection,
		g:             g,
	}
}

func (t *Transformer) Transform() *pdrawioxml.MxFile {
	return gdrawioxml.Build(t.resCollection, t.config, t.g)
}
