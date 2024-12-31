package blitra

// Updates the given element and it's children's layout.
func UpdateLayout(element *Element) {
	element.VisitContainerElementsUp(IntrinsicSizingVisitor)
	element.VisitContainerElementsDownThenUp(
		AvailableSizingVisitor,
		FinalSizingVisitor,
	)
	element.VisitContainerElementsDown(
		PositioningVisitor,
	)
}
