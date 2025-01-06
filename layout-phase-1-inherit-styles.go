package blitra

// Takes an element and copies the parent's inheritable styles to the element.
//
// This visitor runs once. It is part of the first phase, and moves down the
// tree.
//
// TODO: Consider if any other styles should be inherited.
func InheritStylesVisitor(element *Element, _ *LayoutState) error {
	if element.Parent == nil {
		return nil
	}

	if element.Style.TextColor == nil {
		element.Style.TextColor = element.Parent.Style.TextColor
	}
	if element.Style.BackgroundColor == nil {
		element.Style.BackgroundColor = element.Parent.Style.BackgroundColor
	}
	if element.Style.TextWrap == nil {
		element.Style.TextWrap = element.Parent.Style.TextWrap
	}
	if element.Style.Ellipsis == nil {
		element.Style.Ellipsis = element.Parent.Style.Ellipsis
	}

	return nil
}
