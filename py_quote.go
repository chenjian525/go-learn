package sina

import (
	"strconv"
)

const (
	pyquote encoding = iota + 1
)

type EscapeError string

func (e EscapeError) Error() string {
	return "invalid URL escape " + strconv.Quote(string(e))
}

type encoding int

func quote(s string) string {
	return escape(s, pyquote)
}

func unquote(s string) (string, error) {
	return unescape(s, pyquote)
}

func shouldEscape(c byte, mode encoding) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	switch c {
	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@':
		switch mode {
		case pyquote:
			return c != '/'
		default:
			return false
		}
	default:
		return false
	}
}

func escape(s string, mode encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == pyquote {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == pyquote:
			t[j] = '+'
			j++
		case shouldEscape(c, pyquote):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = c
			j++
		}
	}
	return string(t)
}

func isHex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	default:
		return false
	}
}

func unHex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	default:
		return 0
	}
}

func unescape(s string, mode encoding) (string, error) {
	n := 0
	hasPlus := false
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			n++
			if n+2 >= len(s) || !isHex(s[n+1]) || !isHex(s[n+2]) {
				s = s[i:]
				if len(s) > 3 {
					s = s[:3]
				}
				return "", EscapeError(s)
			}
			i += 3
		case '+':
			hasPlus = mode == pyquote
			i++
		default:
			i++
		}
	}

	if n == 0 && !hasPlus {
		return s, nil
	}

	t := make([]byte, len(s)-2*n)
	j := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			t[j] = unHex(s[i+1])<<4 | unHex(s[i+2])
			j++
			i += 3
		case '+':
			if mode == pyquote {
				t[j] = ' '
			} else {
				t[j] = '+'
			}
			i++
			j++
		default:
			t[j] = s[i]
			i++
			j++
		}
	}
	return string(t), nil
}
