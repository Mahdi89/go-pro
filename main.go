// Description of a 4-stage in-order processor based on the
// Manchester Small-Scale Experimental Machine (SSEM) ISA
// The SSEM has a 32-bit word length and a memory of 32 words.

package main

import (
	"bufio"
	//	"bytes"
	//	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type word uint32
type LineAddress uint8
type CRTAddress uint16
type SSEMFunc uint8

// Global memory and state variables
var Memory [100]word
var PC, PC_step LineAddress
var ACC word
var MDR word
var Stopped bool

const (
	JMP SSEMFunc = iota
	JRP
	LDN
	STO
	SUB
	SUB_alt
	TEST
	STOP
)

type SSEMInst struct {
	LineNo LineAddress
	CRTNo  CRTAddress
	Func   SSEMFunc
}

func Mem(fptr string) int {

	// Fetch data/inst. from file
	// TODO read from MEMR chan.
	file, err := os.Open(fptr)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	index := 0
	var h0, h1, h2, h3 uint8
	var i, i_ string

	for scanner.Scan() {
		str := scanner.Text()
		fmt.Sscanf(str, "%4s%4s", &i_, &i)
		fmt.Sscanf(i, "%1x%1x%1x%1x", &h0, &h1, &h2, &h3)
		Memory[index] = word((h0<<8|h1)<<8 | (h2<<8 | h3))
		index += 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()
	return index
}

func Fetch(FD chan<- word) {

	for {
		FD <- Memory[PC]
		if PC < 5 {
			PC++
		}
	}
}

func Decode(FD <-chan word, DE chan<- SSEMInst) {

	for {
		IR := <-FD
		LA := LineAddress(IR & 0x001f)
		CA := CRTAddress((IR & 0x0fff) >> 4)
		Func := SSEMFunc(IR >> 24)

		DE <- SSEMInst{
			LineNo: LA,
			CRTNo:  CA,
			Func:   Func}
	}
}

func Execute(DE <-chan SSEMInst, EW chan<- SSEMInst) {

	for {
		IR := <-DE
		//LA := IR.LineNo
		//CA := IR.CRTNo
		Func := IR.Func

		// Execute the decoded instruction

		switch SSEMFunc(Func) {
		case JMP:
			// Jump - update PC
			// JMP then MemoryRead (); ZeroPC (); AddMDRToPC ()
		case JRP:
			// Relative Jummp
			// JMP then MemoryRead (); ZeroPC (); AddMDRToPC ()
		case LDN:
			// Load negative - read memory (MDR)
			// LDN then ZeroACC (); SUB ()
		case STO:
			// Store ACC result to memory
			// STO then MemoryWrite ()
		case SUB, SUB_alt:
			// MemoryRead (); ACC = (ACC - MDR as word)
		case TEST:
			// Test - if ACC<0 PC++
		case STOP:
			Stopped = true
		default:
			// Nothing
		}

		EW <- IR
	}
}

func WriteBack(EW <-chan SSEMInst, MEMW chan<- SSEMInst) {

	for {
		MDR := <-EW
		// TODO Writre back into memory
		MEMW <- MDR
	}
}
