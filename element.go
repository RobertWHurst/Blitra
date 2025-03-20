package blitra

import (
	"fmt"
	"math"
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

type ElementLayoutState struct {
	ParentAxis Axis
	Grow       int
	Shrink     int
	Basis      int
	Size       int
	Length     int
}

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

	// Based on how big the element wants to be. Includes children.
	IntrinsicSize Size
	// Based on how big the element can be. Constrained by parents and siblings.
	AvailableSize Size

	SourceText      string
	TextReflowWidth *int

	Size     Size
	Position Point
	Text     string
}

type ElementIndex map[string]*Element

// Creates an element tree and element index from a renderable. A renderable is a struct that
// implements the Renderable interface. The element tree is created by calling the Render method
// of the renderable and traversing the result. The element index is a map of element IDs to elements.
// The element index is used to quickly locate elements by ID.
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
				SourceText: v,
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

func (e *Element) Axis() Axis {
	return VOr(e.Style.Axis, HorizontalAxis)
}

func (e *Element) Gap() int {
	return V(e.Style.Gap)
}

func (e *Element) Grow() int {
	return max(V(e.Style.Grow), 0)
}

func (e *Element) Shrink() int {
	return max(VOr(e.Style.Shrink, 1), 0)
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

func (e *Element) HorizontalMargin() int {
	return e.LeftMargin() + e.RightMargin()
}

func (e *Element) VerticalMargin() int {
	return e.TopMargin() + e.BottomMargin()
}

func (e *Element) LeftBorderWidth() int {
	if e.Style.LeftBorder == nil {
		return 0
	}
	return e.Style.LeftBorder.LeftWidth()
}

func (e *Element) RightBorderWidth() int {
	if e.Style.RightBorder == nil {
		return 0
	}
	return e.Style.RightBorder.RightWidth()
}

func (e *Element) TopBorderHeight() int {
	if e.Style.TopBorder == nil {
		return 0
	}
	return e.Style.TopBorder.TopHeight()
}

func (e *Element) BottomBorderHeight() int {
	if e.Style.BottomBorder == nil {
		return 0
	}
	return e.Style.BottomBorder.BottomHeight()
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

func (e *Element) HorizontalEdge() int {
	return e.LeftEdge() + e.RightEdge()
}

func (e *Element) VerticalEdge() int {
	return e.TopEdge() + e.BottomEdge()
}

func (e *Element) AssignedWidth() *int {
	if e.Style.Width == nil {
		return nil
	}
	width := e.clampWidth(V(e.Style.Width))
	return &width
}

func (e *Element) AssignedHeight() *int {
	if e.Style.Height == nil {
		return nil
	}
	height := e.clampHeight(V(e.Style.Height))
	return &height
}

func (e *Element) clampWidth(width int) int {
	minWidth := V(e.Style.MinWidth)
	maxWidth := VOr(e.Style.MaxWidth, math.MaxInt)
	return max(min(width, maxWidth), minWidth)
}

func (e *Element) clampHeight(height int) int {
	minHeight := V(e.Style.MinHeight)
	maxHeight := VOr(e.Style.MaxHeight, math.MaxInt)
	return max(min(height, maxHeight), minHeight)
}

func (e *Element) TextWrap() TextWrap {
	return VOr(e.Style.TextWrap, WordWrap)
}

func (e *Element) Ellipsis() bool {
	return VOr(e.Style.Ellipsis, true)
}

// Executes a visitor function on the root element and each descendant depth-first,
// top-down.
func VisitElementsUp[S any](rootElement *Element, state S, fn func(*Element, S) error) error {
	return VisitElementsDownThenUp(rootElement, state, nil, fn)
}

// Executes a visitor function on the root element and each descendant depth-first,
// bottom-up.
func VisitElementsDown[S any](rootElement *Element, state S, fn func(*Element, S) error) error {
	return VisitElementsDownThenUp(rootElement, state, fn, nil)
}

// Executes two visitor functions on the element and each descendant depth-first.
// The first visitor is executed top-down and the second is executed bottom-up.
func VisitElementsDownThenUp[S any](rootElement *Element, state S, downFn, upFn func(*Element, S) error) error {
	if rootElement.Parent != nil {
		return fmt.Errorf("element is not a root element")
	}

	element := rootElement

loop:
	for element != nil {
		if downFn != nil {
			err := downFn(element, state)
			if err != nil {
				return err
			}
		}

		if element.FirstChild != nil {
			element = element.FirstChild
			continue
		}

		for element.Parent != nil {
			if upFn != nil {
				err := upFn(element, state)
				if err != nil {
					return err
				}
			}
			if element.Next != nil {
				element = element.Next
				continue loop
			}
			element = element.Parent
		}

		if element.Parent == nil {
			if upFn != nil {
				err := upFn(element, state)
				if err != nil {
					return err
				}
			}
			break
		}
	}

	return nil
}

func MergeElementVisitors[S any](visitors ...func(*Element, S) error) func(*Element, S) error {
	return func(e *Element, s S) error {
		for _, visitor := range visitors {
			if err := visitor(e, s); err != nil {
				return err
			}
		}
		return nil
	}
}
