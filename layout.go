package blitra

import (
	"fmt"
	"math"
)

type LayoutState struct {
	needsReflow bool
	isReflowing bool
}

// Updates the given element and it's children's layout.
func UpdateLayout(rootElement *Element) error {
	state := &LayoutState{}
	return flowLayout(rootElement, state)
}

// Updates the given element and it's children's layout.
func flowLayout(rootElement *Element, state *LayoutState) error {

	err := VisitElementsDownThenUp(
		rootElement,
		state,
		MergeElementVisitors(
			calculateIntrinsicSize,
			calculateAvailableSize,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to calculate intrinsic size: %w", err)
	}

	return nil
}

// For intrinsic sizing we want to figure out the "natural" size of the element.
//
// For text, this means the width and height of the text with no wrapping;
// just explicit line breaks.
//
// For container elements, we take the elements styled min size, size, border
// and padding to calculate the width and height.
//
// This visitor affects the element itself only.
//
// The root element is skipped as it gets it's size from the view and is
// static.
func calculateIntrinsicSize(element *Element, _ any) error {
	if element.Parent == nil {
		return nil
	}

	element.ParentAxis = VOr(element.Parent.Style.Axis, HorizontalAxis)
	element.ParentGap = VOr(element.Parent.Style.Gap, 0)
	element.ShrinkFactor = VOr(element.Style.Shrink, 1)
	element.GrowFactor = V(element.Style.Grow)

	if element.Kind == TextElementKind {
		size := Size{Width: math.MaxInt, Height: math.MaxInt}
		_, info, err := ApplyWrap(NoWrap, false, size, element.SourceText)
		if err != nil {
			return fmt.Errorf("failed to calculate intrinsic size: %w", err)
		}
		element.IntrinsicSize = info.Size
		return nil
	}

	var width, height int
	if element.Style.Width != nil {
		width = *element.Style.Width
	} else {
		borderWidth := element.LeftBorderWidth() + element.RightBorderWidth()
		paddingWidth := element.LeftPadding() + element.RightPadding()
		width = borderWidth + paddingWidth
	}
	if element.Style.MinWidth != nil && *element.Style.MinWidth > width {
		width = *element.Style.MinWidth
	}

	if element.Style.Height != nil {
		height = *element.Style.Height
	} else {
		borderHeight := element.TopBorderHeight() + element.BottomBorderHeight()
		paddingHeight := element.TopPadding() + element.BottomPadding()
		height = borderHeight + paddingHeight
	}
	if element.Style.MinHeight != nil && *element.Style.MinHeight > height {
		height = *element.Style.MinHeight
	}

	element.IntrinsicSize = Size{Width: width, Height: height}

	return nil
}

// For available sizing we take the available size of the element and the
// intrinsic size of it's children. We then calculate the size of each child
// based on the available size of element.
//
// This visitor affects the element's children only.
//
// Must be run after calculateIntrinsicSize.
func calculateAvailableSize(element *Element, _ any) error {
	// Get the available inner size of the element.
	innerWidth := element.IntrinsicSize.Width
	innerHeight := element.IntrinsicSize.Height
	if element.Style.Width != nil {
		innerWidth = *element.Style.Width
	}
	if element.Style.Height != nil {
		innerHeight = *element.Style.Height
	}
	if element.Style.MaxWidth != nil && innerWidth > *element.Style.MaxWidth {
		innerWidth = *element.Style.MaxWidth
	}
	if element.Style.MaxHeight != nil && innerHeight > *element.Style.MaxHeight {
		innerHeight = *element.Style.MaxHeight
	}
	innerWidth -= element.LeftEdge() + element.RightEdge()
	innerHeight -= element.TopEdge() + element.BottomEdge()
	innerLength := getByAxis(element.ParentAxis, innerWidth, innerHeight)
	innerSpan := getByAxis(element.ParentAxis, innerHeight, innerWidth)

	// Get the intrinsic size of all the children, plus margins and gaps.
	growElements := []*Element{}
	shrinkElements := []*Element{}

	contentLength := 0
	contentSpan := 0
	for childElement := range element.ChildrenIter {
		intrinsicSize := childElement.IntrinsicSize
		intrinsicLength := getByAxis(element.ParentAxis, intrinsicSize.Width, intrinsicSize.Height)

		childElement.Length = VOr(childElement.Style.Basis, intrinsicLength)
		contentLength += childElement.Length

		childElement.Span = getByAxis(element.ParentAxis, intrinsicSize.Height, intrinsicSize.Width)
		contentSpan = max(contentSpan, childElement.Span)

		if childElement.ShrinkFactor != 0 {
			shrinkElements = append(shrinkElements, childElement)
		}
		if childElement.GrowFactor != 0 {
			growElements = append(growElements, childElement)
		}
	}

	gapLength := element.ParentGap * (element.ChildCount - 1)
	needsShrink := contentLength+gapLength > innerLength && len(shrinkElements) != 0
	needsGrow := contentLength+gapLength < innerLength && len(growElements) != 0

	if !needsShrink && !needsGrow {
		// TODO: maybe just apply the intrinsic size to the children?
		return nil
	}

	factoredLength := 0
	if needsShrink {
		for childElement := range element.ChildrenIter {
			factoredLength += (childElement.Length / element.ShrinkFactor) + childElement.ParentGap
		}

		for len(shrinkElements) != 0 {
			for _, childElement := range shrinkElements {
				targetLen := nil
			}
		}
	}

	if needsGrow {
		for childElement := range element.ChildrenIter {
			factoredLength += (childElement.Length * element.GrowFactor) + childElement.ParentGap
		}
	}

	// for len(elementsToFlow) != 0 {
	// 	for _, childElement := range elementsToFlow {
	// 		growFactor := V(childElement.Style.Grow)
	// 		shrinkFactor := VOr(childElement.Style.Shrink, 1)
	// 		intrinsicWidth := childElement.IntrinsicSize.Width
	// 		intrinsicHeight := childElement.IntrinsicSize.Height

	// 		if axis == HorizontalAxis {
	// 			switch {
	// 			case needsShrink:

	// 			case needsGrow:

	// 			default:
	// 				childElement.AvailableSize.Width = intrinsicWidth
	// 			}
	// 		} else {
	// 			switch {
	// 			case needsShrink:

	// 			case needsGrow:

	// 			default:
	// 				childElement.AvailableSize.Height = intrinsicHeight
	// 			}
	// 		}
	// 	}
	// }

	return nil
}

func getByAxis[V any](axis Axis, horizontalValue V, verticalValue V) V {
	if axis == HorizontalAxis {
		return horizontalValue
	} else {
		return verticalValue
	}
}
