// Description of a 4-stage in-order processor based on the
// Manchester Small-Scale Experimental Machine (SSEM) ISA
// The SSEM has a 32-bit word length and a memory of 32 words.

package main

import "fmt"

type word uint32
type LineAddress uint8
type CRTAddress uint8
type SSEMFunc uint8

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

func Fetch(MEMR <-chan SSEMInst, FD chan<- SSEMInst) {

	//IR := <-MEMR
	// TODO Fetch instruction from memory
	// Increment PC
	FD <- SSEMInst{
		LineNo: 0,
		CRTNo:  0,
		Func:   STOP}
}

func Decode(FD <-chan SSEMInst, DE chan<- SSEMInst) {

	IR := <-FD
	DE <- IR
}

func Execute(DE <-chan SSEMInst, EW chan<- SSEMInst) {

	IR := <-DE
	// Execute the decoded instruction
	MDR := IR
	EW <- MDR
}

func WriteBack(EW <-chan SSEMInst, MEMW chan<- SSEMInst) {

	MDR := <-EW
	// TODO Writre back into memory
	MEMW <- MDR
}

func main(
// Memory and Host side
// channels to be included here
) {
	// Inter stage channels
	MEMR := make(chan SSEMInst, 1)
	FD := make(chan SSEMInst, 1)
	DE := make(chan SSEMInst, 1)
	EW := make(chan SSEMInst, 1)
	MEMW := make(chan SSEMInst, 1)

	// Processor stages
	go Fetch(MEMR, FD)
	go Decode(FD, DE)
	go Execute(DE, EW)
	go WriteBack(EW, MEMW)
}
