package syntax

var (
	// JSONCheckerSet contains a checker set that is used to parse JSON
	JSONCheckerSet = map[string]ITokenChecker{
		"ws": &WhitespaceTokenChecker{},
		"colon": &SymbolTokenChecker{
			ValidSymbols: []string{":"},
		},
		"comma": &SymbolTokenChecker{
			ValidSymbols: []string{","},
		},
		"oobj": &SymbolTokenChecker{
			ValidSymbols: []string{"{"},
		},
		"cobj": &SymbolTokenChecker{
			ValidSymbols: []string{"}"},
		},
		"oarr": &SymbolTokenChecker{
			ValidSymbols: []string{"["},
		},
		"carr": &SymbolTokenChecker{
			ValidSymbols: []string{"]"},
		},
		"bool": &SymbolTokenChecker{
			ValidSymbols: []string{"true", "false"},
		},
		"null": &SymbolTokenChecker{
			ValidSymbols: []string{"null"},
		},
		"strlit": &StringTokenChecker{
			QuoteChars: []rune{'"'},
		},
		"numlit": &NumberTokenChecker{},
	}
)
