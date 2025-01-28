package main

import (
	"fmt"
	"strings"
	"unicode"
)

type parseMode int

const (
	Unquoted parseMode = iota
	SingleQuoted
	DoubleQuoted
)

func parseCommand(in string) ([]string, error) {
	var sb strings.Builder
	tokens := make([]string, 0)

	appendToken := func() {
		if sb.Len() == 0 {
			return
		}
		tokens = append(tokens, sb.String())
		sb.Reset()
	}

	i := 0
	mode := Unquoted
	chs := []rune(in)

	for i < len(chs) {
		ch := chs[i]
		switch mode {
		case Unquoted:
			{
				if ch == '"' {
					mode = DoubleQuoted
				} else if ch == '\'' {
					mode = SingleQuoted
				} else if ch == '\\' {
					if i+1 < len(chs) {
						sb.WriteRune(chs[i+1])
						i += 1
					}
				} else if unicode.IsSpace(ch) {
					appendToken()
				} else {
					sb.WriteRune(ch)
				}
			}
		case SingleQuoted:
			{
				if ch == '\'' {
					mode = Unquoted
				} else {
					sb.WriteRune(ch)
				}
			}
		case DoubleQuoted:
			{
				if ch == '"' {
					mode = Unquoted
				} else if ch == '\\' {
					if i+1 < len(chs) {
						next := chs[i+1]
						if next == '$' || next == '`' || next == '"' || next == '\\' || next == '\n' {
							sb.WriteRune(next)
							i += 1 // skip next
						} else {
							sb.WriteRune(chs[i])
						}
					}
				} else {
					sb.WriteRune(ch)
				}
			}
		default:
			return nil, fmt.Errorf("unknown parse mode = %v", mode)
		}

		i += 1
	}

	// append remaining token
	appendToken()

	if mode != Unquoted {
		return nil, fmt.Errorf("invalid state, mode = %v not handled properly", mode)
	}

	return tokens, nil
}
