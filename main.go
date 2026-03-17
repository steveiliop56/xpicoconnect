package main

import (
	"fmt"
	"log"
	"time"

	xphttpbridgego "github.com/steveiliop56/xphttpbridge-go"
)

type AppState struct {
	BeaconLightState bool
}

const (
	BeaconLightRef = "sim/cockpit/electrical/beacon_lights_on"
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

	log.Print("pico good")

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

	state, err := getBeaconState(xpclient)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to get beacon light state: %v", err))
	}

	log.Printf("initial beacon light state: %v", state)

	appState.BeaconLightState = state

	res, err = ch.SendCommand("beacon", []byte(boolToSwitch(appState.BeaconLightState)))

	if err != nil {
		log.Fatalf("failed to set initial beacon light state: %v", err)
	}

	log.Printf("initial beacon light state set to: %v, response: %v", appState.BeaconLightState, string(res))

	log.Print("starting main loop")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		state, err := getBeaconState(xpclient)

		if err != nil {
			log.Printf("failed to get beacon light state: %v", err)
			continue
		}

		if state != appState.BeaconLightState {
			log.Printf("beacon light state changed: %v", state)
			appState.BeaconLightState = state

			res, err := ch.SendCommand("beacon", []byte(boolToSwitch(appState.BeaconLightState)))

			if err != nil {
				log.Printf("failed to set beacon light state: %v", err)
			}

			log.Printf("beacon light state set to: %v, response: %v", appState.BeaconLightState, string(res))
		}
	}
}

func getBeaconState(xpclient *xphttpbridgego.Client) (bool, error) {
	val, err := xpclient.GetDataRef(BeaconLightRef)
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
