//go:generate mockgen -package resources -destination resource_factory_mock.go -source resource_factory.go
package resources

type ResourceFactory interface {
	CreateResource(id, value, style string) Resource
}
