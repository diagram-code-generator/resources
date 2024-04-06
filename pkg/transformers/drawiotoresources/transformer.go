package drawiotoresources

import (
	"errors"

	drawioxml "github.com/joselitofilho/drawio-parser-go/pkg/parser/xml"

	"github.com/diagram-code-generator/resources/pkg/resources"
)

var ErrInvalidXML = errors.New("invalid XML error")

type Transformer struct {
	mxFile  *drawioxml.MxFile
	factory resources.ResourceFactory
}

func NewTransformer(mxFile *drawioxml.MxFile, factory resources.ResourceFactory) *Transformer {
	return &Transformer{mxFile: mxFile, factory: factory}
}

// Transform parses resources from the MxFile.
func (t *Transformer) Transform() (*resources.ResourceCollection, error) {
	if t.mxFile == nil {
		return nil, ErrInvalidXML
	}

	resc := resources.NewResourceCollection()

	for i := range t.mxFile.Diagram.MxGraphModel.Root.MxCells {
		cell := t.mxFile.Diagram.MxGraphModel.Root.MxCells[i]

		resource := t.factory.CreateResource(cell.ID, cell.Value, cell.Style)
		if resource != nil {
			resc.AddResource(resource)
		}
	}

	resourcesMap := map[string]resources.Resource{}
	for _, resource := range resc.Resources {
		resourcesMap[resource.ID()] = resource
	}

	for i := range t.mxFile.Diagram.MxGraphModel.Root.MxCells {
		cell := t.mxFile.Diagram.MxGraphModel.Root.MxCells[i]
		if cell.Source != "" && cell.Target != "" {
			source := resourcesMap[cell.Source]
			target := resourcesMap[cell.Target]

			if source != nil && target != nil {
				resc.AddRelationship(source, target)
			}
		}
	}

	return resc, nil
}
