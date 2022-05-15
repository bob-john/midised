package main

import (
	"fmt"
	"os"
	"time"
)

var noteOff = make(map[string]struct{})

func Play(ls EventList) {
	var j int
	var dt = time.Minute / (300 * 96)
	var ticker = time.NewTicker(dt)
	var t0 = time.Now()
	defer Stop()
	for range ticker.C {
		if j >= len(ls) {
			return
		}
		var t = time.Since(t0)
		var i int
		for i = j; i < len(ls); i++ {
			var e = ls[i]
			if e.Time <= t {
				if !e.IsRealTime {
					fmt.Fprintln(os.Stdout, e.Message)
				}
				if e.IsNoteOn() {
					noteOff[e.NoteOff()] = struct{}{}
				} else if e.IsNoteOff() {
					delete(noteOff, e.NoteOff())
				}
			} else {
				break
			}
		}
		j = i
	}
}

func Stop() {
	for m := range noteOff {
		fmt.Fprintln(os.Stdout, m)
	}
	noteOff = make(map[string]struct{})
}
