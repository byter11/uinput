package uinput

import (
	"fmt"
	"io"
	"os"
)

// A TouchPad is an input device that uses absolute axis events, meaning that you can specify
// the exact position the cursor should move to. Therefore, it is necessary to define the size
// of the rectangle in which the cursor may move upon creation of the device.
type TouchPad interface {
	// MoveTo will move the cursor to the specified position on the screen
	MoveTo(x int32, y int32) error

	// LeftClick will issue a single left click.
	LeftClick() error

	// RightClick will issue a right click.
	RightClick() error

	// LeftPress will simulate a press of the left mouse button. Note that the button will not be released until
	// LeftRelease is invoked.
	LeftPress() error

	// LeftRelease will simulate the release of the left mouse button.
	LeftRelease() error

	// RightPress will simulate the press of the right mouse button. Note that the button will not be released until
	// RightRelease is invoked.
	RightPress() error

	// RightRelease will simulate the release of the right mouse button.
	RightRelease() error

	// TouchDown will simulate a single touch to a virtual touch device. Use TouchUp to end the touch gesture.
	TouchDown() error

	// TouchUp will end or ,more precisely, unset the touch event issued by TouchDown
	TouchUp() error

	io.Closer
}

type vTouchPad struct {
	name       []byte
	deviceFile *os.File
}

// CreateTouchPad will create a new touch pad device. note that you will need to define the x and y axis boundaries
// (min and max) within which the cursor maybe moved around.
func CreateTouchPad(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32) (TouchPad, error) {
	err := validateDevicePath(path)
	if err != nil {
		return nil, err
	}
	err = validateUinputName(name)
	if err != nil {
		return nil, err
	}

	fd, err := createTouchPad(path, name, minX, maxX, minY, maxY)
	if err != nil {
		return nil, err
	}

	return vTouchPad{name: name, deviceFile: fd}, nil
}

func (vTouch vTouchPad) MoveTo(x int32, y int32) error {
	return sendAbsEvent(vTouch.deviceFile, absX, absY, x, y)
}

func (vTouch vTouchPad) LeftClick() error {
	err := sendBtnEvent(vTouch.deviceFile, []int{evBtnLeft}, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the LeftClick event: %v", err)
	}

	return sendBtnEvent(vTouch.deviceFile, []int{evBtnLeft}, btnStateReleased)
}

func (vTouch vTouchPad) RightClick() error {
	err := sendBtnEvent(vTouch.deviceFile, []int{evBtnRight}, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the RightClick event: %v", err)
	}

	return sendBtnEvent(vTouch.deviceFile, []int{evBtnRight}, btnStateReleased)
}

// LeftPress will simulate a press of the left mouse button. Note that the button will not be released until
// LeftRelease is invoked.
func (vTouch vTouchPad) LeftPress() error {
	return sendBtnEvent(vTouch.deviceFile, []int{evBtnLeft}, btnStatePressed)
}

// LeftRelease will simulate the release of the left mouse button.
func (vTouch vTouchPad) LeftRelease() error {
	return sendBtnEvent(vTouch.deviceFile, []int{evBtnLeft}, btnStateReleased)
}

// RightPress will simulate the press of the right mouse button. Note that the button will not be released until
// RightRelease is invoked.
func (vTouch vTouchPad) RightPress() error {
	return sendBtnEvent(vTouch.deviceFile, []int{evBtnRight}, btnStatePressed)
}

// RightRelease will simulate the release of the right mouse button.
func (vTouch vTouchPad) RightRelease() error {
	return sendBtnEvent(vTouch.deviceFile, []int{evBtnRight}, btnStateReleased)
}

func (vTouch vTouchPad) TouchDown() error {
	return sendBtnEvent(vTouch.deviceFile, []int{evBtnTouch, evBtnToolFinger}, btnStatePressed)
}

func (vTouch vTouchPad) TouchUp() error {
	return sendBtnEvent(vTouch.deviceFile, []int{evBtnTouch, evBtnToolFinger}, btnStateReleased)
}

func (vTouch vTouchPad) Close() error {
	return closeDevice(vTouch.deviceFile)
}

func createTouchPad(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32) (fd *os.File, err error) {
	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not create absolute axis input device: %v", err)
	}

	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register key device: %v", err)
	}
	// register button events (in order to enable left and right click)
	for _, event := range []int{evBtnLeft, evBtnRight, evBtnTouch, evBtnToolFinger} {
		err = ioctl(deviceFile, uiSetKeyBit, uintptr(event))
		if err != nil {
			deviceFile.Close()
			return nil, fmt.Errorf("failed to register button event %v: %v", event, err)
		}
	}

	err = registerDevice(deviceFile, uintptr(evAbs))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute axis input device: %v", err)
	}

	// register x and y axis events
	for _, event := range []int{absX, absY} {
		err = ioctl(deviceFile, uiSetAbsBit, uintptr(event))
		if err != nil {
			deviceFile.Close()
			return nil, fmt.Errorf("failed to register absolute axis event %v: %v", event, err)
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

