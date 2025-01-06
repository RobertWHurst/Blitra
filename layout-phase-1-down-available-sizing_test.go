package blitra_test

import (
	"testing"

	"github.com/RobertWHurst/blitra"
	"github.com/stretchr/testify/assert"
)

func TestAvailableSizingVisitor(t *testing.T) {
	t.Run("Correctly flows available size from parent to child", func(t *testing.T) {
		parent := &blitra.Element{}
		child := &blitra.Element{}

		parent.ChildCount = 1
		parent.FirstChild = child
		parent.LastChild = child

		parent.Style.Width = blitra.P(100)
		parent.Style.Height = blitra.P(60)

		err := blitra.AvailableSizingVisitor(parent, nil)
		assert.Nil(t, err)

		assert.Equal(t, 100, child.AvailableSize.Width)
		assert.Equal(t, 60, child.AvailableSize.Height)
	})

	t.Run("Subtracts margins, padding, and border from available size", func(t *testing.T) {
		parent := &blitra.Element{}
		child := &blitra.Element{}

		parent.ChildCount = 1
		parent.FirstChild = child
		parent.LastChild = child

		parent.Style.Width = blitra.P(100)
		parent.Style.Height = blitra.P(60)
		parent.Style.LeftMargin = blitra.P(1)
		parent.Style.RightMargin = blitra.P(2)
		parent.Style.TopMargin = blitra.P(3)
		parent.Style.BottomMargin = blitra.P(4)
		parent.Style.LeftPadding = blitra.P(5)
		parent.Style.RightPadding = blitra.P(6)
		parent.Style.TopPadding = blitra.P(7)
		parent.Style.BottomPadding = blitra.P(8)
		parent.Style.LeftBorder = blitra.NewBorder("", "", "", "", "|", "", "", "")
		parent.Style.RightBorder = blitra.NewBorder("", "", "", "", "", "||", "", "")
		parent.Style.TopBorder = blitra.NewBorder("", "", "", "", "", "", "|\n|", "")
		parent.Style.BottomBorder = blitra.NewBorder("", "", "", "", "", "", "", "|\n|\n|")

		err := blitra.AvailableSizingVisitor(parent, nil)
		assert.Nil(t, err)

		assert.Equal(t, 100-1-2-5-6-1-2, child.AvailableSize.Width)
		assert.Equal(t, 60-3-4-7-8-2-3, child.AvailableSize.Height)
	})

	t.Run("Correctly flows inherited styles to children", func(t *testing.T) {
		parent := &blitra.Element{}
		child := &blitra.Element{}

		parent.ChildCount = 1
		parent.FirstChild = child
		parent.LastChild = child

		parent.Style.TextColor = blitra.P("#f00")
		parent.Style.BackgroundColor = blitra.P("#00f")
		parent.Style.TextWrap = blitra.P(blitra.CharacterWrap)
		parent.Style.Ellipsis = blitra.P(true)

		err := blitra.AvailableSizingVisitor(parent, nil)
		assert.Nil(t, err)

		assert.Equal(t, "#f00", *child.Style.TextColor)
		assert.Equal(t, "#00f", *child.Style.BackgroundColor)
		assert.Equal(t, blitra.CharacterWrap, *child.Style.TextWrap)
		assert.Equal(t, true, *child.Style.Ellipsis)
	})

	t.Run("Does not allow negative available size", func(t *testing.T) {
		parent := &blitra.Element{}
		child := &blitra.Element{}

		parent.ChildCount = 1
		parent.FirstChild = child
		parent.LastChild = child

		parent.Style.Width = blitra.P(2)
		parent.Style.Height = blitra.P(2)
		parent.Style.LeftMargin = blitra.P(1)
		parent.Style.RightMargin = blitra.P(2)
		parent.Style.TopMargin = blitra.P(3)
		parent.Style.BottomMargin = blitra.P(4)

		err := blitra.AvailableSizingVisitor(parent, nil)
		assert.Nil(t, err)

		assert.Equal(t, 0, child.AvailableSize.Width)
		assert.Equal(t, 0, child.AvailableSize.Height)
	})

	t.Run("Respects max width and height", func(t *testing.T) {
		parent := &blitra.Element{}
		child := &blitra.Element{}

		parent.ChildCount = 1
		parent.FirstChild = child
		parent.LastChild = child

		parent.Style.Width = blitra.P(100)
		parent.Style.Height = blitra.P(60)
		parent.Style.MaxWidth = blitra.P(50)
		parent.Style.MaxHeight = blitra.P(30)

		err := blitra.AvailableSizingVisitor(parent, nil)
		assert.Nil(t, err)

		assert.Equal(t, 50, child.AvailableSize.Width)
		assert.Equal(t, 30, child.AvailableSize.Height)
	})
}
