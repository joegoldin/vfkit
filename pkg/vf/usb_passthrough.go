package vf

import (
	"github.com/Code-Hex/vz/v3"
	"github.com/crc-org/vfkit/pkg/config"
)

type USBPassthrough config.USBPassthrough

func (dev *USBPassthrough) AddToVirtualMachineConfig(vmConfig *VirtualMachineConfiguration) error {
	if err := (*config.USBPassthrough)(dev).Validate(); err != nil {
		return err
	}
	if err := vz.USBPassthroughAvailable(); err != nil {
		return err
	}
	vmConfig.usbPassthroughFilters = append(vmConfig.usbPassthroughFilters, vz.USBPassthroughFilter{
		VendorID:  dev.VendorID,
		ProductID: dev.ProductID,
	})
	return nil
}

// ListenUSBPassthrough registers the requested accessories after the VM starts.
// Accessory Access only hands over devices the user grants to this process.
func (vm *VirtualMachine) ListenUSBPassthrough() (*vz.USBPassthroughSession, error) {
	if len(vm.vfConfig.usbPassthroughFilters) == 0 {
		return nil, nil
	}
	return vm.VirtualMachine.StartUSBPassthrough(vm.vfConfig.usbPassthroughFilters)
}
