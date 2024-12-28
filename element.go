package blitra

import "reflect"

type Element struct {
	Parent   *Element
	Previous *Element
	Next     *Element
	Children []*Element

	Text string
}

func ElementFromAnyWithState(unknown any, state ViewState) *Element {
	unknownVal := reflect.ValueOf(unknown)

	switch unknownVal.Kind() {

	// For slices we create a list element which will contain all the
	// elements in the slice. We then iterate over the slice and create
	// elements for each item in the slice while establishing the correct
	// relationships between the elements.
	case reflect.Slice:
		listElement := &Element{}
		var previousElement *Element
		for i := 0; i < unknownVal.Len(); i++ {
			unknownItem := unknownVal.Index(i).Interface()
			element := ElementFromAnyWithState(unknownItem, state)
			if element != nil {
				element.Previous = previousElement
				previousElement.Next = element
				previousElement = element
				listElement.Children = append(listElement.Children, element)
			}
		}
		return listElement

	// If the unknown value is a struct, then we will check if it implements
	// the Renderable interface. If it does, we will call the Render method
	// and create an element from the returned value.
	case reflect.Struct:
		renderable, ok := unknown.(Renderable)
		if !ok {
			panic("Struct type does not implement the Renderable interface: " + unknownVal.Type().String())
		}
		element := &Element{}
		childElement := ElementFromAnyWithState(renderable.Render(state), state)
		if childElement != nil {
			childElement.Parent = element
			element.Children = append(element.Children, childElement)
		}
		return element

	// If the unknown value is a pointer, then we traverse the pointer to
	// get the underlying value and create an element from it. If the pointer
	// is nil, then we return nil.
	case reflect.Ptr:
		if unknownVal.IsNil() {
			return nil
		}
		return ElementFromAnyWithState(unknownVal.Elem().Interface(), state)

	// If the unknown value is a string, then we create a text element with
	// the string value.
	case reflect.String:
		return &Element{
			Text: unknownVal.String(),
		}
	}

	// If we fell through the switch statement, then we encountered an
	// unsupported type.
	panic("Tried to render an unsupported type: " + unknownVal.Kind().String())
}
