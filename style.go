package blitra

// Style is a struct that contains all the style properties for an element.
type Style struct {
	Grow *bool

	Axis *Axis

	LeftPadding   *int
	RightPadding  *int
	TopPadding    *int
	BottomPadding *int
	Gap           *int

	LeftMargin   *int
	RightMargin  *int
	TopMargin    *int
	BottomMargin *int

	Width    *int
	MinWidth *int
	MaxWidth *int

	Height    *int
	MinHeight *int
	MaxHeight *int

	Align   *Align
	Justify *Justify

	LeftBorder   *Border
	RightBorder  *Border
	TopBorder    *Border
	BottomBorder *Border

	TextWrap *TextWrap
	Ellipsis *bool

	BackgroundColor *string
	TextColor       *string

	DEBUG_ID string
}
