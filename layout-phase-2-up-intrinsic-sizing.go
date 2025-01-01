package blitra

// A visitor that calculates the intrinsic dimensions of an element.
// Must be executed bottom-up; The sizing is derived from the children.
//
// If the element is text, the intrinsic dimensions are calculated based on the
// text size without size constraints.
func IntrinsicSizingVisitor(element *Element) {
	if element.Parent == nil {
		return
	}

	if element.SourceText != nil {
		_, info := ApplyWrap(
			V(element.Style.TextWrap),
			V(element.Style.Ellipsis),
			&element.AvailableSize,
			*element.SourceText,
		)

		element.IntrinsicSize = info.Dimensions
		return
	}

	allChildrenWidth := 0
	allChildrenHeight := 0
	axis := V(element.Style.Axis)

	for childElement := range element.ChildrenIter {
		childWidth := childElement.IntrinsicSize.Width
		childHeight := childElement.IntrinsicSize.Height
		if axis == HorizontalAxis {
			allChildrenWidth += childWidth
			if childHeight > allChildrenHeight {
				allChildrenHeight = childHeight
			}
		} else {
			allChildrenHeight += childHeight
			if childWidth > allChildrenWidth {
				allChildrenWidth = childWidth
			}
		}
	}

	width := allChildrenWidth + element.LeftEdge() + element.RightEdge()
	height := allChildrenHeight + element.TopEdge() + element.BottomEdge()

	if element.ChildCount > 1 {
		if axis == HorizontalAxis {
			width += V(element.Style.Gap) * (element.ChildCount - 1)
		} else {
			height += V(element.Style.Gap) * (element.ChildCount - 1)
		}
	}

	if element.Style.Width != nil {
		width = *element.Style.Width
	}
	if element.Style.MinWidth != nil && width < *element.Style.MinWidth {
		width = *element.Style.MinWidth
	}
	if element.Style.Height != nil {
		height = *element.Style.Height
	}
	if element.Style.MinHeight != nil && height < *element.Style.MinHeight {
		height = *element.Style.MinHeight
	}

	element.IntrinsicSize.Width = width
	element.IntrinsicSize.Height = height
}
