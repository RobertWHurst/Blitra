package blitra_test

import (
	"errors"
	"testing"

	"github.com/RobertWHurst/blitra"
	"github.com/stretchr/testify/assert"
)

func TestElementTreeAndIndexFromRenderable(t *testing.T) {

	t.Run("Can take a renderable that produces a string", func(t *testing.T) {
		id := "test-id"
		text := "Hello, World!"
		renderable := NewTestRenderable(id, text)
		element, elementIndex, err := blitra.ElementTreeAndIndexFromRenderable(renderable, blitra.ViewState{})
		assert.Nil(t, err)

		assert.NotNil(t, element)
		assert.Equal(t, id, element.ID)
		assert.Equal(t, blitra.ContainerElementKind, element.Kind)

		assert.NotNil(t, element.FirstChild)
		assert.Equal(t, blitra.TextElementKind, element.FirstChild.Kind)
		assert.Equal(t, text, element.FirstChild.SourceText)

		assert.NotNil(t, elementIndex)
		assert.Equal(t, 1, len(elementIndex))
		assert.Equal(t, element, elementIndex[id])
	})

	t.Run("Can take a renderable that produces a slice of any, containing strings", func(t *testing.T) {
		id1 := "test-id"
		text1 := "Hello"
		text2 := "World!"
		renderable := NewTestRenderable(id1, []any{text1, text2})

		element, elementIndex, err := blitra.ElementTreeAndIndexFromRenderable(renderable, blitra.ViewState{})
		assert.Nil(t, err)

		assert.NotNil(t, element)
		assert.Equal(t, id1, element.ID)
		assert.Equal(t, blitra.ContainerElementKind, element.Kind)

		assert.NotNil(t, element.FirstChild)
		assert.Equal(t, blitra.TextElementKind, element.FirstChild.Kind)
		assert.Equal(t, text1, element.FirstChild.SourceText)

		assert.NotNil(t, element.FirstChild.Next)
		assert.Equal(t, blitra.TextElementKind, element.FirstChild.Next.Kind)
		assert.Equal(t, text2, element.FirstChild.Next.SourceText)

		assert.NotNil(t, elementIndex)
		assert.Equal(t, 1, len(elementIndex))
		assert.Equal(t, element, elementIndex[id1])
	})

	t.Run("Can take a renderable that produces another renderable", func(t *testing.T) {
		id1 := "test-id"
		id2 := "test-id2"
		text := "Hello, World!"
		renderable := NewTestRenderable(id1, NewTestRenderable(id2, text))

		element, elementIndex, err := blitra.ElementTreeAndIndexFromRenderable(renderable, blitra.ViewState{})
		assert.Nil(t, err)

		assert.NotNil(t, element)
		assert.Equal(t, id1, element.ID)
		assert.Equal(t, blitra.ContainerElementKind, element.Kind)

		assert.NotNil(t, element.FirstChild)
		assert.Equal(t, blitra.ContainerElementKind, element.FirstChild.Kind)
		assert.Equal(t, id2, element.FirstChild.ID)

		assert.NotNil(t, element.FirstChild.FirstChild)
		assert.Equal(t, blitra.TextElementKind, element.FirstChild.FirstChild.Kind)
		assert.Equal(t, text, element.FirstChild.FirstChild.SourceText)

		assert.NotNil(t, elementIndex)
		assert.Equal(t, 2, len(elementIndex))
		assert.Equal(t, element, elementIndex[id1])
		assert.Equal(t, element.FirstChild, elementIndex[id2])
	})

	t.Run("Can take a renderable that produces a slice of any, containing renderables", func(t *testing.T) {
		id1 := "test-id"
		id2 := "test-id2"
		id3 := "test-id3"
		text1 := "Hello"
		text2 := "World!"
		renderable := NewTestRenderable(id1, []any{
			NewTestRenderable(id2, text1),
			NewTestRenderable(id3, text2),
		})

		element, elementIndex, err := blitra.ElementTreeAndIndexFromRenderable(renderable, blitra.ViewState{})
		assert.Nil(t, err)

		assert.NotNil(t, element)
		assert.Equal(t, id1, element.ID)
		assert.Equal(t, blitra.ContainerElementKind, element.Kind)

		assert.NotNil(t, element.FirstChild)
		assert.Equal(t, blitra.ContainerElementKind, element.FirstChild.Kind)
		assert.Equal(t, id2, element.FirstChild.ID)

		assert.NotNil(t, element.FirstChild.FirstChild)
		assert.Equal(t, blitra.TextElementKind, element.FirstChild.FirstChild.Kind)
		assert.Equal(t, text1, element.FirstChild.FirstChild.SourceText)

		assert.NotNil(t, element.FirstChild.Next)
		assert.Equal(t, blitra.ContainerElementKind, element.FirstChild.Next.Kind)
		assert.Equal(t, id3, element.FirstChild.Next.ID)

		assert.NotNil(t, element.FirstChild.Next.FirstChild)
		assert.Equal(t, blitra.TextElementKind, element.FirstChild.Next.FirstChild.Kind)
		assert.Equal(t, text2, element.FirstChild.Next.FirstChild.SourceText)

		assert.NotNil(t, elementIndex)
		assert.Equal(t, 3, len(elementIndex))
		assert.Equal(t, element, elementIndex[id1])
		assert.Equal(t, element.FirstChild, elementIndex[id2])
		assert.Equal(t, element.FirstChild.Next, elementIndex[id3])
	})

	t.Run("Can take a renderable that produces nil producing no element in it's place", func(t *testing.T) {
		id := "test-id"
		renderable := NewTestRenderable(id, nil)
		element, elementIndex, err := blitra.ElementTreeAndIndexFromRenderable(renderable, blitra.ViewState{})
		assert.Nil(t, err)

		assert.NotNil(t, element)
		assert.Equal(t, id, element.ID)
		assert.Equal(t, blitra.ContainerElementKind, element.Kind)

		assert.Nil(t, element.FirstChild)

		assert.NotNil(t, elementIndex)
		assert.Equal(t, 1, len(elementIndex))
		assert.Equal(t, element, elementIndex[id])
	})

	t.Run("Returns an error if the renderable produces a non renderable struct", func(t *testing.T) {
		renderable := NewTestRenderable("", struct{}{})
		_, _, err := blitra.ElementTreeAndIndexFromRenderable(renderable, blitra.ViewState{})
		assert.NotNil(t, err)
	})

	t.Run("Can create a tree of elements from a more complex structure", func(t *testing.T) {
		// A
		// | \
		// B  C --
		// | \  \  \
		// D  E  F  G

		idA := "A"
		idB := "B"
		idC := "C"
		idD := "D"
		idE := "E"
		idF := "F"
		idG := "G"

		textD := "D"
		textE := "E"
		textF := "F"

		renderable := NewTestRenderable(idA, []any{
			NewTestRenderable(idB, []any{
				NewTestRenderable(idD, textD),
				NewTestRenderable(idE, textE),
			}),
			NewTestRenderable(idC, []any{
				NewTestRenderable(idF, textF),
				NewTestRenderable(idG, nil),
			}),
		})

		element, elementIndex, err := blitra.ElementTreeAndIndexFromRenderable(renderable, blitra.ViewState{})
		assert.Nil(t, err)

		eA := element
		eB := element.FirstChild
		eC := eB.Next
		eD := eB.FirstChild
		eE := eD.Next
		eF := eC.FirstChild
		eG := eF.Next

		eDTxt := eD.FirstChild
		eETxt := eE.FirstChild
		eFTxt := eF.FirstChild
		eGTxt := eG.FirstChild

		assert.NotNil(t, element)
		assert.Equal(t, idA, eA.ID)
		assert.Equal(t, blitra.ContainerElementKind, eA.Kind)

		assert.NotNil(t, eB)
		assert.Equal(t, idB, eB.ID)
		assert.Equal(t, blitra.ContainerElementKind, eB.Kind)

		assert.NotNil(t, eC)
		assert.Equal(t, idC, eC.ID)
		assert.Equal(t, blitra.ContainerElementKind, eC.Kind)

		assert.NotNil(t, eD)
		assert.Equal(t, idD, eD.ID)
		assert.Equal(t, blitra.ContainerElementKind, eD.Kind)

		assert.NotNil(t, eE)
		assert.Equal(t, idE, eE.ID)
		assert.Equal(t, blitra.ContainerElementKind, eE.Kind)

		assert.NotNil(t, eF)
		assert.Equal(t, idF, eF.ID)
		assert.Equal(t, blitra.ContainerElementKind, eF.Kind)

		assert.NotNil(t, eG)
		assert.Equal(t, idG, eG.ID)
		assert.Equal(t, blitra.ContainerElementKind, eG.Kind)

		assert.NotNil(t, eDTxt)
		assert.Equal(t, blitra.TextElementKind, eDTxt.Kind)
		assert.Equal(t, textD, eDTxt.SourceText)

		assert.NotNil(t, eETxt)
		assert.Equal(t, blitra.TextElementKind, eETxt.Kind)
		assert.Equal(t, textE, eETxt.SourceText)

		assert.NotNil(t, eFTxt)
		assert.Equal(t, blitra.TextElementKind, eFTxt.Kind)
		assert.Equal(t, textF, eFTxt.SourceText)

		assert.Nil(t, eGTxt)

		assert.NotNil(t, elementIndex)
		assert.Equal(t, 7, len(elementIndex))
		assert.Equal(t, eA, elementIndex[idA])
		assert.Equal(t, eB, elementIndex[idB])
		assert.Equal(t, eC, elementIndex[idC])
		assert.Equal(t, eD, elementIndex[idD])
		assert.Equal(t, eE, elementIndex[idE])
		assert.Equal(t, eF, elementIndex[idF])
	})
}

func TestElementAddChild(t *testing.T) {
	t.Run("Correctly sets up the relationships between a parent and child", func(t *testing.T) {
		parent := &blitra.Element{}
		child := &blitra.Element{}

		parent.AddChild(child)

		assert.Equal(t, 1, parent.ChildCount)
		assert.Equal(t, child, parent.FirstChild)
		assert.Equal(t, child, parent.LastChild)
		assert.Equal(t, parent, child.Parent)
		assert.Nil(t, child.Next)
		assert.Nil(t, child.Previous)
	})

	t.Run("Correctly sets up the relationships between a parent and multiple children", func(t *testing.T) {
		parent := &blitra.Element{}
		child1 := &blitra.Element{}
		child2 := &blitra.Element{}

		parent.AddChild(child1)
		parent.AddChild(child2)

		assert.Equal(t, 2, parent.ChildCount)
		assert.Equal(t, child1, parent.FirstChild)
		assert.Equal(t, child2, parent.LastChild)
		assert.Equal(t, parent, child1.Parent)
		assert.Equal(t, parent, child2.Parent)
		assert.Equal(t, child2, child1.Next)
		assert.Nil(t, child1.Previous)
		assert.Equal(t, child1, child2.Previous)
		assert.Nil(t, child2.Next)
	})
}

func TestElementRemoveChild(t *testing.T) {
	t.Run("Correctly removes a child from a parent", func(t *testing.T) {
		parent := &blitra.Element{}
		child := &blitra.Element{}

		parent.ChildCount = 1
		parent.FirstChild = child
		parent.LastChild = child
		child.Parent = parent

		parent.RemoveChild(child)

		assert.Equal(t, 0, parent.ChildCount)
		assert.Nil(t, parent.FirstChild)
		assert.Nil(t, parent.LastChild)
		assert.Nil(t, child.Parent)
		assert.Nil(t, child.Next)
		assert.Nil(t, child.Previous)
	})

	t.Run("Correctly removes the first child from a parent with three", func(t *testing.T) {
		parent := &blitra.Element{}
		child1 := &blitra.Element{}
		child2 := &blitra.Element{}
		child3 := &blitra.Element{}

		parent.ChildCount = 3
		parent.FirstChild = child1
		parent.LastChild = child3
		child1.Parent = parent
		child1.Next = child2
		child2.Parent = parent
		child2.Previous = child1
		child2.Next = child3
		child3.Parent = parent
		child3.Previous = child2

		parent.RemoveChild(child1)

		assert.Equal(t, 2, parent.ChildCount)
		assert.Equal(t, child2, parent.FirstChild)
		assert.Equal(t, child3, parent.LastChild)
		assert.Equal(t, child3, child2.Next)
		assert.Equal(t, child2, child3.Previous)
		assert.Nil(t, child2.Previous)
		assert.Nil(t, child3.Next)
		assert.Nil(t, child1.Parent)
		assert.Nil(t, child1.Next)
		assert.Nil(t, child1.Previous)
	})
}

func TestElementChildrenIter(t *testing.T) {
	t.Run("Iterates over all the children", func(t *testing.T) {
		parent := &blitra.Element{}
		child1 := &blitra.Element{}
		child2 := &blitra.Element{}
		child3 := &blitra.Element{}

		parent.ChildCount = 3
		parent.FirstChild = child1
		parent.LastChild = child3
		child1.Next = child2
		child2.Previous = child1
		child2.Next = child3
		child3.Previous = child2

		var iteratedChildren []*blitra.Element
		for child := range parent.ChildrenIter {
			iteratedChildren = append(iteratedChildren, child)
		}

		assert.Equal(t, 3, len(iteratedChildren))
		assert.Equal(t, child1, iteratedChildren[0])
		assert.Equal(t, child2, iteratedChildren[1])
		assert.Equal(t, child3, iteratedChildren[2])
	})
}

func TestElementVisitElementsUp(t *testing.T) {
	t.Run("Visits each of the elements descendants up to and including itself", func(t *testing.T) {
		parent := &blitra.Element{}
		child1 := &blitra.Element{}
		child2 := &blitra.Element{}
		child3 := &blitra.Element{}

		parent.ChildCount = 3
		parent.FirstChild = child1
		parent.LastChild = child3
		child1.Parent = parent
		child1.Next = child2
		child2.Parent = parent
		child2.Previous = child1
		child2.Next = child3
		child3.Parent = parent
		child3.Previous = child2

		var visitedElements []*blitra.Element
		err := blitra.VisitElementsUp(parent, nil, func(e *blitra.Element, _ any) error {
			visitedElements = append(visitedElements, e)
			return nil
		})
		assert.Nil(t, err)

		assert.Equal(t, 4, len(visitedElements))
		assert.Equal(t, child1, visitedElements[0])
		assert.Equal(t, child2, visitedElements[1])
		assert.Equal(t, child3, visitedElements[2])
		assert.Equal(t, parent, visitedElements[3])
	})

	t.Run("Stops visiting if the visitor returns an error", func(t *testing.T) {
		parent := &blitra.Element{}
		child1 := &blitra.Element{}
		child2 := &blitra.Element{}
		child3 := &blitra.Element{}

		parent.ChildCount = 3
		parent.FirstChild = child1
		parent.LastChild = child3
		child1.Parent = parent
		child1.Next = child2
		child2.Parent = parent
		child2.Previous = child1
		child2.Next = child3
		child3.Parent = parent
		child3.Previous = child2

		var visitedElements []*blitra.Element
		err := blitra.VisitElementsUp(parent, nil, func(e *blitra.Element, _ any) error {
			visitedElements = append(visitedElements, e)
			if e == child2 {
				return errors.New("Test error")
			}
			e.ID = "Visited"
			return nil
		})
		assert.Error(t, err)

		assert.Equal(t, parent.ID, "")
		assert.Equal(t, child1.ID, "Visited")
		assert.Equal(t, child2.ID, "")
		assert.Equal(t, child3.ID, "")
	})
}

func TestElementVisitElementsDown(t *testing.T) {
	t.Run("Visits each of the elements descendants down to and including itself", func(t *testing.T) {
		parent := &blitra.Element{}
		child1 := &blitra.Element{}
		child2 := &blitra.Element{}
		child3 := &blitra.Element{}

		parent.ChildCount = 3
		parent.FirstChild = child1
		parent.LastChild = child3
		child1.Parent = parent
		child1.Next = child2
		child2.Parent = parent
		child2.Previous = child1
		child2.Next = child3
		child3.Parent = parent
		child3.Previous = child2

		var visitedElements []*blitra.Element
		err := blitra.VisitElementsDown(parent, nil, func(e *blitra.Element, _ any) error {
			visitedElements = append(visitedElements, e)
			return nil
		})
		assert.Nil(t, err)

		assert.Equal(t, 4, len(visitedElements))
		assert.Equal(t, parent, visitedElements[0])
		assert.Equal(t, child1, visitedElements[1])
		assert.Equal(t, child2, visitedElements[2])
		assert.Equal(t, child3, visitedElements[3])
	})

	t.Run("Stops visiting if the visitor returns an error", func(t *testing.T) {
		parent := &blitra.Element{}
		child1 := &blitra.Element{}
		child2 := &blitra.Element{}
		child3 := &blitra.Element{}

		parent.ChildCount = 3
		parent.FirstChild = child1
		parent.LastChild = child3
		child1.Parent = parent
		child1.Next = child2
		child2.Parent = parent
		child2.Previous = child1
		child2.Next = child3
		child3.Parent = parent
		child3.Previous = child2

		var visitedElements []*blitra.Element
		err := blitra.VisitElementsDown(parent, nil, func(e *blitra.Element, _ any) error {
			visitedElements = append(visitedElements, e)
			if e == child2 {
				return errors.New("Test error")
			}
			e.ID = "Visited"
			return nil
		})
		assert.Error(t, err)

		assert.Equal(t, parent.ID, "Visited")
		assert.Equal(t, child1.ID, "Visited")
		assert.Equal(t, child2.ID, "")
		assert.Equal(t, child3.ID, "")
	})
}

func TestElementVisitElementsDownThenUp(t *testing.T) {
	t.Run("Visits each of the elements descendants down to and including itself, then up to and including itself", func(t *testing.T) {
		parent := &blitra.Element{ID: "parent"}
		child1 := &blitra.Element{ID: "1"}
		child2 := &blitra.Element{ID: "2"}
		child3 := &blitra.Element{ID: "3"}
		child4 := &blitra.Element{ID: "4"}
		child5 := &blitra.Element{ID: "5"}
		child6 := &blitra.Element{ID: "6"}
		child7 := &blitra.Element{ID: "7"}
		child8 := &blitra.Element{ID: "8"}
		child9 := &blitra.Element{ID: "9"}
		child10 := &blitra.Element{ID: "10"}
		child11 := &blitra.Element{ID: "11"}
		child12 := &blitra.Element{ID: "12"}

		parent.FirstChild = child1
		parent.LastChild = child3
		child1.Parent = parent
		child2.Parent = parent
		child3.Parent = parent
		child1.Next = child2
		child2.Next = child3
		child2.Previous = child1
		child3.Previous = child2
		child1.FirstChild = child4
		child1.LastChild = child6
		child4.Parent = child1
		child5.Parent = child1
		child6.Parent = child1
		child4.Next = child5
		child5.Next = child6
		child5.Previous = child4
		child6.Previous = child5
		child2.FirstChild = child7
		child2.LastChild = child9
		child7.Parent = child2
		child8.Parent = child2
		child9.Parent = child2
		child7.Next = child8
		child8.Next = child9
		child8.Previous = child7
		child9.Previous = child8
		child3.FirstChild = child10
		child3.LastChild = child12
		child10.Parent = child3
		child11.Parent = child3
		child12.Parent = child3
		child10.Next = child11
		child11.Next = child12
		child11.Previous = child10
		child12.Previous = child11

		var visitedElements []*blitra.Element
		var visitedElementsDown []*blitra.Element
		var visitedElementsUp []*blitra.Element
		err := blitra.VisitElementsDownThenUp(parent, nil, func(e *blitra.Element, _ any) error {
			visitedElements = append(visitedElements, e)
			visitedElementsDown = append(visitedElementsDown, e)
			return nil
		}, func(e *blitra.Element, _ any) error {
			visitedElements = append(visitedElements, e)
			visitedElementsUp = append(visitedElementsUp, e)
			return nil
		})
		assert.Nil(t, err)

		assert.Equal(t, 26, len(visitedElements))
		assert.Equal(t, parent, visitedElements[0])
		assert.Equal(t, child1, visitedElements[1])
		assert.Equal(t, child4, visitedElements[2])
		assert.Equal(t, child4, visitedElements[3])
		assert.Equal(t, child5, visitedElements[4])
		assert.Equal(t, child5, visitedElements[5])
		assert.Equal(t, child6, visitedElements[6])
		assert.Equal(t, child6, visitedElements[7])
		assert.Equal(t, child1, visitedElements[8])
		assert.Equal(t, child2, visitedElements[9])
		assert.Equal(t, child7, visitedElements[10])
		assert.Equal(t, child7, visitedElements[11])
		assert.Equal(t, child8, visitedElements[12])
		assert.Equal(t, child8, visitedElements[13])
		assert.Equal(t, child9, visitedElements[14])
		assert.Equal(t, child9, visitedElements[15])
		assert.Equal(t, child2, visitedElements[16])
		assert.Equal(t, child3, visitedElements[17])
		assert.Equal(t, child10, visitedElements[18])
		assert.Equal(t, child10, visitedElements[19])
		assert.Equal(t, child11, visitedElements[20])
		assert.Equal(t, child11, visitedElements[21])
		assert.Equal(t, child12, visitedElements[22])
		assert.Equal(t, child12, visitedElements[23])
		assert.Equal(t, child3, visitedElements[24])
		assert.Equal(t, parent, visitedElements[25])

		assert.Equal(t, 13, len(visitedElementsDown))
		assert.Equal(t, parent, visitedElementsDown[0])
		assert.Equal(t, child1, visitedElementsDown[1])
		assert.Equal(t, child4, visitedElementsDown[2])
		assert.Equal(t, child5, visitedElementsDown[3])
		assert.Equal(t, child6, visitedElementsDown[4])
		assert.Equal(t, child2, visitedElementsDown[5])
		assert.Equal(t, child7, visitedElementsDown[6])
		assert.Equal(t, child8, visitedElementsDown[7])
		assert.Equal(t, child9, visitedElementsDown[8])
		assert.Equal(t, child3, visitedElementsDown[9])
		assert.Equal(t, child10, visitedElementsDown[10])
		assert.Equal(t, child11, visitedElementsDown[11])
		assert.Equal(t, child12, visitedElementsDown[12])

		assert.Equal(t, 13, len(visitedElementsUp))
		assert.Equal(t, child4, visitedElementsUp[0])
		assert.Equal(t, child5, visitedElementsUp[1])
		assert.Equal(t, child6, visitedElementsUp[2])
		assert.Equal(t, child1, visitedElementsUp[3])
		assert.Equal(t, child7, visitedElementsUp[4])
		assert.Equal(t, child8, visitedElementsUp[5])
		assert.Equal(t, child9, visitedElementsUp[6])
		assert.Equal(t, child2, visitedElementsUp[7])
		assert.Equal(t, child10, visitedElementsUp[8])
		assert.Equal(t, child11, visitedElementsUp[9])
		assert.Equal(t, child12, visitedElementsUp[10])
		assert.Equal(t, child3, visitedElementsUp[11])
		assert.Equal(t, parent, visitedElementsUp[12])
	})
}

func TestLeftMargin(t *testing.T) {
	t.Run("Gets the left margin from the style", func(t *testing.T) {
		style := blitra.Style{LeftMargin: blitra.P(10)}
		element := blitra.Element{Style: style}
		assert.Equal(t, 10, element.LeftMargin())
	})

	t.Run("Gets Zero if the left margin is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.LeftMargin())
	})
}

func TestRightMargin(t *testing.T) {
	t.Run("Gets the right margin from the style", func(t *testing.T) {
		style := blitra.Style{RightMargin: blitra.P(10)}
		element := blitra.Element{Style: style}
		assert.Equal(t, 10, element.RightMargin())
	})

	t.Run("Gets Zero if the right margin is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.RightMargin())
	})
}

func TestTopMargin(t *testing.T) {
	t.Run("Gets the top margin from the style", func(t *testing.T) {
		style := blitra.Style{TopMargin: blitra.P(10)}
		element := blitra.Element{Style: style}
		assert.Equal(t, 10, element.TopMargin())
	})

	t.Run("Gets Zero if the top margin is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.TopMargin())
	})
}

func TestBottomMargin(t *testing.T) {
	t.Run("Gets the bottom margin from the style", func(t *testing.T) {
		style := blitra.Style{BottomMargin: blitra.P(10)}
		element := blitra.Element{Style: style}
		assert.Equal(t, 10, element.BottomMargin())
	})

	t.Run("Gets Zero if the bottom margin is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.BottomMargin())
	})
}

func TestLeftBorderWidth(t *testing.T) {
	t.Run("Gets the left border width from the style", func(t *testing.T) {
		style := blitra.Style{LeftBorder: blitra.NewBorder("", "", "", "", "||", "", "", "")}
		element := blitra.Element{Style: style}
		assert.Equal(t, 2, element.LeftBorderWidth())
	})

	t.Run("Gets Zero if the left border width is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.LeftBorderWidth())
	})
}

func TestRightBorderWidth(t *testing.T) {
	t.Run("Gets the right border width from the style", func(t *testing.T) {
		style := blitra.Style{RightBorder: blitra.NewBorder("", "", "", "", "", "||", "", "")}
		element := blitra.Element{Style: style}
		assert.Equal(t, 2, element.RightBorderWidth())
	})

	t.Run("Gets Zero if the right border width is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.RightBorderWidth())
	})
}

func TestTopBorderHeight(t *testing.T) {
	t.Run("Gets the top border height from the style", func(t *testing.T) {
		style := blitra.Style{TopBorder: blitra.NewBorder("", "", "", "", "", "", "|\n|", "")}
		element := blitra.Element{Style: style}
		assert.Equal(t, 2, element.TopBorderHeight())
	})

	t.Run("Gets Zero if the top border height is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.TopBorderHeight())
	})
}

func TestBottomBorderHeight(t *testing.T) {
	t.Run("Gets the bottom border height from the style", func(t *testing.T) {
		style := blitra.Style{BottomBorder: blitra.NewBorder("", "", "", "", "", "", "", "|\n|")}
		element := blitra.Element{Style: style}
		assert.Equal(t, 2, element.BottomBorderHeight())
	})

	t.Run("Gets Zero if the bottom border height is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.BottomBorderHeight())
	})
}

func TestLeftPadding(t *testing.T) {
	t.Run("Gets the left padding from the style", func(t *testing.T) {
		style := blitra.Style{LeftPadding: blitra.P(10)}
		element := blitra.Element{Style: style}
		assert.Equal(t, 10, element.LeftPadding())
	})

	t.Run("Gets Zero if the left padding is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.LeftPadding())
	})
}

func TestRightPadding(t *testing.T) {
	t.Run("Gets the right padding from the style", func(t *testing.T) {
		style := blitra.Style{RightPadding: blitra.P(10)}
		element := blitra.Element{Style: style}
		assert.Equal(t, 10, element.RightPadding())
	})

	t.Run("Gets Zero if the right padding is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.RightPadding())
	})
}

func TestTopPadding(t *testing.T) {
	t.Run("Gets the top padding from the style", func(t *testing.T) {
		style := blitra.Style{TopPadding: blitra.P(10)}
		element := blitra.Element{Style: style}
		assert.Equal(t, 10, element.TopPadding())
	})

	t.Run("Gets Zero if the top padding is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.TopPadding())
	})
}

func TestBottomPadding(t *testing.T) {
	t.Run("Gets the bottom padding from the style", func(t *testing.T) {
		style := blitra.Style{BottomPadding: blitra.P(10)}
		element := blitra.Element{Style: style}
		assert.Equal(t, 10, element.BottomPadding())
	})

	t.Run("Gets Zero if the bottom padding is not set", func(t *testing.T) {
		style := blitra.Style{}
		element := blitra.Element{Style: style}
		assert.Equal(t, 0, element.BottomPadding())
	})
}

func TestLeftEdge(t *testing.T) {
	t.Run("Get a sum of the left margin, left padding, and left border width", func(t *testing.T) {
		style := blitra.Style{
			LeftMargin:  blitra.P(10),
			LeftPadding: blitra.P(10),
			LeftBorder:  blitra.NewBorder("", "", "", "", "||", "", "", ""),
		}
		element := blitra.Element{Style: style}
		assert.Equal(t, 22, element.LeftEdge())
	})
}

func TestRightEdge(t *testing.T) {
	t.Run("Get a sum of the right margin, right padding, and right border width", func(t *testing.T) {
		style := blitra.Style{
			RightMargin:  blitra.P(10),
			RightPadding: blitra.P(10),
			RightBorder:  blitra.NewBorder("", "", "", "", "", "||", "", ""),
		}
		element := blitra.Element{Style: style}
		assert.Equal(t, 22, element.RightEdge())
	})
}

func TestTopEdge(t *testing.T) {
	t.Run("Get a sum of the top margin, top padding, and top border height", func(t *testing.T) {
		style := blitra.Style{
			TopMargin:  blitra.P(10),
			TopPadding: blitra.P(10),
			TopBorder:  blitra.NewBorder("", "", "", "", "", "", "|\n|", ""),
		}
		element := blitra.Element{Style: style}
		assert.Equal(t, 22, element.TopEdge())
	})
}

func TestBottomEdge(t *testing.T) {
	t.Run("Get a sum of the bottom margin, bottom padding, and bottom border height", func(t *testing.T) {
		style := blitra.Style{
			BottomMargin:  blitra.P(10),
			BottomPadding: blitra.P(10),
			BottomBorder:  blitra.NewBorder("", "", "", "", "", "", "", "|\n|"),
		}
		element := blitra.Element{Style: style}
		assert.Equal(t, 22, element.BottomEdge())
	})
}

type TestRenderable struct {
	id          string
	returnValue any
}

var _ blitra.Renderable = TestRenderable{}

func NewTestRenderable(id string, returnValue any) TestRenderable {
	return TestRenderable{id: id, returnValue: returnValue}
}

func (r TestRenderable) ID() string {
	return r.id
}

func (r TestRenderable) Style() blitra.Style {
	return blitra.Style{}
}

func (r TestRenderable) Render(state blitra.ViewState) any {
	return r.returnValue
}
