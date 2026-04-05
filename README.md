# XPicoConnect

A Go library that bridges a Raspberry Pi Pico (over USB serial) with [X-Plane](https://www.x-plane.com/) via [XPHTTPBridge](https://github.com/steveiliop56/xphttpbridge). Build physical cockpit controls and instruments using a Pico microcontroller and have them talk to the simulator in real time.

## How It Works

```
┌───────────┐  serial   ┌──────────────────┐  HTTP   ┌─────────────┐
│   Pico    │◄─────────►│  XPicoConnect    │◄───────►│  X-Plane    │
│  (USB)    │           │  (Go program)    │         │  (XPHTTPBridge)
└───────────┘           └──────────────────┘         └─────────────┘
```

XPicoConnect manages the serial connection to a Pico running the companion Arduino sketch and connects to X-Plane through XPHTTPBridge. You can then:

- Send commands to the Pico and receive responses (e.g. toggle LEDs, drive servos).
- Bind X-Plane datarefs to callbacks that fire every poll cycle (e.g. read airspeed, altitude).
- Bind Pico commands to callbacks so the Pico can request data from the host (e.g. switch inputs).

## Serial Protocol

Communication uses a simple newline-delimited text protocol:

| Direction | Format | Example |
|---|---|---|
| Command | `command:value\n` | `led:on\n` |
| Response | `command:status:result\n` | `led:ok:on\n` |

A full-duplex handshake (`fdx`) is performed at startup to confirm both sides are ready.

## Getting Started

### Prerequisites

- Go 1.26+
- A Raspberry Pi Pico (or compatible board) connected via USB
- [XPHTTPBridge](https://github.com/steveiliop56/xphttpbridge) running alongside X-Plane

### Install

```
go get github.com/steveiliop56/xpicoconnect
```

### Flash the Pico

Upload `base_pico.ino` (or one of the project-specific sketches in `projects/`) to your Pico using the Arduino IDE.

### Configuration

Create a `config.ini` file:

```ini
poll_time = 100

[serial]
baudrate = 115200
port = /dev/ttyACM0
buffer_size = 256
timeout = 10000

[xphttpbridge]
address = localhost
port = 8080
```

### Minimal Example

```go
package main

import "github.com/steveiliop56/xpicoconnect"

func main() {
    xpc := xpicoconnect.NewXPicoConnector(xpicoconnect.XPicoConnectorConfig{})

    if err := xpc.ReadInConfig("config.ini"); err != nil {
        panic(err)
    }

    if err := xpc.Initialize(); err != nil {
        panic(err)
    }

    // React when the Pico sends a command
    xpc.BindPicoCommand(xpicoconnect.PicoCommandBind{
        Command: "switch",
        Callback: func(value []byte) ([]byte, error) {
            // handle the command, return an encoded response
            return commands.EncodeResponse("switch", "ok", []byte("toggled")), nil
        },
    })

    // Poll an X-Plane dataref and forward it to the Pico
    xpc.BindBridgeRef(xpicoconnect.BridgeRefBind{
        Ref:     "sim/cockpit2/gauges/indicators/airspeed_kts_pilot",
        IsSlice: false,
        Callback: func(value any) {
            // send the value to the Pico, update a display, etc.
        },
    })

    // Blocks until SIGINT/SIGTERM
    xpc.Listen()
}
```

## Project Structure

```
├── base_pico.ino          # Base Arduino sketch for the Pico
├── commands/              # Encode/decode helpers for the serial protocol
├── hat/                   # Raspberry Pi Sense HAT LED matrix animations
├── python/                # Python helpers (Sense HAT animations, protocol utils)
├── projects/              # Example projects
├── serial.go              # Serial port reader
├── types.go               # Core types and config structs
└── xpicoconnector.go      # Main library entry point
```

## Example Projects

The `projects/` directory contains complete, working examples. Each project has its own `main.go`, Arduino sketch (`.ino`), `config.ini`, and `go.mod`. See `projects/README.md` for details on running them locally.

## License

[MIT](LICENSE)
