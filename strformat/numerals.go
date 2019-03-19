package strformat

import (
	"math"
	"strconv"
	"strings"
)

/*
	StrFormat package -> Numerals
	ver 1.0 - 2019-03-18
	by Zecchan Silverlake

	This package contains useful function to print number to text
*/

type Numeral struct {
	// SplitDigit is how many digits will be taken until the conversion repeat
	SplitDigit int
	// Conversion will convert number to text. Value is taken from a digit, example: 1 -> one. Only works for 1-9
	Conversion map[int]string
	// ZeroConversion converts zero value (0) to this string
	ZeroConversion string
	// LiteralConversion will convert number to text if value is found in a group or group mod operation, example: 11 -> eleven
	LiteralConversion map[int]string
	// DigitNames is the name of the digits, example: 100 -> hundred (group at mod 100 is hundred)
	DigitNames map[int]string
	// GroupNames is the name of the group based on splitted digits, example: 1 -> thousand (group at index 1 is thousand)
	GroupNames map[int]string
	// PointConversion is the name of the point
	PointConversion string
	// Correction will correct substring into mapped string, example: two ty -> twenty
	Correction map[string]string
	// CurrencyName is the name of the currency
	CurrencyName string
	// CurrencyPointName is the name of the currency decimal point, example: cent
	CurrencyPointName string
	// CurrencyPointConversion is the joiner of the currency decimal point, example: and
	CurrencyPointConversion string
	// CurrencyPointLength is the length of the currency decimal point
	CurrencyPointLength int
}

func (n *Numeral) ConvertCurrency(value float64) string {
	strVal := strconv.FormatFloat(value, 'f', n.CurrencyPointLength, 64)
	spl := strings.Split(strVal, ".")
	digits := spl[0]
	points := ""
	if len(spl) == 2 {
		points = spl[1]
	}

	digitInt, e := strconv.Atoi(digits)
	res := ""
	if e == nil {
		res = n.Convert(float64(digitInt), 0)
	}
	res = strings.Trim(res, " \t")

	if points != "" {
		points = points[0:int(math.Min(float64(len(points)), float64(n.CurrencyPointLength)))]
		for len(points) < n.CurrencyPointLength {
			points += "0"
		}
		ptInt, e := strconv.Atoi(points)
		if e == nil && ptInt != 0 {
			res += " " + n.CurrencyPointConversion + " " + n.Convert(float64(ptInt), 0)
			if n.CurrencyPointName != "" {
				res += " " + n.CurrencyPointName
			}
		}
	}

	if n.CurrencyName != "" {
		res += " " + n.CurrencyName
	}

	return res + n.CurrencyName
}

func (n *Numeral) Convert(value float64, prec int) string {
	strVal := strconv.FormatFloat(value, 'f', prec, 64)
	spl := strings.Split(strVal, ".")
	digits := spl[0]
	points := ""
	if len(spl) == 2 {
		points = spl[1]
	}

	res := ""
	group := ""
	gIdx := 0
	for i := len(digits) - 1; i >= 0; i-- {
		group = digits[i:i+1] + group
		if len(group) == n.SplitDigit || i == 0 {
			groupStr := n.groupConvert(group)
			groupName, ok := n.GroupNames[gIdx]
			if ok && groupName != "" && groupStr != "" {
				groupStr += " " + groupName
			}
			if groupStr != "" {
				res = groupStr + " " + res
			}
			group = ""
			gIdx++
		}
	}
	res = strings.Trim(res, " \t")

	if points != "" {
		ptWord := ""
		for i := 0; i < len(points); i++ {
			pt, e := strconv.Atoi(points[i : i+1])
			if e == nil {
				ptWord += " " + n.Conversion[pt]
			} else {
				ptWord += " " + n.ZeroConversion
			}
		}

		ptWord = strings.Trim(ptWord, " ")
		if ptWord != "" {
			res += " " + n.PointConversion + " " + ptWord
		}
	}
	return res
}

func (n *Numeral) groupConvert(group string) string {
	res := ""
	grVal, e := strconv.Atoi(group)
	if e != nil {
		return res
	}
	for i := n.SplitDigit - 1; i >= 1; i-- {
		base := int(math.Pow(10, float64(i)))
		rem := grVal % base
		quo := (grVal / base) % 10
		quoLit := quo * base

		// left side
		ql, ok := n.LiteralConversion[quoLit]
		if ok {
			res += " " + ql
		} else {
			if quo != 0 {
				res += " " + n.Conversion[quo] + " " + n.DigitNames[base]
			}
		}

		// right side
		if i == 1 {
			res += " " + n.Conversion[rem]
		} else {
			v, ok := n.LiteralConversion[rem]
			if ok {
				res += " " + v
				break
			}
		}
	}

	for key, cor := range n.Correction {
		res = strings.Replace(res, key, cor, -1)
	}

	return strings.Trim(res, " \t")
}

// NumeralCreateIndonesian creates numeral struct for Indonesian language
func NumeralCreateIndonesian() *Numeral {
	num := Numeral{
		SplitDigit:     3,
		ZeroConversion: "nol",
		Conversion: map[int]string{
			1: "satu",
			2: "dua",
			3: "tiga",
			4: "empat",
			5: "lima",
			6: "enam",
			7: "tujuh",
			8: "delapan",
			9: "sembilan",
		},
		LiteralConversion: map[int]string{
			10:  "sepuluh",
			11:  "sebelas",
			12:  "dua belas",
			13:  "tiga belas",
			14:  "empat belas",
			15:  "lima belas",
			16:  "enam belas",
			17:  "tujuh belas",
			18:  "delapan belas",
			19:  "sembilan belas",
			100: "seratus",
		},
		DigitNames: map[int]string{
			10:  "puluh",
			100: "ratus",
		},
		GroupNames: map[int]string{
			0: "",
			1: "ribu",
			2: "juta",
			3: "miliar",
			4: "trilyun",
		},
		Correction: map[string]string{
			"satu ribu": "seribu",
		},
		PointConversion:         "koma",
		CurrencyName:            "rupiah",
		CurrencyPointConversion: "dan",
		CurrencyPointName:       "sen",
		CurrencyPointLength:     2,
	}
	return &num
}

// NumeralCreateEnglish creates numeral struct for English language
func NumeralCreateEnglish() *Numeral {
	num := Numeral{
		SplitDigit:     3,
		ZeroConversion: "zero",
		Conversion: map[int]string{
			1: "one",
			2: "two",
			3: "three",
			4: "four",
			5: "five",
			6: "six",
			7: "seven",
			8: "eight",
			9: "nine",
		},
		LiteralConversion: map[int]string{
			10: "ten",
			11: "eleven",
			12: "twelve",
			13: "thirteen",
			14: "fourteen",
			15: "fifteen",
			16: "sixteen",
			17: "seventeen",
			18: "eighteen",
			19: "nineteen",
		},
		DigitNames: map[int]string{
			10:  "ty",
			100: "hundred",
		},
		GroupNames: map[int]string{
			0: "",
			1: "thousand",
			2: "million",
			3: "billion",
			4: "trillion",
		},
		Correction: map[string]string{
			"two ty":   "twenty",
			"three ty": "thirty",
			"four ty":  "fourty",
			"five ty":  "fifty",
			"six ty":   "sixty",
			"seven ty": "seventy",
			"eight ty": "eighty",
			"nine ty":  "ninety",
		},
		PointConversion:         "point",
		CurrencyName:            "",
		CurrencyPointConversion: "and",
		CurrencyPointName:       "cents",
		CurrencyPointLength:     2,
	}
	return &num
}
