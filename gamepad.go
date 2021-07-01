package uinput

import (
	"io"
	"os"
	"fmt"
)

type Gamepad interface {
	// SetAxis will set the position of the gamepad's left axis.
	SetAxis(x, y int32) error

	// SetAxisR will set the postition of the gamepad's right axis.
	SetAxisR(x, y int32) error

	// BtnDown will send a BtnDown event to an existing gamepad device.
	// The button can be any of the predefined codes from gamepadcodes.go.
	BtnDown(btn int) error

	// BtnUp will send a BtnUp event to an existing gamepad device.
	// The button can be any of the predefined codes from gamepadcodes.go.
	BtnUp(btn int) error

	io.Closer
}

type vGamepad struct {
	name       []byte
	deviceFile *os.File
}

// CreateGamepad will create a new gamepad using the given uinput
// device path of the uinput device.
func CreateGamepad(path string, name []byte) (Gamepad, error) {
	err := validateDevicePath(path)
	if err != nil {
		return nil, err
	}
	err = validateUinputName(name)
	if err != nil {
		return nil, err
	}

	fd, err := createVGamepadDevice(path, name, -32767, 32767, -32767, 32767)
	if err != nil {
		return nil, err
	}

	return vGamepad{name: name, deviceFile: fd}, nil
}

func(vg vGamepad) SetAxis(x, y int32) error {
	if err := sendAbsEvent(vg.deviceFile, absX, absY, x, y); err != nil{
		return fmt.Errorf("failed to move axis")
	}
	return nil
}

func(vg vGamepad) SetAxisR(x, y int32) error {
	if err := sendAbsEvent(vg.deviceFile, absRX, absRY, x, y); err != nil{
		return fmt.Errorf("failed to move right axis")
	}
	return nil
}

func (vg vGamepad) BtnDown(key int) error {
	if !BtnCodeInRange(key) {
		return fmt.Errorf("failed to perform BtnDown. Code %d is not in range", key)
	}
	return sendBtnEvent(vg.deviceFile, []int{key}, btnStatePressed)
}

func (vg vGamepad) BtnUp(key int) error {
	if !BtnCodeInRange(key) {
		return fmt.Errorf("failed to perform BtnUp. Code %d is not in range", key)
	}

	return sendBtnEvent(vg.deviceFile, []int{key}, btnStateReleased)
}

// Close will close the device and free resources.
// It's usually a good idea to use defer to call this function.
func (vg vGamepad) Close() error {
	return closeDevice(vg.deviceFile)
}

func createVGamepadDevice(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32) (fd *os.File, err error) {
	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual gamepad device: %v", err)
	}

	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register virtual gamepad device: %v", err)
	}

	// register key events
	for i := BtnMin; i <= BtnMax; i++ {
		err = ioctl(deviceFile, uiSetKeyBit, uintptr(i))
		if err != nil {
			deviceFile.Close()
			fmt.Print(i)
			return nil, fmt.Errorf("failed to register Btn number %d: %v", i, err)
		}
	}

	err = registerDevice(deviceFile, uintptr(evAbs))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register Absolute axis input device: %v", err)
	}

	for _, event := range []int{absX, absY, absZ, absRX, absRY, absRZ} {
		err = ioctl(deviceFile, uiSetAbsBit, uintptr(event))
		if err != nil {
			deviceFile.Close()
			return nil, fmt.Errorf("failed to register relative event %v: %v", event, err)
		}
	}

	var absMin [absSize]int32
	absMin[absX] = minX
	absMin[absY] = minY

	var absMax [absSize]int32
	absMax[absX] = maxX
	absMax[absY] = maxY

	return createUsbDevice(deviceFile,
		uinputUserDev{
			Name: toUinputName(name),
			ID: inputID{
				Bustype: busUsb,
				Vendor:  0x4711,
				Product: 0x0817,
				Version: 1},
			Absmin: absMin,
			Absmax: absMax})
}


func BtnCodeInRange(btn int) bool {
	return btn >= BtnMin && btn <= BtnMax
}