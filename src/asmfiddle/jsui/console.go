package jsui

import (
	"asmfiddle"

	"github.com/gopherjs/gopherjs/js"
)

type console struct {
	pre *js.Object
	log string
}

func (c *console) init() error {
	// store the console div
	c.pre = js.Global.Get("document").Call("getElementById", "console")
	return nil
}

func (c *console) Print(msg string) {
	// write to console div
	js.Global.Get("console").Call("log", msg, c.pre)
	c.log += msg + "\n"
	c.pre.Set("innerText", c.log)
}

func NewConsole() (asmfiddle.Console, error) {
	c := &console{}
	return c, c.init()
}
