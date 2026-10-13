package config

import (
	"fmt"
	"strconv"
	"strings"
)

// USBPassthrough requests user-approved USB accessories through Accessory Access.
// An omitted vendor/product filter expresses interest in all USB accessories.
type USBPassthrough struct {
	VendorID  *uint16 `json:"vendorId,omitempty"`
	ProductID *uint16 `json:"productId,omitempty"`
}

func (dev *USBPassthrough) Validate() error {
	if dev.ProductID != nil && dev.VendorID == nil {
		return fmt.Errorf("usb-passthrough productId requires vendorId")
	}
	return nil
}

func (dev *USBPassthrough) FromOptions(options []option) error {
	for _, opt := range options {
		var field **uint16
		switch opt.key {
		case "vendorId":
			field = &dev.VendorID
		case "productId":
			field = &dev.ProductID
		default:
			return fmt.Errorf("unknown option for usb-passthrough devices: %s", opt.key)
		}
		if *field != nil {
			return fmt.Errorf("duplicate usb-passthrough option: %s", opt.key)
		}
		base := 10
		value := opt.value
		if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
			base = 16
			value = value[2:]
		}
		id, err := strconv.ParseUint(value, base, 16)
		if err != nil {
			return fmt.Errorf("invalid usb-passthrough %s: %w", opt.key, err)
		}
		parsed := uint16(id)
		*field = &parsed
	}
	return dev.Validate()
}

func (dev *USBPassthrough) ToCmdLine() ([]string, error) {
	if err := dev.Validate(); err != nil {
		return nil, err
	}
	arg := "usb-passthrough"
	if dev.VendorID != nil {
		arg += fmt.Sprintf(",vendorId=0x%04x", *dev.VendorID)
	}
	if dev.ProductID != nil {
		arg += fmt.Sprintf(",productId=0x%04x", *dev.ProductID)
	}
	return []string{"--device", arg}, nil
}
