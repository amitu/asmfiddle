package machine

import (
	"asmfiddle"
	"fmt"
)

type ram []int

func (r ram) Set(pos, val int) {
	// where are our peripherals?
	// keyboard: 1000
	//      on off | int address | context
	//Set(1000, 1) // activate keyboard
	//Set(1004, 1) // use interrupt 1 for keyboard events
	//Set(1008, 1)

	// mouse 1020
	//       on off | int address | context
	//Set(1020, 1) // activate mouse
	//Set(1024, 2) // use interrupt 1 for mouse events
	//Set(1028, 2)
	// how many interrupts?
}

func (r ram) Get(int) int {
	return 0
}

type cpu struct {
	keyboard asmfiddle.Keyboard
	mouse    asmfiddle.Mouse
	lcd      asmfiddle.LCD
	fs       asmfiddle.FileSystem

	ram       ram
	registers *registers
	stack     *stack
}

type registers struct {
	data [17]int
}

func (r *registers) EIP() int {
	return r.data[2]
}

func (r *registers) IncrEIP(incr int) {
	r.data[2] += incr
}

func (r *registers) String() string {
	return fmt.Sprintf(`
	// Stack pointer, points to the top of the stack
	ESP = %d
	// Base pointer, points to the base of the stack
	EBP = %d
	// Instruction pointer, this is modified with the
	// jump commands, never directly
	EIP = %d

	EAX = %d
	EBX = %d
	ECX = %d
	EDX = %d

	ESI = %d
	EDI = %d

	R08 = %d
	R09 = %d
	R10 = %d
	R11 = %d
	R12 = %d
	R13 = %d
	R14 = %d
	R15 = %d`, r.data[0], r.data[1], r.data[2], r.data[3], r.data[4], r.data[5],
		r.data[6], r.data[7], r.data[8], r.data[9], r.data[10], r.data[11], r.data[12],
		r.data[13], r.data[14], r.data[15], r.data[16])
}

func (r *registers) Set(i, val int) {
	r.data[i] = val
}

type stack struct {
	sp    int
	stack []int
}

func (s *stack) Push(val int) {
	s.stack = append(s.stack, val)
	s.sp += 1
}

func (s *stack) Pop() int {
	return s.stack[s.sp]
}

func NewCPU(
	keyboard asmfiddle.Keyboard, mouse asmfiddle.Mouse, lcd asmfiddle.LCD,
	fs asmfiddle.FileSystem, size int,
) asmfiddle.Machine {
	c := &cpu{
		keyboard:  keyboard,
		mouse:     mouse,
		lcd:       lcd,
		fs:        fs,
		ram:       make([]int, size),
		registers: &registers{},
		stack:     &stack{},
	}
	c.loadfs()

	// register event handlers
	// c.keyboard.OnKey(asmfiddle.KeyboardHandler(c.onKey))
	// c.mouse.OnMouse(asmfiddle.MouseHandler(c.onMouse))

	return c
}

func (c *cpu) onKey(asmfiddle.KeyEvent) {

}

func (c *cpu) onMouse(asmfiddle.MouseEvent) {

}

func (c *cpu) RAM() []int {
	return []int(c.ram)
}

func (c *cpu) Registers() asmfiddle.Registers {
	return nil
}

func (c *cpu) Stack() ([]int, int) {
	return nil, 0
}

func (c *cpu) loadfs() {
	// iterate through c.fs, if any file has a name <int>.ext or <int>
	// load the content of that file in c.ram.

	// eg 2000.txt contains "hello world", so write ram[2000:2012] = "hello world\0"
}

func (c *cpu) Run() {
	for {
		if c.registers.EIP() > len(c.ram) {
			return
		}

		op := OpCode(c.ram[c.registers.EIP()]) // fetch
		switch op {                            // decode
		case OpMovRI:
			// execute
			argr := c.ram[c.registers.EIP()+1]
			argi := c.ram[c.registers.EIP()+2]
			c.registers.IncrEIP(3)
			c.registers.Set(argr, argi)
		case OpHalt:
			return
		}
	}
}

type OpCode int

// I -> immediate
// R -> register
// M -> memory

const (
	OpMovRI OpCode = iota
	OpMovRR
	OpMovRM
	OpMovMI
	OpMovMR
	OpMovMM
	OpPush
	OpPop
	OpPushf
	OpPopf
	OpCall
	OpRet
	OpInc
	OpDec
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpRem
	OpNot
	OpXor
	OpOr
	OpAnd
	OpShl
	OpShr
	OpCmp
	OpJmp
	OpJE
	OpJne
	OpJg
	OpJge
	OpJl
	OpJle
	OpPrn
	OpInt
	OpHalt
)
