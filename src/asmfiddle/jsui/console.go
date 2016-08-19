package jsui

import (
	"asmfiddle"

	"github.com/gopherjs/gopherjs/js"
)

type console struct {
	pre *js.Object
}

func (c *console) init() error {
	// store the console div
	return nil
}

func (c *console) Print(msg string) {
	// write to console div
	js.Global.Get("console").Call("log", msg)
}

func NewConsole() (asmfiddle.Console, error) {
	c := &console{}
	return c, c.init()
}
