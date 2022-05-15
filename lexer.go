package main

import (
	"bufio"
	"io"
	"unicode"
)

type Lexer struct {
	r   *bufio.Reader
	eof bool
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{bufio.NewReader(r), false}
}

func (l *Lexer) Int() (n int, ok bool) {
	l.Skim()
	for {
		r, eof := l.ReadRune()
		if eof {
			return
		}
		if unicode.IsDigit(r) {
			n = n*10 + int(r-'0')
			ok = true
		} else {
			l.r.UnreadRune()
			return
		}
	}
}

func (l *Lexer) Char(c rune) (string, bool) {
	l.Skim()
	r, eof := l.ReadRune()
	if eof {
		return "", false
	}
	if r == c {
		return string(r), true
	}
	l.r.UnreadRune()
	return "", false
}

func (l *Lexer) Range() (Range, bool) {
	var r Range
	if i, ok := l.Int(); ok {
		r[0] = i
	}
	if _, ok := l.Char(':'); !ok {
		return r, false
	}
	if i, ok := l.Int(); ok {
		r[1] = i
	}
	return r, true
}

func (l *Lexer) Skim() {
	for {
		r, eof := l.ReadRune()
		if eof {
			return
		}
		if !unicode.IsSpace(r) {
			l.r.UnreadRune()
			return
		}
	}
}

func (l *Lexer) Done() bool {
	return l.eof
}

func (l *Lexer) ReadRune() (rune, bool) {
	r, _, err := l.r.ReadRune()
	if err != nil && err != io.EOF {
		panic(err)
	}
	l.eof = err == io.EOF
	return r, l.eof
}
