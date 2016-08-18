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
	cpum := NewCPU(nil, nil, nil, nil, console, 20, 0)
	c, ok := cpum.(*cpu)
	if !ok {
		t.Fatal("invalid")
	}

	c.ram = ram([]int{int(OpMovRI), 4, 42, int(OpHalt)})
	c.Run()

	if c.registers.EBX() != 42 {
		t.Fatal("test failed")
	}
}

func TestCpu_OpPrnII(t *testing.T) {
	console := &FakeConsole{}
	cpum := NewCPU(nil, nil, nil, nil, console, 20, 0)

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

func TestCpu_Stack(t *testing.T) {
	cpum := NewCPU(nil, nil, nil, nil, nil, 20, 4)
	c, _ := cpum.(*cpu)

	c.Push(10)
	if c.Pop() != 10 {
		t.Fatal("Pop failed")
	}

	c.Push(1)
	c.Push(2)
	c.Push(3)
	c.Push(4)

	if c.Pop() != 4 {
		t.Fatal("Pop failed")
	}
	if c.Pop() != 3 {
		t.Fatal("Pop failed")
	}
	if c.Pop() != 2 {
		t.Fatal("Pop failed")
	}
	if c.Pop() != 1 {
		t.Fatal("Pop failed")
	}
}
