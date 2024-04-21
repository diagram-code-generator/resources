package resourcestodrawio

import (
	"crypto/rand"
	"errors"
	"testing"

	drawioxml "github.com/joselitofilho/drawio-parser-go/pkg/parser/xml"

	gdrawioxml "github.com/diagram-code-generator/resources/pkg/parser/graphviz/drawioxml"
	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/goccy/go-graphviz"
	"github.com/stretchr/testify/require"
)

func TestTransformer_Transform(t *testing.T) {
	const lambdaStyle = "outlineConnect=0;dashed=0;verticalLabelPosition=bottom;verticalAlign=top;align=center;html=1;" +
		"shape=mxgraph.aws3.lambda;fillColor=#F58534;gradientColor=none;aspect=fixed;"

	errDummy := errors.New("dummy error")

	type fields struct {
		config        *gdrawioxml.Config
		resCollection *resources.ResourceCollection
		g             *graphviz.Graphviz
	}

	lambda1 := resources.NewGenericResource("1", "lambda1", "lambda")
	lambda2 := resources.NewGenericResource("2", "lambda2", "lambda")
	lambda3 := resources.NewGenericResource("3", "lambda3", "lambda")

	tests := []struct {
		name   string
		fields fields
		setup  func() (tearDown func())
		want   *drawioxml.MxFile
	}{
		{
			name: "happy path",
			fields: fields{
				config: &gdrawioxml.Config{
					NodeStyles: map[string]string{"lambda": lambdaStyle},
				},
				resCollection: &resources.ResourceCollection{
					Resources: []resources.Resource{lambda1, lambda2, lambda3},
					Relationships: []resources.Relationship{
						{Source: lambda1, Target: lambda2},
					},
				},
				g: graphviz.New(),
			},
			setup: func() (tearDown func()) {
				gdrawioxml.RandRead = func(b []byte) (n int, err error) {
					require.Len(t, b, 15)
					return 15, nil
				}

				return func() {
					gdrawioxml.RandRead = rand.Read
				}
			},
			want: &drawioxml.MxFile{Diagram: drawioxml.Diagram{MxGraphModel: drawioxml.MxGraphModel{
				Root: drawioxml.Root{MxCells: []drawioxml.MxCell{
					{ID: "0"},
					{ID: "1", Parent: "0"},
					{
						ID: "aaaaaaaaaaaaaaa-1", Value: "lambda1", Style: lambdaStyle, Parent: "1", Vertex: "1",
						Geometry: &drawioxml.Geometry{
							X: "96.8597", Y: "-85.8", Width: 40, Height: 40, As: "geometry",
						},
					},
					{
						ID: "aaaaaaaaaaaaaaa-2", Value: "lambda2", Style: lambdaStyle, Parent: "1", Vertex: "1",
						Geometry: &drawioxml.Geometry{
							X: "96.8597", Y: "-13.8", Width: 40, Height: 40, As: "geometry",
						},
					},
					{
						ID: "aaaaaaaaaaaaaaa-3", Value: "lambda3", Style: lambdaStyle, Parent: "1", Vertex: "1",
						Geometry: &drawioxml.Geometry{
							X: "308.8597", Y: "-85.8", Width: 40, Height: 40, As: "geometry",
						},
					},
					{
						ID: "aaaaaaaaaaaaaaa-4", Parent: "1", Edge: "1",
						Source: "aaaaaaaaaaaaaaa-1", Target: "aaaaaaaaaaaaaaa-2",
						Geometry: &drawioxml.Geometry{As: "geometry"},
					},
				}},
			}}},
		},
		{
			name: "when randRead fails should return an empty base ID",
			fields: fields{
				config: &gdrawioxml.Config{},
				resCollection: &resources.ResourceCollection{
					Resources: []resources.Resource{lambda1, lambda2, lambda3},
					Relationships: []resources.Relationship{
						{Source: lambda1, Target: lambda2},
					},
				},
				g: graphviz.New(),
			},
			setup: func() (tearDown func()) {
				gdrawioxml.RandRead = func(b []byte) (n int, err error) {
					require.Len(t, b, 15)
					return 15, errDummy
				}

				return func() {
					gdrawioxml.RandRead = rand.Read
				}
			},
			want: &drawioxml.MxFile{Diagram: drawioxml.Diagram{MxGraphModel: drawioxml.MxGraphModel{
				Root: drawioxml.Root{MxCells: []drawioxml.MxCell{
					{ID: "0"},
					{ID: "1", Parent: "0"},
					{
						ID: "-1", Value: "lambda1", Parent: "1", Vertex: "1",
						Geometry: &drawioxml.Geometry{
							X: "96.8597", Y: "-85.8", Width: 40, Height: 40, As: "geometry",
						},
					},
					{
						ID: "-2", Value: "lambda2", Parent: "1", Vertex: "1",
						Geometry: &drawioxml.Geometry{
							X: "96.8597", Y: "-13.8", Width: 40, Height: 40, As: "geometry",
						},
					},
					{
						ID: "-3", Value: "lambda3", Parent: "1", Vertex: "1",
						Geometry: &drawioxml.Geometry{
							X: "308.8597", Y: "-85.8", Width: 40, Height: 40, As: "geometry",
						},
					},
					{
						ID: "-4", Parent: "1", Edge: "1",
						Source: "-1", Target: "-2",
						Geometry: &drawioxml.Geometry{As: "geometry"},
					},
				}},
			}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tearDown := tt.setup()
			defer tearDown()

			tr := NewTransformer(tt.fields.resCollection, tt.fields.config, tt.fields.g)

			got := tr.Transform()

			require.Equal(t, tt.want, got)
		})
	}
}
