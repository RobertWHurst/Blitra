package blitra_test

import (
	blitra "_/home/robert/Developer/Personal/Blitra"
	"testing"
)

func TestElementFromAnyWithState(t *testing.T) {
	t.Run("Can take a string and wrap it in a element containing the text", func(t *testing.T) {
		element := blitra.ElementFromAnyWithState("Hello, World!", blitra.ViewState{})
		if element.Text != "Hello, World!" {
			t.Errorf("Expected element to contain text 'Hello, World!', but got '%s'", element.Text)
		}
	})
}
