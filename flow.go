package blitra

import "fmt"

func Flow(el *Element) error {
	reflow := true
	for reflow {
		reflow = false
		if err := VisitElementsUp(el, nil, intrinsicSizeVisitor); err != nil {
			return fmt.Errorf("Failed to calculate intrinsic sizing: %w", err)
		}
		if err := VisitElementsDown(el, &reflow, availableSizeVisitor); err != nil {
			return fmt.Errorf("Failed to calculate available sizing: %w", err)
		}
	}
	if err := VisitElementsDown(el, nil, positioningVisitor); err != nil {
		return fmt.Errorf("Failed to calculate positioning: %w", err)
	}
	return nil
}

func intrinsicSizeVisitor(el *Element, _ any) error {
	switch el.Kind {
	case TextElementKind:
		return calcIntrinsicTextSize(el)
	case ContainerElementKind:
		return calcIntrinsicContainerSize(el)
	default:
		return fmt.Errorf("unknown element kind: %v", el.Kind)
	}
}

func availableSizeVisitor(el *Element, reflow *bool) error {
	switch el.Kind {
	case TextElementKind:
		return finalizeText(el, reflow)
	case ContainerElementKind:
		return calcAvailableContainerSizesForChildren(el)
	default:
		return fmt.Errorf("unknown element kind: %v", el.Kind)
	}
}

func positioningVisitor(el *Element, _ any) error {
	switch el.Kind {
	case TextElementKind:
		return calcTextPosition(el)
	case ContainerElementKind:
		return calcContainerPositionsForChildren(el)
	default:
		return fmt.Errorf("unknown element kind: %v", el.Kind)
	}
}
