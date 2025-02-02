package main

import (
	"os"
	"fmt"
	"errors"
	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/slog"
	"math/rand"
)

type Machine struct {
	Reg_V [16]uint8
	I uint16
	DisplayBuf [32]uint64
	Mem [4096]uint8	
	Stack [16]uint16
	SP uint8
	PC uint16
	DT uint8
	ST uint8
}

var KeyMap = map[uint8]int32 {
	0: rl.KeyKp0,
	1: rl.KeyKp7,
	2: rl.KeyKp8,
	3: rl.KeyKp9,
	4: rl.KeyKp4,
	5: rl.KeyKp5,
	6: rl.KeyKp6,
	7: rl.KeyKp1,
	8: rl.KeyKp2,
	9: rl.KeyKp3,
	0xA: rl.KeyZ,
	0xB: rl.KeyX,
	0xC: rl.KeyC,
	0xD: rl.KeyV,
	0xE: rl.KeyB,
	0xF: rl.KeyN,
}

var RKeyMap = map[int32]uint8 {
	rl.KeyKp0: 0,   
	rl.KeyKp1: 7,   
	rl.KeyKp2: 8,   
	rl.KeyKp3: 9,   
	rl.KeyKp4: 4,   
	rl.KeyKp5: 5,   
	rl.KeyKp6: 6,   
	rl.KeyKp7: 1,   
	rl.KeyKp8: 2,   
	rl.KeyKp9: 3,   
	rl.KeyZ: 0xA,   
	rl.KeyX: 0xB,   
	rl.KeyC: 0xC,   
	rl.KeyV: 0xD,   
	rl.KeyB: 0xE,   
	rl.KeyN: 0xF,   
}

func (m *Machine) Init(rom_path string) {
	fp, err := os.Open(rom_path)
	check(err)	
	
	n, err := fp.Read(m.Mem[0x200:])
	check(err)	

	m.PC = 0x200
	m.SP = 0x0
	fmt.Printf("%d bytes read\n", n)

	digits := []uint8 {
		0xF0,
        0x90,
        0x90,
        0x90,
        0xF0,
        0x20,
        0x60,
        0x20,
        0x20,
        0x70,
        0xF0,
        0x10,
        0xF0,
        0x80,
        0xF0,
        0xF0,
        0x10,
        0xF0,
        0x10,
        0xF0,
        0x90,
        0x90,
        0xF0,
        0x10,
        0x10,
        0xF0,
        0x80,
        0xF0,
        0x10,
        0xF0,
        0xF0,
        0x80,
        0xF0,
        0x90,
        0xF0,
        0xF0,
        0x10,
        0x20,
        0x40,
        0x40,
        0xF0,
        0x90,
        0xF0,
        0x90,
        0xF0,
        0xF0,
        0x90,
        0xF0,
        0x10,
        0xF0,
        0xF0,
        0x90,
        0xF0,
        0x90,
        0x90,
        0xE0,
        0x90,
        0xE0,
        0x90,
        0xE0,
        0xF0,
        0x80,
        0x80,
        0x80,
        0xF0,
        0xE0,
        0x90,
        0x90,
        0x90,
        0xE0,
        0xF0,
        0x80,
        0xF0,
        0x80,
        0xF0,
        0xF0,
        0x80,
        0xF0,
        0x80,
        0x80,
	
	}

	for i := range digits {
        m.Mem[i] = digits[i]
	}
}



func (m *Machine) clearDisplay() {
	for i := 0; i < len(m.DisplayBuf); i++ {
		m.DisplayBuf[i] = 0;
	}
}

func (m *Machine) drawSprite(vx uint8, vy uint8, n uint8) {
	x := m.Reg_V[vx]
	y := m.Reg_V[vy]
	start := m.I 
	// fmt.Printf("sprite offset: x: %x, y: %x\n", x, y)

	var i uint8
	var j uint8
	m.Reg_V[15] = 0 
	for i = 0; i < n; i++ {
		if y + i > 31 {
			// below screen
			continue
		}
		line := m.Mem[start + uint16(i)]

		for j = 0; j < 8; j++ {
			var bit uint64
			var bitmask uint8 
			bitmask = 1 << 7
			bit = uint64(line & (bitmask >> j)) 
			if bit != 0{
				bit = (1 << 63)
			}
			rshift := x + j
			bit >>= rshift
			// fmt.Printf("display bitmask: %x\n", bit)
			if m.DisplayBuf[y + i] & bit != 0{
				m.Reg_V[15] = 1
			}
			m.DisplayBuf[y + i] ^= bit
		}
	}

	/* Debug sprite */
	/*
	fmt.Printf("start: %x, x: %x, y: %x\n", start, x, y)
	for i = 0; i < n; i++ {
		line := m.Mem[start + uint16(i)]
		fmt.Printf("line: %02x\n", line)

		for j = 0; j < 8; j++ {
			var bit uint64
			var bitmask uint8 
			bitmask = 1 << 7
			bit = uint64(line & (bitmask >> j)) 
			if bit != 0{
				fmt.Printf("*")
			} else {
				fmt.Printf(" ")
			}			
		}
		fmt.Println("")
	}
	*/

	/* Debug Display buffer */
	// for i = 0; i < 32; i++ {
	// 	line := m.DisplayBuf[i]
	// 	for j = 0; j < 64; j++ {
	// 		var bit uint64
	// 		bit = uint64(line & (1 << 63 >> j))
	// 		if bit != 0{
	// 			fmt.Printf("*")
	// 		} else {
	// 			fmt.Printf(" ")
	// 		}
	// 	}
	// 	fmt.Printf("\n")
	// }
			
}

func (m *Machine) TimerUpdate() {
	if m.DT > 0 {
		m.DT--
	}

	if m.ST > 0 {
		m.ST--
		AudioUpdate()
		if m.ST <= 0 {
			AudioStop()
		}
	}
}

func (m *Machine) Clocktick() error {
		
	
	instr1, instr2 := m.Mem[m.PC], m.Mem[m.PC + 1]

	slog.Debug("%04x : ", m.PC)

	switch {		
	/*
		Load instructions
	*/
	case instr1 >> 4 == 6:
		slog.Debug("Load immediate instruction")
		reg := instr1 & 0xF
		val := instr2
		m.Reg_V[reg] = val
	case instr1 >> 4 == 0xA:
		slog.Debug("Load index instruction")
		var val uint16
		val = uint16(instr1) & 0xF
		val <<= 8
		val |= uint16(instr2)
		m.I = val //+ 0x200
	
	case (instr1 >> 4 == 0xF) && (instr2 == 0x29):
		slog.Debug("Load digit sprite")
		reg := instr1 & 0xF
		m.I = uint16(m.Reg_V[reg]) * 5
		
	/*
		Store instructions
	*/
	case (instr1 >> 4 == 0x8) && (instr2 & 0xF == 0):
		slog.Debug("Store direct instruction")
		regx := instr1 & 0xF
		regy := instr2 >> 4
		m.Reg_V[regx] = m.Reg_V[regy]

	case (instr1 >> 4 == 0x8) && (instr2 & 0xF == 1):
		slog.Debug("Store ORed instruction")
		regx := instr1 & 0xF
		regy := instr2 >> 4
		m.Reg_V[regx] |= m.Reg_V[regy] 

	case (instr1 >> 4 == 0x8) && (instr2 & 0xF == 2):
		slog.Debug("Store ANDed instruction")
		regx := instr1 & 0xF
		regy := instr2 >> 4
		m.Reg_V[regx] &= m.Reg_V[regy]

	case (instr1 >> 4 == 0x8) && (instr2 & 0xF == 3):
		slog.Debug("Store XORed instruction")
		regx := instr1 & 0xF
		regy := instr2 >> 4
		m.Reg_V[regx] ^= m.Reg_V[regy]

	case (instr1 >> 4 == 0x8) && (instr2 & 0xF == 4):
		slog.Debug("Store ADDed instruction")
		regx := instr1 & 0xF
		regy := instr2 >> 4
		tmp := uint16(m.Reg_V[regx]) + uint16(m.Reg_V[regy])
		m.Reg_V[regx] += m.Reg_V[regy]
		if tmp > 0xFF {
			m.Reg_V[0xF] = 1
		} else {
			m.Reg_V[0xF] = 0
		}
		
	case (instr1 >> 4 == 0x8) && (instr2 & 0xF == 5):
		slog.Debug("Store SUBed instruction")
		regx := instr1 & 0xF
		regy := instr2 >> 4
		flag := m.Reg_V[regx] >= m.Reg_V[regy]
		m.Reg_V[regx] -= m.Reg_V[regy]
		if flag {
			m.Reg_V[15] = 1
		} else {
			m.Reg_V[15] = 0
		}
		slog.Debug("Register file", m.Reg_V)

	case (instr1 >> 4 == 0x8) && (instr2 & 0xF == 6):
		slog.Debug("Store Right shift instruction")
		regx := instr1 & 0xF
		// regy := instr2 >> 4
		flag := m.Reg_V[regx] % 2
		m.Reg_V[regx] >>= 1
		m.Reg_V[15] = flag

	case (instr1 >> 4 == 0x8) && (instr2 & 0xF == 7):
		slog.Debug("Store SUBN instruction")
		regx := instr1 & 0xF
		regy := instr2 >> 4
		y, x := m.Reg_V[regy], m.Reg_V[regx]
		m.Reg_V[regx] = m.Reg_V[regy] - m.Reg_V[regx]
		if y >= x {
			m.Reg_V[15] = 1
		} else {
			m.Reg_V[15] = 0
		}


	case (instr1 >> 4 == 0x8) && (instr2 & 0xF == 0xE):
		slog.Debug("Store Left shift instruction")
		regx := instr1 & 0xF
		// regy := instr2 >> 4
		flag := m.Reg_V[regx] >> 7
		m.Reg_V[regx] <<= 1
		m.Reg_V[15] = flag
		
	/*
		Add immediate
	*/
	case instr1 >> 4 == 0x7:
		reg := instr1 & 0xF
		val := instr2
		m.Reg_V[reg] += val
		slog.Debug("Add immediate instruction")

	/*
		Jumps
	*/

	case instr1 >> 4 == 0x1:
		addr := uint16(instr2)
		addr2 := (uint16(instr1) & 0xF) << 8
		addr |= addr2
		slog.Debug("Jump instruction")
		// fmt.Printf("Jmp address: %x:%x => %04x\n", addr2, instr2, addr)
		m.PC = addr
		return nil

	/*
		Conditional Jumps
	*/

	case instr1 >> 4 == 0x3:
		reg := instr1 & 0xF
		val := instr2
		if m.Reg_V[reg] == val {
			m.PC += 2
		}
		slog.Debug("Skip if equal imm")

	case instr1 >> 4 == 0x4:
		reg := instr1 & 0xF
		val := instr2
		if m.Reg_V[reg] != val {
			m.PC += 2
		}
		slog.Debug("Skip if not equal imm")

	case instr1 >> 4 == 0x5:
		reg1 := instr1 & 0xF
		reg2 := (instr2 & 0xF0) >> 4
		if m.Reg_V[reg1] == m.Reg_V[reg2] {
			m.PC += 2
		}
		slog.Debug("Skip if equal")

	case instr1 >> 4 == 0x9:
		reg1 := instr1 & 0xF
		reg2 := (instr2 & 0xF0) >> 4
		if m.Reg_V[reg1] != m.Reg_V[reg2] {
			m.PC += 2
		}
		slog.Debug("Skip if equal")

	/* 
		Subroutine
	*/
	case instr1 >> 4 == 0x2:
		addr1 := instr1 & 0xF
		addr2 := instr2
		var addr uint16
		addr = uint16(addr2) | (uint16(addr1) & 0xF << 8)
		m.Stack[m.SP] = m.PC
		m.SP += 1
		m.PC = addr
		slog.Debug("Subroutine call")
		slog.Debug("sp: %x addr: %x, return: %x\n", m.SP, m.PC, m.Stack[m.SP - 1])
		return nil

	case instr1 == 0 && instr2 == 0xEE:
		slog.Debug("Subroutine return")
		
		m.SP -= 1
		ret_addr := m.Stack[m.SP]
		m.PC = ret_addr

		fmt.Printf("SP: %x, ret: %x\n", m.SP, ret_addr)

	/* 
		Special instructions
	*/
	case instr1 == 0 && instr2 == 0xE0:
		slog.Debug("Clear screen instruction")
		m.clearDisplay()

	case instr1 >> 4 == 0xC:
		slog.Debug("Get Random byte instruction")
		vx := instr1 & 0xF
		n := instr2
		rand_byte := rand.Intn(0xFF)
		m.Reg_V[vx] = uint8(rand_byte) & n

	case instr1 >> 4 == 0xD:
		slog.Debug("Draw sprite instruction")
		
		vx := instr1 & 0xF
		vy := (instr2 & 0xF0) >> 4
		n := instr2 & 0xF
		// fmt.Printf("instr1: %x, instr2: %x, vx: %x, vy: %x, n: %x\n", 
		// 	instr1, instr2, vx, vy, n)
		m.drawSprite(vx, vy, n)

	case (instr1 >> 4 == 0xF) && (instr2 == 0x65):
		slog.Debug("Load regs from memory instruction")
		
		x := uint8(instr1) & 0xF
		var i uint8
		for i = 0; i <= x; i++ {
			m.Reg_V[i] = m.Mem[m.I + uint16(i)]
		}
		
	case (instr1 >> 4 == 0xF) && (instr2 == 0x55):
		slog.Debug("Store regs to memory instruction")
		
		x := uint8(instr1) & 0xF
		var i uint8
		for i = 0; i <= x; i++ {
			m.Mem[m.I + uint16(i)] = m.Reg_V[i]
		}
	
	case (instr1 >> 4 == 0xF) && (instr2 == 0x33):
		slog.Debug("Store BCD of Vx to memory instruction")
		
		regx := instr1 & 0xF
		val := m.Reg_V[regx]
		h, t, o := val / 100, (val % 100) / 10, (val % 10)
		m.Mem[m.I] = h
		m.Mem[m.I+1] = t
		m.Mem[m.I+2] = o

	case (instr1 >> 4 == 0xF) && (instr2 == 0x1E):
		slog.Debug("Add Vx to I instruction")
		regx := instr1 & 0xF
		val := m.Reg_V[regx]
		m.I += uint16(val)

	/*
		Keypad instructions
	*/
	case (instr1 >> 4 == 0xE) && (instr2 == 0x9E):
		slog.Debug("Key SKP instruction")
		reg := instr1 & 0xF
		keycode := KeyMap[m.Reg_V[reg]]
		if rl.IsKeyDown(keycode) == true {
			m.PC += 2
		}
		
	case (instr1 >> 4 == 0xE) && (instr2 == 0xA1):
		slog.Debug("Key SKNP instruction")
		reg := instr1 & 0xF
		keycode := KeyMap[m.Reg_V[reg]]
		if rl.IsKeyDown(keycode) != true {
			m.PC += 2
		}
	
	case (instr1 >> 4 == 0xF) && (instr2 == 0x0A):
		slog.Debug("GetKey instruction")
		reg := instr1 & 0xF
		var key int32
		for key != 0{
			key = rl.GetKeyPressed()			
		}
		
		ikey := RKeyMap[key]
		m.Reg_V[reg] = ikey

	/*
		Timer Instructions
	*/
	case (instr1 >> 4 == 0xF) && (instr2 == 0x15):
		slog.Debug("Set Timer instruction")
		reg := instr1 & 0xF
		m.DT = m.Reg_V[reg]
		
	case (instr1 >> 4 == 0xF) && (instr2 == 0x07):
		slog.Debug("Get Timer instruction")
		reg := instr1 & 0xF
		m.Reg_V[reg] = m.DT

	case (instr1 >> 4 == 0xF) && (instr2 == 0x18):
		slog.Debug("Set Sound Timer instruction")
		reg := instr1 & 0xF
		m.ST = m.Reg_V[reg]
		AudioPlay()
		fmt.Println("Set Sound Timer instruction", m.Reg_V[reg])

	default:
		fmt.Printf("Unknown instruction: %02x:%02x\n", instr1, instr2)
		return errors.New("Unknown instruction")
	}

	m.PC += 2
	
	
	return nil
	
}

func (m *Machine) Deinit() {
	
}