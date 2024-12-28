package tui

type Element interface {
	LayoutMeta() *LayoutMeta
	Children() []Element
}
