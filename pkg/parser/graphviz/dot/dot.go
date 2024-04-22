package dot

import (
	"github.com/emicklei/dot"

	"github.com/diagram-code-generator/resources/pkg/resources"
)

const OrientationLeftRight = "LR"

var (
	DefaultNodeAttrs = map[string]any{
		"shape":    "plaintext",
		"imagepos": "tc",
		"labelloc": "b",
		"height":   0.9,
	}

	DefaultEdgeAttrs = map[string]any{
		"arrowhead": "vee",
		"arrowtail": "normal",
	}
)

type DotDiagram struct {
	config *Config
	g      *dot.Graph
}

func NewDotDiagram(config *Config) *DotDiagram {
	return &DotDiagram{
		config: config,
		g:      dot.NewGraph(dot.Directed),
	}
}

func (d *DotDiagram) Build(resc *resources.ResourceCollection) string {
	if d.config == nil {
		d.config = &Config{}
	}
	config := d.config

	if config.Orientation != "" {
		d.g.Attr("rankdir", config.Orientation)
	}

	d.g.NodeInitializer(func(n dot.Node) {
		var nodeAttrs map[string]any = DefaultNodeAttrs

		if len(config.NodeAttrs) > 0 {
			nodeAttrs = config.NodeAttrs
		}

		for name, value := range nodeAttrs {
			n.Attrs(name, value)
		}
	})

	d.g.EdgeInitializer(func(e dot.Edge) {
		var edgeAttrs map[string]any = DefaultEdgeAttrs

		if len(config.EdgeAttrs) > 0 {
			edgeAttrs = config.EdgeAttrs
		}

		for name, value := range edgeAttrs {
			e.Attrs(name, value)
		}
	})

	nodes := map[string]dot.Node{}
	edges := map[string]struct{}{}

	d.applyStyleForNodes(resc, nodes)

	d.applyStyleForArrows(resc, edges, nodes)

	return d.g.String()
}

func (d *DotDiagram) applyStyleForNodes(resc *resources.ResourceCollection, nodes map[string]dot.Node) {
	style := d.config.Style
	if style == nil {
		style = &Style{}
	}

	resourceImageMap := d.config.ResourceImageMap

	for i := range resc.Resources {
		res := resc.Resources[i]

		node := d.g.Node(res.Value()).Attr("image", resourceImageMap[res.ResourceType()])

		if color, ok := style.Nodes[res]; ok {
			node = node.Attr("fontcolor", color)
		}

		nodes[res.Value()] = node
	}

	for k, v := range style.Nodes {
		nodes[k.Value()] = d.g.Node(k.Value()).
			Attr("fontcolor", v).
			Attr("image", resourceImageMap[k.ResourceType()])
	}
}

func (d *DotDiagram) applyStyleForArrows(
	resc *resources.ResourceCollection, edges map[string]struct{}, nodes map[string]dot.Node,
) {
	style := d.config.Style
	if style == nil {
		style = &Style{}
	}

	for _, rel := range resc.Relationships {
		if rel.Source == nil || rel.Target == nil {
			continue
		}

		edgeKey := rel.Source.Value() + "###" + rel.Target.Value()
		if _, ok := edges[edgeKey]; ok {
			continue
		}

		sourceNode := nodes[rel.Source.Value()]
		targetNode := nodes[rel.Target.Value()]

		if color, ok := getArrowColor(style, rel); ok {
			d.g.Edge(sourceNode, targetNode).Attr("color", color)
		} else {
			d.g.Edge(sourceNode, targetNode)
		}

		edges[edgeKey] = struct{}{}
	}

	d.applyCustomArrowStyles(style, edges, nodes)
}

func (d *DotDiagram) applyCustomArrowStyles(style *Style, edges map[string]struct{}, nodes map[string]dot.Node) {
	for source, targets := range style.Arrows {
		for i := range targets {
			for target, color := range targets[i] {
				edgeKey := source + "###" + target

				if _, ok := edges[edgeKey]; !ok {
					d.g.Edge(nodes[source], nodes[target]).Attr("color", color)

					edges[edgeKey] = struct{}{}
				}
			}
		}
	}
}

func getArrowColor(style *Style, rel resources.Relationship) (string, bool) {
	if list, exists := style.Arrows[rel.Source.Value()]; exists {
		for _, colors := range list {
			if color, ok := colors[rel.Target.Value()]; ok {
				return color, true
			}
		}
	}

	return "", false
}
