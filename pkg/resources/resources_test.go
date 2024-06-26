package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	kinesisType = "kinesis"
	lambdaType  = "lambda"
	sqsType     = "sqs"
)

func TestResourceCollection_AddResource(t *testing.T) {
	type fields struct {
		Resources     []Resource
		Relationships []Relationship
	}

	type args struct {
		resource Resource
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Resource
	}{
		{
			name: "add a Lamda resource",
			fields: fields{
				Resources:     []Resource{},
				Relationships: []Relationship{},
			},
			args: args{
				resource: NewGenericResource("1", "MyLambda", lambdaType),
			},
			want: []Resource{&GenericResource{
				id: "1", value: "MyLambda", resourceType: lambdaType,
			}},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			rc := NewResourceCollection()
			rc.Resources = tc.fields.Resources
			rc.Relationships = tc.fields.Relationships

			rc.AddResource(tc.args.resource)

			require.Equal(t, tc.want, rc.Resources)
		})
	}
}

func TestResourceCollection_AddRelationship(t *testing.T) {
	type fields struct {
		Resources     []Resource
		Relationships []Relationship
	}

	type args struct {
		source Resource
		target Resource
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Relationship
	}{
		{
			name: "add relationship between a Lambda and SQS",
			fields: fields{
				Resources:     []Resource{},
				Relationships: []Relationship{},
			},
			args: args{
				source: NewGenericResource("1", "MyLambda", lambdaType),
				target: NewGenericResource("2", "MyQueue", sqsType),
			},
			want: []Relationship{{
				Source: &GenericResource{id: "1", value: "MyLambda", resourceType: lambdaType},
				Target: &GenericResource{id: "2", value: "MyQueue", resourceType: sqsType},
			}},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			rc := NewResourceCollection()
			rc.Resources = tc.fields.Resources
			rc.Relationships = tc.fields.Relationships

			rc.AddRelationship(tc.args.source, tc.args.target)

			require.Equal(t, tc.want, rc.Relationships)
		})
	}
}
