package tui

type Renderable interface {
	Render() Element
}

func RenderRenderable(r any) []Renderable {

}
