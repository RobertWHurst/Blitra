package blitra_test

import (
	"testing"

	"github.com/RobertWHurst/blitra"
	"github.com/stretchr/testify/assert"
)

func TestApplyWrap(t *testing.T) {
	t.Run("Returns an error if an invalid wrap mode is provided", func(t *testing.T) {
		_, _, err := blitra.ApplyWrap(42, false, blitra.Size{}, "Hello, World!")
		assert.Error(t, err)
	})
}

func TestApplyWordWrap(t *testing.T) {
	t.Run("Ensures the text does not exceed the maximum width", func(t *testing.T) {
		text := "Hello, World!"
		maxDimensions := blitra.Size{
			Width:  8,
			Height: 5,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.WordWrap, false, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "Hello,\nWorld!", wrappedText)
		assert.Equal(t, blitra.Size{Width: 6, Height: 2}, info.Size)
	})

	t.Run("Will try to preserve word boundaries", func(t *testing.T) {
		text := "I'm a little amazed everyday with the things people create in computer science."
		maxDimensions := blitra.Size{
			Width:  20,
			Height: 5,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.WordWrap, false, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "I'm a little amazed\neveryday with the\nthings people create\nin computer science.", wrappedText)
		assert.Equal(t, blitra.Size{Width: 20, Height: 4}, info.Size)
	})

	t.Run("Will split words that are too long", func(t *testing.T) {
		text := "This was totally Supercalifragilisticexpialidocious and I'm not sure how to handle it."
		maxDimensions := blitra.Size{
			Width:  10,
			Height: 20,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.WordWrap, false, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "This was\ntotally\nSupercali-\nfragilist-\nicexpiali-\ndocious\nand I'm\nnot sure\nhow to\nhandle it.", wrappedText)
		assert.Equal(t, blitra.Size{Width: 10, Height: 10}, info.Size)
	})

	t.Run("Do nothing if given an empty string", func(t *testing.T) {
		text := ""
		maxDimensions := blitra.Size{
			Width:  10,
			Height: 5,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.WordWrap, false, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "", wrappedText)
		assert.Equal(t, blitra.Size{}, info.Size)
	})

	t.Run("Will work with a single row and column", func(t *testing.T) {
		text := "Hello, World!"
		maxDimensions := blitra.Size{
			Width:  1,
			Height: 1,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.WordWrap, false, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "H", wrappedText)
		assert.Equal(t, blitra.Size{Width: 1, Height: 1}, info.Size)
	})

	t.Run("Will work with a single row and column ignoring ellipsis", func(t *testing.T) {
		text := "Hello, World!"
		maxDimensions := blitra.Size{
			Width:  1,
			Height: 1,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.WordWrap, true, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "H", wrappedText)
		assert.Equal(t, blitra.Size{Width: 1, Height: 1}, info.Size)
	})
}

func TestApplyCharacterWrap(t *testing.T) {
	t.Run("Ensures the text does not exceed the maximum width", func(t *testing.T) {
		text := "It's not as common to use character wrap, but it's still useful."
		maxDimensions := blitra.Size{
			Width:  10,
			Height: 15,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.CharacterWrap, false, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "It's not\nas common\nto use ch-\naracter\nwrap, but\nit's still\nuseful.", wrappedText)
		assert.Equal(t, blitra.Size{Width: 10, Height: 7}, info.Size)
	})
}

func TestApplyNoWrap(t *testing.T) {
	t.Run("Ensures the text does not exceed the maximum width", func(t *testing.T) {
		text := "Hello, World!"
		maxDimensions := blitra.Size{
			Width:  5,
			Height: 5,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.NoWrap, false, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "Hello", wrappedText)
		assert.Equal(t, blitra.Size{Width: 5, Height: 1}, info.Size)
	})

	t.Run("Will use an ellipsis if the text is too long and ellipsis is enabled", func(t *testing.T) {
		text := "Hello, World!"
		maxDimensions := blitra.Size{
			Width:  5,
			Height: 5,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.NoWrap, true, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "Hell…", wrappedText)
		assert.Equal(t, blitra.Size{Width: 5, Height: 1}, info.Size)
	})

	t.Run("Allows for explicit line breaks", func(t *testing.T) {
		text := "Hello,\nWorld!"
		maxDimensions := blitra.Size{
			Width:  5,
			Height: 5,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.NoWrap, false, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "Hello\nWorld", wrappedText)
		assert.Equal(t, blitra.Size{Width: 5, Height: 2}, info.Size)
	})

	t.Run("Will combine ellipsis and explicit line breaks", func(t *testing.T) {
		text := "Hello,\nWorld!"
		maxDimensions := blitra.Size{
			Width:  5,
			Height: 5,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.NoWrap, true, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "Hell…\nWorl…", wrappedText)
		assert.Equal(t, blitra.Size{Width: 5, Height: 2}, info.Size)
	})

	t.Run("Will not exceed the maximum height", func(t *testing.T) {
		text := "Hello,\nWorld!"
		maxDimensions := blitra.Size{
			Width:  10,
			Height: 1,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.NoWrap, false, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "Hello,", wrappedText)
		assert.Equal(t, blitra.Size{Width: 6, Height: 1}, info.Size)
	})

	t.Run("Will insert an ellipsis on the last line when the height is exceeded", func(t *testing.T) {
		text := "Hello,\nWorld!"
		maxDimensions := blitra.Size{
			Width:  10,
			Height: 1,
		}
		wrappedText, info, err := blitra.ApplyWrap(blitra.NoWrap, true, maxDimensions, text)
		assert.NoError(t, err)
		assert.Equal(t, "Hello,…", wrappedText)
		assert.Equal(t, blitra.Size{Width: 7, Height: 1}, info.Size)
	})
}
