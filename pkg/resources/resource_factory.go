//go:generate mockgen -package mocks -source resource_factory.go -destination mocks/resource_factory_mock.go
package resources

type ResourceFactory interface {
	CreateResource(id, value, style string) Resource
}
