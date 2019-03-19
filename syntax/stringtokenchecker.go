package syntax

import "strconv"

// StringTokenChecker is a token checker for string literals
type StringTokenChecker struct {
	QuoteChars []rune

	buffer       string
	rawBuffer    string
	isValid      bool
	afterEscape  bool
	openQuote    rune
	stringStart  bool
	stringEnd    bool
	beginUnicode bool
	unicodeHex   string
}

func (w *StringTokenChecker) Reset() {
	w.buffer = ""
	w.rawBuffer = ""
	w.isValid = true
	w.afterEscape = false
	w.beginUnicode = false
	w.stringStart = false
	w.stringEnd = false
	w.unicodeHex = ""
	if len(w.QuoteChars) == 0 {
		w.QuoteChars = []rune{'"'}
	}
}

func (w *StringTokenChecker) isAQuote(chr rune) bool {
	for _, c := range w.QuoteChars {
		if c == chr {
			return true
		}
	}
	return false
}

func (w *StringTokenChecker) isNumeric(chr rune, allowHex bool) bool {
	num := "0123456789"
	if allowHex {
		num += "abcdefABCDEF"
	}
	for _, c := range num {
		if c == chr {
			return true
		}
	}
	return false
}

func (w *StringTokenChecker) Feed(chr rune) (string, string, bool, bool) {
	if !w.isValid || w.stringEnd {
		return "", "", false, false
	}

	w.rawBuffer += string(chr)

	if !w.stringStart {
		if !w.isAQuote(chr) {
			w.isValid = false
			return "", "", false, false
		}
		w.openQuote = chr
		w.stringStart = true
		return w.rawBuffer, w.buffer, true, false
	}

	if w.afterEscape {
		if w.beginUnicode {
			if !w.isNumeric(chr, true) {
				w.isValid = false
				return w.rawBuffer, "", false, false
			}
			w.unicodeHex += string(chr)
			if len(w.unicodeHex) == 4 {
				w.beginUnicode = false
				w.afterEscape = false
				uni := "'\\u" + w.unicodeHex + "'"
				unic, err := strconv.Unquote(uni)
				if err == nil {
					w.buffer += unic
				} else {
					w.buffer += "?"
				}
			}
			return w.rawBuffer, w.buffer, true, false
		}
		if chr == w.openQuote || chr == '\\' || chr == '/' {
			w.buffer += string(chr)
			w.afterEscape = false
			return w.rawBuffer, w.buffer, true, false
		}
		if chr == 'n' {
			w.buffer += "\n"
			w.afterEscape = false
			return w.rawBuffer, w.buffer, true, false
		}
		if chr == 'r' {
			w.buffer += "\r"
			w.afterEscape = false
			return w.rawBuffer, w.buffer, true, false
		}
		if chr == 'b' {
			w.buffer += "\b"
			w.afterEscape = false
			return w.rawBuffer, w.buffer, true, false
		}
		if chr == 'f' {
			w.buffer += "\f"
			w.afterEscape = false
			return w.rawBuffer, w.buffer, true, false
		}
		if chr == 't' {
			w.buffer += "\t"
			w.afterEscape = false
			return w.rawBuffer, w.buffer, true, false
		}
		if chr == 'u' {
			w.beginUnicode = true
			w.unicodeHex = ""
			return w.rawBuffer, w.buffer, true, false
		}
		w.isValid = false
		return "", "", false, false
	}

	if chr == '\\' {
		w.afterEscape = true
		return w.rawBuffer, w.buffer, true, false
	}

	if chr == w.openQuote {
		w.stringEnd = true
		return w.rawBuffer, w.buffer, true, true
	}

	w.buffer += string(chr)
	return w.rawBuffer, w.buffer, true, false
}
