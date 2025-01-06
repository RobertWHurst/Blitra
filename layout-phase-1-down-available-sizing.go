package blitra

// Takes the available size of an element and flows it down to its children.
// It subtracts the margins, padding, and border from the available size when
// doing so.
//
// The available size is not distributed to the children, but instead the
// each child is given the full available size. This is important to allow for
// an accurate intrinsic size calculation. In the final sizing phase the
// any over-sizing will be corrected.
//
// Available sizing is also restricted by the element's style width and max
// width, and height and max height.
//
// This visitor only runs once. It is part of the first phase, and moves
// down the tree.
//
// TODO: Introduce absolute positioning.
func AvailableSizingVisitor(element *Element, state *LayoutState) error {
	if element.ChildCount == 0 {
		return nil
	}

	width := element.AvailableSize.Width
	height := element.AvailableSize.Height

	if element.Style.Width != nil {
		width = *element.Style.Width
	}
	if element.Style.MaxWidth != nil && width > *element.Style.MaxWidth {
		width = *element.Style.MaxWidth
	}

	if element.Style.Height != nil {
		height = *element.Style.Height
	}
	if element.Style.MaxHeight != nil && height > *element.Style.MaxHeight {
		height = *element.Style.MaxHeight
	}

	innerWidth := width - element.LeftEdge() - element.RightEdge()
	innerHeight := height - element.TopEdge() - element.BottomEdge()

	for childElement := range element.ChildrenIter {
		childElement.AvailableSize.Width = max(innerWidth, 0)
		childElement.AvailableSize.Height = max(innerHeight, 0)
	}

	return nil
}
