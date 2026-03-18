package main

import (
	"log"
)

const (
	StrobeLightRef = "sim/cockpit/electrical/strobe_lights_on"
	BatteryOnRef   = "sim/cockpit/electrical/battery_on"
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
		PollTime: 100,
	}

	connector, err := NewXPicoConnector(connectorCfg)

	if err != nil {
		log.Fatalf("failed to setup connector: %v", err)
	}

	connector.BindPicoCommand(PicoCommandBind{
		Command: "switch_0",
		Callback: func(value []byte) ([]byte, error) {
			log.Printf("received command: switch_0 with value: %s", string(value))

			if string(value) == "1" {
				err := connector.GetXPBridge().SetDataRef(StrobeLightRef, true)
				if err != nil {
					return nil, err
				}
			} else {
				err := connector.GetXPBridge().SetDataRef(StrobeLightRef, false)
				if err != nil {
					return nil, err
				}
			}

			encodedRes := encodeResponse("switch_0", "ok", []byte{})
			return encodedRes, nil
		},
	})

	connector.Listen()
}
