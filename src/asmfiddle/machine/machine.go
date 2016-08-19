// Memory Layout (peripherals etc)
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
// 2460: LEDs, 32 LEDs
// 2464: Input PINs (pins will be connected to switches)
//         0: Control
//            0: pin 0 toggle will interrupt
//            1: pin 1 toggle will interrupt
//            ...
//
//         2468: Pin High: if this is high then we are interested when
//                         when pin goes up, else down
//
//         2472: Pin Toggle: if this is 0 then both on and off will
//                         trigger interrupt, else just high or low
//
//         2476: pin state
//            0: pin 0
//            1: pin 1
//            ...
//
// 2480 Ethernet:
//        0: control
//           0: on | off
//        2484: mac address
//        2492: Buffer Out: WO
//
// 2496 IPNet: will store IP address on Ctrl + 1, Ctrl + 2 buffer out, which is IP packet
//        0: control
//        2500: IP address
//        2504: out buffer
//
// 2508: TCPNet:
// 		  0: control
//               0: on | off
//               1: client | server
//               2: domain | IP
//               3: port ready
//               4: new connection
//               5: error
// 				16-32 will store Port at Control (lower 16 bits)
//                      read returns actual port
//                      write is illegal
//        2512: remote ip
//        2512: On buffer: buffer size | remote ip | remote port | local ip | local port | data
//
// 2516: HTTPNet:
//      0: control
//         0: on | off
//         1: client | server
//         2: domain | ip
//      +1: remote ip
// 		2520: out Buffer
//
// DNSNet:
//          0: control
//               0: start dns lookup
//          +1: address of null terminated domain name
//          +2: ip address after resolution

// 3000: Interrupt: invalid memory access
// 3004: Interrupt: Keyboard Interrupt
// 3008: Interrupt: mouse
// 3012: timer 0
// 3016: timer 1
// 3020: timer 2
// 3024: timer 3
// 3028: user defined int 0
// 3032: user defined int 1
// 3036: user defined int 2
// 3040: user defined int3
// 3044: pin interrupt
// 3048: net* in interrupt
// 3052: net* out interrupt
// 3056: dns resolved interrupt
//
// 4000: 4kb Net* In
// 8000: main program start
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

	if loc >= 1000 {
		if loc >= 1000+len(c.ram) {
			c.registers.SetEIP(3000)
		}
		c.ram[loc-1000] = val
		return
	}

	c.special[loc] = val
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
	leds     asmfiddle.LEDBank
	switches asmfiddle.SwitchBank

	ram       ram
	registers *registers
	stack     []int

	special [760]int
}

func (c *cpu) SetRAM(data []int) {
	c.ram = data
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

func (c *cpu) Push(val int) {
	c.stack[c.registers.data[0]] = val
	c.registers.data[0] += 1
}

func (c *cpu) Pop() int {
	c.registers.data[0] -= 1
	return c.stack[c.registers.data[0]]
}

func NewCPU(
	keyboard asmfiddle.Keyboard, mouse asmfiddle.Mouse, lcd asmfiddle.LCD,
	fs asmfiddle.FileSystem, console asmfiddle.Console, leds asmfiddle.LEDBank,
	switches asmfiddle.SwitchBank, ramsize int, stacksize int,
) asmfiddle.Machine {
	c := &cpu{
		keyboard: keyboard,
		mouse:    mouse,
		lcd:      lcd,
		fs:       fs,
		console:  console,
		leds:     leds,
		switches: switches,

		ram:       make([]int, ramsize),
		registers: &registers{},
		stack:     make([]int, stacksize),
	}
	c.loadfs()
	c.registers.SetEIP(4000) // main

	// register event handlers
	// c.keyboard.OnKey(asmfiddle.KeyboardHandler(c.onKey))
	// c.mouse.OnMouse(asmfiddle.MouseHandler(c.onMouse))

	return c
}

func (c *cpu) stackDump() string {
	return fmt.Sprintf(`
		stack: %v
		pointer: %d`, c.stack, c.registers.ESP())
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
		case OpMovRM:
			argr, argm := c.readTwo()
			c.registers.Set(argr, c.Get(argm))
		case OpMovRR:
			r1, r2 := c.readTwo()
			c.registers.Set(r1, c.registers.data[r2])
		case OpMovMI:
			argm, argi := c.readTwo()
			c.Set(argm, argi)
		case OpMovMM:
			m1, m2 := c.readTwo()
			c.Set(m1, c.Get(m2))
		case OpMovMR:
			m, r := c.readTwo()
			c.Set(m, c.registers.data[r])

		case OpPushI:
			c.Push(c.readOne())
		case OpPushM:
			c.Push(c.Get(c.readOne()))
		case OpPushR:
			c.Push(c.registers.data[c.readOne()])

		case OpPopR:
			c.registers.data[c.readOne()] = c.Pop()
		case OpPopM:
			c.Set(c.readOne(), c.Pop())

		case OpIncR:
			c.registers.data[c.readOne()]++
		case OpIncM:
			m := c.readOne()
			c.Set(m, c.Get(m)+1)

		case OpDecR:
			c.registers.data[c.readOne()]--
		case OpDecM:
			m := c.readOne()
			c.Set(m, c.Get(m)-1)

		case OpAddRI: // TODO: overflow
			argr, argi := c.readTwo()
			c.registers.Set(argr, argi+c.registers.data[argr])
		case OpAddRR:
			r1, r2 := c.readTwo()
			c.registers.Set(r1, c.registers.data[r1]+c.registers.data[r2])
		case OpAddRM:
			argr, argm := c.readTwo()
			c.registers.Set(argr, c.registers.data[argr]+c.Get(argm))
		case OpAddMI:
			argm, argi := c.readTwo()
			c.Set(argm, c.Get(argm)+argi)
		case OpAddMR:
			m, r := c.readTwo()
			c.Set(m, c.registers.data[r]+c.Get(m))
		case OpAddMM:
			m1, m2 := c.readTwo()
			c.Set(m1, c.Get(m2)+c.Get(m1))

		case OpSubRI: // TODO: overflow
			argr, argi := c.readTwo()
			c.registers.Set(argr, c.registers.data[argr]-argi)
		case OpSubRR:
			r1, r2 := c.readTwo()
			c.registers.Set(r1, c.registers.data[r1]-c.registers.data[r2])
		case OpSubRM:
			argr, argm := c.readTwo()
			c.registers.Set(argr, c.registers.data[argr]-c.Get(argm))
		case OpSubMI:
			argm, argi := c.readTwo()
			c.Set(argm, c.Get(argm)-argi)
		case OpSubMR:
			m, r := c.readTwo()
			c.Set(m, c.Get(m)-c.registers.data[r])
		case OpSubMM:
			m1, m2 := c.readTwo()
			c.Set(m1, c.Get(m1)-c.Get(m2))

		case OpMulRI: // TODO: overflow
			argr, argi := c.readTwo()
			c.registers.Set(argr, argi*c.registers.data[argr])
		case OpMulRR:
			r1, r2 := c.readTwo()
			c.registers.Set(r1, c.registers.data[r1]*c.registers.data[r2])
		case OpMulRM:
			argr, argm := c.readTwo()
			c.registers.Set(argr, c.registers.data[argr]*c.Get(argm))
		case OpMulMI:
			argm, argi := c.readTwo()
			c.Set(argm, c.Get(argm)*argi)
		case OpMulMR:
			m, r := c.readTwo()
			c.Set(m, c.registers.data[r]*c.Get(m))
		case OpMulMM:
			m1, m2 := c.readTwo()
			c.Set(m1, c.Get(m2)*c.Get(m1))

		case OpDivRI: // TODO: overflow etc
			argr, argi := c.readTwo()
			c.registers.Set(argr, c.registers.data[argr]/argi)
		case OpDivRR:
			r1, r2 := c.readTwo()
			c.registers.Set(r1, c.registers.data[r1]/c.registers.data[r2])
		case OpDivRM:
			argr, argm := c.readTwo()
			c.registers.Set(argr, c.registers.data[argr]/c.Get(argm))
		case OpDivMI:
			argm, argi := c.readTwo()
			c.Set(argm, c.Get(argm)/argi)
		case OpDivMR:
			m, r := c.readTwo()
			c.Set(m, c.Get(m)/c.registers.data[r])
		case OpDivMM:
			m1, m2 := c.readTwo()
			c.Set(m1, c.Get(m1)/c.Get(m2))

		case OpPrnII:
			argi := c.readOne()
			c.console.Print(fmt.Sprintf("%d", argi))
		case OpPrnIR:
			r := c.readOne()
			c.console.Print(fmt.Sprintf("%d", c.registers.data[r]))
		case OpPrnIM:
			m := c.readOne()
			c.console.Print(fmt.Sprintf("%d", c.Get(m)))
		case OpHalt:
			return
		}
	}
}

type OpCode int

func (o OpCode) String() string {
	switch o {
	case OpMovRI:
		return "MovRI"
	case OpPrnII:
		return "PrnII"
	case OpHalt:
		return "Halt"
	default:
		return fmt.Sprintf("Unknown: %d", o)
	}
}

// I -> immediate
// R -> register
// M -> memory

const (
	// Invalid should never come in code
	OpInv OpCode = iota
	OpMovRI
	OpMovRR
	OpMovRM
	OpMovMI
	OpMovMR
	OpMovMM

	OpPushI
	OpPushR
	OpPushM

	OpPopR
	OpPopM

	OpCallI
	OpCallM
	OpCallR

	OpRet

	OpIncR
	OpIncM

	OpDecR
	OpDecM

	OpAddRI
	OpAddRR
	OpAddRM
	OpAddMI
	OpAddMR
	OpAddMM

	OpSubRI
	OpSubRM
	OpSubRR
	OpSubMI
	OpSubMM
	OpSubMR

	OpMulRI
	OpMulRR
	OpMulRM
	OpMulMI
	OpMulMR
	OpMulMM

	// Divides arg0 by arg1, storing the quotient in arg0
	OpDivRI
	OpDivRR
	OpDivRM
	OpDivMI
	OpDivMM
	OpDivMR

	// Same as the '%' (modulus) operator in C. Calculates arg0 mod arg1 and
	// stores the result in the remainder register.
	OpModRI
	OpModRM
	OpModRR
	OpModMI
	OpModMR
	OpModMM

	// Retrieves the value stored in the remainder register, storing it in arg
	OpRemR
	OpRemM

	OpNotR
	OpNotM

	OpXorRI
	OpXorRM
	OpXorRR
	OpXorMI
	OpXorMR
	OpXorMM

	OpOrRI
	OpOrRM
	OpOrRR
	OpOrMI
	OpOrMM
	OpOrMR

	OpAndRI
	OpAndRM
	OpAndRR
	OpAndMI
	OpAndMM
	OpAndMR

	OpShlRI
	OpShlRR
	OpShlRM
	OpShlMI
	OpShlMR
	OpShlMM

	OpShrRI
	OpShrRR
	OpShrRM
	OpShrMI
	OpShrMR
	OpShrMM

	OpCmpII
	OpCmpIR
	OpCmpIM
	OpCmpRI
	OpCmpRR
	OpCmpRM
	OpCmpMI
	OpCmpMR
	OpCmpMM

	OpJmpI
	OpJmpR
	OpJmpM

	OpJEI
	OpJER
	OpJEM

	OpJneI
	OpJneR
	OpJneM

	OpJgI
	OpJgR
	OpJgM

	OpJgeI
	OpJgeR
	OpJgeM

	OpJlI
	OpJlR
	OpJlM

	OpJleI
	OpJleR
	OpJleM

	OpPrnII
	OpPrnIR
	OpPrnIM
	OpPrnSI
	OpPrnSR
	OpPrnSM

	OpIntI
	OpIntM
	OpIntR

	// pause the program till some reg/memory is zero/nonzero
	// only handle interrupts
	OpPauseRZ
	OpPauseRNz
	OpPauseMZ
	OpPauseMNz

	// terminate the program
	OpHalt
)
