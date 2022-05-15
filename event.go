package main

import (
	"bufio"
	"encoding/hex"
	"io"
	"strings"
	"time"
)

type Event struct {
	Time       time.Duration
	Tick       int
	Message    string
	Raw        []byte
	IsRealTime bool
}

func (e Event) IsNoteOn() bool {
	return e.Raw[0]&0xF0 == 0x90 && e.Raw[2] > 0
}

func (e Event) IsNoteOff() bool {
	return e.Raw[0]&0xF0 == 0x80 && (e.Raw[0]&0xF0 == 0x90 && e.Raw[2] == 0)
}

func (e Event) NoteOff() string {
	return hex.EncodeToString([]byte{0x80 | e.Raw[0]&0x0F, e.Raw[1], 0})
}

type EventList []*Event

func ReadEventList(r io.Reader) (ls EventList) {
	s := bufio.NewScanner(r)
	var t0 time.Time
	var tick int
	for s.Scan() {
		l := s.Text()
		if len(l) == 0 {
			continue
		}
		c := strings.Split(l, " ")
		if len(c) == 0 {
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, c[0])
		check(err)
		if t0.IsZero() {
			t0 = t
		}
		if len(c) == 1 {
			continue
		}
		m, err := hex.DecodeString(c[1])
		check(err)
		if m[0] == 0xF8 {
			tick++
		}
		ls = append(ls, &Event{t.Sub(t0), tick, hex.EncodeToString(m), m, m[0]&0xF0 == 0xF0})
	}
	return
}

func (ls EventList) SetChannel(c int) {
	for _, e := range ls {
		if c > 0 && e.Raw[0] >= 0x80 && e.Raw[0] <= 0xEF {
			e.Raw[0] &= 0xF0
			e.Raw[0] |= byte(c-1) & 0xF
			e.Message = hex.EncodeToString(e.Raw)
		}
	}
}
