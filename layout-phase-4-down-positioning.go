package blitra

func PositioningVisitor(element *Element) {
	if element.ChildCount == 0 {
		return
	}

	axis := V(element.Style.Axis)
	alignment := V(element.Style.Align)
	justification := V(element.Style.Justify)
	remainingChildWidth := element.Size.Width - element.IntrinsicSize.Width
	remainingChildHeight := element.Size.Height - element.IntrinsicSize.Height
	innerX := element.Position.X + element.GetLeftEdgeWidth()
	innerY := element.Position.Y + element.GetTopEdgeHeight()
	innerHeight := element.Size.Height - element.GetEdgeHeight()
	innerWidth := element.Size.Width - element.GetEdgeWidth()

	// Justification requires us to calculate a gap offset to properly space the
	// children. It also requires us to adjust the innerX or innerY position for
	// some modes.
	gapOffset := 0
	switch justification {
	case SpaceBetweenJustify:
		if element.ChildCount != 1 {
			if axis == HorizontalAxis {
				gapOffset = remainingChildWidth / (element.ChildCount - 1)
			} else {
				gapOffset = remainingChildHeight / (element.ChildCount - 1)
			}
		}
	case SpaceAroundJustify:
		if axis == HorizontalAxis {
			gapOffset = remainingChildWidth / element.ChildCount
			innerX += gapOffset / 2
		} else {
			gapOffset = remainingChildHeight / element.ChildCount
			innerY += gapOffset / 2
		}
	case SpaceEvenlyJustify:
		if axis == HorizontalAxis {
			gapOffset = remainingChildWidth / (element.ChildCount + 1)
			innerX += gapOffset
		} else {
			gapOffset = remainingChildHeight / (element.ChildCount + 1)
			innerY += gapOffset
		}
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

		// Justification resolves X for horizontal axis, and Y for vertical axis.
		switch justification {
		case StartJustify:
			if axis == HorizontalAxis {
				childElement.Position.X = innerX
			} else {
				childElement.Position.Y = innerY
			}
		case CenterJustify:
			if axis == HorizontalAxis {
				xOffset := (innerWidth - childElement.Size.Width) / 2
				childElement.Position.X = innerX + xOffset
			} else {
				yOffset := (innerHeight - childElement.Size.Height) / 2
				childElement.Position.Y = innerY + yOffset
			}
		case EndJustify:
			if axis == HorizontalAxis {
				xOffset := innerWidth - childElement.Size.Width
				childElement.Position.X = innerX + xOffset
			} else {
				yOffset := innerHeight - childElement.Size.Height
				childElement.Position.Y = innerY + yOffset
			}
		case SpaceBetweenJustify, SpaceAroundJustify, SpaceEvenlyJustify:
			if axis == HorizontalAxis {
				childElement.Position.X = innerX
			} else {
				childElement.Position.Y = innerY
			}
		}

		if axis == HorizontalAxis {
			innerX += childElement.Size.Width + V(element.Style.Gap) + gapOffset
		} else {
			innerY += childElement.Size.Height + V(element.Style.Gap) + gapOffset
		}
	}
}
