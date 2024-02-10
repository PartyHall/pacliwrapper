package pacliwrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
)

type PaVolume struct {
	Value        int64  `json:"value"`
	ValuePercent string `json:"value_percent"`
	Decibel      string `json:"db"`
}

type PaDevice struct {
	Index int    `json:"index"`
	State string `json:"state"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Driver string   `json:"driver"`
	Mute   bool     `json:"mute"`
	Volume PaVolume `json:"volume"`
}

func (d PaDevice) GetVolume() int {
	return int((float64(d.Volume.Value) / 65536.0) * 100)
}

func (d PaDevice) SetMute(mute bool) error {
	val := "0"
	if mute {
		val = "1"
	}

	cmd := exec.Command("pactl", "set-sink-mute", fmt.Sprintf("%v", d.Index), val)
	return cmd.Run()
}

func (d *PaDevice) SetVolume(percentage int) error {
	cmd := exec.Command("pactl", "set-sink-volume", fmt.Sprintf("%v", d.Index), fmt.Sprintf("%v%%", percentage))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing pactl command: %v, stderr: %s", err, stderr.String())
	}

	d.Volume.Value = int64(float64(percentage) / 100 * 65536)

	return nil
}

type PaCliWrapper struct {
	Devices    []PaDevice
	MainDevice *PaDevice
}

func New() (*PaCliWrapper, error) {
	pa := &PaCliWrapper{}

	return pa, nil
}

func (pcw *PaCliWrapper) Refresh() error {
	cmd := exec.Command("pactl", "-f", "json", "list", "sinks")

	data, err := cmd.Output()
	if err != nil {
		return err
	}

	var parsed []map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}

	cmd = exec.Command("pactl", "get-default-sink")
	data, err = cmd.Output()
	if err != nil {
		return err
	}

	defaultSink := strings.TrimSpace(string(data))

	pcw.MainDevice = nil
	pcw.Devices = []PaDevice{}
	for _, device := range parsed {
		idx, _ := device["index"].(float64)
		state, _ := device["state"].(string)
		name, _ := device["name"].(string)
		desc, _ := device["description"].(string)
		driver, _ := device["driver"].(string)
		mute, _ := device["mute"].(bool)

		bv, _ := device["volume"].(map[string]interface{})

		keys := reflect.ValueOf(bv).MapKeys()
		if len(keys) == 0 {
			fmt.Println("weird device")
			continue
		}

		bv_data := bv[(keys[0].String())].(map[string]interface{})

		bv_val, _ := bv_data["value"].(float64)
		bv_perc, _ := bv_data["value_percent"].(string)
		bv_db, _ := bv_data["db"].(string)

		currDevice := PaDevice{
			Index:       int(idx),
			State:       strings.TrimSpace(state),
			Name:        strings.TrimSpace(name),
			Description: strings.TrimSpace(desc),
			Driver:      strings.TrimSpace(driver),
			Mute:        mute,
			Volume: PaVolume{
				Value:        int64(bv_val),
				ValuePercent: strings.TrimSpace(bv_perc),
				Decibel:      strings.TrimSpace(bv_db),
			},
		}

		if currDevice.Name == defaultSink {
			pcw.MainDevice = &currDevice
		}

		pcw.Devices = append(pcw.Devices, currDevice)
	}

	return nil
}

func (pcw *PaCliWrapper) SetDefaultOutput(device PaDevice) error {
	cmd := exec.Command("pactl", "set-default-sink", fmt.Sprintf("%v", device.Index))
	return cmd.Run()
}
