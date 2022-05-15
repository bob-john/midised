package main

import (
	"bytes"
	"time"

	"github.com/pkg/errors"
)

type Remix struct {
	Begin, End int
}

var ErrRemixSyntax = errors.New("remix: syntax error")

func ParseRemix(s string) (*Remix, error) {
	p := NewParser(bytes.NewReader([]byte(s)))
	r := p.Parse()
	return &r, nil
}

func (r *Remix) Apply(s EventList) (d EventList) {
	d = make([]*Event, 0, cap(s))
	var begin, end = r.Begin * 24, r.End * 24
	var t0 time.Duration
	for _, e := range s {
		if e.Tick < begin {
			continue
		}
		if t0 == 0 {
			t0 = e.Time
		}
		if end > 0 && e.Tick >= end {
			break
		}
		e.Time -= t0
		e.Tick -= begin
		d = append(d, e)
	}
	return
}
