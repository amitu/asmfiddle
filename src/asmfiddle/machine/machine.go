// Memory Layout
//
// Peripherals
// 0000000 LCD  [320 x 240]
//         0: Control
//            0: on or off
//            1: clear the LCD
//            2: clear color
//            3: show test message on LCD
//         4-2404: Video Memory
//
// 0002408 Keyboard:
//         0: Control
//            0: on or off
//         2412: Code
//  	       2416: Down
//
// 0002420 Mouse:
// 	    0: Control
// 	    	   0: on or off
// 	    	2424: pos x
// 	    	2428: pos y
// 	    	2432: Mask
//
// 2436: Timer
//         0: Control
//             0: timer 0 is active
//             1: timer 1 is active
//             2: timer 2 is active
//             3: timer 3 is active
//          2440: current time
//          2444: timer 0 deadline
//          2448: timer 1 deadline
//          2452: timer 2 deadline
//          2456: timer 3 deadline
//
// 3000: Interrupt: invalid memory access
// 3004: Interrupt: Keyboard Interrupt
// 3008: Interrupt: mouse
// 3012: timer 0
// 3016: timer 1
// 3020: timer 2
// 3024: timer 3
//
// 4000: main program start
package machine

import (
	"asmfiddle"

	"fmt"
	"time"
)

type ram []int

func (c *cpu) Set(loc, val int) {
	if loc%4 != 0 {
		// set some flag to indicate what has happened
		c.registers.SetEIP(3000)
	}

	if loc == 2412 || loc == 2416 || loc == 2424 || loc == 2428 || loc == 2432 || loc == 2440 {
		c.registers.SetEIP(3000)
	}

	loc = loc / 4
}

func (c *cpu) Get(loc int) int {
	if loc%4 != 0 {
		// set some flag to indicate what has happened
		c.registers.SetEIP(3000)
	}

	loc = loc / 4

	if loc == 2444 {
		return int((time.Now().UnixNano() / 1000000) & 0xFFFFFFFF)
	}

	if loc >= 1000 {
		if loc >= 1000+len(c.ram) {
			c.registers.SetEIP(3000)
		}
		return c.ram[loc-1000]
	}

	return c.special[loc]
}

type cpu struct {
	keyboard asmfiddle.Keyboard
	mouse    asmfiddle.Mouse
	lcd      asmfiddle.LCD
	console  asmfiddle.Console
	fs       asmfiddle.FileSystem

	ram       ram
	registers *registers
	stack     *stack

	special [756]int
}

type registers struct {
	data [18]int
}

func (r *registers) ESP() int {
	return r.data[0]
}

func (r *registers) EBP() int {
	return r.data[1]
}

func (r *registers) EIP() int {
	return r.data[2]
}

func (r *registers) EAX() int {
	return r.data[3]
}

func (r *registers) EBX() int {
	return r.data[4]
}

func (r *registers) ECX() int {
	return r.data[5]
}

func (r *registers) EDX() int {
	return r.data[6]
}

func (r *registers) ESI() int {
	return r.data[7]
}

func (r *registers) EDI() int {
	return r.data[8]
}

func (r *registers) R08() int {
	return r.data[9]
}

func (r *registers) R09() int {
	return r.data[10]
}

func (r *registers) R10() int {
	return r.data[11]
}

func (r *registers) R11() int {
	return r.data[12]
}

func (r *registers) R12() int {
	return r.data[13]
}

func (r *registers) R13() int {
	return r.data[14]
}

func (r *registers) R14() int {
	return r.data[15]
}

func (r *registers) R15() int {
	return r.data[16]
}

func (r *registers) FLAGS() int {
	return r.data[17]
}

func (r *registers) SetEIP(val int) {
	r.data[2] = val
}

func (r *registers) IncrEIP(incr int) {
	r.data[2] += incr * 4
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
	R15 = %d

	FLAGS = %b`, r.data[0], r.data[1], r.data[2], r.data[3], r.data[4], r.data[5],
		r.data[6], r.data[7], r.data[8], r.data[9], r.data[10], r.data[11], r.data[12],
		r.data[13], r.data[14], r.data[15], r.data[16], r.data[17])
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
	fs asmfiddle.FileSystem, console asmfiddle.Console, size int,
) asmfiddle.Machine {
	c := &cpu{
		keyboard: keyboard,
		mouse:    mouse,
		lcd:      lcd,
		fs:       fs,
		console:  console,

		ram:       make([]int, size),
		registers: &registers{},
		stack:     &stack{},
	}
	c.loadfs()
	c.registers.SetEIP(4000) // main

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
	return c.registers
}

func (c *cpu) Stack() ([]int, int) {
	return nil, 0
}

func (c *cpu) loadfs() {
	// iterate through c.fs, if any file has a name <int>.ext or <int>
	// load the content of that file in c.ram.

	// eg 2000.txt contains "hello world", so write ram[2000:2012] = "hello world\0"
}

func (c *cpu) readOp() OpCode {
	op := OpCode(c.Get(c.registers.EIP()))
	c.registers.IncrEIP(1)
	return op
}

func (c *cpu) readOne() int {
	one := c.Get(c.registers.EIP())
	c.registers.IncrEIP(1)
	return one
}

func (c *cpu) readTwo() (int, int) {
	one := c.Get(c.registers.EIP())
	two := c.Get(c.registers.EIP() + 4)
	c.registers.IncrEIP(2)
	return one, two
}

func (c *cpu) Run() {
	for {
		if c.registers.EIP() > 4000+len(c.ram)*4 {
			return
		}

		op := c.readOp() // fetch
		switch op {      // decode
		case OpMovRI:
			// execute
			argr, argi := c.readTwo()
			c.registers.Set(argr, argi)
		case OpPrnII:
			argi := c.readOne()
			c.console.Print(fmt.Sprintf("%d", argi))
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
	OpPrnII
	OpPrnIR
	OpPrnIM
	OpPrnSI
	OpPrnSR
	OpPrnSM
	OpInt
	OpHalt
)
