package blitra

func PositioningVisitor(element *Element) {
	if element.ChildCount == 0 {
		return
	}

	axis := V(element.Style.Axis)
	alignment := V(element.Style.Align)
	justification := V(element.Style.Justify)
	gap := V(element.Style.Gap)
	innerX := element.Position.X + element.LeftEdge()
	innerY := element.Position.Y + element.TopEdge()
	innerWidth := element.Size.Width - element.LeftEdge() - element.RightEdge()
	innerHeight := element.Size.Height - element.TopEdge() - element.BottomEdge()

	// We start with the negative gap as we only want to add the gap between
	// children.
	allChildrenWidth := -gap
	allChildrenHeight := -gap
	for childElement := range element.ChildrenIter {
		childWidth := childElement.Size.Width
		childHeight := childElement.Size.Height

		if axis == HorizontalAxis {
			allChildrenWidth += childWidth + gap
			if childHeight > allChildrenHeight {
				allChildrenHeight = childHeight
			}
		} else {
			allChildrenHeight += childHeight + gap
			if childWidth > allChildrenWidth {
				allChildrenWidth = childWidth
			}
		}
	}

	justificationOffset := 0
	justificationGap := 0
	switch justification {
	case StartJustify:
		// NOTE: we don't need to do anything for start justification.
	case CenterJustify:
		if axis == HorizontalAxis {
			justificationOffset = (innerWidth - allChildrenWidth) / 2
		} else {
			justificationOffset = (innerHeight - allChildrenHeight) / 2
		}
	case EndJustify:
		if axis == HorizontalAxis {
			justificationOffset = innerWidth - allChildrenWidth
		} else {
			justificationOffset = innerHeight - allChildrenHeight
		}
	case SpaceBetweenJustify:
		if element.ChildCount != 1 {
			if axis == HorizontalAxis {
				justificationGap = (innerWidth - allChildrenWidth) / (element.ChildCount - 1)
			} else {
				justificationGap = (innerHeight - allChildrenHeight) / (element.ChildCount - 1)
			}
		}
	case SpaceAroundJustify:
		if axis == HorizontalAxis {
			justificationGap = (innerWidth - allChildrenWidth) / element.ChildCount
		} else {
			justificationGap = (innerHeight - allChildrenHeight) / element.ChildCount
		}
		justificationOffset = justificationGap / 2
	case SpaceEvenlyJustify:
		if axis == HorizontalAxis {
			justificationGap = (innerWidth - allChildrenWidth) / (element.ChildCount + 1)
		} else {
			justificationGap = (innerHeight - allChildrenHeight) / (element.ChildCount + 1)
		}
		justificationOffset = justificationGap
	}

	if axis == HorizontalAxis {
		innerX += justificationOffset
	} else {
		innerY += justificationOffset
	}

	// Position the children.
	for childElement := range element.ChildrenIter {

		// Alignment resolves Y for horizontal axis, and X for vertical axis.
		switch alignment {
		case StretchAlign:
			if axis == HorizontalAxis {
				childElement.Position.Y = innerY
			} else {
				childElement.Position.X = innerX
			}
		case StartAlign:
			if axis == HorizontalAxis {
				childElement.Position.Y = innerY
			} else {
				childElement.Position.X = innerX
			}
		case CenterAlign:
			if axis == HorizontalAxis {
				yOffset := (innerHeight - childElement.Size.Height) / 2
				childElement.Position.Y = innerY + yOffset
			} else {
				xOffset := (innerWidth - childElement.Size.Width) / 2
				childElement.Position.X = innerX + xOffset
			}
		case EndAlign:
			if axis == HorizontalAxis {
				yOffset := innerHeight - childElement.Size.Height
				childElement.Position.Y = innerY + yOffset
			} else {
				xOffset := innerWidth - childElement.Size.Width
				childElement.Position.X = innerX + xOffset
			}
		}

		if axis == HorizontalAxis {
			childElement.Position.X = innerX
			innerX += childElement.Size.Width + gap + justificationGap
		} else {
			childElement.Position.Y = innerY
			innerY += childElement.Size.Height + gap + justificationGap
		}
	}
}
