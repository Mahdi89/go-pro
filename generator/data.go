package main

type Processor struct {
	Type        string
	TypeWidth   int `yaml:"typeWidth"`
	Deserialize string
	Function    string
	Replicate   int
}

type Router struct {
	Type      string
	TypeWidth int `yaml:"typeWidth"`
	Serialize string
	Function  string
}

type Data struct {
	Processor Processor
	Router    Router
}

type ProcessorSpec struct {
	Index        int
	ContextIndex int
}

func (d Data) Processors() []ProcessorSpec {
	ret := make([]ProcessorSpec, d.Processor.Replicate, d.Processor.Replicate)
	for i := range ret {
		ret[i] = ProcessorSpec{Index: i, ContextIndex: i >> 4}
	}
	return ret
}

type RouterSpec struct {
	Last        bool
	OutputIndex int
	InputA      int
	InputB      int
}
