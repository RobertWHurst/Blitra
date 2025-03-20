package blitra

import (
	"fmt"
)

func calcIntrinsicTextSize(el *Element) error {
	if el.Parent == nil {
		return nil
	}

	size := MaxSize
	if el.TextReflowWidth != nil {
		size.Width = *el.TextReflowWidth
	}
	_, wrapInfo, err := ApplyWrap(el.TextWrap(), false, size, el.SourceText)
	if err != nil {
		return fmt.Errorf("Failed to calculate intrinsic text size: %w", err)
	}
	el.IntrinsicSize = wrapInfo.Size
	return nil
}

func finalizeText(el *Element, reflow *bool) error {
	text, wrapInfo, err := ApplyWrap(el.TextWrap(), el.Ellipsis(), el.AvailableSize, el.SourceText)
	if err != nil {
		return fmt.Errorf("Failed to calculate available text size: %w", err)
	}

	if wrapInfo.IsVerticallyTruncated {
		if el.TextReflowWidth == nil || *el.TextReflowWidth != el.AvailableSize.Width {
			el.TextReflowWidth = &el.AvailableSize.Width
			*reflow = true
		}
	} else {
		el.Text = text
		el.TextReflowWidth = nil
	}

	return nil
}

func calcTextPosition(el *Element) error {
	if el.Parent == nil {
		return nil
	}

	el.Position.X = el.Parent.Position.X + el.Parent.LeftEdge()
	el.Position.Y = el.Parent.Position.Y + el.Parent.TopEdge()

	return nil
}
