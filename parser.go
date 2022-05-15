package main

import (
	"io"
)

type Parser struct {
	l *Lexer
}

func NewParser(r io.Reader) *Parser {
	return &Parser{NewLexer(r)}
}

func (p *Parser) Parse() (r Remix) {
	for !p.l.Done() {
		if _, ok := p.l.Char('C'); ok {
			if v, ok := p.l.Range(); ok {
				r.Channels = v
			} else {
				panic("parser: expect range 'first:last' after 'C'")
			}
		} else if _, ok := p.l.Char('B'); ok {
			if v, ok := p.l.Range(); ok {
				r.Beats = v
			} else {
				panic("parser: expect range 'first:last' after 'B'")
			}
		} else if !p.l.Done() {
			panic("parser: expect 'C' or 'B'")
		}
	}
	return
}
