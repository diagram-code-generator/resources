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
	type args struct {
		resc   *resources.ResourceCollection
		config *Config
	}

	lambdaResource := resources.NewGenericResource("1", "MyLambda", "lambda")
	sqsResource := resources.NewGenericResource("2", "my-queue", "sqs")

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
		name string
		args args
		want string
	}{
		{
			name: "happy path",
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: sqsResource,
					}},
				},
				config: &Config{ResourceImageMap: reourceImageMap},
			},
			want: string(happyPath),
		},
		{
			name: "custom node and edge attrs",
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: sqsResource,
					}},
				},
				config: &Config{NodeAttrs: nodeAttrs, EdgeAttrs: edgeAttrs, ResourceImageMap: reourceImageMap},
			},
			want: string(customNodeEdgeAttrs),
		},
		{
			name: "left-right orientation",
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: sqsResource,
					}},
				},
				config: &Config{Orientation: "LR", ResourceImageMap: reourceImageMap},
			},
			want: string(lrOrientation),
		},
		{
			name: "empty graph",
			args: args{
				resc:   resources.NewResourceCollection(),
				config: &Config{ResourceImageMap: map[string]string{}},
			},
			want: "digraph  {\n\t\n\t\n}\n",
		},
		{
			name: "skip edge for nil sources",
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: nil,
						Target: sqsResource,
					}},
				},
				config: &Config{ResourceImageMap: reourceImageMap},
			},
			want: string(sourceOrTargetNil),
		},
		{
			name: "skip edge for nil targets",
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: nil,
					}},
				},
				config: &Config{ResourceImageMap: reourceImageMap},
			},
			want: string(sourceOrTargetNil),
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := Build(tc.args.resc, tc.args.config)

			require.Equal(t, tc.want, got)
		})
	}
}

func TestBuildWithStyle(t *testing.T) {
	type args struct {
		resc   *resources.ResourceCollection
		config *Config
	}

	lambdaResource := resources.NewGenericResource("1", "MyLambda", "lambda")
	sqsResource := resources.NewGenericResource("2", "my-queue", "sqs")
	kinesisResource := resources.NewGenericResource("3", "MyStream", "kinesis")
	databaseResource := resources.NewGenericResource("4", "doc", "database")

	reourceImageMap := map[string]string{
		"lambda": "images/lambda.svg",
		"sqs":    "images/sqs.svg",
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy path",
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
			want: string(withStyleHappyPath),
		},
		{
			name: "ignore duplicated relationship",
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
			want: string(withStyleHappyPath),
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := Build(tc.args.resc, tc.args.config)

			require.Equal(t, tc.want, got)
		})
	}
}
