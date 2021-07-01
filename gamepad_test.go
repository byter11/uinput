package uinput

import (
	"testing"
	"time"
	// "time"
)

// This test will confirm that basic gamepad button events are working.
func TestBtnsInValidRangeWork(t *testing.T) {
	vk, err := CreateGamepad("/dev/uinput", []byte("Test Basic Gamepad"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	// for i := 0; i<100; i++ {
	err = vk.BtnDown(BtnMin)
	if err != nil {
		t.Fatalf("Failed to send button down event %d. Last error was: %s\n", BtnMin, err)
	}

	err = vk.BtnUp(BtnMin)
	if err != nil {
		t.Fatalf("Failed to send button up event. Last error was: %s\n", err)
	}
	err = vk.SetAxis(-32767, 32767)
	if err != nil {
		t.Fatalf("Failed to send axis event")
	}
	// time.Sleep(time.Second)
	err = vk.SetAxis(0, 0)
	if err != nil {
		t.Fatalf("Failed to send axis event")
	}

	err = vk.SetAxisR(-32767, 32767)
	if err != nil {
		t.Fatalf("Failed to send axis event")
	}
	// time.Sleep(time.Second)
	err = vk.SetAxisR(0, 0)
	if err != nil {
		t.Fatalf("Failed to send axis event")
	}
	// time.Sleep(time.Second)
	// }

	err = vk.BtnDown(BtnMax)
	if err != nil {
		t.Fatalf("Failed to send button down event. Last error was: %s\n", err)
	}

	err = vk.BtnUp(BtnMax)
	if err != nil {
		t.Fatalf("Failed to send button up event. Last error was: %s\n", err)
	}

	err = vk.Close()

	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}

func TestDpad(t *testing.T) {
	vk, err := CreateGamepad("/dev/uinput", []byte("Test Basic Gamepad"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	for i := 0; i < 100; i++ {
		err = vk.BtnDown(BtnDpadDown)
		if err != nil {
			t.Fatalf("Failed to send button down event %d. Last error was: %s\n", BtnMin, err)
		}
		time.Sleep(time.Second)
		err = vk.BtnUp(BtnDpadDown)
		if err != nil {
			t.Fatalf("Failed to send button up event. Last error was: %s\n", err)
		}
		time.Sleep(time.Second)
		// err = vk.BtnDown(BtnTriggerHappy1)
		// if err != nil {
		// 	t.Fatalf("Failed to send button down event %d. Last error was: %s\n", BtnMin, err)
		// }
	}
}
