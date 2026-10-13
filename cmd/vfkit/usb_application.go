//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c -fno-objc-arc
#cgo LDFLAGS: -framework AppKit
#include "usb_application.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"runtime/cgo"

	"github.com/crc-org/vfkit/pkg/config"
)

func init() {
	// AppKit must run on the process's initial thread, including console VMs.
	runtime.LockOSThread()
}

func needsUSBApplication(vmConfig *config.VirtualMachine) bool {
	for _, dev := range vmConfig.Devices {
		if _, ok := dev.(*config.USBPassthrough); ok {
			return true
		}
	}
	return false
}

func validateUSBApplication(vmConfig *config.VirtualMachine) error {
	if needsUSBApplication(vmConfig) {
		for _, gpu := range vmConfig.VirtioGPUDevices() {
			if gpu.UsesGUI {
				return fmt.Errorf("usb-passthrough currently requires a console VM; remove --gui")
			}
		}
	}
	return nil
}

func runUSBApplication(run func() error) error {
	if !bool(C.prepareUSBApplication()) {
		return fmt.Errorf("USB accessories require a regular macOS application with a Dock icon")
	}
	ready := make(chan struct{})
	handle := cgo.NewHandle(ready)
	defer handle.Delete()
	result := make(chan error, 1)
	go func() {
		<-ready
		result <- run()
		C.stopUSBApplication()
	}()
	C.runUSBApplication(C.uintptr_t(handle))
	return <-result
}

//export usbApplicationReady
func usbApplicationReady(handle C.uintptr_t) {
	close(cgo.Handle(handle).Value().(chan struct{}))
}
