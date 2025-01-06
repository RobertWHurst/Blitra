package blitra

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// Indicates how text should be wrapped.
type TextWrap int

const (
	// Will cause wrap to occur at word boundaries, unless a word is too long to
	// fit on a line. In that case, the word will be split with a hyphen.
	WordWrap TextWrap = iota
	// Will cause wrap to occur at character boundaries. Partial words will be
	// split with a hyphen.
	CharacterWrap
	// Will prevent all wrapping.
	NoWrap
)

type WrapInfo struct {
	Size        Size
	HasEllipsis bool
}

func ApplyWrap(mode TextWrap, useEllipsis bool, size Size, text string) (string, WrapInfo, error) {
	switch mode {
	case WordWrap:
		return ApplyWordOrCharWrap(true, useEllipsis, size, text)
	case CharacterWrap:
		return ApplyWordOrCharWrap(false, useEllipsis, size, text)
	case NoWrap:
		return ApplyNoWrap(useEllipsis, size, text)
	}
	return "", WrapInfo{}, fmt.Errorf("unknown TextWrap mode: %s", reflect.TypeOf(mode).String())
}

func ApplyWordOrCharWrap(useWordWrap bool, useEllipsis bool, size Size, text string) (string, WrapInfo, error) {
	maxWidth := size.Width
	maxHeight := size.Height
	if len(text) == 0 || maxWidth < 1 || maxHeight < 1 {
		return "", WrapInfo{}, nil
	}

	// Convert the text to a slice of runes, and append a newline rune
	// to force the last word to be processed.
	chars := []rune(text + "\n")

	// With word wrapping, the minimum partial word length is 3. With
	// character wrapping, the minimum partial word length is 2.
	var minPartialWordLen int
	if useWordWrap {
		minPartialWordLen = 3
	} else {
		minPartialWordLen = 2
	}

	// If the maximum width is to small to accommodate hyphenation or
	// ellipsis, disable them.
	useHyphens := true
	if maxWidth < minPartialWordLen {
		minPartialWordLen = maxWidth
	}
	if maxWidth < 2 {
		useEllipsis = false
		useHyphens = false
	}

	var (
		charIndex   int
		lines       [][]rune
		line        []rune
		word        []rune
		width       int
		hasEllipsis bool
	)

charLoop:
	for charIndex < len(chars) || len(word) != 0 {
		// Add the current word to one or more lines.
		for len(word) != 0 {
			// If the word fits in the limits of the current line, add it
			// and continue to the next word.
			if len(line) != 0 && len(line)+1+len(word) <= maxWidth || len(line) == 0 && len(word) <= maxWidth {
				if len(line) != 0 {
					line = append(line, ' ')
				}
				line = append(line, word...)
				word = []rune{}
				charIndex += 1
				if charIndex >= len(chars) {
					lineLen := len(line)
					if lineLen > width {
						width = lineLen
					}
					lines = append(lines, line)
					break charLoop
				}
				continue charLoop
			}

			// If this is the last possible line, add as much of the word as
			// possible.
			if len(lines) == maxHeight-1 {
				partialWordLen := maxWidth - len(line)
				if useEllipsis {
					partialWordLen -= 1
				}
				partialWord := word[0:partialWordLen]
				if useEllipsis {
					partialWord = append(partialWord, '…')
					hasEllipsis = true
				}
				if len(line) != 0 {
					line = append(line, ' ')
				}
				line = append(line, partialWord...)
				lineLen := len(line)
				if lineLen > width {
					width = lineLen
				}
				lines = append(lines, line)
				break charLoop
			}

			// If word wrapping is enabled, and the word is shorter than
			// the maximum width, start a new line.
			if useWordWrap && len(word) < maxWidth {
				lineLen := len(line)
				if lineLen > width {
					width = lineLen
				}
				lines = append(lines, line)
				line = []rune{}
				continue
			}

			// If the remaining space in the line is less than the minimum
			// partial word length, push the line and start a new one.
			if len(line) != 0 {
				availableLineLen := maxWidth - len(line)
				if useHyphens {
					availableLineLen -= 1
				}
				if len(line) != 0 {
					availableLineLen -= 1
				}
				if availableLineLen < minPartialWordLen {
					lineLen := len(line)
					if lineLen > width {
						width = lineLen
					}
					lines = append(lines, line)
					line = []rune{}
					continue
				}
			}

			// take a partial of the word, add it to the line, push the line,
			// start a new line
			partialWordLen := maxWidth - len(line)
			if useHyphens {
				partialWordLen -= 1
			}
			if len(line) != 0 {
				partialWordLen -= 1
			}
			if len(word) > minPartialWordLen*2 {
				for len(word)-partialWordLen < minPartialWordLen {
					partialWordLen -= 1
				}
			}
			partialWord := append([]rune{}, word[0:partialWordLen]...)
			word = word[partialWordLen:]
			if useHyphens {
				partialWord = append(partialWord, '-')
			}
			if len(line) != 0 {
				line = append(line, ' ')
			}
			line = append(line, partialWord...)
			lineLen := len(line)
			if lineLen > width {
				width = lineLen
			}
			lines = append(lines, line)
			line = []rune{}
		}

		// Collect the current word.
		for ; charIndex < len(chars); charIndex += 1 {
			char := chars[charIndex]
			if unicode.IsSpace(char) {
				continue charLoop
			}
			word = append(word, char)
		}

		if len(word) == 0 {
			charIndex += 1
		}
	}

	// Convert the lines (rune slices) to textLines (string slices).
	textLines := []string{}
	for _, line := range lines {
		textLines = append(textLines, string(line))
	}

	// Return the textLines as a single string, and the wrap info.
	return strings.Join(textLines, "\n"), WrapInfo{
		Size: Size{
			Width:  width,
			Height: len(textLines),
		},
		HasEllipsis: hasEllipsis,
	}, nil
}

func ApplyNoWrap(useEllipsis bool, size Size, text string) (string, WrapInfo, error) {
	maxWidth := size.Width
	maxHeight := size.Height
	if len(text) == 0 || maxWidth < 1 || maxHeight < 1 {
		return "", WrapInfo{}, nil
	}

	lines := [][]rune{}
	line := []rune{}
	width := 0
	inTailOfTruncatedLine := false
	lineHasEllipsis := false
	hasEllipsis := false
	for _, r := range text + "\n" {
		isLineBreak := r == '\n'

		// Start a new line if we reach a line break
		if isLineBreak {
			inTailOfTruncatedLine = false

			if len(lines) == maxHeight-1 {
				if useEllipsis && !lineHasEllipsis {
					line = append(line, '…')
					hasEllipsis = true
				}
				lineLen := len(line)
				if lineLen > width {
					width = lineLen
				}
				lines = append(lines, line)
				break
			}
			lineLen := len(line)
			if lineLen > width {
				width = lineLen
			}
			lines = append(lines, line)
			line = []rune{}
			lineHasEllipsis = false
			continue
		}

		// Skip the current rune if we are in the tail of a truncated line
		if inTailOfTruncatedLine {
			continue
		}

		// If we have exceeded the maximum number of columns, we will ignore the
		// rest of the text until we reach a line break.
		if len(line) == maxWidth-1 {
			inTailOfTruncatedLine = true
			if useEllipsis {
				line = append(line, '…')
				lineHasEllipsis = true
				hasEllipsis = true
				continue
			}
		}

		line = append(line, r)
	}

	textLines := []string{}
	for _, line := range lines {
		textLines = append(textLines, string(line))
	}

	return strings.Join(textLines, "\n"), WrapInfo{
		Size: Size{
			Width:  width,
			Height: len(textLines),
		},
		HasEllipsis: hasEllipsis,
	}, nil
}
