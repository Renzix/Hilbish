package readline

import (
	"strings"
)

func (rl *Instance) emacsForwardWord(tokeniser tokeniser) (adjust int) {
	split, index, pos := tokeniser(rl.line, rl.pos)
	if len(split) == 0 {
		return
	}

	word := strings.TrimSpace(split[index])

	switch {
	case len(split) == 0:
		return
	case pos == len(word) && index != len(split)-1:
		extrawhitespace := len(strings.TrimLeft(split[index], " ")) - len(word)
		word = split[index+1]
		adjust = len(word) + extrawhitespace
	default:
		adjust = len(word) - pos
	}
	return
}

func (rl *Instance) emacsBackwardWord(tokeniser tokeniser) (adjust int) {
	split, index, pos := tokeniser(rl.line, rl.pos)
	if len(split) == 0 {
		return
	}

	switch {
	case len(split) == 0:
		return
	case pos == 0 && index != 0:
		adjust = len(split[index-1])
	default:
		adjust = pos
	}
	return
}

func (rl *Instance) ResetCount() {
	rl.count = 1
	rl.resetInfoText()
}
