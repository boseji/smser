//  `smser` - UTC Serial Time Sender
//
// Small Go program that opens a serial port, sends the current UTC time
// as a single ISO8601 timestamp prefixed with `~`, then exits.
//
// Copyright 2026 boseji(https://github.com/boseji)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	serial "go.bug.st/serial"
)

func main() {
	defaultPort := "/dev/ttyACM0"
	defaultBaud := 115200
	defaultEpoch := false

	port := flag.String("port", defaultPort,
		"Serial port path/name (e.g., /dev/ttyACM0 or COM3)")
	baud := flag.Int("baud", defaultBaud, "Baud rate (e.g., 115200)")
	epoch := flag.Bool("e", defaultEpoch, "Send time in Unix Epoch format")
	readBack := flag.Bool("r", false, "Read back the Time in UNIX Epoch"+
		"    This disables the set functionality. ")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [options]\n\nOptions:\n",
			"smser",
		)
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(),
			"\nSends a single UTC ISO8601 timestamp to the serial port,"+
				" \n  prefixed with '~', then exits."+
				"\n For Reading back the Time in UTC Epoch it sends,"+
				"\n '?' and then waits for receiving the 10 chars of UTC. ")
	}
	flag.Parse()

	mode := &serial.Mode{
		BaudRate: *baud,
	}

	p, err := serial.Open(*port, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	line := ""

	// Only Enable if we are in Set Mode
	if !*readBack {
		nowUTC := time.Now().UTC()
		iso := fmt.Sprintf("%d", nowUTC.Unix()) // Normally in Unix Epoch mode
		if !*epoch {
			iso = nowUTC.Format(time.RFC3339Nano) // e.g. 2026-07-31T12:34:56.123456789Z
		}
		line = "~" + iso + "\n"

		if _, err := p.Write([]byte(line)); err != nil {
			log.Fatal(err)
		}

		fmt.Print(line) // optional; remove if you want zero stdout
		fmt.Printf("Sent to %q @ %d BAUD\n", *port, *baud)
	} else {
		buf := make([]byte, 10) // Buffer to Read Back
		// Setup The Read Timeout
		if err := p.SetReadTimeout(1 * time.Second); err != nil {
			log.Fatal(err)
		}

		// Send the Read Command First
		if _, err := p.Write([]byte("?\n")); err != nil {
			log.Fatal(err)
		}

		// Read Back the Response
		if n, err := p.Read(buf); err != nil || n < 10 {
			log.Fatal(err, n)
		}

		// Get the Actual Time Stamp
		line = fmt.Sprintf("%s", buf)

		// Get the Unix Time Integer
		unix, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		t := time.Unix(unix, 0)
		fmt.Printf("Read Back time As : %s\n", t.UTC())
	}
	fmt.Println("Done!")
}
