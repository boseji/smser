# `smser` - UTC Serial Time Sender

Small Go program that opens a serial port, sends the current UTC time as a single ISO8601 timestamp prefixed with `~`, then exits.

## Build

```bash
go build -trimpath .
```

## Run

Default port is `/dev/ttyACM0` and default baud rate is `115200`.

```bash
./smser --port /dev/ttyACM0 --baud 115200
```

You can also use the defaults:

```bash
./smser
```

## Output format

A single line like:

`~2026-07-31T12:34:56.123456789Z`

## License

```text
Copyright 2026 boseji(https://github.com/boseji)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
