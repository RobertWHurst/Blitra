package blitra

import "reflect"

// Used by Biltra to represent each renderable as an element. Elements are
// arranged in a tree structure and are traversed to compute layout and render
// the final output.
type Element struct {
	// The parent element. Is a root element if nil.
	Parent *Element
	// The previous sibling element. If nil, this is the first child.
	Previous *Element
	// The next sibling element. If nil, this is the last child.
	Next *Element
	// The first child element.
	FirstChild *Element
	// The last child element.
	LastChild *Element
	// The number of children.
	ChildCount int

	// For text elements, the initial text value is stored here. It will
	// be converted to runes and stored in the Runes field after rendering.
	SourceText *string

	IntrinsicSize Size
	AvailableSize Size
	Size          Size
	Position      Point
	Text          string

	// The style of the element. Used to compute layout and visual style.
	Style Style
}

func ElementFromRenderable(renderable Renderable, state ViewState) *Element {
	rootElement := &Element{
		Style: renderable.Style(),
	}

	type pending struct {
		parent *Element
		next   *pending
		result any
	}

	head := &pending{
		parent: rootElement,
		result: renderable.Render(state),
	}
	tail := head
	for head != nil {

		switch v := head.result.(type) {
		case nil:
			head = head.next

		case string:
			element := &Element{
				SourceText: &v,
			}
			head.parent.AddChild(element)
			head = head.next

		case []any:
			for _, subV := range v {
				tail.next = &pending{
					parent: head.parent,
					result: subV,
				}
				tail = tail.next
			}
			head = head.next

		default:
			renderable, ok := v.(Renderable)
			if !ok {
				panic("Struct type does not implement the Renderable interface: " + reflect.TypeOf(v).String())
			}
			element := &Element{
				Style: renderable.Style(),
			}
			head.parent.AddChild(element)
			tail.next = &pending{
				parent: element,
				result: renderable.Render(state),
			}
			tail = tail.next
			head = head.next
		}
	}

	return rootElement
}

// Adds a child element. Sets up all the necessary relationship pointers.
func (e *Element) AddChild(childElement *Element) {
	if childElement.Parent != nil {
		childElement.Parent.RemoveChild(childElement)
	}
	childElement.Parent = e
	if e.FirstChild == nil {
		e.FirstChild = childElement
	}
	if e.LastChild != nil {
		e.LastChild.Next = childElement
		childElement.Previous = e.LastChild
	}
	e.LastChild = childElement
	e.ChildCount += 1
}

// Removes a child element. Updates all the necessary relationship pointers.
func (e *Element) RemoveChild(childElement *Element) {
	if childElement.Parent != e {
		return
	}
	if e.FirstChild == childElement {
		e.FirstChild = childElement.Next
	}
	if e.LastChild == childElement {
		e.LastChild = childElement.Previous
	}
	if childElement.Previous != nil {
		childElement.Previous.Next = childElement.Next
	}
	if childElement.Next != nil {
		childElement.Next.Previous = childElement.Previous
	}
	childElement.Parent = nil
	childElement.Previous = nil
	childElement.Next = nil
	e.ChildCount -= 1
}

// Iterates over the children of the element.
func (e *Element) ChildrenIter(yield func(*Element) bool) {
	element := e.FirstChild
	for element != nil {
		if !yield(element) {
			return
		}
		element = element.Next
	}
}

// Executes a visitor function on the element and each descendant depth-first,
// top-down.
func (e *Element) VisitContainerElementsUp(fn func(*Element)) {
	e.traverseElements(nil, &fn)
}

// Executes a visitor function on the element and each descendant depth-first,
// bottom-up.
func (e *Element) VisitContainerElementsDown(fn func(*Element)) {
	e.traverseElements(&fn, nil)
}

// Executes two visitor functions on the element and each descendant depth-first.
// The first visitor is executed top-down and the second is executed bottom-up.
func (e *Element) VisitContainerElementsDownThenUp(downFn, upFn func(*Element)) {
	e.traverseElements(&downFn, &upFn)
}

// Gets the width of the left edge (left margin, left padding, and left border).
func (e *Element) GetLeftEdgeWidth() int {
	return V(e.Style.LeftMargin) + V(e.Style.LeftPadding) + VMap(e.Style.LeftBorder, func(b Border) int {
		return len([]rune(b.Left))
	})
}

// Gets the width of the right edge (right margin, right padding, and right border).
func (e *Element) GetRightEdgeWidth() int {
	return V(e.Style.RightMargin) + V(e.Style.RightPadding) + VMap(e.Style.RightBorder, func(b Border) int {
		return len([]rune(b.Right))
	})
}

// Gets the height of the top edge (top margin, top padding, and top border).
func (e *Element) GetTopEdgeHeight() int {
	return V(e.Style.TopMargin) + V(e.Style.TopPadding) + VMap(e.Style.TopBorder, func(b Border) int {
		return len([]rune(b.Top))
	})
}

// Gets the height of the bottom edge (bottom margin, bottom padding, and bottom border).
func (e *Element) GetBottomEdgeHeight() int {
	return V(e.Style.BottomMargin) + V(e.Style.BottomPadding) + VMap(e.Style.BottomBorder, func(b Border) int {
		return len([]rune(b.Bottom))
	})
}

// Gets the width taken up by the element's edges (margin, padding, border).
func (e *Element) GetEdgeWidth() int {
	width := 0
	width += V(e.Style.LeftMargin)
	width += V(e.Style.RightMargin)
	width += V(e.Style.LeftPadding)
	width += V(e.Style.RightPadding)
	width += VMap(e.Style.LeftBorder, func(b Border) int {
		return len([]rune(b.Left))
	})
	width += VMap(e.Style.RightBorder, func(b Border) int {
		return len([]rune(b.Right))
	})
	return width
}

// Gets the height taken up by the element's edges (margin, padding, border).
func (e *Element) GetEdgeHeight() int {
	height := 0
	height += V(e.Style.TopMargin)
	height += V(e.Style.BottomMargin)
	height += V(e.Style.TopPadding)
	height += V(e.Style.BottomPadding)
	height += VMap(e.Style.TopBorder, func(b Border) int {
		return len([]rune(b.Top))
	})
	height += VMap(e.Style.BottomBorder, func(b Border) int {
		return len([]rune(b.Bottom))
	})
	return height
}

// Gets the width taken up by the gaps between the element's children.
func (e *Element) GetGapWidth() int {
	if V(e.Style.Axis) == VerticalAxis || e.Style.Gap == nil {
		return 0
	}
	return V(e.Style.Gap) * (e.ChildCount - 1)
}

// Gets the height taken up by the gaps between the element's children.
func (e *Element) GetGapHeight() int {
	if V(e.Style.Axis) == HorizontalAxis || e.Style.Gap == nil {
		return 0
	}
	return V(e.Style.Gap) * (e.ChildCount - 1)
}

func (e *Element) traverseElements(downFn, upFn *func(*Element)) {
	element := e

loop:
	for element != nil {
		if downFn != nil {
			(*downFn)(element)
		}

		if element.FirstChild != nil {
			element = element.FirstChild
			continue
		}

		for element.Parent != nil {
			if upFn != nil {
				(*upFn)(element)
			}
			if element.Next != nil {
				element = element.Next
				continue loop
			}
			element = element.Parent
		}

		if element.Parent == nil {
			if upFn != nil {
				(*upFn)(element)
			}
			break
		}
	}
}
