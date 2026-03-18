package main

import (
	"log"
)

const (
	BeaconLightRef = "sim/cockpit/electrical/beacon_lights_on"
)

func main() {
	connectorCfg := XPicoConnectorConfig{
		SerialConfig: SerialConfig{
			Baudrate:   115200,
			Port:       "/dev/ttyACM0",
			BufferSize: 256,
			Timeout:    10000,
		},
		XPHTTPBridgeConfig: XPHTTPBridgeConfig{
			Address: "127.0.0.1",
			Port:    49000,
		},
		PollTime: 50,
	}

	connector, err := NewXPicoConnector(connectorCfg)

	if err != nil {
		log.Fatalf("failed to setup connector: %v", err)
	}

	connector.BindPicoCommand(PicoCommandBind{
		Command: "beacon_switch",
		Callback: func(value []byte) ([]byte, error) {
			log.Printf("received command: beacon_switch with value: %s", string(value))

			if string(value) == "1" {
				err := connector.GetXPBridge().SetDataRef(BeaconLightRef, true)
				if err != nil {
					return nil, err
				}
			} else {
				err := connector.GetXPBridge().SetDataRef(BeaconLightRef, false)
				if err != nil {
					return nil, err
				}
			}

			encodedRes := encodeResponse("beacon_switch", "ok", []byte{})
			return encodedRes, nil
		},
	})

	prevBeaconState := false

	connector.BindBridgeRef(BridgeRefBind{
		Ref:     BeaconLightRef,
		IsSlice: false,
		Callback: func(value any) {
			state, ok := value.(bool)
			if !ok {
				log.Printf("unexpected type for beacon light ref: %T", value)
				return
			}
			if state != prevBeaconState {
				log.Printf("beacon light state changed: %v", state)
				prevBeaconState = state
				if state {
					_, err := connector.SendPicoCommand("beacon", []byte("on"))
					if err != nil {
						log.Printf("failed to send beacon command: %v", err)
					}
				} else {
					_, err := connector.SendPicoCommand("beacon", []byte("off"))
					if err != nil {
						log.Printf("failed to send beacon command: %v", err)
					}
				}
			}
		},
	})

	connector.Listen()
}
