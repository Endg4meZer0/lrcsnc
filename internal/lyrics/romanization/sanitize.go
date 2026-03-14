package romanization

import (
	"strings"
	"unicode"
)

// sanitizeBefore makes sure the str string is not going to cause
// any trouble later on the actual romanization step.
func sanitizeBefore(str string) string {
	out := strings.Builder{}

	// Convert common full-width characters to half-width
	for _, _r := range str {
		var r = _r
		if (_r >= 'Ａ' && _r <= 'Ｚ') ||
			(_r >= 'ａ' && _r <= 'ｚ') ||
			(_r >= '０' && _r <= '９') {
			r -= 0xFEE0
		} else {
			var ok bool
			r, ok = romajiMap[_r]
			if !ok {
				r = _r
			}
		}
		out.WriteRune(r)
	}

	return out.String()
}

// sanitizeAfter performs additional stuff after the romanization is done
func sanitizeAfter(str string) string {
	out := strings.Builder{}

	// Remove redundant spaces, add required spaces
	// and capitalize letters after punctuation
	// for better readability

	capitalizeNext := true
	spaceCheckActive := false
	for i, r := range str {
		if spaceCheckActive && !unicode.IsSpace(r) {
			out.WriteRune(' ')
		}
		spaceCheckActive = false

		if capitalizeNext && unicode.IsUpper(r) {
			capitalizeNext = false
		}

		if capitalizeNext && unicode.IsLower(r) {
			r = unicode.ToUpper(r)
			capitalizeNext = false
		}

		if r == ' ' && i != len(str)-1 &&
			(str[i+1] == '.' || str[i+1] == '!' || str[i+1] == '?') {
			continue
		}

		if r == '.' || r == '!' || r == '?' || r == '"' {
			capitalizeNext = true
			spaceCheckActive = true
		}

		out.WriteRune(r)
	}

	return out.String()
}

var romajiMap map[rune]rune = map[rune]rune{
	// Full-width space
	'　': ' ',
	// Full-width punctuation and symbols
	'，': ',',
	'．': '.',
	'；': ';',
	'：': ':',
	'！': '!',
	'？': '?',
	'＂': '"',
	'＇': '\'',
	'｀': '`',
	'＾': '^',
	'～': '~',
	'￣': '~',
	'－': '-',
	'＿': '_',
	'＆': '&',
	'＠': '@',
	'＃': '#',
	'％': '%',
	'＋': '+',
	'＊': '*',
	'＝': '=',
	'＜': '<',
	'＞': '>',
	'（': '(',
	'）': ')',
	'［': '[',
	'］': ']',
	'｛': '{',
	'｝': '}',
	'｟': '(',
	'｠': ')',
	'｜': '|',
	'￤': '|',
	'／': '/',
	'＼': '\\',
	'￢': '¬',
	'＄': '$',
	'￡': '£',
	'￠': '¢',
	'￦': '₩',
	'￥': '¥',
}
