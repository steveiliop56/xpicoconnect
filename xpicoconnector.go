package xpicoconnect

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	xphttpbridgego "github.com/steveiliop56/xphttpbridge-go"
	"go.bug.st/serial"
)

func NewXPicoConnector(config XPicoConnectorConfig) *XPicoConnector {
	return &XPicoConnector{
		config: config,
		state: XPicoConnectorState{
			picoCommandBinders: make(map[string]PicoCommandBind),
			bridgeRefBinders:   make(map[string]BridgeRefBind),
		},
	}
}

func (xpc *XPicoConnector) Initialize() error {
	_, err := xpc.setupSerial()
	if err != nil {
		return fmt.Errorf("failed to setup serial connection: %v", err)
	}

	_, err = xpc.setupXPHTTPBridge()
	if err != nil {
		return fmt.Errorf("failed to setup XPHTTPBridge connection: %v", err)
	}

	return nil
}

func (xpc *XPicoConnector) WithSerial(serial *serial.Port) (*XPicoConnector, error) {
	xpc.port = serial
	return xpc, nil
}

func (xpc *XPicoConnector) WithXPHTTPBridge(xpbridge *xphttpbridgego.Client) (*XPicoConnector, error) {
	xpc.xpbridge = xpbridge
	return xpc, nil
}

func (xpc *XPicoConnector) setupSerial() (serial.Port, error) {
	if xpc.port != nil {
		// End previous pico session if it exists, ignore errors since it might fail if there's no session
		xpc.SendPicoCommand("end", []byte("foo"))
		err := xpc.testPicoFDX()
		if err != nil {
			return nil, err
		}
		return *xpc.port, nil
	}

	mode := &serial.Mode{
		BaudRate: xpc.config.SerialConfig.Baudrate,
	}

	port, err := serial.Open(xpc.config.SerialConfig.Port, mode)

	if err != nil {
		return nil, err
	}

	xpc.port = &port

	xpc.setupReader()

	// End previous pico session if it exists, ignore errors since it might fail if there's no session
	xpc.SendPicoCommand("end", []byte("foo"))

	err = xpc.testPicoFDX()

	if err != nil {
		return nil, err
	}

	return port, nil
}

func (xpc *XPicoConnector) setupXPHTTPBridge() (*xphttpbridgego.Client, error) {
	if xpc.xpbridge != nil {
		err := xpc.ensureBridgeHealthy()
		if err != nil {
			return nil, err
		}
		return xpc.xpbridge, nil
	}

	clientCfg := xphttpbridgego.Config{
		Address: xpc.config.XPHTTPBridgeConfig.Address,
		Port:    xpc.config.XPHTTPBridgeConfig.Port,
	}

	client := xphttpbridgego.NewClient(clientCfg)

	xpc.xpbridge = client

	err := xpc.ensureBridgeHealthy()

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (xpc *XPicoConnector) setupReader() {
	xpc.readerChan = make(chan []byte, xpc.config.SerialConfig.BufferSize)
	go func() {
		for {
			line, err := readSerialLine(xpc.config.SerialConfig.BufferSize, *xpc.port)
			if err != nil {
				continue
			}
			xpc.readerChan <- line
		}
	}()
}

// Send a command to the pico, wait for the response,
// then wait for a command from the pico and send a response back
func (xpc *XPicoConnector) testPicoFDX() error {
	encoded := EncodeCommand("fdx", []byte("foo"))

	_, err := (*xpc.port).Write(encoded)
	if err != nil {
		return err
	}

	xpc.state.isCommandPending = true

	var res []byte
	select {
	case res = <-xpc.readerChan:
	case <-time.After(time.Duration(xpc.config.SerialConfig.Timeout) * time.Millisecond):
		xpc.state.isCommandPending = false
		return fmt.Errorf("timed out waiting for response from pico")
	}

	_, err = DecodeResponse(res)

	if err != nil {
		return err
	}

	select {
	case res = <-xpc.readerChan:
	case <-time.After(time.Duration(xpc.config.SerialConfig.Timeout) * time.Millisecond):
		xpc.state.isCommandPending = false
		return fmt.Errorf("timed out waiting for FDX response from pico")
	}

	xpc.state.isCommandPending = false

	cmd, value, err := DecodeCommand(res)
	if err != nil {
		return err
	}

	if cmd != "fdx" {
		return fmt.Errorf("expected command 'fdx', got '%s'", cmd)
	}

	if value != "foo" {
		return fmt.Errorf("expected response 'foo', got '%s'", value)
	}

	res = EncodeResponse("fdx", "ok", []byte("bar"))

	_, err = (*xpc.port).Write(res)
	if err != nil {
		return err
	}

	return nil
}

func (xpc *XPicoConnector) ensureBridgeHealthy() error {
	return xpc.xpbridge.Ping()
}

func (xpc *XPicoConnector) SendPicoCommand(command string, value []byte) (string, error) {
	encoded := EncodeCommand(command, value)

	_, err := (*xpc.port).Write(encoded)
	if err != nil {
		return "", err
	}

	xpc.state.isCommandPending = true

	var res []byte
	select {
	case res = <-xpc.readerChan:
	case <-time.After(time.Duration(xpc.config.SerialConfig.Timeout) * time.Millisecond):
		xpc.state.isCommandPending = false
		return "", fmt.Errorf("timed out waiting for response from pico")
	}

	xpc.state.isCommandPending = false

	return DecodeResponse(res)
}

func (xpc *XPicoConnector) BindPicoCommand(bind PicoCommandBind) {
	xpc.state.picoCommandBinders[bind.Command] = bind
}

func (xpc *XPicoConnector) BindBridgeRef(bind BridgeRefBind) {
	xpc.state.bridgeRefBinders[bind.Ref] = bind
}

func (xpc *XPicoConnector) DestroyPicoBind(command string) {
	delete(xpc.state.picoCommandBinders, command)
}

func (xpc *XPicoConnector) DestroyBridgeBind(ref string) {
	delete(xpc.state.bridgeRefBinders, ref)
}

func (xpc *XPicoConnector) GetPort() *serial.Port {
	return xpc.port
}

func (xpc *XPicoConnector) GetXPBridge() *xphttpbridgego.Client {
	return xpc.xpbridge
}

func (xpc *XPicoConnector) Close() error {
	xpc.state.bridgeRefBinders = make(map[string]BridgeRefBind)
	xpc.state.picoCommandBinders = make(map[string]PicoCommandBind)
	xpc.state.isCommandPending = false
	xpc.SendPicoCommand("end", []byte("foo"))
	return (*xpc.port).Close()
}

func (xpc *XPicoConnector) Listen() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(time.Duration(xpc.config.PollTime) * time.Millisecond)
	defer ticker.Stop()
	defer stop()

	for {
		select {
		case <-ticker.C:
			for _, bind := range xpc.state.bridgeRefBinders {
				if bind.IsSlice {
					res, err := xpc.xpbridge.GetDataRefSlice(bind.Ref)
					if err != nil {
						continue
					}
					bind.Callback(res)
					continue
				}
				res, err := xpc.xpbridge.GetDataRef(bind.Ref)
				if err != nil {
					continue
				}
				bind.Callback(res)
			}
		case line := <-xpc.readerChan:
			log.Printf("received line from pico: %s\n", string(line))
			if xpc.state.isCommandPending {
				continue
			}
			command, value, err := DecodeCommand(line)
			if err != nil {
				continue
			}
			binder, exists := xpc.state.picoCommandBinders[command]
			if !exists {
				continue
			}
			res, err := binder.Callback([]byte(value))
			if err != nil {
				res = EncodeResponse(command, "not_ok", []byte("callback_failed"))
			}
			_, err = (*xpc.port).Write(res)
			if err != nil {
				continue
			}
		case <-ctx.Done():
			log.Printf("%v received, shutting down...\n", context.Cause(ctx))
			goto shutdown
		}
	}

shutdown:
	err := xpc.Close()
	if err != nil {
		fmt.Printf("error closing XPicoConnector: %v\n", err)
	}
}
