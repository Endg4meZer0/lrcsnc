package dynreplacer

import "strings"

// DynamicReplacer is a utility used to improve template functionality
// of output/client module. It is a very simplified middle-ground between
// strings.Replacer and a proper text/template.
//
// It is hardcoded to use curly brackets as start and end delimeters for keys.
type DynamicReplacer struct {
	m map[string]func() string
}

const startDelim = '{'
const endDelim = '}'

func NewDynamicReplacer(_m map[string]func() string) *DynamicReplacer {
	return &DynamicReplacer{
		m: _m,
	}
}

func (dr *DynamicReplacer) Replace(template string) string {
	var result strings.Builder

	// An optimization to reduce the number of allocations.
	result.Grow(len(template) * 2)

	i := 0

	for i < len(template) {
		startIndex := strings.IndexRune(template[i:], startDelim)
		if startIndex == -1 {
			result.WriteString(template[i:])
			break
		}
		startIndex += i
		result.WriteString(template[i:startIndex])

		endIndex := strings.IndexRune(template[startIndex:], endDelim)
		if endIndex == -1 {
			result.WriteString(template[startIndex:])
			break
		}
		endIndex += startIndex
		key := template[startIndex+1 : endIndex]

		if fn, ok := dr.m[key]; ok {
			res := fn()
			result.WriteString(res)
		} else {
			result.WriteString(template[startIndex : endIndex+1])
		}
		i = endIndex + 1
	}

	return result.String()
}
