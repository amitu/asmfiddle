package main

import (
	"asmfiddle/jsui"
	"asmfiddle/machine"
)

func main() {
	lcd, err := jsui.NewLCD()
	if err != nil {
		panic(err)
	}

	console, err := jsui.NewConsole()
	if err != nil {
		panic(err)
	}

	//kb := jsui.NewKB()
	//mouse := jsui.NewMouse()
	//leds := jsui.NewLEDs()
	//switches := jsui.NewSwitches()

	m := machine.NewCPU(nil, nil, lcd, nil, console, nil, nil, 20, 0)
	m.SetRAM([]int{
		int(machine.OpPrnII), 42,
		int(machine.OpHalt),
	})
	m.Run()
}
