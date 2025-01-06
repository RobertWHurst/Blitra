package blitra_test

import (
	"testing"

	"github.com/RobertWHurst/blitra"
	"github.com/stretchr/testify/assert"
)

func TestIntrinsicSizingVisitor(t *testing.T) {
	t.Run("Calculates the intrinsic size from text", func(t *testing.T) {
		element := &blitra.Element{}

		element.Parent = &blitra.Element{}
		element.Kind = blitra.TextElementKind
		element.AvailableSize = blitra.Size{Width: 10, Height: 5}
		element.SourceText = "Hello, World!"

		err := blitra.IntrinsicSizingVisitor(element, nil)
		assert.Nil(t, err)

		assert.Equal(t, 6, element.IntrinsicSize.Width)
		assert.Equal(t, 2, element.IntrinsicSize.Height)
	})

	t.Run("Calculates the intrinsic size from children", func(t *testing.T) {
		parent := &blitra.Element{}
		child := &blitra.Element{}

		parent.Kind = blitra.ContainerElementKind
		parent.Parent = &blitra.Element{}
		parent.Style.Axis = blitra.P(blitra.HorizontalAxis)
		parent.FirstChild = child
		parent.LastChild = child
		parent.ChildCount = 1
		parent.AvailableSize = blitra.Size{Width: 10, Height: 5}

		child.Kind = blitra.TextElementKind
		child.Parent = parent
		child.AvailableSize = blitra.Size{Width: 10, Height: 5}
		child.IntrinsicSize = blitra.Size{Width: 6, Height: 2}

		err := blitra.IntrinsicSizingVisitor(parent, nil)
		assert.Nil(t, err)

		assert.Equal(t, 6, parent.IntrinsicSize.Width)
		assert.Equal(t, 2, parent.IntrinsicSize.Height)
	})

	t.Run("Calculates the intrinsic size from multiple children in a horizontal axis", func(t *testing.T) {
		parent := &blitra.Element{}
		child1 := &blitra.Element{}
		child2 := &blitra.Element{}

		parent.Kind = blitra.ContainerElementKind
		parent.Parent = &blitra.Element{}
		parent.Style.Axis = blitra.P(blitra.HorizontalAxis)
		parent.FirstChild = child1
		parent.LastChild = child2
		parent.ChildCount = 2
		parent.AvailableSize = blitra.Size{Width: 10, Height: 5}

		child1.Kind = blitra.TextElementKind
		child1.Parent = parent
		child1.Next = child2
		child1.AvailableSize = blitra.Size{Width: 10, Height: 5}
		child1.IntrinsicSize = blitra.Size{Width: 8, Height: 4}

		child2.Kind = blitra.TextElementKind
		child2.Parent = parent
		child2.Previous = child1
		child2.AvailableSize = blitra.Size{Width: 10, Height: 5}
		child2.IntrinsicSize = blitra.Size{Width: 4, Height: 3}

		err := blitra.IntrinsicSizingVisitor(parent, nil)
		assert.Nil(t, err)

		assert.Equal(t, 12, parent.IntrinsicSize.Width)
		assert.Equal(t, 4, parent.IntrinsicSize.Height)
	})

	t.Run("Calculates the intrinsic size from multiple children in a vertical axis", func(t *testing.T) {
		parent := &blitra.Element{}
		child1 := &blitra.Element{}
		child2 := &blitra.Element{}

		parent.Kind = blitra.ContainerElementKind
		parent.Parent = &blitra.Element{}
		parent.Style.Axis = blitra.P(blitra.VerticalAxis)
		parent.FirstChild = child1
		parent.LastChild = child2
		parent.ChildCount = 2
		parent.AvailableSize = blitra.Size{Width: 10, Height: 5}

		child1.Kind = blitra.TextElementKind
		child1.Parent = parent
		child1.Next = child2
		child1.AvailableSize = blitra.Size{Width: 10, Height: 5}
		child1.IntrinsicSize = blitra.Size{Width: 8, Height: 4}

		child2.Kind = blitra.TextElementKind
		child2.Parent = parent
		child2.Previous = child1
		child2.AvailableSize = blitra.Size{Width: 10, Height: 5}
		child2.IntrinsicSize = blitra.Size{Width: 4, Height: 3}

		err := blitra.IntrinsicSizingVisitor(parent, nil)
		assert.Nil(t, err)

		assert.Equal(t, 8, parent.IntrinsicSize.Width)
		assert.Equal(t, 7, parent.IntrinsicSize.Height)
	})
}
