package hat

import (
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/steveiliop56/xpicoconnect/commands"
)

type HatConfig struct {
	PyPath         string
	BinPath        string
	AnimationDelay float32
}

type Hat struct {
	config HatConfig
}

func NewHat(config HatConfig) *Hat {
	return &Hat{
		config: config,
	}
}

func (h *Hat) callPy(cmd string) (string, error) {
	abspy, err := filepath.Abs(h.config.PyPath)
	if err != nil {
		return "", err
	}
	basepath := path.Dir(abspy)
	pypath := path.Base(abspy)
	var stdout strings.Builder
	var stderr strings.Builder
	ecmd := exec.Command(h.config.BinPath, pypath, cmd)
	ecmd.Dir = basepath
	ecmd.Stdout = &stdout
	ecmd.Stderr = &stderr
	err = ecmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed with %v and out %v", err, stderr.String())
	}
	return stdout.String(), nil
}

func (h *Hat) buildAnimation(animation string, direction string, delay float32) string {
	var builder strings.Builder
	builder.WriteString(animation)
	builder.WriteString(":")
	builder.WriteString(direction)
	builder.WriteString(",")
	fmt.Fprintf(&builder, "%f", delay)
	return builder.String()
}

func (h *Hat) runAnimation(animation string) error {
	res, err := h.callPy(animation)
	if err != nil {
		return err
	}
	_, err = commands.DecodeResponse([]byte(res))
	if err != nil {
		return err
	}
	return nil
}

func (h *Hat) Test() error {
	anim := h.buildAnimation("test", "left", h.config.AnimationDelay)
	return h.runAnimation(anim)
}

func (h *Hat) Transmit() error {
	anim := h.buildAnimation("transmit", "up", h.config.AnimationDelay)
	return h.runAnimation(anim)
}

func (h *Hat) Receive() error {
	anim := h.buildAnimation("receive", "down", h.config.AnimationDelay)
	return h.runAnimation(anim)
}

func (h *Hat) Main() error {
	anim := h.buildAnimation("main", "left", h.config.AnimationDelay)
	return h.runAnimation(anim)
}

func (h *Hat) Shutdown() error {
	anim := h.buildAnimation("shutdown", "right", h.config.AnimationDelay)
	return h.runAnimation(anim)
}
