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
	Min: -1000,
	Max: 1000,
}

func main() {
	connector := xpicoconnect.NewXPicoConnector(xpicoconnect.XPicoConnectorConfig{})

	err := connector.ReadInConfig("config.ini")

	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	err = connector.Initialize()

	if err != nil {
		log.Fatalf("failed to setup connector: %v", err)
	}

	// reset the servo to the middle position on startup
	connector.SendPicoCommand("vsi_servo", fmt.Appendf([]byte{}, "%d", 90))

	log.Println("subscribing to events")

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

	log.Println("starting listener")

	connector.Listen()
}
