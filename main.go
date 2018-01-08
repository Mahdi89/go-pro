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
type CRTAddress uint8
type SSEMFunc uint8

var Memory [100]word

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

func Fetch(MEMR <-chan SSEMInst, FD chan<- string) {

	// Fetch data/inst. from file
	// TODO read from MEMR chan.
	file, err := os.Open("./data/gcd.raw")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		FD <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func Decode(FD <-chan string, DE chan<- SSEMInst) {

	var s, i string
	var h0, h1, h2, h3 byte

	for {

		IR := <-FD

		fmt.Sscanf(IR, "%4s%4s", &s, &i)
		fmt.Sscanf(i, "%1x%1x%1x%1x", &h0, &h1, &h2, &h3)

		fmt.Println(h0, h1, h2, h3)

		switch SSEMFunc(h0) {
		case JMP:
			DE <- SSEMInst{
				LineNo: 0,
				CRTNo:  0,
				Func:   JMP}
		case JRP:
			DE <- SSEMInst{
				LineNo: 0,
				CRTNo:  0,
				Func:   JRP}
		case LDN:
			DE <- SSEMInst{
				LineNo: 0,
				CRTNo:  0,
				Func:   LDN}
		case STO:
			DE <- SSEMInst{
				LineNo: 0,
				CRTNo:  0,
				Func:   STO}
		case SUB:
			DE <- SSEMInst{
				LineNo: 0,
				CRTNo:  0,
				Func:   SUB}
		case SUB_alt:
			DE <- SSEMInst{
				LineNo: 0,
				CRTNo:  0,
				Func:   SUB_alt}
		case TEST:
			DE <- SSEMInst{
				LineNo: 0,
				CRTNo:  0,
				Func:   TEST}
		case STOP:
			DE <- SSEMInst{
				LineNo: 0,
				CRTNo:  0,
				Func:   STOP}
		default:
			DE <- SSEMInst{
				LineNo: 0,
				CRTNo:  0,
				Func:   STOP}
		}

	}
}

func Execute(DE <-chan SSEMInst, EW chan<- SSEMInst) {

	for {
		IR := <-DE
		// Execute the decoded instruction
		MDR := IR
		EW <- MDR
	}
}

func WriteBack(EW <-chan SSEMInst, MEMW chan<- SSEMInst) {

	for {
		MDR := <-EW
		// TODO Writre back into memory
		MEMW <- MDR
	}
}

func main(
// Memory and Host side
// channels to be included here
) {
	// Inter stage channels
	MEMR := make(chan SSEMInst, 1)
	FD := make(chan string, 1)
	DE := make(chan SSEMInst, 1)
	EW := make(chan SSEMInst, 1)
	MEMW := make(chan SSEMInst, 1)

	for {
		// Processor stages
		go Fetch(MEMR, FD)
		go Decode(FD, DE)
		go Execute(DE, EW)
		go WriteBack(EW, MEMW)
	}

}
