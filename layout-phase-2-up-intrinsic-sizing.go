package blitra

// Calculates the intrinsic size of an element. If the element is text, then it
// calculates the wrapped text size by the available size from the previous
// phase, using it as the intrinsic size. If the element is a container, then
// it calculates the intrinsic size based on the intrinsic size of its children,
// plus the margins, padding, border, and gap.
//
// This visitor could run twice. It will run once during the second phase,
// and possibly again during a reflow (occurs if dynamically sized elements
// such as text require more space than they have after the final sizing phase).
func IntrinsicSizingVisitor(element *Element, state *LayoutState) error {
	if element.Parent == nil {
		return nil
	}

	// If the element has a reflow size, then we use that as the intrinsic size.
	if element.ReflowSize != nil {
		element.IntrinsicSize = *element.ReflowSize
		element.ReflowSize = nil
		return nil
	}

	// If text, then we calculate the size the text would take if given
	// all available space. This becomes the intrinsic size of the text.
	if element.Kind == TextElementKind {
		_, info, err := ApplyWrap(
			V(element.Style.TextWrap),
			V(element.Style.Ellipsis),
			element.AvailableSize,
			element.SourceText,
		)
		if err != nil {
			return err
		}

		element.IntrinsicSize = info.Size
		return nil
	}

	// If a container, then we calculate it's intrinsic size based on it's
	// children. We do so by summing their intrinsic sizes, then adding the
	// margins, padding, border, and gap.

	allChildrenWidth := 0
	allChildrenHeight := 0
	axis := V(element.Style.Axis)
	gap := V(element.Style.Gap)

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
			width += gap * (element.ChildCount - 1)
		} else {
			height += gap * (element.ChildCount - 1)
		}
	}

	element.IntrinsicSize.Width = max(width, 0)
	element.IntrinsicSize.Height = max(height, 0)

	return nil
}
