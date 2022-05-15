package main

import (
	"fmt"
	"os"
	"time"
)

func Play(ls EventList) {
	var j int
	var dt = time.Minute / (300 * 96)
	var ticker = time.NewTicker(dt)
	var t0 = time.Now()
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
			} else {
				break
			}
		}
		j = i
	}
}
