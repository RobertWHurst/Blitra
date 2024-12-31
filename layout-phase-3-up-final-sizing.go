package blitra

// A visitor that calculates the final dimensions of an element. It uses the
// intrinsic and available sizing calculated by the previous two visitors.
// Must be executed bottom-up; The sizing is derived from the children as well
// as previous visitor phases.
//
// If the element is text, the final dimensions are calculated based on the
// text size after being wrapped to fit the available width.
func FinalSizingVisitor(element *Element) {

	// If this element has text, calculate its final dimensions based on the
	// text size after being wrapped to fit the available width.
	if element.SourceText != nil {
		text, info := ApplyWrap(
			V(element.Style.TextWrap),
			V(element.Style.Ellipsis),
			&element.AvailableSize,
			*element.SourceText,
		)
		element.Text = text
		element.Size = info.Dimensions
		return
	}

	// Grab the element axis so we can use it when calculating the children.
	// sizes.
	axis := V(element.Style.Axis)
	alignment := V(element.Style.Align)

	// Calculate the total size of all children that do not grow, and count the
	// ones that do. Also calculate the children element's final size.
	allChildrenWidth := 0
	allChildrenHeight := 0
	growCount := 0
	for childElement := range element.ChildrenIter {

		// Grab the child's intrinsic dimensions and grow property.
		childGrows := V(childElement.Style.Grow)
		childElement.Size = childElement.IntrinsicSize

		// If the child is set to stretch, set it's size along the axis to the
		// intrinsic size of the parent.
		if alignment == StretchAlign {
			if axis == HorizontalAxis {
				childElement.Size.Height = element.IntrinsicSize.Height - element.GetEdgeHeight() - element.GetGapHeight()
			} else {
				childElement.Size.Width = element.IntrinsicSize.Width - element.GetEdgeWidth() - element.GetGapWidth()
			}
		}

		// If the child grows, increment the grow count.
		if childGrows {
			growCount += 1
		}

		// If the child does not grow, get it's size along the axis. and add it
		// to the total size along the axis for all children.
		// Regardless of if the child grows, take it's lateral size if larger
		// than the current total lateral size, replace it.
		if axis == HorizontalAxis {
			if !childGrows {
				allChildrenWidth += childElement.Size.Width
			}
			if childElement.Size.Height > allChildrenHeight {
				allChildrenHeight = childElement.Size.Height
			}
		} else {
			if !childGrows {
				allChildrenHeight += childElement.Size.Height
			}
			if childElement.Size.Width > allChildrenWidth {
				allChildrenWidth = childElement.Size.Width
			}
		}
	}

	// Calculate remainder space for children that grow.
	var growWidth int
	var growHeight int
	var clipWidth int
	var clipHeight int
	if axis == HorizontalAxis {
		remainderWidth := element.AvailableSize.Width - allChildrenWidth
		if remainderWidth < 0 {
			clipWidth = -(remainderWidth / element.ChildCount)
		} else if growCount > 0 {
			growWidth = remainderWidth / growCount
		}
	} else {
		remainderHeight := element.AvailableSize.Height - allChildrenHeight
		if remainderHeight < 0 {
			clipHeight = -(remainderHeight / element.ChildCount)
		} else if growCount > 0 {
			growHeight = remainderHeight / growCount
		}
	}

	for childElement := range element.ChildrenIter {
		childGrows := V(childElement.Style.Grow)

		if childGrows {
			childElement.Size.Width += growWidth
			childElement.Size.Height += growHeight
		}
		childElement.Size.Width -= clipWidth
		childElement.Size.Height -= clipHeight

		if childElement.Size.Width < 0 {
			childElement.Size.Width = 0
		}
		if childElement.Size.Height < 0 {
			childElement.Size.Height = 0
		}
		if childElement.Size.Width > childElement.AvailableSize.Width {
			childElement.Size.Width = childElement.AvailableSize.Width
		}
		if childElement.Size.Height > childElement.AvailableSize.Height {
			childElement.Size.Height = childElement.AvailableSize.Height
		}
	}

	if element.Parent == nil {
		element.Size.Width = element.AvailableSize.Width
		element.Size.Height = element.AvailableSize.Height
	}
}
