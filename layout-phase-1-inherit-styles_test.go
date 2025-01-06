package blitra_test

import (
	"testing"

	"github.com/RobertWHurst/blitra"
	"github.com/stretchr/testify/assert"
)

func TestInheritStylesVisitor(t *testing.T) {
	t.Run("Applies an element's inheritable styles to it's children", func(t *testing.T) {
		parent := &blitra.Element{
			Style: blitra.Style{
				TextColor:       blitra.P("#f00"),
				BackgroundColor: blitra.P("#00f"),
				TextWrap:        blitra.P(blitra.CharacterWrap),
				Ellipsis:        blitra.P(true),
			},
		}
		child := &blitra.Element{
			Parent: parent,
			Style:  blitra.Style{},
		}

		err := blitra.InheritStylesVisitor(child, nil)
		assert.Nil(t, err)

		assert.Equal(t, parent.Style.TextColor, child.Style.TextColor)
		assert.Equal(t, parent.Style.BackgroundColor, child.Style.BackgroundColor)
		assert.Equal(t, parent.Style.TextWrap, child.Style.TextWrap)
		assert.Equal(t, parent.Style.Ellipsis, child.Style.Ellipsis)
	})
}
