package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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
			fmt.Fprintln(os.Stderr, "midiseq: too many arguments")
			os.Exit(2)
		}
		timestamp()

	case "play":
		var (
			cmd     = flag.NewFlagSet("play", flag.ExitOnError)
			remix   = cmd.String("remix", "", "remix instructions")
			channel = cmd.Int("channel", 0, "midi output channel")
		)
		cmd.Parse(os.Args[2:])
		if cmd.NArg() != 1 {
			fmt.Fprintln(os.Stderr, "midiseq: invalid argument count")
			os.Exit(2)
		}
		play(cmd.Arg(0), *remix, *channel)

	default:
		if flag.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "midiseq: missing a command, try -h")
		} else {
			fmt.Fprintf(os.Stderr, "midiseq: unknown command %q\n", flag.Arg(0))
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

func play(name string, remix string, channel int) {
	// var r, err = ParseRemix(remix)
	// check(err)
	f, err := os.Open(name)
	check(err)
	defer f.Close()
	var ls = ReadEventList(f)
	f.Close()
	if channel > 0 {
		ls.SetChannel(channel)
	}
	Play(ls)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
