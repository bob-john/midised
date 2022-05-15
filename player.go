package main

import (
	"fmt"
	"os"
	"time"
)

func Play(ls EventList) {
	var j int
	var t time.Duration
	var dt = time.Minute / (300 * 96)
	var ticker = time.NewTicker(dt)
	for range ticker.C {
		if j >= len(ls) {
			return
		}
		t += dt
		var i int
		for i = j; i < len(ls); i++ {
			var e = ls[i]
			if e.Time <= t {
				fmt.Fprintln(os.Stdout, e.Message)
			} else {
				break
			}
		}
		j = i
	}
}
