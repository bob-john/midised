package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Remix struct {
	Begin, End int
}

var ErrRemixSyntax = errors.New("remix: syntax error")

func ParseRemix(s string) (*Remix, error) {
	c := strings.Split(s, ":")
	if len(c) != 2 {
		return nil, errors.WithStack(ErrRemixSyntax)
	}
	b, err := strconv.Atoi(c[0])
	if err != nil {
		return nil, errors.WithStack(err)
	}
	e, err := strconv.Atoi(c[1])
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Remix{b, e}, nil
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
		if e.Tick >= end {
			break
		}
		e.Time -= t0
		e.Tick -= begin
		d = append(d, e)
	}
	return
}
