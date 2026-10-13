package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUSBPassthroughRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		input   string
		command string
		json    string
	}{
		{"usb-passthrough", "usb-passthrough", `{"kind":"usbpassthrough"}`},
		{"usb-passthrough,vendorId=4660", "usb-passthrough,vendorId=0x1234", `{"kind":"usbpassthrough","vendorId":4660}`},
		{"usb-passthrough,productId=0x0000,vendorId=0x1234", "usb-passthrough,vendorId=0x1234,productId=0x0000", `{"kind":"usbpassthrough","vendorId":4660,"productId":0}`},
	} {
		t.Run(tc.input, func(t *testing.T) {
			dev, err := deviceFromCmdLine(tc.input)
			require.NoError(t, err)
			args, err := dev.ToCmdLine()
			require.NoError(t, err)
			require.Equal(t, []string{"--device", tc.command}, args)
			data, err := json.Marshal(dev)
			require.NoError(t, err)
			require.JSONEq(t, tc.json, string(data))
			decoded, err := unmarshalDevice(data)
			require.NoError(t, err)
			require.Equal(t, dev, decoded)
		})
	}
}

func TestUSBPassthroughInvalidOptions(t *testing.T) {
	for _, options := range []string{
		"productId=1", "vendorId=", "vendorId=-1", "vendorId=65536",
		"vendorId=0x10000", "vendorId=no", "vendorId=1,vendorId=2",
		"vendorId=1,productId=2,productId=3", "vendorId=1,productId=65536",
		"path=/dev/usb", "all", "vendorId=1,serial=abc",
	} {
		t.Run(options, func(t *testing.T) {
			_, err := deviceFromCmdLine("usb-passthrough," + options)
			require.Error(t, err)
		})
	}
}

func TestUSBPassthroughInvalidJSON(t *testing.T) {
	for _, data := range []string{
		`{"kind":"usbpassthrough","productId":1}`,
		`{"kind":"usbpassthrough","vendorId":65536}`,
	} {
		_, err := unmarshalDevice([]byte(data))
		require.Error(t, err)
	}
}
