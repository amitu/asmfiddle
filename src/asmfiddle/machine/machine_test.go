package machine

import (
	"fmt"
	"testing"
)

func TestCpu_Run(t *testing.T) {
	cpum := NewCPU(nil, nil, nil, nil, 20)
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


type FakeConsole struct {
	last string
}

func (c *FakeConsole) Print(msg string) {
	c.last = msg
}

func TestCpu_OpMovRI(t *testing.T) {
	cpum := NewCPU(nil, nil, nil, nil, nil, 20, 0)
	c, ok := cpum.(*cpu)
	if !ok {
		t.Fatal("invalid")
	}

	c.ram = ram([]int{
		int(OpMovRI), 4, 42,
		int(OpHalt),
	})
	c.Run()

	if c.registers.EBX() != 42 {
		t.Fatal("test failed")
	}
}

func TestCpu_OpMovRM(t *testing.T) {
	cpum := NewCPU(nil, nil, nil, nil, nil, 30, 0)
	c, _ := cpum.(*cpu)
	c.ram = ram([]int{
		int(OpMovRM), 4, 4000,
		int(OpHalt),
	})
	c.Run()

	if c.registers.EBX() != int(OpMovRM) {
		t.Fatal("test failed")
	}
}

func TestCpu_OpMovRR(t *testing.T) {
	cpum := NewCPU(nil, nil, nil, nil, nil, 30, 0)
	c, _ := cpum.(*cpu)
	c.ram = ram([]int{
		int(OpMovRI), 4, 42,
		int(OpMovRR), 5, 4,
		int(OpHalt),
	})
	c.Run()

	if c.registers.ECX() != 42 {
		t.Fatal("test failed")
	}
}

func TestCpu_OpMovMI(t *testing.T) {
	cpum := NewCPU(nil, nil, nil, nil, nil, 30, 0)
	c, _ := cpum.(*cpu)
	c.ram = ram([]int{
		int(OpMovMI), 4000, 42,
		int(OpHalt),
	})
	c.Run()

	if c.Get(4000) != 42 {
		t.Fatal("test failed")
	}
}

func TestCpu_OpMovMM(t *testing.T) {
	cpum := NewCPU(nil, nil, nil, nil, nil, 30, 0)
	c, _ := cpum.(*cpu)
	c.ram = ram([]int{
		int(OpMovMM), 4000, 4004,
		int(OpHalt),
	})
	c.Run()

	if c.Get(4000) != 4000 {
		t.Fatal("test failed")
	}
}

func TestCpu_OpMovMR(t *testing.T) {
	cpum := NewCPU(nil, nil, nil, nil, nil, 30, 0)
	c, _ := cpum.(*cpu)
	c.ram = ram([]int{
		int(OpMovMR), 4000, 2,
		int(OpHalt),
	})
	c.Run()

	if c.Get(4000) != 4012 {
		t.Fatalf("test failed, got: %d", c.Get(4000))
	}
}

func TestCpu_OpPushI(t *testing.T) {
	cpum := NewCPU(nil, nil, nil, nil, nil, 30, 1)
	c, _ := cpum.(*cpu)
	c.ram = ram([]int{
		int(OpPushI), 42,
		int(OpHalt),
	})
	c.Run()

	if c.Pop() != 42 {
		t.Fatalf("test failed")
	}
}

func TestCpu_OpPrnII(t *testing.T) {
	console := &FakeConsole{}
	cpum := NewCPU(nil, nil, nil, nil, console, 20, 1)

	c, _ := cpum.(*cpu)
	c.ram = ram([]int{
		int(OpPrnII), 42,
		int(OpHalt),
	})
	c.Run()

	if console.last != "42" {
		t.Fatal("test failed")
	}
}

func TestCpu_OpPrnIR(t *testing.T) {
	console := &FakeConsole{}
	cpum := NewCPU(nil, nil, nil, nil, console, 20, 1)

	c, _ := cpum.(*cpu)
	c.ram = ram([]int{
		int(OpPrnIR), 2,
		int(OpHalt),
	})
	c.Run()

	if console.last != "4008" {
		t.Fatal("test failed")
	}
}

func TestCpu_OpPrnIM(t *testing.T) {
	console := &FakeConsole{}
	cpum := NewCPU(nil, nil, nil, nil, console, 20, 0)

	c, ok := cpum.(*cpu)
	if !ok {
		t.Fatal("invalid")
	}

	c.ram = ram([]int{
		int(OpPrnIM), 4000,
		int(OpHalt),
	})
	c.Run()

	if console.last != fmt.Sprintf("%d", OpPrnIM) {
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

