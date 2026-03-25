package main

import (
	"log"
	"strconv"

	"github.com/steveiliop56/xpicoconnect"
)

const (
	BeaconLightRef = "sim/cockpit/electrical/beacon_lights_on"
	LandLightRef   = "sim/cockpit/electrical/landing_lights_on"
	TaxiLightRef   = "sim/cockpit/electrical/taxi_light_on"
	NavLightRef    = "sim/cockpit/electrical/nav_lights_on"
	StrobeLightRef = "sim/cockpit/electrical/strobe_lights_on"
)

type SwitchStates struct {
	Beacon bool
	Land   bool
	Taxi   bool
	Nav    bool
	Strobe bool
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

	connector.BindPicoCommand(xpicoconnect.PicoCommandBind{
		Command: "beacon_switch",
		Callback: func(value []byte) ([]byte, error) {
			return handleSwitchChange(connector, "beacon_switch", value, BeaconLightRef)
		},
	})

	connector.BindPicoCommand(xpicoconnect.PicoCommandBind{
		Command: "land_switch",
		Callback: func(value []byte) ([]byte, error) {
			return handleSwitchChange(connector, "land_switch", value, LandLightRef)
		},
	})

	connector.BindPicoCommand(xpicoconnect.PicoCommandBind{
		Command: "taxi_switch",
		Callback: func(value []byte) ([]byte, error) {
			return handleSwitchChange(connector, "taxi_switch", value, TaxiLightRef)
		},
	})

	connector.BindPicoCommand(xpicoconnect.PicoCommandBind{
		Command: "nav_switch",
		Callback: func(value []byte) ([]byte, error) {
			return handleSwitchChange(connector, "nav_switch", value, NavLightRef)
		},
	})

	connector.BindPicoCommand(xpicoconnect.PicoCommandBind{
		Command: "strobe_switch",
		Callback: func(value []byte) ([]byte, error) {
			return handleSwitchChange(connector, "strobe_switch", value, StrobeLightRef)
		},
	})

	connector.Listen()
}

func handleSwitchChange(connector *xpicoconnect.XPicoConnector, command string, value []byte, ref string) ([]byte, error) {
	log.Printf("received command: %s with value: %s", command, string(value))

	valueInt, err := strconv.Atoi(string(value))

	if err != nil {
		log.Printf("invalid value for command %s: %s, expected 0 or 1", command, string(value))
		return nil, err
	}

	err = connector.GetXPBridge().SetDataRef(ref, valueInt)

	if err != nil {
		log.Printf("failed to set dataref %s to value %s: %v", ref, string(value), err)
		return nil, err
	}

	encodedRes := xpicoconnect.EncodeResponse(command, "ok", []byte{})
	return encodedRes, nil
}
