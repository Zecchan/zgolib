package syntax

// CommentTokenChecker is a token checker for identifiers
type CommentTokenChecker struct {
	buffer    string
	rawBuffer string
	isValid   bool

	isMultiline        bool
	expectEndMultiline bool
	isCommentEnded     bool
	AllowMultiline     bool
}

func (w *CommentTokenChecker) Reset() {
	w.buffer = ""
	w.rawBuffer = ""
	w.isValid = true
	w.isMultiline = false
	w.expectEndMultiline = false
	w.isCommentEnded = false
}

func (w *CommentTokenChecker) Feed(chr rune) (string, string, bool, bool) {
	if !w.isValid || w.isCommentEnded {
		return "", "", false, false
	}
	if w.rawBuffer == "/" {
		w.rawBuffer += string(chr)
		if chr != '/' && chr != '*' {
			w.isValid = false
			return "", "", false, false
		}
		if chr == '*' && !w.AllowMultiline {
			w.isValid = false
			return "", "", false, false
		}
		w.isMultiline = chr == '*'
		return w.rawBuffer, w.buffer, true, false
	}

	if w.rawBuffer == "" {
		w.rawBuffer += string(chr)
		if chr != '/' {
			w.isValid = false
			return "", "", false, false
		}
		return w.rawBuffer, w.buffer, true, false
	}

	w.rawBuffer += string(chr)

	if chr == '\n' && !w.isMultiline {
		w.isCommentEnded = true
		return w.rawBuffer, w.buffer, true, true
	}

	if w.expectEndMultiline {
		if chr == '/' {
			w.buffer = w.buffer[0 : len(w.buffer)-1]
			w.isCommentEnded = true
			return w.rawBuffer, w.buffer, true, true
		}
		w.expectEndMultiline = false
	}
	if chr == '*' {
		w.expectEndMultiline = true
	}

	w.buffer += string(chr)
	return w.rawBuffer, w.buffer, true, false
}
