package blitra

// A visitor that calculates the available dimensions of an element.
// Must be executed top-down; The sizing is derived from the element and
// applied to the children.
//
// The available dimensions are calculated based on the parent's available
// dimensions within it's padding, margin, border, and gaps.
//
// Inherited styles are also passed down to the children in this phase.
func AvailableSizingVisitor(element *Element) {
	if element.ChildCount == 0 {
		return
	}

	width := element.AvailableSize.Width
	height := element.AvailableSize.Height

	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}

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
		childElement.AvailableSize.Width = innerWidth
		childElement.AvailableSize.Height = innerHeight

		if childElement.Style.TextColor == nil {
			childElement.Style.TextColor = element.Style.TextColor
		}
		if childElement.Style.BackgroundColor == nil {
			childElement.Style.BackgroundColor = element.Style.BackgroundColor
		}
		if childElement.Style.TextWrap == nil {
			childElement.Style.TextWrap = element.Style.TextWrap
		}
		if childElement.Style.Ellipsis == nil {
			childElement.Style.Ellipsis = element.Style.Ellipsis
		}
	}
}
