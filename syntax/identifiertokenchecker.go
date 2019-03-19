package syntax

import "strings"

// IdentifierTokenChecker is a token checker for identifiers
type IdentifierTokenChecker struct {
	buffer    string
	rawBuffer string
	isValid   bool

	ValidFirstCharacters string
	ValidCharacters      string
}

func (w *IdentifierTokenChecker) Reset() {
	w.buffer = ""
	w.rawBuffer = ""
	w.isValid = true
}

func (w *IdentifierTokenChecker) Feed(chr rune) (string, string, bool, bool) {
	if !w.isValid {
		return "", "", false, false
	}
	w.rawBuffer += string(chr)

	if w.buffer == "" {
		if !strings.Contains(w.ValidFirstCharacters, string(chr)) {
			w.isValid = false
			return "", "", false, false
		}
	} else {
		if !strings.Contains(w.ValidCharacters, string(chr)) {
			w.isValid = false
			return "", "", false, false
		}
	}

	w.buffer += string(chr)
	return w.rawBuffer, w.buffer, true, true
}
