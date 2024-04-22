package dot

import (
	_ "embed"
	"testing"

	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/happy_path.dot
	happyPath []byte
	//go:embed testdata/custom_node_edge_attrs.dot
	customNodeEdgeAttrs []byte
	//go:embed testdata/default_config.dot
	defaultConfigDot []byte
	//go:embed testdata/lr_orientation.dot
	lrOrientation []byte
	//go:embed testdata/source_or_target_nil.dot
	sourceOrTargetNil []byte
)

var (
	//go:embed testdata/with_style/happy_path.dot
	withStyleHappyPath []byte
)

func TestBuild(t *testing.T) {
	type fields struct {
		config *Config
	}

	type args struct {
		resc *resources.ResourceCollection
	}

	lambdaResource := resources.NewGenericResource("1", "MyLambda", "lambda")
	sqsResource := resources.NewGenericResource("2", "my-queue", "sqs")
	kinesisResource := resources.NewGenericResource("3", "MyStream", "kinesis")
	databaseResource := resources.NewGenericResource("4", "doc", "database")

	reourceImageMap := map[string]string{
		"lambda": "images/lambda.svg",
		"sqs":    "images/sqs.svg",
	}

	nodeAttrs := make(map[string]any)
	for k, v := range DefaultNodeAttrs {
		nodeAttrs[k] = v
	}
	delete(nodeAttrs, "height")

	edgeAttrs := make(map[string]any)
	for k, v := range DefaultEdgeAttrs {
		edgeAttrs[k] = v
	}
	edgeAttrs["arrowtail"] = "dot"

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "happy path",
			fields: fields{
				config: &Config{Splines: SplineOrtho, ResourceImageMap: reourceImageMap},
			},
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: sqsResource,
					}},
				},
			},
			want: string(happyPath),
		},
		{
			name: "with style",
			fields: fields{
				config: &Config{
					Style: &Style{
						Nodes: map[resources.Resource]string{lambdaResource: "green"},
						Arrows: map[string][]map[string]string{
							"MyLambda": {{"my-queue": "red"}, {"MyStream": "green"}},
							"MyStream": {{"MyLambda": "green"}},
						},
					},
					ResourceImageMap: reourceImageMap,
				},
			},
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource, kinesisResource, databaseResource},
					Relationships: []resources.Relationship{
						{
							Source: lambdaResource,
							Target: sqsResource,
						},
						{
							Source: lambdaResource,
							Target: kinesisResource,
						},
						{
							Source: lambdaResource,
							Target: databaseResource,
						},
					},
				},
			},
			want: string(withStyleHappyPath),
		},
		{
			name: "ignore duplicated relationship by applying style",
			fields: fields{
				config: &Config{
					Style: &Style{
						Nodes: map[resources.Resource]string{lambdaResource: "green"},
						Arrows: map[string][]map[string]string{
							"MyLambda": {{"my-queue": "red"}, {"MyStream": "green"}},
							"MyStream": {{"MyLambda": "green"}},
						},
					},
					ResourceImageMap: reourceImageMap,
				},
			},
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource, kinesisResource, databaseResource},
					Relationships: []resources.Relationship{
						{
							Source: lambdaResource,
							Target: sqsResource,
						},
						{
							Source: lambdaResource,
							Target: sqsResource,
						},
						{
							Source: lambdaResource,
							Target: kinesisResource,
						},
						{
							Source: lambdaResource,
							Target: databaseResource,
						},
					},
				},
			},
			want: string(withStyleHappyPath),
		},
		{
			name: "custom node and edge attrs",
			fields: fields{
				config: &Config{NodeAttrs: nodeAttrs, EdgeAttrs: edgeAttrs, ResourceImageMap: reourceImageMap},
			},
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: sqsResource,
					}},
				},
			},
			want: string(customNodeEdgeAttrs),
		},
		{
			name: "left-right orientation",
			fields: fields{
				config: &Config{Orientation: "LR", ResourceImageMap: reourceImageMap},
			},
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: sqsResource,
					}},
				},
			},
			want: string(lrOrientation),
		},
		{
			name: "empty graph",
			fields: fields{
				config: &Config{ResourceImageMap: map[string]string{}},
			},
			args: args{
				resc: resources.NewResourceCollection(),
			},
			want: "digraph  {\n\t\n\t\n}\n",
		},
		{
			name: "skip edge for nil sources",
			fields: fields{
				config: &Config{ResourceImageMap: reourceImageMap},
			},
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: nil,
						Target: sqsResource,
					}},
				},
			},
			want: string(sourceOrTargetNil),
		},
		{
			name: "skip edge for nil targets",
			fields: fields{
				config: &Config{ResourceImageMap: reourceImageMap},
			},
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: nil,
					}},
				},
			},
			want: string(sourceOrTargetNil),
		},
		{
			name: "default config",
			fields: fields{
				config: nil,
			},
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: sqsResource,
					}},
				},
			},
			want: string(defaultConfigDot),
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := NewDotDiagram(tc.fields.config).Build(tc.args.resc)

			require.Equal(t, tc.want, got)
		})
	}
}
