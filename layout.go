package blitra

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
	MergeElementVisitors(
		FinalSizingVisitor,
		PositioningVisitor,
	)

	err := VisitElementsDownThenUp(
		rootElement,
		state,
		MergeElementVisitors(
			InheritStylesVisitor,
			AvailableSizingVisitor,
		),
		IntrinsicSizingVisitor,
	)
	if err != nil {
		return err
	}

	err = VisitElementsDown(
		rootElement,
		state,
		MergeElementVisitors(
			FinalSizingVisitor,
			PositioningVisitor,
		),
	)
	if err != nil {
		return err
	}

	if state.needsReflow {
		err = VisitElementsUp(rootElement, state, IntrinsicSizingVisitor)
		if err != nil {
			return err
		}

		err := VisitElementsDown(rootElement, state, MergeElementVisitors(
			FinalSizingVisitor,
			PositioningVisitor,
		))
		if err != nil {
			return err
		}
	}

	return nil
}
