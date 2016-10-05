package lex

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	Str string
	Pos int
}

func NewLexer(s string) *Lexer {
	l := new(Lexer)
	l.Str = s
	return l
}

func (l *Lexer) Position() int {
	return l.Pos
}

func (l *Lexer) Len() int {
	return len(l.Str)
}

func (l *Lexer) Next() rune {
	if l.Pos >= len(l.Str) {
		return '\x00'
	}
	runeValue, width := utf8.DecodeRuneInString(l.Str[l.Pos:])
	l.Pos += width

	return runeValue
}

func (l *Lexer) Peek() rune {
	if l.Pos >= len(l.Str) {
		return '\x00'
	}

	runeValue, _ := utf8.DecodeRuneInString(l.Str[l.Pos:])
	return runeValue
}

func (l *Lexer) Rewind() {
	// XXX: we need to rewind whatever the length of the last rune was
	if l.Pos > 0 {
		l.Pos--
	}
}

// State 0: No quote
// State 1: Token break (whitespace)
// State 2: backslash
// State 3: In Quote

func Tokenize(s string) ([]string, error) {
	result := make([]string, 0)
	token := make([]rune, 0)
	state := 0
	quote := '\x00'

	chunk := func() {
		if len(token) > 0 {
			result = append(result, string(token))
			token = token[0:-1] // zero length slice
		}
	}

	for c := range s {
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
			case unicode.IsSpace(c) && quote == '\x00':
				token = append(token, c)
			case c == '\'' && quote == '\'':
				token = append(token, c)
			case c == '"' && quote == '"':
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
