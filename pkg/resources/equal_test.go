package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceCollection_Equal(t *testing.T) {
	type fields struct {
		Resources     []Resource
		Relationships []Relationship
	}

	type args struct {
		other any
	}

	res1 := &GenericResource{id: "1"}
	res2 := &GenericResource{id: "2"}
	res3 := &GenericResource{id: "3"}

	resources1 := []Resource{res1, res2, res3}
	resources2 := []Resource{res3, res1, res2}

	relationships1 := []Relationship{
		{Source: res1, Target: res2},
		{Source: res2, Target: res3},
	}
	relationships2 := []Relationship{
		{Source: res2, Target: res3},
		{Source: res1, Target: res2},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "resource items in the same order",
			fields: fields{
				Resources:     resources1,
				Relationships: relationships1,
			},
			args: args{
				other: &ResourceCollection{
					Resources:     resources1,
					Relationships: relationships2,
				},
			},
			want: true,
		},
		{
			name: "relationship items in the same order",
			fields: fields{
				Resources:     resources1,
				Relationships: relationships1,
			},
			args: args{
				other: &ResourceCollection{
					Resources:     resources2,
					Relationships: relationships1,
				},
			},
			want: true,
		},
		{
			name: "resource and relationship items in the different order",
			fields: fields{
				Resources:     resources1,
				Relationships: relationships1,
			},
			args: args{
				other: &ResourceCollection{
					Resources:     resources2,
					Relationships: relationships1,
				},
			},
			want: true,
		},
		{
			name: "ResourceCollection as reference",
			fields: fields{
				Resources:     resources1,
				Relationships: relationships1,
			},
			args: args{
				other: ResourceCollection{
					Resources:     resources2,
					Relationships: relationships2,
				},
			},
			want: true,
		},
		{
			name: "resource items with same size and different values",
			fields: fields{
				Resources:     resources1,
				Relationships: relationships1,
			},
			args: args{
				other: &ResourceCollection{
					Resources:     []Resource{res3, res1, &GenericResource{id: "0"}},
					Relationships: relationships2,
				},
			},
			want: false,
		},
		{
			name: "relationship items with same size and different values",
			fields: fields{
				Resources:     resources1,
				Relationships: relationships1,
			},
			args: args{
				other: &ResourceCollection{
					Resources: resources2,
					Relationships: []Relationship{
						{Source: res2, Target: res1},
						{Source: res3, Target: res2},
					},
				},
			},
			want: false,
		},
		{
			name: "resource items are different",
			fields: fields{
				Resources:     resources1,
				Relationships: relationships1,
			},
			args: args{
				other: &ResourceCollection{
					Resources:     []Resource{},
					Relationships: relationships2,
				},
			},
			want: false,
		},
		{
			name: "relationship items are different",
			fields: fields{
				Resources:     resources1,
				Relationships: relationships1,
			},
			args: args{
				other: &ResourceCollection{
					Resources:     resources2,
					Relationships: []Relationship{},
				},
			},
			want: false,
		},
		{
			name: "comparing it to an instance that is not a ResourceCollection",
			fields: fields{
				Resources:     resources1,
				Relationships: relationships1,
			},
			args: args{
				other: "string",
			},
			want: false,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			rc := &ResourceCollection{
				Resources:     tc.fields.Resources,
				Relationships: tc.fields.Relationships,
			}

			got := rc.Equal(tc.args.other)

			require.Equal(t, tc.want, got)
		})
	}
}
