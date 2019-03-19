package syntax

import (
	"errors"
	"strconv"
)

// Tokenizer enable you to tokenize a script
type Tokenizer struct {
	// Checkers is a token checker which type is described by the map key
	Checkers map[string]ITokenChecker

	IgnoreTokenTypes []string
}

func (t *Tokenizer) feedChar(char rune, line, col int) []Token {
	res := []Token{}
	for typ, chk := range t.Checkers {
		raw, val, valid, complete := chk.Feed(char)
		if valid {
			token := Token{
				Type:       typ,
				Value:      val,
				RawValue:   raw,
				IsValid:    valid,
				IsComplete: complete,
				Line:       line,
				Column:     col - len(raw) + 1,
			}
			res = append(res, token)
		}
	}
	return res
}
func (t *Tokenizer) resetCheckers() {
	for _, chk := range t.Checkers {
		chk.Reset()
	}
}

func (t *Tokenizer) ignored(token Token) bool {
	if t.IgnoreTokenTypes == nil {
		return false
	}
	for _, typ := range t.IgnoreTokenTypes {
		if typ == token.Type {
			return true
		}
	}
	return false
}

// Tokenize will parse specified script into Tokens
func (t *Tokenizer) Tokenize(script string) ([]Token, error) {
	line := 1
	col := 0
	ret := []Token{}
	var prevRes []Token

	script += "\a"
	t.resetCheckers()
	for ix, chr := range script {

		col++

		res := t.feedChar(chr, line, col)

		if len(res) > 0 {
			prevRes = res
		} else {
			if len(prevRes) == 1 {
				if !prevRes[0].IsComplete {
					return nil, errors.New("Line " + strconv.Itoa(line) + " Col " + strconv.Itoa(col) + ": Invalid character '" + string(chr) + "' after \"" + prevRes[0].RawValue + "\"")
				}
				if !t.ignored(prevRes[0]) {
					ret = append(ret, prevRes[0])
				}
				t.resetCheckers()
				prevRes = t.feedChar(chr, line, col)
				if len(prevRes) == 0 && ix != len(script)-1 {
					return nil, errors.New("Line " + strconv.Itoa(line) + " Col " + strconv.Itoa(col) + ": Unexpected character '" + string(chr) + "'")
				}
			} else if len(prevRes) > 1 {
				var completedToken Token
				var hasACompleteToken bool

				for _, tkn := range prevRes {
					if tkn.IsComplete {
						if hasACompleteToken {
							hasACompleteToken = false
							break
						}
						hasACompleteToken = true
						completedToken = tkn
					}
				}

				if !hasACompleteToken {
					return nil, errors.New("Line " + strconv.Itoa(line) + " Col " + strconv.Itoa(prevRes[0].Column) + ": Ambiguous token type for \"" + prevRes[0].RawValue + "\"")
				}

				if !t.ignored(completedToken) {
					ret = append(ret, completedToken)
				}
				t.resetCheckers()
				prevRes = t.feedChar(chr, line, col)
				if len(prevRes) == 0 && ix != len(script)-1 {
					return nil, errors.New("Line " + strconv.Itoa(line) + " Col " + strconv.Itoa(col) + ": Unexpected character '" + string(chr) + "'")
				}
			} else {
				return nil, errors.New("Line " + strconv.Itoa(line) + " Col " + strconv.Itoa(col) + ": Unexpected character '" + string(chr) + "'")
			}
		}

		if chr == '\n' {
			line++
			col = 0
		}
	}
	return ret, nil
}

// ITokenChecker must implement a checker function that is used to determine whether a substring is valid for a given token
type ITokenChecker interface {
	Feed(chr rune) (string, string, bool, bool)
	Reset()
}

// Token is a token data represented by its Type and Value
type Token struct {
	Value      string
	RawValue   string
	Type       string
	IsValid    bool
	IsComplete bool
	Line       int
	Column     int
}

// NewlineTokenChecker is a token checker for newlines
type NewlineTokenChecker struct {
	ExcludeTab     bool
	ExcludeSpace   bool
	ExcludeNewline bool

	buffer    string
	rawBuffer string
	isValid   bool
	isEnded   bool
}

func (w *NewlineTokenChecker) Reset() {
	w.buffer = ""
	w.rawBuffer = ""
	w.isValid = true
	w.isEnded = false
}
func (w *NewlineTokenChecker) Feed(chr rune) (string, string, bool, bool) {
	if !w.isValid || w.isEnded {
		return "", "", false, false
	}

	w.rawBuffer += string(chr)

	if w.buffer == "" {
		if chr != '\r' && chr != '\n' {
			w.isValid = false
			return "", "", false, false
		}
		w.buffer += string(chr)
		return w.rawBuffer, w.buffer, true, chr == '\n'
	}
	if chr != '\n' {
		w.isValid = false
		return "", "", false, false
	}
	w.buffer += string(chr)
	return w.rawBuffer, w.buffer, true, true
}

// WhitespaceTokenChecker is a token checker for whitespaces
type WhitespaceTokenChecker struct {
	ExcludeTab     bool
	ExcludeSpace   bool
	ExcludeNewline bool

	buffer    string
	rawBuffer string
	isValid   bool
}

func (w *WhitespaceTokenChecker) Reset() {
	w.buffer = ""
	w.rawBuffer = ""
	w.isValid = true
}

func (w *WhitespaceTokenChecker) Feed(chr rune) (string, string, bool, bool) {
	if !w.isValid {
		return "", "", false, false
	}

	w.rawBuffer += string(chr)

	if !w.ExcludeTab && chr == '\t' {
		w.buffer += string(chr)
		return w.rawBuffer, w.buffer, true, true
	}

	if !w.ExcludeSpace && chr == ' ' {
		w.buffer += string(chr)
		return w.rawBuffer, w.buffer, true, true
	}

	if !w.ExcludeNewline && (chr == '\n' || chr == '\r') {
		w.buffer += string(chr)
		return w.rawBuffer, w.buffer, true, true
	}

	w.isValid = false
	return "", "", false, false
}

// SymbolTokenChecker is a token checker for symbols
type SymbolTokenChecker struct {
	ValidSymbols []string
	buffer       string
	rawBuffer    string
	isValid      bool
}

func (w *SymbolTokenChecker) Reset() {
	w.buffer = ""
	w.rawBuffer = ""
	w.isValid = true
	if w.ValidSymbols == nil {
		w.ValidSymbols = []string{}
	}
}

func (w *SymbolTokenChecker) Feed(chr rune) (string, string, bool, bool) {
	if !w.isValid {
		return "", "", false, false
	}
	w.rawBuffer += string(chr)

	w.buffer += string(chr)
	subvalid := ""
	for _, sym := range w.ValidSymbols {
		if sym == w.buffer {
			return w.rawBuffer, sym, true, true
		}
		if len(sym) > len(w.buffer) {
			subsym := sym[0:len(w.buffer)]
			if subsym == w.buffer {
				subvalid = subsym
			}
		}
	}
	if subvalid != "" {
		return w.rawBuffer, subvalid, true, false
	}

	w.isValid = false
	return "", "", false, false
}
