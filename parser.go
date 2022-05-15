package main

import "io"

type Parser struct {
	l *Lexer
}

func NewParser(r io.Reader) *Parser {
	return &Parser{NewLexer(r)}
}

func (p *Parser) Parse() (r Remix) {
	if i, ok := p.l.Int(); ok {
		r.Begin = i
	}
	if _, ok := p.l.Colon(); !ok {
		panic("parser: expect ':'")
	}
	if i, ok := p.l.Int(); ok {
		r.End = i
	} else if p.l.IsEOF() {
		r.End = -1
	} else {
		panic("parser: expect integer")
	}
	return
}
