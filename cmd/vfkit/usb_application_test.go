//go:build darwin

package main

import (
	"testing"

	"github.com/crc-org/vfkit/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestUSBApplicationConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name     string
		devices  []config.VirtioDevice
		needsApp bool
		invalid  bool
	}{
		{name: "ordinary console"},
		{name: "ordinary GUI", devices: []config.VirtioDevice{&config.VirtioGPU{UsesGUI: true}}},
		{name: "USB console", devices: []config.VirtioDevice{&config.USBPassthrough{}}, needsApp: true},
		{name: "USB GPU without window", devices: []config.VirtioDevice{&config.USBPassthrough{}, &config.VirtioGPU{}}, needsApp: true},
		{name: "USB GUI unsupported", devices: []config.VirtioDevice{&config.USBPassthrough{}, &config.VirtioGPU{UsesGUI: true}}, needsApp: true, invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			vm := &config.VirtualMachine{Devices: tc.devices}
			require.Equal(t, tc.needsApp, needsUSBApplication(vm))
			err := validateUSBApplication(vm)
			if tc.invalid {
				require.ErrorContains(t, err, "requires a console VM")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
