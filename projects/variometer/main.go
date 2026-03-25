package main

import (
	"log"

	"github.com/steveiliop56/xpicoconnect"
)

const (
	VerticalSpeedRef = "sim/cockpit2/gauges/indicators/vvi_fpm_pilot"
)

type VSI struct {
	Min int
	Max int
}

var C172VSI = VSI{
	Min: -2000,
	Max: 2000,
}

func main() {
	connectorCfg := xpicoconnect.XPicoConnectorConfig{
		SerialConfig: xpicoconnect.SerialConfig{
			Baudrate:   115200,
			Port:       "/dev/ttyACM0",
			BufferSize: 256,
			Timeout:    10000,
		},
		XPHTTPBridgeConfig: xpicoconnect.XPHTTPBridgeConfig{
			Address: "127.0.0.1",
			Port:    49000,
		},
		PollTime: 50,
	}

	connector := xpicoconnect.NewXPicoConnector(connectorCfg)
	err := connector.Initialize()

	if err != nil {
		log.Fatalf("failed to setup connector: %v", err)
	}

	connector.BindBridgeRef(xpicoconnect.BridgeRefBind{
		Ref:     VerticalSpeedRef,
		IsSlice: false,
		Callback: func(value any) {
			verticalSpeed := value.(float64)
			log.Printf("Vertical speed: %f fpm", verticalSpeed)
			servo_pos := (verticalSpeed-float64(C172VSI.Min))/(float64(C172VSI.Max)-float64(C172VSI.Min))*(180-0) + 0
			if servo_pos < 0 {
				servo_pos = 0
			} else if servo_pos > 180 {
				servo_pos = 180
			}
			log.Printf("Servo position: %f degrees", servo_pos)
		},
	})

	connector.Listen()
}
