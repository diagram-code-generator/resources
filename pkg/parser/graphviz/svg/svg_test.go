package drawioxml

import (
	"encoding/xml"
	"testing"

	"github.com/diagram-code-generator/resources/pkg/resources"
	"github.com/stretchr/testify/require"
)

func TestSVGDiagram_Build(t *testing.T) {
	type args struct {
		resCollection *resources.ResourceCollection
	}

	lambda1 := resources.NewGenericResource("1", "lambda1", "lambda")
	lambda2 := resources.NewGenericResource("2", "lambda2", "lambda")

	tests := []struct {
		name string
		args args
		want *SVG
	}{
		{
			name: "happy path",
			args: args{
				resCollection: &resources.ResourceCollection{
					Resources: []resources.Resource{lambda1, lambda2},
					Relationships: []resources.Relationship{
						{Source: lambda1, Target: lambda2},
					},
				},
			},
			want: &SVG{
				XMLName: xml.Name{
					Space: "http://www.w3.org/2000/svg",
					Local: "svg",
				},
				G: G{
					ID:    "graph0",
					Class: "graph",
					Nodes: []Node{
						{
							Title: "1$$lambda1$$lambda",
							Text:  Text{Content: "1$$lambda1$$lambda", X: "96.8597", Y: "-85.8"},
						},
						{
							Title: "2$$lambda2$$lambda",
							Text:  Text{Content: "2$$lambda2$$lambda", X: "96.8597", Y: "-13.8"},
						},
						{
							Title: "1$$lambda1$$lambda->2$$lambda2$$lambda",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewSVGDiagram()

			got := d.Build(tt.args.resCollection)

			require.Equal(t, tt.want, got)
		})
	}
}
