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

	index := 0
	for {
		FD <- Memory[index]
		index += 1
	}
}

func Decode(FD <-chan word, DE chan<- SSEMInst) {

	var LA LineAddress
	var CA CRTAddress
	var Fun SSEMFunc
	for {
		IR := <-FD
		LA = LineAddress(IR & 0x001f)
		CA = CRTAddress((IR & 0x0fff) >> 4)
		Fun = SSEMFunc(IR >> 24)

		switch SSEMFunc(Fun) {
		case JMP:
			DE <- SSEMInst{
				LineNo: LA,
				CRTNo:  CA,
				Func:   JMP}
		case JRP:
			DE <- SSEMInst{
				LineNo: LA,
				CRTNo:  CA,
				Func:   JRP}
		case LDN:
			DE <- SSEMInst{
				LineNo: LA,
				CRTNo:  CA,
				Func:   LDN}
		case STO:
			DE <- SSEMInst{
				LineNo: LA,
				CRTNo:  CA,
				Func:   STO}
		case SUB:
			DE <- SSEMInst{
				LineNo: LA,
				CRTNo:  CA,
				Func:   SUB}
		case SUB_alt:
			DE <- SSEMInst{
				LineNo: LA,
				CRTNo:  CA,
				Func:   SUB_alt}
		case TEST:
			DE <- SSEMInst{
				LineNo: LA,
				CRTNo:  CA,
				Func:   TEST}
		case STOP:
			DE <- SSEMInst{
				LineNo: LA,
				CRTNo:  CA,
				Func:   STOP}
		default:
			DE <- SSEMInst{
				LineNo: LA,
				CRTNo:  CA,
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
