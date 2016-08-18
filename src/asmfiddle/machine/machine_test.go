package machine

import "testing"

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
