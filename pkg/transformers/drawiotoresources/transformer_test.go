package drawiotoresources

import (
	"testing"

	"github.com/diagram-code-generator/resources/pkg/resources"
	"github.com/diagram-code-generator/resources/pkg/resources/mocks"
	drawioxml "github.com/joselitofilho/drawio-parser-go/pkg/parser/xml"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransform(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		mxFile  *drawioxml.MxFile
		factory resources.ResourceFactory
	}

	appEngineResource := resources.NewGenericResource("APPENGINE_ID", "appengine", "mx-appengine")
	dataFlowResource := resources.NewGenericResource("DATAFLOW_ID", "dataflow", "mx-dataflow")

	tests := []struct {
		name      string
		args      args
		setup     func(*mocks.MockResourceFactory)
		want      *resources.ResourceCollection
		targetErr error
	}{
		{
			name: "Single Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{{
									ID:    "APPENGINE_ID",
									Value: "appengine",
									Style: "mx-appengine",
								}},
							},
						},
					},
				},
				factory: mocks.NewMockResourceFactory(ctrl),
			},
			setup: func(mrf *mocks.MockResourceFactory) {
				mrf.EXPECT().
					CreateResource("APPENGINE_ID", "appengine", "mx-appengine").
					Return(appEngineResource)
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{appEngineResource},
			},
		},
		{
			name: "Two Connected Resources",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{{
									ID:    "APPENGINE_ID",
									Value: "appengine",
									Style: "mx-appengine",
								}, {
									ID:    "DATAFLOW_ID",
									Value: "dataflow",
									Style: "mx-dataflow",
								}, {
									ID: "3", Source: "APPENGINE_ID", Target: "DATAFLOW_ID",
								}},
							},
						},
					},
				},
				factory: mocks.NewMockResourceFactory(ctrl),
			},
			setup: func(mrf *mocks.MockResourceFactory) {
				gomock.InOrder(
					mrf.EXPECT().
						CreateResource("APPENGINE_ID", "appengine", "mx-appengine").
						Return(appEngineResource),
					mrf.EXPECT().
						CreateResource("DATAFLOW_ID", "dataflow", "mx-dataflow").
						Return(dataFlowResource),
					mrf.EXPECT().
						CreateResource("3", "", "").
						Return(nil),
				)
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{appEngineResource, dataFlowResource},
				Relationships: []resources.Relationship{{Source: appEngineResource, Target: dataFlowResource}},
			},
		},
		{
			name: "when XML is invalid should return an error",
			args: args{
				mxFile:  nil,
				factory: &mocks.MockResourceFactory{},
			},
			want:      nil,
			targetErr: ErrInvalidXML,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(tc.args.factory.(*mocks.MockResourceFactory))
			}

			got, err := NewTransformer(tc.args.mxFile, tc.args.factory).Transform()

			require.ErrorIs(t, err, tc.targetErr)
			require.Equal(t, tc.want, got)
		})
	}
}
