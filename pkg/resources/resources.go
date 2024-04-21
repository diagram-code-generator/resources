package resources

// Resource defines the basic contract for all resource types, specifying methods to retrieve the ID, value, and
// resource type of a resource.
type Resource interface {
	ID() string
	Value() string
	ResourceType() string
}

// GenericResource represents a generic implementation of a resource, providing methods to retrieve the ID, value, and
// resource type.
type GenericResource struct {
	id           string
	value        string
	resourceType string
}

// NewGenericResource creates a new instance of GenericResource.
func NewGenericResource(id, value string, resourceType string) *GenericResource {
	return &GenericResource{id: id, value: value, resourceType: resourceType}
}
func (r *GenericResource) ID() string           { return r.id }
func (r *GenericResource) Value() string        { return r.value }
func (r *GenericResource) ResourceType() string { return r.resourceType }

// Relationship represents the relationship between two resources. It consists of a source and a target resource.
type Relationship struct {
	Source Resource
	Target Resource
}

// ResourceCollection represents a collection of resources and their relationships. It includes slices to store
// resources and relationships.
type ResourceCollection struct {
	Resources     []Resource
	Relationships []Relationship
}

// NewResourceCollection creates a new ResourceCollection.
func NewResourceCollection() *ResourceCollection {
	return &ResourceCollection{}
}

// AddResource adds a resource to the collection.
func (rc *ResourceCollection) AddResource(resource Resource) {
	rc.Resources = append(rc.Resources, resource)
}

// AddRelationship adds a relationship to the collection.
func (rc *ResourceCollection) AddRelationship(source, target Resource) {
	relationship := Relationship{Source: source, Target: target}
	rc.Relationships = append(rc.Relationships, relationship)
}
