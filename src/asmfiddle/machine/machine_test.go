package machine

import "testing"

type FakeConsole struct {
	last string
}

func (c *FakeConsole) Print(msg string) {
	c.last = msg
}

func TestCpu_OpMovRI(t *testing.T) {
	console := &FakeConsole{}
	cpum := NewCPU(nil, nil, nil, nil, console, 20)
	c, ok := cpum.(*cpu)
	if !ok {
		t.Fatal("invalid")
	}

	c.ram = ram([]int{int(OpMovRI), 4, 42, int(OpHalt)})
	c.Run()

	if c.registers.data[4] != 42 {
		t.Fatal("test failed")
	}
}

func TestCpu_OpPrnII(t *testing.T) {
	console := &FakeConsole{}
	cpum := NewCPU(nil, nil, nil, nil, console, 20)

	c, ok := cpum.(*cpu)
	if !ok {
		t.Fatal("invalid")
	}

	c.ram = ram([]int{int(OpPrnII), 42, int(OpHalt)})
	c.Run()

	if console.last != "42" {
		t.Fatal("test failed")
	}
}
