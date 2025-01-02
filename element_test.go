package blitra_test

import (
	"testing"

	blitra "github.com/RobertWHurst/blitra"
)

func TestElementFromAnyWithState(t *testing.T) {
	// t.Run("Can take a string and wrap it in a element containing the text", func(t *testing.T) {
	// 	element := blitra.BuildElementTreeFromRenderable("Hello, World!", blitra.ViewState{})
	// 	if *element.Text != "Hello, World!" {
	// 		t.Errorf("Expected element to contain text 'Hello, World!', but got '%s'", *element.Text)
	// 	}
	// })

	// t.Run("Can take a slice of strings and wrap them in a list element containing elements with the text", func(t *testing.T) {
	// 	element := blitra.BuildElementTreeFromRenderable([]string{"Hello", "World!"}, blitra.ViewState{})
	// 	if len(element.Children) != 2 {
	// 		t.Errorf("Expected element to contain 2 children, but got %d", len(element.Children))
	// 	}
	// 	if *element.Children[0].Text != "Hello" {
	// 		t.Errorf("Expected first child element to contain text 'Hello', but got '%s'", *element.Children[0].Text)
	// 	}
	// 	if *element.Children[1].Text != "World!" {
	// 		t.Errorf("Expected second child element to contain text 'World!', but got '%s'", *element.Children[1].Text)
	// 	}
	// })

	// t.Run("Can take a struct that implements the Renderable interface and wrap the returned value in a child element", func(t *testing.T) {
	// 	element := blitra.BuildElementTreeFromRenderable(TestRenderable{ReturnValue: "Hello World!"}, blitra.ViewState{})
	// 	if len(element.Children) != 1 {
	// 		t.Errorf("Expected element to contain 1 child, but got %d", len(element.Children))
	// 	}
	// 	if *element.Children[0].Text != "Hello World!" {
	// 		t.Errorf("Expected child element to contain text 'Hello World!', but got '%s'", *element.Children[0].Text)
	// 	}
	// })

	// t.Run("Can take a pointer to a struct that implements the Renderable interface", func(t *testing.T) {
	// 	element := blitra.BuildElementTreeFromRenderable(&TestRenderable{ReturnValue: "Hello World!"}, blitra.ViewState{})
	// 	if len(element.Children) != 1 {
	// 		t.Errorf("Expected element to contain 1 child, but got %d", len(element.Children))
	// 	}
	// 	if *element.Children[0].Text != "Hello World!" {
	// 		t.Errorf("Expected child element to contain text 'Hello World!', but got '%s'", *element.Children[0].Text)
	// 	}
	// })

	// t.Run("Panics if the struct does not implement the Renderable interface", func(t *testing.T) {
	// 	defer func() {
	// 		if r := recover(); r == nil {
	// 			t.Errorf("Expected to panic, but did not")
	// 		}
	// 	}()
	// 	_ = blitra.BuildElementTreeFromRenderable(struct{}{}, blitra.ViewState{})
	// })

	// t.Run("Can take typed pointer to nil and return nil", func(t *testing.T) {
	// 	var renderable *TestRenderable
	// 	element := blitra.BuildElementTreeFromRenderable(renderable, blitra.ViewState{})
	// 	if element != nil {
	// 		t.Errorf("Expected element to be nil, but got %+v", element)
	// 	}
	// })

	// t.Run("Can take nil and return nil", func(t *testing.T) {
	// 	element := blitra.BuildElementTreeFromRenderable(nil, blitra.ViewState{})
	// 	if element != nil {
	// 		t.Errorf("Expected element to be nil, but got %+v", element)
	// 	}
	// })

	// t.Run("Can create a tree of elements from a more complex structure", func(t *testing.T) {

	// 	child21 := TestRenderable{ReturnValue: "Child 2.1"}
	// 	child11 := TestRenderable{ReturnValue: child21}
	// 	child12 := TestRenderable{ReturnValue: "Child 1.2"}
	// 	root := TestRenderable{ReturnValue: []any{child11, child12}}

	// 	element := blitra.BuildElementTreeFromRenderable(root, blitra.ViewState{})

	// 	if len(element.Children) != 1 {
	// 		t.Errorf("Expected element to contain 1 child, but got %d", len(element.Children))
	// 	}
	// 	rootElement := element.Children[0]

	// 	if len(rootElement.Children) != 2 {
	// 		t.Errorf("Expected element to contain 2 children, but got %d", len(element.Children))
	// 	}
	// 	child11Element := rootElement.Children[0]

	// 	if len(child11Element.Children) != 1 {
	// 		t.Errorf("Expected element to contain 1 child, but got %d", len(child11Element.Children))
	// 	}
	// 	child21Element := child11Element.Children[0]

	// 	if len(child21Element.Children) != 1 {
	// 		t.Errorf("Expected element to contain 1 child, but got %d", len(child21Element.Children))
	// 	}
	// 	child21TextElement := child21Element.Children[0]
	// 	if *child21TextElement.Text != "Child 2.1" {
	// 		t.Errorf("Expected element to contain text 'Child 2.1', but got '%s'", *child21TextElement.Text)
	// 	}

	// 	child12Element := rootElement.Children[1]
	// 	if len(child12Element.Children) != 1 {
	// 		t.Errorf("Expected element to contain 1 child, but got %d", len(child12Element.Children))
	// 	}

	// 	child12TextElement := child12Element.Children[0]
	// 	if *child12TextElement.Text != "Child 1.2" {
	// 		t.Errorf("Expected element to contain text 'Child 1.2', but got '%s'", *child12TextElement.Text)
	// 	}
	// })
}

type TestRenderable struct {
	ReturnValue any
}

var _ blitra.Renderable = TestRenderable{}

func (r TestRenderable) ID() string {
	return "TestRenderable"
}

func (r TestRenderable) Style() blitra.Style {
	return blitra.Style{}
}

func (r TestRenderable) Render(state blitra.ViewState) any {
	return r.ReturnValue
}
