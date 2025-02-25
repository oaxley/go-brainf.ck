package main

import (
	"fmt"
	"log"
	"os"
)

// ----- constants
const CodeSize = 1 << 15 // max code size is 32,768 bytes
const DataSize = 1 << 15 // max data size is 32,768 bytes

const (
	DataValueInc byte = 43 // '+' increase value at the data pointer
	DataValueDec byte = 45 // '-' decrease value at the data pointer
	ReadChar     byte = 44 // ',' read a character from the user
	WriteChar    byte = 46 // '.' write a character to the screen
	DataPtrDec   byte = 60 // '<' move data pointer to the left
	DataPtrInc   byte = 62 // '>' move data pointer to the right
	JumpFwd      byte = 91 // '[' jump forward is data value is 0
	JumpBck      byte = 93 // ']' jump backward if data value is not 0
)

// ----- structures
type (
	VMCore struct {
		PC     uint16 // program counter
		DP     uint16 // data pointer
		Code   []byte // code
		Data   []byte // data
		Length int    // code size

		jumps map[uint16]uint16 // hashmap for the jumps
	}
)

// ----- functions

// load a new ROM in the VM
func (core *VMCore) Load(filename string) uint16 {
	// open the file for reading
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to open the file '%s' for reading!", filename)
	}
	defer fh.Close()

	// initialize the VM
	core.PC = 0
	core.DP = 0
	core.Code = make([]byte, CodeSize)
	core.Data = make([]byte, DataSize)

	// read the file character per character
	n, err := fh.Read(core.Code)
	if err != nil {
		log.Fatalln("Unable to read the code!")
	} else {
		core.Length = n
	}

	// compute the jumps
	core.ComputeJumps(core.Length)

	return uint16(core.Length)
}

// compute the jumps
func (core *VMCore) ComputeJumps(size int) {

	core.jumps = make(map[uint16]uint16, 0)

	var s stack
	s = make(stack, 0)

	for counter := 0; counter < size; counter++ {
		opcode := core.Code[counter]

		// jump forward
		if opcode == byte(JumpFwd) {
			s = s.Push(uint16(counter))
		}

		// jump backward
		if opcode == byte(JumpBck) {
			// empty stack -> error
			if s.Count() == 0 {
				log.Fatalln("Error: unbalanced number of '[' and ']' in the source code!")
			}

			// retrieve the last value from the stack
			var last uint16
			s, last = s.Pop()

			// insert in the map
			core.jumps[last] = uint16(counter + 1)
			core.jumps[uint16(counter)] = last + 1
		}
	}
}

// increase the value at the data pointer
func (core *VMCore) DataValueIncrease() {
	value := core.Data[core.DP]
	value = (value + 1) & 0xFF
	core.Data[core.DP] = value
}

// decrease the value at the data pointer
func (core *VMCore) DataValueDecrease() {
	value := core.Data[core.DP]
	value = (value - 1) & 0xFF
	core.Data[core.DP] = value
}

// increase the data pointer
func (core *VMCore) IncreaseDP() {
	core.DP = (core.DP + 1) & (DataSize - 1)
}

// decrease the data pointer
func (core *VMCore) DecreaseDP() {
	core.DP = (core.DP - 1) & (DataSize - 1)
}

// jump forward
func (core *VMCore) JumpFwd() {
	if core.Data[core.DP] == 0 {
		value := core.PC - 1
		core.PC = core.jumps[value]
	}

}

// jump backward
func (core *VMCore) JumpBck() {
	if core.Data[core.DP] != 0 {
		value := core.PC - 1
		core.PC = core.jumps[value]
	}
}

func (core *VMCore) Execute() {

	for core.PC < uint16(core.Length) {
		// read the next opcode and increment the Program Counter
		opcode := core.Code[core.PC]
		core.PC = core.PC + 1

		// opcode lookup
		switch opcode {
		case DataValueInc:
			core.DataValueIncrease()
		case DataValueDec:
			core.DataValueDecrease()

		case DataPtrInc:
			core.IncreaseDP()
		case DataPtrDec:
			core.DecreaseDP()

		case JumpFwd:
			core.JumpFwd()
		case JumpBck:
			core.JumpBck()

		case WriteChar:
			fmt.Printf("%c", core.Data[core.DP])
			continue

		default:
			continue
		}
	}
}
