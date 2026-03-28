package main

import (
	"fmt"
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
		// We don't need to poll too fast for the VSI since it doesn't change that rapidly,
		// and it gives us more time to process the data and update the servo position
		PollTime: 200,
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
			servo_pos := int((verticalSpeed-float64(C172VSI.Min))/(float64(C172VSI.Max)-float64(C172VSI.Min))*(180-0) + 0)
			log.Printf("Servo position: %d degrees", servo_pos)
			connector.SendPicoCommand("vsi_servo", fmt.Appendf([]byte{}, "%d", servo_pos))
		},
	})

	connector.Listen()
}
