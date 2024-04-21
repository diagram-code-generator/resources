<div align="center">

# resources

[![GitHub tag](https://img.shields.io/github/release/diagram-code-generator/resources?include_prereleases=&sort=semver&color=2ea44f&style=for-the-badge)](https://github.com/diagram-code-generator/resources/releases/)
[![Code coverage](https://img.shields.io/badge/Coverage-100.0%25-2ea44f?style=for-the-badge)](#)
[![Go Report Card](https://goreportcard.com/badge/github.com/diagram-code-generator/resources?style=for-the-badge)](https://goreportcard.com/report/github.com/diagram-code-generator/resources)
[![GoDoc](https://img.shields.io/badge/reference-blue?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/diagram-code-generator/resources "Go to Golang docs")

[![Made with Golang](https://img.shields.io/badge/Golang-1.21.6-blue?logo=go&logoColor=white&style=for-the-badge)](https://go.dev "Go to Golang homepage")
[![Using Diagrams](https://img.shields.io/badge/diagrams.net-orange?logo=&logoColor=white&style=for-the-badge)](https://app.diagrams.net/ "Go to Diagrams homepage")

[![BuyMeACoffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black)](https://www.buymeacoffee.com/joselitofilho)

</div>

# Overview

The `resources` repository is a collection designed to define a structure for resources in a system, providing 
interfaces, generic implementations, and helper methods for resource and resource collection management.
It gives you different ways to handle resources and collections, with handy shortcuts for managing relationships and 
spotting differences between sets.

## How to Use

```bash
$ go get github.com/diagram-code-generator/resources@latest
```

## Core Components

### Resource Interface
The `Resource` interface defines the basic contract for a resource in the system. Any resource must implement this 
interface to be considered valid within the context of this framework.

### GenericResource Generic Implementation
`GenericResource` is a generic implementation that can be extended to create specific resources. It provides a basic 
implementation of common methods required for resource manipulation.

### ResourceCollection
`ResourceCollection` is a structure for storing and managing collections of resources. It offers methods for adding, 
removing, and manipulating resources within the collection.

## Example Usage

```Go
package main

import (
	"fmt"

	"github.com/diagram-code-generator/resources/pkg/resources"
)

func main() {
	// Create some resources
	resource1 := resources.NewGenericResource("1", "Value 1", "Type A")
	resource2 := resources.NewGenericResource("2", "Value 2", "Type B")
	resource3 := resources.NewGenericResource("3", "Value 3", "Type A")

	// Create a resource collection
	collection1 := resources.NewResourceCollection()

	// Add resources to the collection
	collection1.AddResource(resource1)
	collection1.AddResource(resource2)
	collection1.AddResource(resource3)

	// Establish relationships between resources
	collection1.AddRelationship(resource1, resource2)
	collection1.AddRelationship(resource2, resource3)

	// Print information about the resources and relationships
	fmt.Println("Resources:")
	for _, resource := range collection1.Resources {
		fmt.Printf("ID: %s, Value: %s, Type: %s\n", resource.ID(), resource.Value(), resource.ResourceType())
	}

	fmt.Println("\nRelationships:")
	for _, relationship := range collection1.Relationships {
		fmt.Printf("Source: %s, Target: %s\n", relationship.Source.ID(), relationship.Target.ID())
	}

    // Create a resource collection
	collection2 := resources.NewResourceCollection()

    // Print differences between the collections
	availableTypes := []resources.ResourceType{"Type A", "Type B"}
	resources.PrintDiff(collection1, collection2, availableTypes)
}

```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to create an 
[issue][issues] or submit a pull request. Your contribution is much appreciated. See [Contributing](CONTRIBUTING.md).

[![open - Contributing](https://img.shields.io/badge/open-contributing-blue?style=for-the-badge)](CONTRIBUTING.md "Go to contributing")

## License

This project is licensed under the [MIT License](LICENSE).

[issues]: https://github.com/diagram-code-generator/resources/issues