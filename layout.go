package blitra

// Updates the given element and it's children's layout.
func UpdateLayout(element *Element) {
	element.VisitContainerElementsDownThenUp(AvailableSizingVisitor, IntrinsicSizingVisitor)
	element.VisitContainerElementsDown(func(e *Element) {
		FinalSizingVisitor(e)
		PositioningVisitor(e)
	})
}
