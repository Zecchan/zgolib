package grammar

import "github.com/zecchan/zgolib/syntax"

var (
	// CACFGCheckerSet contains a checker set that is used to parse CACFG
	CACFGCheckerSet = map[string]syntax.ITokenChecker{
		"ws": &syntax.WhitespaceTokenChecker{
			ExcludeNewline: true,
		},
		"nl": &syntax.NewlineTokenChecker{},
		"grdef": &syntax.SymbolTokenChecker{
			ValidSymbols: []string{"->"},
		},
		"otkn": &syntax.SymbolTokenChecker{
			ValidSymbols: []string{"<"},
		},
		"ctkn": &syntax.SymbolTokenChecker{
			ValidSymbols: []string{">"},
		},
		"oanc": &syntax.SymbolTokenChecker{
			ValidSymbols: []string{"("},
		},
		"canc": &syntax.SymbolTokenChecker{
			ValidSymbols: []string{")"},
		},
		"oflt": &syntax.SymbolTokenChecker{
			ValidSymbols: []string{"{"},
		},
		"cflt": &syntax.SymbolTokenChecker{
			ValidSymbols: []string{"}"},
		},
		"ref": &syntax.SymbolTokenChecker{
			ValidSymbols: []string{"*"},
		},
		"ident": &syntax.IdentifierTokenChecker{
			ValidFirstCharacters: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			ValidCharacters:      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_",
		},
		"comment": &syntax.CommentTokenChecker{
			AllowMultiline: false,
		},
	}
)
