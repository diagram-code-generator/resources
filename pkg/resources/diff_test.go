package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	availableTypes = []string{kinesisType, lambdaType, sqsType}
)

func TestFindDifferences(t *testing.T) {
	type args struct {
		rc1 *ResourceCollection
		rc2 *ResourceCollection
	}

	lambda1Resource := NewGenericResource("1", "myReceiver", lambdaType)
	sqs1Resource := NewGenericResource("2", "my-queue", sqsType)

	lambda2Resource := NewGenericResource("1", "myProcessor", lambdaType)
	sqs2Resource := NewGenericResource("2", "my-q", sqsType)

	tests := []struct {
		name                       string
		args                       args
		wantAddedResourcesByType   map[string][]Resource
		wantRemovedResourcesByType map[string][]Resource
		wantAddedRelationships     []Relationship
		wantRemovedRelationships   []Relationship
	}{
		{
			name: "happy path",
			args: args{
				rc1: &ResourceCollection{
					Resources: []Resource{lambda1Resource, sqs1Resource},
					Relationships: []Relationship{
						{Source: lambda1Resource, Target: sqs1Resource},
					},
				},
				rc2: &ResourceCollection{
					Resources: []Resource{lambda2Resource, sqs2Resource},
					Relationships: []Relationship{
						{Source: lambda2Resource, Target: sqs2Resource},
					},
				},
			},
			wantAddedResourcesByType: map[string][]Resource{
				lambdaType: {lambda2Resource},
				sqsType:    {sqs2Resource},
			},
			wantRemovedResourcesByType: map[string][]Resource{
				lambdaType: {lambda1Resource},
				sqsType:    {sqs1Resource},
			},
			wantAddedRelationships: []Relationship{
				{Source: lambda2Resource, Target: sqs2Resource},
			},
			wantRemovedRelationships: []Relationship{
				{Source: lambda1Resource, Target: sqs1Resource},
			},
		},
		{
			name: "empty",
			args: args{
				rc1: &ResourceCollection{},
				rc2: &ResourceCollection{},
			},
			wantAddedResourcesByType:   map[string][]Resource{},
			wantRemovedResourcesByType: map[string][]Resource{},
			wantAddedRelationships:     nil,
			wantRemovedRelationships:   nil,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			gotAddedResourcesByType, gotRemovedResourcesByType, gotAddedRelationships, gotRemovedRelationships :=
				FindDifferences(tc.args.rc1, tc.args.rc2)

			require.Equal(t, tc.wantAddedResourcesByType, gotAddedResourcesByType)
			require.Equal(t, tc.wantRemovedResourcesByType, gotRemovedResourcesByType)
			require.Equal(t, tc.wantAddedRelationships, gotAddedRelationships)
			require.Equal(t, tc.wantRemovedRelationships, gotRemovedRelationships)
		})
	}
}

func TestPrintDiff(t *testing.T) {
	type args struct {
		rc1 *ResourceCollection
		rc2 *ResourceCollection
	}

	tests := []struct {
		name           string
		args           args
		logMsgExpected string
	}{
		{
			name: "happy path",
			args: args{
				rc1: &ResourceCollection{
					Resources: []Resource{
						NewGenericResource("1", "myLambda", lambdaType),
						NewGenericResource("2", "my-queue", sqsType),
						NewGenericResource("3", "MyStream", kinesisType),
					},
					Relationships: []Relationship{
						{
							Source: NewGenericResource("1", "myLambda", lambdaType),
							Target: NewGenericResource("2", "my-queue", sqsType),
						},
						{
							Source: NewGenericResource("1", "myLambda", lambdaType),
							Target: NewGenericResource("3", "MyStream", sqsType),
						},
					},
				},
				rc2: &ResourceCollection{
					Resources: []Resource{
						NewGenericResource("1", "myLam", lambdaType),
						NewGenericResource("2", "my-q", sqsType),
						NewGenericResource("3", "MyStream", kinesisType),
					},
					Relationships: []Relationship{
						{
							Target: NewGenericResource("2", "my-q", sqsType),
							Source: NewGenericResource("1", "myLam", lambdaType),
						},
						{
							Source: NewGenericResource("1", "myLambda", lambdaType),
							Target: NewGenericResource("3", "MyStream", sqsType),
						},
					},
				},
			},
			logMsgExpected: "",
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(_ *testing.T) {
			PrintDiff(tc.args.rc1, tc.args.rc2, availableTypes)
		})
	}
}
