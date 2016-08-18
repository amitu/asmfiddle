package machine

import "asmfiddle"

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
	// Stack pointer, points to the top of the stack
	ESP int
	// Base pointer, points to the base of the stack
	EBP int
	// Instruction pointer, this is modified with the
	// jump commands, never directly
	EIP int

	EAX int
	EBX int
	ECX int
	EDX int

	ESI int
	EDI int

	R08 int
	R09 int
	R10 int
	R11 int
	R12 int
	R13 int
	R14 int
	R15 int
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
	c.keyboard.OnKey(asmfiddle.KeyboardHandler(c.onKey))
	c.mouse.OnMouse(asmfiddle.MouseHandler(c.onMouse))

	return c
}

func (c *cpu) onKey(asmfiddle.KeyEvent) {

}

func (c *cpu) onMouse(asmfiddle.MouseEvent) {

}

func (c *cpu) RAM() []byte {
	return []byte(c.ram)
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
		op := OpCode(c.ram[c.registers.EIP]) // fetch
		switch op {                          // decode
		case OpMovRI:
			// execute
			arg0 := c.ram[c.registers.EIP+1]
			arg1 := c.ram[c.registers.EIP+2]
			c.registers.EIP += 3
		default:
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
)
