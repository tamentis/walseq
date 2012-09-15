/*
 * Copyright (c) 2012 Bertrand Janin <b@grun.gy>
 * 
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	MIN_TIMELINE uint64 = 0x00000001
	MAX_TIMELINE uint64 = 0xFFFFFFFF
	MIN_LOGICAL  uint64 = 0x00000000
	MAX_LOGICAL  uint64 = 0xFF000000
	MIN_PHYSICAL uint64 = 0x00000000
	MAX_PHYSICAL uint64 = 0x000000FE

	SEG_FMT = "%08X%08X%08X"
)

// get_args parses the command-line arguments and returns the two segments.
func getArgs() (string, string) {
	var showHelp = flag.Bool("h", false, "help")
	flag.Parse()

	if *showHelp {
		fmt.Println("usage: walseq [-h] [seg_start [seg_end]]")
		os.Exit(-1)
	}

	args := flag.Args()

	defaultStart := fmt.Sprintf(SEG_FMT, MIN_TIMELINE, MIN_LOGICAL, MIN_PHYSICAL)
	defaultStop := fmt.Sprintf(SEG_FMT, MAX_TIMELINE, MAX_LOGICAL, MAX_PHYSICAL)

	if len(args) == 0 {
		return defaultStart, defaultStop
	}

	segStart := args[0]
	if len(segStart) != 24 {
		fmt.Fprintln(os.Stderr, "seg_start should be 24 characters long")
		os.Exit(100)
	}

	if len(args) == 1 {
		return segStart, defaultStop
	}

	segStop := args[1]
	if len(segStop) != 24 {
		fmt.Fprintln(os.Stderr, "seg_stop should be 24 characters long")
		os.Exit(101)
	}

	return segStart, segStop
}

// parseSegUnit converts an hexadecimal string to a unsigned 64 bit integer,
// doing all the error handling.
func parseSegUint(seg string) uint64 {
	i, err := strconv.ParseUint(seg, 16, 64)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to parse segment '%s'", seg)
		os.Exit(102)
	}

	return i
}

// segToIntegers converts a 24 char WAL log segment code into a three integer
// tuple.
func segToIntegers(seg string) (uint64, uint64, uint64) {
	timeline := parseSegUint(seg[0:8])
	logical := parseSegUint(seg[8:16])
	physical := parseSegUint(seg[16:24])

	return timeline, logical, physical
}

func main() {
	segStart, segStop := getArgs()

	if segStop < segStart {
		fmt.Fprintln(os.Stderr, "End segment should be larger (older) than "+
			"the start segment")
		os.Exit(103)
	}

	startTimeline, startLogical, startPhysical := segToIntegers(segStart)
	stopTimeline, stopLogical, stopPhysical := segToIntegers(segStop)

	initialStopPhysical := stopPhysical
	initialStopLogical := stopLogical

	for timeline := startTimeline; timeline <= stopTimeline; timeline++ {
		// Only the last loop gets actual stopLogical, others get MAX_LOGICAL
		if timeline < stopTimeline {
			stopLogical = MAX_LOGICAL
		} else {
			stopLogical = initialStopLogical
		}

		for logical := startLogical; logical <= stopLogical; logical++ {
			// Only the last loop gets actual stopPhysical, others get MAX_PHYSICAL.
			if logical < stopLogical {
				stopPhysical = MAX_PHYSICAL
			} else {
				stopPhysical = initialStopPhysical
			}

			for physical := startPhysical; physical <= stopPhysical; physical++ {
				fmt.Printf(SEG_FMT+"\n", timeline, logical, physical)
			}

			// After the first loop, start back at 0.
			startPhysical = MIN_PHYSICAL
		}

		// After the first loop, start back at 0.
		startLogical = MIN_LOGICAL
	}
}
