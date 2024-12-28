package blitra

// Anything that can be rendered to an element.
type Renderable interface {
	// Should return a value that can be converted to an element via the
	// `ElementFrom` function.
	Render(viewState ViewState) any
}
