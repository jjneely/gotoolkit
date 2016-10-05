package lex

import (
	"fmt"
	"unicode"
)

// State 0: No quote
// State 1: Token break (whitespace)
// State 2: backslash
// State 3: In Quote

func Tokenize(s string) ([]string, error) {
	var (
		result []string
		token  []rune
		state  int
		quote  rune
	)

	chunk := func() {
		if len(token) > 0 {
			result = append(result, string(token))
			token = nil // zero length slice
		}
	}

	for _, c := range s {
		switch state {
		case 0:
			switch {
			case unicode.IsSpace(c):
				chunk()
			case c == '\\':
				state = 2
			case c == '\'' || c == '"':
				state = 3
				quote = c
			default:
				token = append(token, c)
			}
		case 2:
			switch {
			case unicode.IsSpace(c):
				fallthrough
			case c == '\'':
				fallthrough
			case c == '"':
				token = append(token, c)
			default:
				token = append(token, '\\')
				token = append(token, c)
			}
			if quote == '\x00' {
				state = 0
			} else {
				state = 3
			}
		case 3:
			switch c {
			case '\\':
				state = 2
			case quote:
				chunk()
				quote = '\x00'
				state = 0
			default:
				token = append(token, c)
			}
		}
	}

	// End of String
	switch state {
	case 2:
		token = append(token, '\\')
	case 3:
		return nil, fmt.Errorf("Missing closing quote")
	}

	chunk()
	return result, nil
}
