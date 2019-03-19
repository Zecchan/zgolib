package syntax

// NumberTokenChecker is a token checker for numeric representation
type NumberTokenChecker struct {
	buffer    string
	rawBuffer string
	isValid   bool

	hasComma    bool
	numberFound bool
	eFound      bool
	afterE      bool
}

func (w *NumberTokenChecker) Reset() {
	w.buffer = ""
	w.rawBuffer = ""
	w.isValid = true
	w.hasComma = false
	w.numberFound = false
	w.eFound = false
	w.afterE = false
}

func (w *NumberTokenChecker) isDigit(chr rune) bool {
	num := "0123456789"
	for _, c := range num {
		if c == chr {
			return true
		}
	}
	return false
}

func (w *NumberTokenChecker) Feed(chr rune) (string, string, bool, bool) {
	if !w.isValid {
		return "", "", false, false
	}

	w.rawBuffer += string(chr)

	if chr == '-' && w.buffer != "" {
		w.isValid = false
		return "", "", false, false
	}

	if chr == '.' && (w.hasComma || !w.numberFound) {
		w.isValid = false
		return "", "", false, false
	}

	if w.afterE {
		if chr == '-' || chr == '+' {
			w.buffer += string(chr)
			w.afterE = false
			return w.rawBuffer, w.buffer, true, false
		}
		w.isValid = false
		return "", "", false, false
	}

	if w.isDigit(chr) {
		w.numberFound = true
		w.buffer += string(chr)
		return w.rawBuffer, w.buffer, true, true
	}

	if chr == '-' || chr == '.' {
		w.buffer += string(chr)
		return w.rawBuffer, w.buffer, true, false
	}

	if (chr == 'e' || chr == 'E') && !w.eFound {
		w.eFound = true
		w.afterE = true
		w.buffer += string(chr)
		return w.rawBuffer, w.buffer, true, false
	}

	w.isValid = false
	return "", "", false, false
}
