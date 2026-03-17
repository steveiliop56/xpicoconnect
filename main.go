package main

import (
	"fmt"
	"log"
	"time"

	xphttpbridgego "github.com/steveiliop56/xphttpbridge-go"
)

type AppState struct {
	StrobeLightState bool
	BatteryOnState   bool
}

const (
	StrobeLightRef = "sim/cockpit/electrical/strobe_lights_on"
	BatteryOnRef   = "sim/cockpit/electrical/battery_on"
)

func main() {
	appState := AppState{}

	log.Print("setup pico bridge")

	pb, err := NewPicoBridge(PicoBridgeConfig{
		PortName:   "/dev/ttyACM0",
		BaudRate:   115200,
		BufferSize: 256,
	})

	if err != nil {
		log.Fatalf("failed to setup pico bridge: %v", err)
	}

	readchan := pb.StartReader()

	log.Print("setup command handler")

	ch := NewCommandHandler(pb)

	log.Print("ping pico")

	res, err := ch.SendCommand("ping", []byte{})

	if err != nil {
		log.Fatalf("failed to ping pico: %v", err)
	}

	if string(res) != "pong" {
		log.Fatalf("unexpected response to ping: %v", string(res))
	}

	pingBack := <-readchan

	cmd, val, err := ch.DecodeCommand(pingBack)

	if err != nil {
		log.Fatalf("failed to decode ping response: %v", err)
	}

	log.Printf("received command: %v, value: %v", cmd, val)

	log.Print("pico good, bidirectional communication established")

	ch.SendCommand("led", []byte("on"))

	log.Print("setup xphttpbridge")

	xpclient := xphttpbridgego.NewClient(xphttpbridgego.Config{
		Port:    49000,
		Address: "127.0.0.1",
	})

	err = xpclient.Ping()

	if err != nil {
		log.Fatalf("failed to ping xphttpbridge: %v", err)
	}

	log.Print("bridge good")

	strobeState, err := getBoolState(xpclient, StrobeLightRef)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to get strobe light state: %v", err))
	}

	batteryState, err := getBoolState(xpclient, BatteryOnRef)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to get battery on state: %v", err))
	}

	log.Printf("initial battery on state: %v", batteryState)
	appState.BatteryOnState = batteryState

	log.Printf("initial strobe light state: %v", strobeState)

	appState.StrobeLightState = strobeState

	if appState.StrobeLightState && appState.BatteryOnState {
		res, err = ch.SendCommand("strobe", []byte(boolToSwitch(appState.StrobeLightState)))

		if err != nil {
			log.Fatalf("failed to set initial strobe light state: %v", err)
		}

		log.Printf("initial strobe light state set to: %v, response: %v", true, string(res))
	} else {
		log.Printf("strobe light is off due to battery state: %v", appState.BatteryOnState)
	}

	log.Print("starting main loop")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			batteryState, err := getBoolState(xpclient, BatteryOnRef)

			if err != nil {
				log.Printf("failed to get battery on state: %v", err)
				continue
			}

			if batteryState != appState.BatteryOnState {
				log.Printf("battery state changed: %v", batteryState)
				appState.BatteryOnState = batteryState
				if !batteryState {
					log.Printf("battery turned off, ensuring strobe light is off")
					if appState.StrobeLightState {
						res, err := ch.SendCommand("strobe", []byte("off"))
						if err != nil {
							log.Printf("failed to turn off strobe light: %v", err)
						} else {
							log.Printf("strobe light turned off due to battery state, response: %v", string(res))
						}
					}
					appState.StrobeLightState = false
				}
				continue
			}

			strobeState, err := getBoolState(xpclient, StrobeLightRef)

			if err != nil {
				log.Printf("failed to get strobe light state: %v", err)
				continue
			}

			if strobeState != appState.StrobeLightState {
				log.Printf("strobe light state changed: %v", strobeState)
				appState.StrobeLightState = strobeState

				if !appState.BatteryOnState {
					log.Printf("battery is off, skipping strobe light state change")
					continue
				}

				res, err := ch.SendCommand("strobe", []byte(boolToSwitch(appState.StrobeLightState)))

				if err != nil {
					log.Printf("failed to set strobe light state: %v", err)
					continue
				}

				log.Printf("strobe light state set to: %v, response: %v", appState.StrobeLightState, string(res))
			}
		case data := <-readchan:
			if ch.commandPending {
				log.Printf("command pending, ignoring data: %v", string(data))
				continue
			}

			cmd, val, err := ch.DecodeCommand(data)

			if err != nil {
				log.Printf("failed to decode command: %v", err)
				continue
			}

			log.Printf("received command: %v, value: %v", cmd, val)

			switch cmd {
			case "ping":
				_, err := ch.SendCommand("ping", []byte("pong"))
				if err != nil {
					log.Printf("failed to respond to ping: %v", err)
				}
			default:
				log.Printf("unknown command: %v", cmd)
			}
		}
	}
}

func getBoolState(xpclient *xphttpbridgego.Client, refName string) (bool, error) {
	val, err := xpclient.GetDataRef(refName)
	if err != nil {
		return false, err
	}
	state, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("value not bool")
	}
	return state, nil
}

func boolToSwitch(b bool) string {
	if b {
		return "on"
	}
	return "off"
}
