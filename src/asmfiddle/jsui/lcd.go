package jsui

import "asmfiddle"

type lcd struct {
}

func (l *lcd) init() error {
	// TODO: create canvas, store its reference
	return nil
}

func (l *lcd) Write(video []int) {
	// TODO: write to canvas
}

func NewLCD() (asmfiddle.LCD, error) {
	display := &lcd{}
	return display, display.init()
}
