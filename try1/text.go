package tui

type TextElement struct {
	layoutMeta LayoutMeta
	text       string
}

func (t TextElement) LayoutMeta() *LayoutMeta {
	return &t.layoutMeta
}

func (t TextElement) Children() []Element {
	return nil
}
