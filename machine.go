package main

import (
	"os"
	"fmt"
	"errors"
)

type Machine struct {
	Reg_V [16]uint8
	I uint16
	DisplayBuf [32]uint64
	Mem [4096]uint8
	fp *os.File
	PC uint16
}

func (m *Machine) Init(rom_path string) {
	fp, err := os.Open(rom_path)
	check(err)	
	
	n, err := fp.Read(m.Mem[0x200:])
	check(err)	

	m.PC = 0x200
	fmt.Printf("%d bytes read\n", n)
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
	fmt.Printf("sprite offset: x: %x, y: %x\n", x, y)

	var i uint8
	var j uint8
	for i = 0; i < n; i++ {
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
			m.DisplayBuf[y + i] ^= bit
		}
	}

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

func (m *Machine) Clocktick() error {
		
	
	instr1, instr2 := m.Mem[m.PC], m.Mem[m.PC + 1]

	fmt.Printf("%04x : ", m.PC)

	switch {
	case instr1 == 0 && instr2 == 0xE0:
		fmt.Println("Clear screen instruction")
		m.clearDisplay()
	case instr1 >> 4 == 6:
		fmt.Println("Load immediate instruction")
		reg := instr1 & 0xF
		val := instr2
		m.Reg_V[reg] = val
	case instr1 >> 4 == 0xA:
		fmt.Println("Load index instruction")
		var val uint16
		val = uint16(instr1) & 0xF
		val <<= 16
		val |= uint16(instr2)
		m.I = val + 0x200
		fmt.Printf("Load index val: %x\n", val)
	case instr1 >> 4 == 0xD:
		fmt.Println("Draw sprite instruction")
		
		vx := instr1 & 0xF
		vy := (instr2 & 0xF0) >> 4
		n := instr2 & 0xF
		fmt.Printf("instr1: %x, instr2: %x, vx: %x, vy: %x, n: %x\n", 
			instr1, instr2, vx, vy, n)
		m.drawSprite(vx, vy, n)

	case instr1 >> 4 == 0x1:
		addr := uint16(instr2)
		addr2 := (uint16(instr1) & 0xF) << 8
		addr |= addr2
		fmt.Println("Jump instruction")
		// fmt.Printf("Jmp address: %x:%x => %04x\n", addr2, instr2, addr)
		m.PC = addr
		return nil

	case instr1 >> 4 == 0x7:
		reg := instr1 & 0xF
		val := instr2
		m.Reg_V[reg] += val
		fmt.Println("Add immediate instruction")
	default:
		fmt.Println("Unknown instruction")
		return errors.New("Unknown instruction")
	}

	m.PC += 2
	
	
	return nil
	
}

func (m *Machine) Deinit() {
	m.fp.Close()
}