package grammar

import (
	"errors"
	"strconv"

	"github.com/zecchan/zgolib/syntax"
)

type Grammar struct {
	Definitions []GrammarDefinition
}

type GrammarDefinition struct {
	Name      string
	Structure []GrammarAtom
	Line      int
}

type GrammarAtom struct {
	Type GrammarAtomType
	Name string
	Line int
}

type GrammarAtomType int

const (
	AtomTypeToken      GrammarAtomType = 1
	AtomTypeSymbol     GrammarAtomType = 2
	AtomTypeSymbolRef  GrammarAtomType = 3
	AtomTypeFlatSymbol GrammarAtomType = 4
	AtomTypeAnchor     GrammarAtomType = 5
)

func (g *Grammar) Parse(script string) error {
	t := syntax.Tokenizer{}
	t.Checkers = CACFGCheckerSet
	t.IgnoreTokenTypes = []string{"ws"}
	tkns, err := t.Tokenize(script)
	if err != nil {
		return err
	}

	g.Definitions = []GrammarDefinition{}

	tbuf := []syntax.Token{}

	tkns = append(tkns, syntax.Token{Type: "nl"})

	for _, tkn := range tkns {
		if tkn.Type == "nl" {
			if len(tbuf) > 2 {
				gs, err := g.toDefinition(tbuf)
				if err != nil {
					return err
				}
				g.Definitions = append(g.Definitions, gs)
				tbuf = []syntax.Token{}
				continue
			}
			if len(tbuf) == 0 {
				tbuf = []syntax.Token{}
				continue
			}
			return errors.New("Line" + strconv.Itoa(tbuf[0].Line) + ": Expected definition after \"->\"")
		}
		if len(tbuf) == 0 {
			if tkn.Type == "comment" {
				continue
			}
			if tkn.Type != "ident" {
				return errors.New("Line" + strconv.Itoa(tkn.Line) + ": Expected identifier, \"" + tkn.RawValue + "\" found")
			}
		}
		if len(tbuf) == 1 && tkn.Type != "grdef" {
			return errors.New("Line" + strconv.Itoa(tkn.Line) + ": Expected \"->\" after \"" + tbuf[0].RawValue + "\", \"" + tkn.RawValue + "\" found")
		}
		if tkn.Type == "nl" {
			gs, err := g.toDefinition(tbuf)
			if err != nil {
				return err
			}
			g.Definitions = append(g.Definitions, gs)

			tbuf = []syntax.Token{}
			continue
		}
		tbuf = append(tbuf, tkn)
	}

	for _, def := range g.Definitions {
		for _, atm := range def.Structure {
			if atm.Type == AtomTypeFlatSymbol || atm.Type == AtomTypeSymbol || atm.Type == AtomTypeSymbolRef {
				if !g.HasDefinition(atm.Name) {
					return errors.New("Line " + strconv.Itoa(atm.Line) + ": Undefined symbol \"" + atm.Name + "\".")
				}
			}
		}
	}

	return nil
}

func (g *Grammar) toDefinition(tokens []syntax.Token) (GrammarDefinition, error) {
	gs := GrammarDefinition{}
	gs.Name = tokens[0].Value
	gs.Line = tokens[0].Line
	gs.Structure = []GrammarAtom{}
	err := gs.ParseTokens(tokens[2:len(tokens)])
	if err == nil {
		return gs, nil
	}
	return gs, err
}

func (g *Grammar) HasDefinition(name string) bool {
	for _, def := range g.Definitions {
		if def.Name == name {
			return true
		}
	}
	return false
}

func (d *GrammarDefinition) ParseTokens(tokens []syntax.Token) error {
	var state string = ""

	for _, tkn := range tokens {
		if state == "" {
			switch tkn.Type {
			case "ref":
				state = "reference"
			case "oflt":
				state = "flat_0"
			case "oanc":
				state = "anc_0"
			case "otkn":
				state = "tkn_0"
			case "ident":
				a := GrammarAtom{
					Name: tkn.Value,
					Type: AtomTypeSymbol,
					Line: tkn.Line,
				}
				d.Structure = append(d.Structure, a)
			default:
				return errors.New("Line " + strconv.Itoa(tkn.Line) + " Column " + strconv.Itoa(tkn.Column) + ": Unexpected token \"" + tkn.RawValue + "\"")
			}
		} else if state == "reference" {
			if tkn.Type != "ident" {
				return errors.New("Line " + strconv.Itoa(tkn.Line) + " Column " + strconv.Itoa(tkn.Column) + ": Expected an identifier, \"" + tkn.RawValue + "\" found.")
			}
			a := GrammarAtom{
				Name: tkn.Value,
				Type: AtomTypeSymbolRef,
				Line: tkn.Line,
			}
			d.Structure = append(d.Structure, a)
			state = ""
		} else if state == "flat_0" {
			if tkn.Type != "ident" {
				return errors.New("Line " + strconv.Itoa(tkn.Line) + " Column " + strconv.Itoa(tkn.Column) + ": Expected an identifier, \"" + tkn.RawValue + "\" found.")
			}
			a := GrammarAtom{
				Name: tkn.Value,
				Type: AtomTypeFlatSymbol,
				Line: tkn.Line,
			}
			d.Structure = append(d.Structure, a)
			state = "flat_1"
		} else if state == "flat_1" {
			if tkn.Type != "clft" {
				return errors.New("Line " + strconv.Itoa(tkn.Line) + " Column " + strconv.Itoa(tkn.Column) + ": Expecting \"}\", \"" + tkn.RawValue + "\" found.")
			}
			state = ""
		} else if state == "anc_0" {
			if tkn.Type != "ident" {
				return errors.New("Line " + strconv.Itoa(tkn.Line) + " Column " + strconv.Itoa(tkn.Column) + ": Expected an identifier, \"" + tkn.RawValue + "\" found.")
			}
			a := GrammarAtom{
				Name: tkn.Value,
				Type: AtomTypeAnchor,
				Line: tkn.Line,
			}
			d.Structure = append(d.Structure, a)
			state = "anc_1"
		} else if state == "anc_1" {
			if tkn.Type != "canc" {
				return errors.New("Line " + strconv.Itoa(tkn.Line) + " Column " + strconv.Itoa(tkn.Column) + ": Expecting \")\", \"" + tkn.RawValue + "\" found.")
			}
			state = ""
		} else if state == "tkn_0" {
			if tkn.Type != "ident" {
				return errors.New("Line " + strconv.Itoa(tkn.Line) + " Column " + strconv.Itoa(tkn.Column) + ": Expected an identifier, \"" + tkn.RawValue + "\" found.")
			}
			a := GrammarAtom{
				Name: tkn.Value,
				Type: AtomTypeToken,
				Line: tkn.Line,
			}
			d.Structure = append(d.Structure, a)
			state = "tkn_1"
		} else if state == "tkn_1" {
			if tkn.Type != "ctkn" {
				return errors.New("Line " + strconv.Itoa(tkn.Line) + " Column " + strconv.Itoa(tkn.Column) + ": Expecting \">\", \"" + tkn.RawValue + "\" found.")
			}
			state = ""
		} else if tkn.Type != "comment" {
			return errors.New("Line " + strconv.Itoa(tkn.Line) + " Column " + strconv.Itoa(tkn.Column) + ": Expected an identifier, \"" + tkn.RawValue + "\" found.")
		}
	}

	var refCount = 0
	for _, atom := range d.Structure {
		if atom.Type == AtomTypeSymbolRef {
			refCount = 1
		}
	}
	if refCount >= 1 {
		if len(d.Structure) > 1 {
			return errors.New("Line " + strconv.Itoa(d.Structure[0].Line) + ": A reference atom must be the only member of a definition.")
		}
	}
	if len(d.Structure) == 0 {
		return errors.New("Line " + strconv.Itoa(d.Structure[0].Line) + ": Grammar definition cannot be empty.")
	}
	if len(d.Structure) == 1 {
		var atm = d.Structure[0]
		if atm.Type == AtomTypeFlatSymbol || atm.Type == AtomTypeSymbol || atm.Type == AtomTypeSymbolRef {
			if atm.Name == d.Name {
				return errors.New("Line " + strconv.Itoa(d.Structure[0].Line) + ": Definition of \"" + d.Name + "\" cannot have only a single symbol that refer to itself.")
			}
		}
	}

	return nil
}
