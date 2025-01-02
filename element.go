package blitra

import (
	"fmt"
	"reflect"
)

// represents the kind of element.
type ElementKind int

const (
	// represents a container element.
	ContainerElementKind ElementKind = iota
	// represents a text element.
	TextElementKind
)

// Used by Biltra to represent each renderable as an element. Elements are
// arranged in a tree structure and are traversed to compute layout and render
// the final output.
type Element struct {
	Kind  ElementKind
	ID    string
	Style Style

	Parent     *Element
	Previous   *Element
	Next       *Element
	FirstChild *Element
	LastChild  *Element
	ChildCount int

	IntrinsicSize Size
	AvailableSize Size
	SourceText    string

	Size     Size
	Position Point
	Text     string
}

type ElementIndex map[string]*Element

func ElementTreeAndIndexFromRenderable(renderable Renderable, state ViewState) (*Element, ElementIndex, error) {
	elementIndex := map[string]*Element{}
	rootElement := &Element{
		Kind:  ContainerElementKind,
		ID:    renderable.ID(),
		Style: renderable.Style(),
	}
	elementIndex[rootElement.ID] = rootElement

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
				Kind:       TextElementKind,
				ID:         v,
				SourceText: v,
			}
			elementIndex[v] = element
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
				return nil, nil, fmt.Errorf("struct type does not implement the Renderable interface: %s", reflect.TypeOf(v).String())
			}
			element := &Element{
				Kind:  ContainerElementKind,
				ID:    renderable.ID(),
				Style: renderable.Style(),
			}
			elementIndex[element.ID] = element
			head.parent.AddChild(element)
			tail.next = &pending{
				parent: element,
				result: renderable.Render(state),
			}
			tail = tail.next
			head = head.next
		}
	}

	return rootElement, elementIndex, nil
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

func (e *Element) LeftMargin() int {
	return V(e.Style.LeftMargin)
}

func (e *Element) RightMargin() int {
	return V(e.Style.RightMargin)
}

func (e *Element) TopMargin() int {
	return V(e.Style.TopMargin)
}

func (e *Element) BottomMargin() int {
	return V(e.Style.BottomMargin)
}

func (e *Element) LeftBorderWidth() int {
	return VMap(e.Style.LeftBorder, func(b Border) int {
		return len([]rune(b.Left))
	})
}

func (e *Element) RightBorderWidth() int {
	return VMap(e.Style.RightBorder, func(b Border) int {
		return len([]rune(b.Right))
	})
}

func (e *Element) TopBorderHeight() int {
	return VMap(e.Style.TopBorder, func(b Border) int {
		return len([]rune(b.Top))
	})
}

func (e *Element) BottomBorderHeight() int {
	return VMap(e.Style.BottomBorder, func(b Border) int {
		return len([]rune(b.Bottom))
	})
}

func (e *Element) LeftPadding() int {
	return V(e.Style.LeftPadding)
}

func (e *Element) RightPadding() int {
	return V(e.Style.RightPadding)
}

func (e *Element) TopPadding() int {
	return V(e.Style.TopPadding)
}

func (e *Element) BottomPadding() int {
	return V(e.Style.BottomPadding)
}

func (e *Element) LeftEdge() int {
	return e.LeftMargin() + e.LeftPadding() + e.LeftBorderWidth()
}

func (e *Element) RightEdge() int {
	return e.RightMargin() + e.RightPadding() + e.RightBorderWidth()
}

func (e *Element) TopEdge() int {
	return e.TopMargin() + e.TopPadding() + e.TopBorderHeight()
}

func (e *Element) BottomEdge() int {
	return e.BottomMargin() + e.BottomPadding() + e.BottomBorderHeight()
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
