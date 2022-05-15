package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "Commands:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\ttimestamp\tadd timestamp, by eg. `midicat in | midiseq timestamp > record.txt`\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\tplay [FILE]\toutput message according to the timestamps, by eg. `midiseq play record.txt | midicat out`\n")
	}
	flag.Parse()
	switch flag.Arg(0) {
	case "timestamp":
		if flag.NArg() != 1 {
			fmt.Print("midiseq: too many arguments\n")
			os.Exit(2)
		}
		timestamp()

	case "play":
		var (
			cmd     = flag.NewFlagSet("play", flag.ExitOnError)
			channel = cmd.Int("channel", 0, "midi output channel")
		)
		cmd.Parse(os.Args[2:])
		fmt.Println(cmd.Args())
		if cmd.NArg() != 1 {
			fmt.Print("midiseq: invalid argument count\n")
			os.Exit(2)
		}
		play(cmd.Arg(0), *channel)

	default:
		if flag.NArg() == 0 {
			fmt.Print("midiseq: missing a command, try -h\n")
		} else {
			fmt.Printf("midiseq: unknown command %q\n", flag.Arg(0))
		}
		os.Exit(2)
	}
}

func timestamp() {
	var s = bufio.NewScanner(os.Stdin)
	for s.Scan() {
		fmt.Fprintf(os.Stdout, "%s %s\n", time.Now().Format(time.RFC3339Nano), s.Text())
	}
}

func play(name string, channel int) {
	type event struct {
		time    time.Duration
		delta   time.Duration
		message string
	}
	f, err := os.Open(name)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	var t0, lt time.Time
	var events []event
	for s.Scan() {
		l := s.Text()
		c := strings.Split(l, " ")
		if len(c) == 0 {
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, c[0])
		check(err)
		if t0.IsZero() {
			t0, lt = t, t
		}
		if len(c) == 1 {
			continue
		}
		m, err := hex.DecodeString(c[1])
		check(err)
		if m[0]&0xF0 == 0xF0 {
			continue
		}
		if channel > 0 && m[0] >= 0x80 && m[0] <= 0xEF {
			m[0] &= 0xF0
			m[0] |= byte(channel-1) & 0xF
		}
		events = append(events, event{t.Sub(t0), t.Sub(lt), hex.EncodeToString(m)})
		lt = t
	}
	var j int
	var t time.Duration
	var dt = time.Minute / (300 * 96)
	var ticker = time.NewTicker(dt)
	for range ticker.C {
		if j >= len(events) {
			return
		}
		t += dt
		var i int
		for i = j; i < len(events); i++ {
			var e = events[i]
			if e.time <= t {
				fmt.Fprintln(os.Stdout, e.message)
			} else {
				break
			}
		}
		j = i
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
