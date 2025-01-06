package blitra

// Calculates the final size for an element's children based on the intrinsic
// and available of the element. If the element is text, then it will re-wrap
// the text to fit the available width.
func FinalSizingVisitor(element *Element, state *LayoutState) error {

	// If this element has text, calculate its final dimensions based on the
	// text size after being wrapped to fit the available width.
	if element.Kind == TextElementKind {
		text, info, err := ApplyWrap(
			V(element.Style.TextWrap),
			V(element.Style.Ellipsis),
			Size{Width: element.Size.Width, Height: element.AvailableSize.Height},
			element.SourceText,
		)
		if err != nil {
			return err
		}

		if info.Size.Width > element.Size.Width || info.Size.Height > element.Size.Height {
			element.ReflowSize = &info.Size
			state.needsReflow = true
			state.elementsToReflow = append(state.elementsToReflow, element)
		}
		element.Text = text

		return nil
	}

	// Grab the element axis so we can use it when calculating the children.
	// sizes.
	axis := V(element.Style.Axis)
	alignment := V(element.Style.Align)
	gap := V(element.Style.Gap)

	edgeWidth := element.LeftEdge() + element.RightEdge()
	edgeHeight := element.TopEdge() + element.BottomEdge()
	innerIntrinsicWidth := element.IntrinsicSize.Width - edgeWidth
	innerIntrinsicHeight := element.IntrinsicSize.Height - edgeHeight
	innerAvailableWidth := element.AvailableSize.Width - edgeWidth
	innerAvailableHeight := element.AvailableSize.Height - edgeHeight

	// Calculate the total size of all children that do not grow, and count the
	// ones that do. Also calculate the children element's final size.
	allChildrenWidth := -gap
	allChildrenHeight := -gap
	growCount := 0
	for childElement := range element.ChildrenIter {

		// Grab the child's intrinsic dimensions and grow property.
		childGrows := V(childElement.Style.Grow)
		childElement.Size = childElement.IntrinsicSize

		// If the child is set to stretch, set it's size along the axis to the
		// intrinsic size of the parent.
		if alignment == StretchAlign {
			if axis == HorizontalAxis {
				childElement.Size.Height = max(innerIntrinsicHeight, 0)
			} else {
				childElement.Size.Width = max(innerIntrinsicWidth, 0)
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
			allChildrenWidth += gap
			allChildrenWidth += childElement.Size.Width + gap
			if childElement.Size.Height > allChildrenHeight {
				allChildrenHeight = childElement.Size.Height
			}
		} else {
			allChildrenHeight += childElement.Size.Height + gap
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
		remainderWidth := innerAvailableWidth - allChildrenWidth
		if remainderWidth < 0 {
			clipWidth = -(remainderWidth / element.ChildCount)
		} else if growCount > 0 {
			growWidth = remainderWidth / growCount
		}
	} else {
		remainderHeight := innerAvailableHeight - allChildrenHeight
		if remainderHeight < 0 {
			clipHeight = -(remainderHeight / element.ChildCount)
		} else if growCount > 0 {
			growHeight = remainderHeight / growCount
		}
	}
	for childElement := range element.ChildrenIter {
		childGrows := V(childElement.Style.Grow)

		if childGrows {
			if axis == HorizontalAxis {
				childElement.Size.Width += growWidth
			} else {
				childElement.Size.Height += growHeight
			}
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

	return nil
}
