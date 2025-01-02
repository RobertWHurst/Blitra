package blitra_test

import (
	"testing"

	blitra "github.com/RobertWHurst/blitra"
)

func TestApplyWrap(t *testing.T) {
	t.Run("Panics if an invalid wrap mode is given", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected to panic, but did not")
			}
		}()
		blitra.ApplyWrap(42, false, blitra.Size{}, "Hello, World!")
	})
}

func TestApplyWordWrap(t *testing.T) {
	t.Run("Ensures the text does not exceed the maximum width", func(t *testing.T) {
		text := "Hello, World!"
		maxDimensions := blitra.Size{
			Width: 8,
		}
		result, _ := blitra.ApplyWrap(blitra.WordWrap, false, maxDimensions, text)
		if result != "Hello,\nWorld!" {
			t.Errorf("Expected result to be 'Hello,\nWorld!', but got '%s'", result)
		}
	})

	t.Run("Will try to preserve word boundaries", func(t *testing.T) {
		text := "I'm a little amazed everyday with the things people create in computer science."
		maxDimensions := blitra.Size{
			Width: 20,
		}
		result, _ := blitra.ApplyWrap(blitra.WordWrap, false, maxDimensions, text)
		expected := "I'm a little amazed\n" +
			"everyday with the\n" +
			"things people create\n" +
			"in computer science."
		if result != expected {
			t.Errorf("Expected result to be '"+expected+"', but got '%s'", result)
		}
	})

	t.Run("Will split words that are too long", func(t *testing.T) {
		text := "This was totally Supercalifragilisticexpialidocious and I'm not sure how to handle it."
		maxDimensions := blitra.Size{
			Width: 10,
		}
		result, _ := blitra.ApplyWrap(blitra.WordWrap, false, maxDimensions, text)
		expected := "This was\n" +
			"totally\n" +
			"Supercali-\n" +
			"fragilist-\n" +
			"icexpiali-\n" +
			"docious\n" +
			"and I'm\n" +
			"not sure\n" +
			"how to\n" +
			"handle it."
		if result != expected {
			t.Errorf("Expected result to be '"+expected+"', but got '%s'", result)
		}
	})

	t.Run("Will work with an empty string", func(t *testing.T) {
		text := ""
		maxDimensions := blitra.Size{
			Width: 10,
		}
		result, _ := blitra.ApplyWrap(blitra.WordWrap, false, maxDimensions, text)
		if result != "" {
			t.Errorf("Expected result to be '', but got '%s'", result)
		}
	})

	t.Run("Will work with a single row and column", func(t *testing.T) {
		text := "Hello, World!"
		maxDimensions := blitra.Size{
			Width:  1,
			Height: 1,
		}
		result, _ := blitra.ApplyWrap(blitra.WordWrap, false, maxDimensions, text)
		if result != "H" {
			t.Errorf("Expected result to be 'H', but got '%s'", result)
		}
	})

	t.Run("Will work with a single row and column ignoring ellipsis", func(t *testing.T) {
		text := "Hello, World!"
		maxDimensions := blitra.Size{
			Width:  1,
			Height: 1,
		}
		result, _ := blitra.ApplyWrap(blitra.WordWrap, true, maxDimensions, text)
		if result != "H" {
			t.Errorf("Expected result to be 'H', but got '%s'", result)
		}
	})
}

func TestApplyCharacterWrap(t *testing.T) {
	t.Run("Ensures the text does not exceed the maximum width", func(t *testing.T) {
		text := "It's not as common to use character wrap, but it's still useful."
		maxDimensions := blitra.Size{
			Width: 10,
		}
		result, _ := blitra.ApplyWrap(blitra.CharacterWrap, false, maxDimensions, text)
		expected := "It's not\n" +
			"as common\n" +
			"to use ch-\n" +
			"aracter\n" +
			"wrap, but\n" +
			"it's still\n" +
			"useful."
		if result != expected {
			t.Errorf("Expected result to be '"+expected+"', but got '%s'", result)
		}
	})

}

func TestApplyNoWrap(t *testing.T) {
	t.Run("Ensures the text does not exceed the maximum width", func(t *testing.T) {
		text := "Hello, World!"
		maxDimensions := blitra.Size{
			Width: 5,
		}
		result, _ := blitra.ApplyWrap(blitra.NoWrap, false, maxDimensions, text)
		if result != "Hello" {
			t.Errorf("Expected result to be 'Hello', but got '%s'", result)
		}
	})

	t.Run("Will use an ellipsis if the text is too long and ellipsis is enabled", func(t *testing.T) {
		text := "Hello, World!"
		maxDimensions := blitra.Size{
			Width: 5,
		}
		result, _ := blitra.ApplyWrap(blitra.NoWrap, true, maxDimensions, text)
		if result != "Hell…" {
			t.Errorf("Expected result to be 'Hell…', but got '%s'", result)
		}
	})

	t.Run("Allows for explicit line breaks", func(t *testing.T) {
		text := "Hello,\nWorld!"
		maxDimensions := blitra.Size{
			Width: 5,
		}
		result, _ := blitra.ApplyWrap(blitra.NoWrap, false, maxDimensions, text)
		if result != "Hello\nWorld" {
			t.Errorf("Expected result to be 'Hello\nWorld', but got '%s'", result)
		}
	})

	t.Run("Will combine ellipsis and explicit line breaks", func(t *testing.T) {
		text := "Hello,\nWorld!"
		maxDimensions := blitra.Size{
			Width: 5,
		}
		result, _ := blitra.ApplyWrap(blitra.NoWrap, true, maxDimensions, text)
		if result != "Hell…\nWorl…" {
			t.Errorf("Expected result to be 'Hell…\nWorl…', but got '%s'", result)
		}
	})

	t.Run("Will not exceed the maximum height", func(t *testing.T) {
		text := "Hello,\nWorld!"
		maxDimensions := blitra.Size{
			Width:  10,
			Height: 1,
		}
		result, _ := blitra.ApplyWrap(blitra.NoWrap, false, maxDimensions, text)
		if result != "Hello," {
			t.Errorf("Expected result to be 'Hello,', but got '%s'", result)
		}
	})

	t.Run("Will not exceed the maximum height even without line breaks", func(t *testing.T) {
		text := "Hello,\nWorld!"
		maxDimensions := blitra.Size{
			Height: 1,
		}
		result, _ := blitra.ApplyWrap(blitra.NoWrap, false, maxDimensions, text)
		if result != "Hello," {
			t.Errorf("Expected result to be 'Hello,', but got '%s'", result)
		}
	})

	t.Run("Will insert an ellipsis on the last line when the height is exceeded", func(t *testing.T) {
		text := "Hello,\nWorld!"
		maxDimensions := blitra.Size{
			Width:  10,
			Height: 1,
		}
		result, _ := blitra.ApplyWrap(blitra.NoWrap, true, maxDimensions, text)
		if result != "Hello,…" {
			t.Errorf("Expected result to be 'Hello,…', but got '%s'", result)
		}
	})
}
