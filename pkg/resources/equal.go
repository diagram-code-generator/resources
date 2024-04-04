package resources

import "fmt"

// Equal checks if two ResourceCollection instances are equal.
func (rc *ResourceCollection) Equal(other any) bool {
	switch o := other.(type) {
	case ResourceCollection:
		return rc.equal(&o)
	case *ResourceCollection:
		return rc.equal(o)
	}

	return false
}

func (rc *ResourceCollection) equal(other *ResourceCollection) bool {
	if ok := rc.resourcesEqualValues(other.Resources); ok {
		if ok := rc.relationshipsEqualValues(other.Relationships); ok {
			return true
		}
	}

	return false
}

func (rc *ResourceCollection) resourcesEqualValues(otherResources []Resource) bool {
	if len(rc.Resources) != len(otherResources) {
		return false
	}

	rcIDs := make(map[string]struct{})
	for _, resource := range rc.Resources {
		rcIDs[resource.ID()] = struct{}{}
	}

	for _, resource := range otherResources {
		if _, ok := rcIDs[resource.ID()]; !ok {
			return false
		}
	}

	return true
}

func (rc *ResourceCollection) relationshipsEqualValues(otherRelationships []Relationship) bool {
	if len(rc.Relationships) != len(otherRelationships) {
		return false
	}

	rcRelationships := make(map[string]struct{})
	for _, rel := range rc.Relationships {
		key := fmt.Sprintf("%s_%s", rel.Source.ID(), rel.Target.ID())
		rcRelationships[key] = struct{}{}
	}

	for _, rel := range otherRelationships {
		key := fmt.Sprintf("%s_%s", rel.Source.ID(), rel.Target.ID())
		if _, ok := rcRelationships[key]; !ok {
			return false
		}
	}

	return true
}
