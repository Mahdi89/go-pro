package main

func Deserialize(inputChan <-chan uint32, outputChan chan<- uint32) {
	for {
		outputChan <- <-inputChan
	}
}

func Serialize(inputChan <-chan uint32, outputChan chan<- uint32) {
	for {
		outputChan <- uint32(0)
		outputChan <- uint32(0)
	}
}

func Processor(outputChan chan<- uint32) {

	// Executes the instructions in its associated memory
	// and outputs the results, forwards them to some target.
	for {

		outputChan <- uint32(0)
	}
}

func Router(inputChan <-chan uint32, outputChan0 chan<- uint32, outputChan1 chan<- uint32) {

	// Check the routing info and forwards the packet
	// to either one of its output ports
	for {

		outputChan0 <- <-inputChan
	}
}
