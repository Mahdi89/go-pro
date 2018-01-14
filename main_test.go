package main

import (
	//	"fmt"
	"testing"
)

func TestProcessor(t *testing.T) {

	// Inter stage channels
	FD := make(chan word, 1)
	DE := make(chan SSEMInst, 1)
	EW := make(chan SSEMInst, 1)
	MEMW := make(chan SSEMInst, 1)

	// Load the memory
	go Mem("./data/gcd.raw")

	// Processor stages
	go Fetch(FD)
	go Decode(FD, DE)
	go Execute(DE, EW)
	go WriteBack(EW, MEMW)

	ret := <-MEMW
	// Calculate a program 2 + 2
	expected := SSEMInst{
		LineNo: 0,
		CRTNo:  0,
		Func:   JMP,
	}
	if expected != ret {
		t.Errorf("Expected price %v got %v", expected, ret)
	}

}
