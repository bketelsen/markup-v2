package markup

// Componer is the interface that describes a component.
type Componer interface {
	// Render should return a string describing the component with HTML5
	// standard.
	// It support Golang template/text API.
	// Pipeline is based on the component struct.
	// See https://golang.org/pkg/text/template for more informations.
	Render() string
}
