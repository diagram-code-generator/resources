package dot

import "github.com/diagram-code-generator/resources/pkg/resources"

type Style struct {
	Nodes  map[resources.Resource]string
	Arrows map[string][]map[string]string
}

type Config struct {
	Direction        DiagramDirection
	Splines          DiagramSpline
	NodeAttrs        map[string]any
	EdgeAttrs        map[string]any
	ResourceImageMap map[string]string
	Style            *Style
}
