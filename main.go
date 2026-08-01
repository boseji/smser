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
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [options]\n\nOptions:\n",
			"smser",
		)
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(),
			"\nSends a single UTC ISO8601 timestamp to the serial port,"+
				" \n  prefixed with '~', then exits.")
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

	nowUTC := time.Now().UTC()
	iso := fmt.Sprintf("%d", nowUTC.Unix()) // Normally in Unix Epoch mode
	if !*epoch {
		iso = nowUTC.Format(time.RFC3339Nano) // e.g. 2026-07-31T12:34:56.123456789Z
	}
	line := "~" + iso + "\n"

	if _, err := p.Write([]byte(line)); err != nil {
		log.Fatal(err)
	}

	fmt.Print(line) // optional; remove if you want zero stdout
	fmt.Printf("Sent to %q @ %d BAUD\n", *port, *baud)
	fmt.Println("Done!")
}
