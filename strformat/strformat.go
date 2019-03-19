package strformat

/*
	StrFormat package
	ver 1.2 - 2019-03-19
	by Zecchan Silverlake

	This package contains useful function to manipulate strings
*/

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StringFormatter is used to format a string template, it comes with custom format.
// Note: Call Init() before adding CustomFormat or it will panic
type StringFormatter struct {
	CustomFormat  map[string]func(string) string
	UseCustomTime bool
	CustomTime    time.Time
}

// FormatString formats a specified string using Format
func (sf *StringFormatter) FormatString(str string) string {
	var regex, err = regexp.Compile("(%date\\([yMdHhmsa]+\\)%)")
	var cts = time.Now()
	if sf.UseCustomTime {
		cts = sf.CustomTime
	}
	if err == nil {
		str = regex.ReplaceAllStringFunc(str, func(a string) string {
			a = strings.TrimPrefix(a, "%date(")
			a = strings.TrimSuffix(a, ")%")

			var HH = cts.Hour()
			var hh = cts.Hour() % 12
			var aa = "AM"
			if HH > 12 {
				aa = "PM"
			}
			if hh == 0 {
				hh = 12
			}

			a = strings.Replace(a, "y", strconv.Itoa(cts.Year()), 1)
			a = strings.Replace(a, "M", PadLeft(strconv.Itoa(int(cts.Month())), "0", 2), 1)
			a = strings.Replace(a, "d", PadLeft(strconv.Itoa(cts.Day()), "0", 2), 1)
			a = strings.Replace(a, "H", PadLeft(strconv.Itoa(HH), "0", 2), 1)
			a = strings.Replace(a, "h", PadLeft(strconv.Itoa(hh), "0", 2), 1)
			a = strings.Replace(a, "i", PadLeft(strconv.Itoa(cts.Minute()), "0", 2), 1)
			a = strings.Replace(a, "s", PadLeft(strconv.Itoa(cts.Second()), "0", 2), 1)
			a = strings.Replace(a, "a", aa, 1)
			return a
		})
	}

	if sf.CustomFormat != nil {
		for k, v := range sf.CustomFormat {
			if strings.Contains(str, k) {
				str = v(str)
			}
		}
	}
	return str
}

// Init initializes this StringFormatter
func (sf *StringFormatter) Init() {
	sf.CustomFormat = map[string]func(string) string{}
}

// PadLeft pads the left of a string with specified char so that the string will have a length of totalLength
func PadLeft(str string, padChar string, totalLength int) string {
	var c = padChar[0:1]
	var cnt = totalLength - len(str)
	if cnt > 0 {
		str = strings.Repeat(c, cnt) + str
	}
	return str
}

// PadRight pads the right of a string with specified char so that the string will have a length of totalLength
func PadRight(str string, padChar string, totalLength int) string {
	var c = padChar[0:1]
	var cnt = totalLength - len(str)
	if cnt > 0 {
		str += strings.Repeat(c, cnt)
	}
	return str
}

// Filter filters out characters that is not defined in charset
func Filter(str string, charset string) string {
	var re string
	for _, rn := range strings.Split(str, "") {
		if strings.Contains(charset, rn) {
			re += rn
		}
	}
	return re
}

// Capitalize will capitalize each word excluding a and of
func Capitalize(str string) string {
	if str == "" {
		return ""
	}
	var spl = strings.Split(str, " ")
	res := ""
	for _, word := range spl {
		if word != "" && word != "a" && word != "of" {
			nword := strings.ToUpper(word[0:1])
			if len(word) > 1 {
				nword += word[1:]
			}
			word = nword
		}
		if word != "" {
			if res != "" || word == "" {
				res += " "
			}
			res += word
		}
	}
	return res
}

const (
	// CharsetNumber contains all numbers
	CharsetNumber = "0123456789"
	// CharsetAlphaLowercase contains lowercase alpha
	CharsetAlphaLowercase = "abcdefghijklmnopqrstuvwxyz"
	// CharsetAlphaUppercase contains uppercase alpha
	CharsetAlphaUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// CharsetAlphaNumeric contains lowercase alpha, uppercase alpha and numbers
	CharsetAlphaNumeric = CharsetNumber + CharsetAlphaLowercase + CharsetAlphaUppercase
)
