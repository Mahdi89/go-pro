package main

import (
	"testing"
)

func TestProcessor(t *testing.T) {

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

	ret := <-MEMW
	// Calculate a program 2 + 2
	expected := SSEMInst{
		LineNo: 0,
		CRTNo:  0,
		Func:   STOP,
	}

	if expected != ret {
		t.Errorf("Expected price %v got %v", expected, ret)
	}

}
