package blitra

import "slices"

func calcIntrinsicContainerSize(el *Element) error {
	if el.Parent == nil {
		return nil
	}

	assignedWidth := el.AssignedWidth()
	assignedHeight := el.AssignedHeight()

	if assignedWidth != nil {
		el.IntrinsicSize.Width = *assignedWidth + el.HorizontalMargin()
	}
	if assignedHeight != nil {
		el.IntrinsicSize.Height = *assignedHeight + el.VerticalMargin()
	}
	if assignedWidth != nil && assignedHeight != nil {
		return nil
	}

	if assignedWidth == nil {
		el.IntrinsicSize.Width = el.HorizontalEdge()
	}
	if assignedHeight == nil {
		el.IntrinsicSize.Height = el.VerticalEdge()
	}

	if el.ChildCount == 0 {
		return nil
	}

	axis := el.Axis()
	for cEl := range el.ChildrenIter {
		if axis == HorizontalAxis {
			if assignedWidth == nil {
				el.IntrinsicSize.Width += cEl.IntrinsicSize.Width
			}
			if assignedHeight == nil && cEl.IntrinsicSize.Height > el.IntrinsicSize.Height {
				el.IntrinsicSize.Height = cEl.IntrinsicSize.Height
			}
		} else {
			if assignedWidth == nil && cEl.IntrinsicSize.Width > el.IntrinsicSize.Width {
				el.IntrinsicSize.Width = cEl.IntrinsicSize.Width
			}
			if assignedHeight == nil {
				el.IntrinsicSize.Height += cEl.IntrinsicSize.Height
			}
		}
	}

	if el.ChildCount > 1 {
		if axis == HorizontalAxis {
			if assignedWidth == nil {
				el.IntrinsicSize.Width += el.Gap() * (el.ChildCount - 1)
			}
		} else {
			if assignedHeight == nil {
				el.IntrinsicSize.Height += el.Gap() * (el.ChildCount - 1)
			}
		}
	}

	if assignedWidth == nil {
		el.IntrinsicSize.Width = el.clampWidth(el.IntrinsicSize.Width)
	}
	if assignedHeight == nil {
		el.IntrinsicSize.Height = el.clampHeight(el.IntrinsicSize.Height)
	}

	return nil
}

func calcAvailableContainerSizesForChildren(el *Element) error {
	if el.ChildCount == 0 {
		return nil
	}

	axis := el.Axis()

	availableAxisLength := 0
	availableAxisSpan := 0
	if axis == HorizontalAxis {
		availableAxisLength = el.AvailableSize.Width - el.HorizontalEdge()
		availableAxisSpan = el.AvailableSize.Height - el.VerticalEdge()
	} else {
		availableAxisLength = el.AvailableSize.Height - el.VerticalEdge()
		availableAxisSpan = el.AvailableSize.Width - el.HorizontalEdge()
	}

	// assign initial available size from intrinsic, check for grow/shrink
	// and keep track of elements with assigned width/height.
	// most importantly, calculate the desired axis length and span.
	elsWithAssignedWidth := map[*Element]bool{}
	elsWithAssignedHeight := map[*Element]bool{}
	growableEls := []*Element{}
	shrinkableEls := []*Element{}

	growDivisor := 0
	shrinkDivisor := 0
	desiredAxisLength := 0
	desiredAxisSpan := 0
	for cEl := range el.ChildrenIter {
		if grow := cEl.Grow(); grow > 0 {
			growableEls = append(growableEls, cEl)
			growDivisor += grow
		}
		if shrink := cEl.Shrink(); shrink > 0 {
			shrinkableEls = append(shrinkableEls, cEl)
			shrinkDivisor += shrink
		}

		assignedWidth := cEl.AssignedWidth()
		assignedHeight := cEl.AssignedHeight()
		if assignedWidth != nil {
			cEl.AvailableSize.Width = *assignedWidth + cEl.HorizontalMargin()
			elsWithAssignedWidth[cEl] = true
		} else {
			cEl.AvailableSize.Width = el.IntrinsicSize.Width
		}
		if assignedHeight != nil {
			cEl.AvailableSize.Height = *assignedHeight + cEl.VerticalMargin()
			elsWithAssignedHeight[cEl] = true
		} else {
			cEl.AvailableSize.Height = el.IntrinsicSize.Height
		}

		if axis == HorizontalAxis {
			if cEl.AvailableSize.Height > desiredAxisSpan {
				desiredAxisSpan = cEl.AvailableSize.Height
			}
			desiredAxisLength += cEl.AvailableSize.Width
		} else {
			if cEl.AvailableSize.Width > desiredAxisSpan {
				desiredAxisSpan = cEl.AvailableSize.Width
			}
			desiredAxisLength += cEl.AvailableSize.Height
		}
	}
	if el.ChildCount > 1 {
		desiredAxisLength += el.Gap() * (el.ChildCount - 1)
	}
	axisLengthDelta := availableAxisLength - desiredAxisLength
	desiredAxisSpan = min(availableAxisSpan, desiredAxisSpan)

	// assign span
	for cEl := range el.ChildrenIter {
		if axis == HorizontalAxis {
			if !elsWithAssignedHeight[cEl] {
				cEl.AvailableSize.Height = desiredAxisSpan
			}
		} else {
			if !elsWithAssignedWidth[cEl] {
				cEl.AvailableSize.Width = desiredAxisSpan
			}
		}
	}

	// figure out if we are growing/shrinking
	growOrShrink := 0
	if axisLengthDelta > 0 {
		if growDivisor > 0 {
			growOrShrink = 1
		}
	} else if axisLengthDelta < 0 {
		if shrinkDivisor > 0 {
			growOrShrink = -1
		}
	}

	// If no shrink/grow then leave elements with their intrinsic length
	// and assigned span.
	if growOrShrink == 0 {
		return nil
	}

	// grow/shrink until we have no more length to distribute or we run out
	// of elements to grow/shrink.
	var targetEls []*Element
	var targetDivisor int
	if growOrShrink == 1 {
		targetEls = growableEls[:]
		targetDivisor = growDivisor
	} else {
		targetEls = shrinkableEls[:]
		targetDivisor = shrinkDivisor
	}
	for axisLengthDelta != 0 && len(targetEls) > 0 {
		fractionalLength := axisLengthDelta / targetDivisor

		for i := 0; i < len(targetEls); i += 1 {
			cEl := targetEls[i]

			var targetScalar int
			if growOrShrink == 1 {
				targetScalar = cEl.Grow()
			} else {
				targetScalar = cEl.Shrink()
			}
			targetLength := fractionalLength * targetScalar
			if axis == HorizontalAxis {
				currentWidth := cEl.AvailableSize.Width
				desiredWidth := currentWidth + targetLength
				cEl.AvailableSize.Width = cEl.clampWidth(desiredWidth)
				axisLengthDelta -= cEl.AvailableSize.Width - currentWidth
				if desiredWidth != cEl.AvailableSize.Width {
					targetEls = slices.Delete(targetEls, i, i+1)
					i -= 1
					targetDivisor -= targetScalar
				}
			} else {
				currentHeight := cEl.AvailableSize.Height
				desiredHeight := currentHeight + targetLength
				cEl.AvailableSize.Height = cEl.clampHeight(desiredHeight)
				axisLengthDelta -= cEl.AvailableSize.Height - currentHeight
				if desiredHeight != cEl.AvailableSize.Height {
					targetEls = slices.Delete(targetEls, i, i+1)
					i -= 1
					targetDivisor -= targetScalar
				}
			}
		}
	}

	return nil
}

func calcContainerPositionsForChildren(el *Element) error {
	if el.ChildCount == 0 {
		return nil
	}

	axis := el.Axis()

	// final sizing
	// TODO: handle alignment related size adjustments
	for cEl := range el.ChildrenIter {
		cEl.Size = cEl.AvailableSize
	}

	// axis spanwise positioning
	// TODO: handle alignment
	for cEl := range el.ChildrenIter {
		if axis == HorizontalAxis {
			cEl.Position.Y = el.Position.Y + el.TopPadding()
		} else {
			cEl.Position.X = el.Position.X + el.LeftPadding()
		}
	}

	// axis lengthwise positioning
	// TODO: handle justification
	axisLengthwiseOffset := 0
	if axis == HorizontalAxis {
		axisLengthwiseOffset = el.LeftPadding()
	} else {
		axisLengthwiseOffset = el.TopPadding()
	}
	for cEl := range el.ChildrenIter {
		if axis == HorizontalAxis {
			cEl.Position.X = el.Position.X + axisLengthwiseOffset
			axisLengthwiseOffset += cEl.Size.Width + el.Gap()
		} else {
			cEl.Position.Y = el.Position.Y + axisLengthwiseOffset
			axisLengthwiseOffset += cEl.Size.Height + el.Gap()
		}
	}

	return nil
}
