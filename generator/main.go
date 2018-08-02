package main

import (
	"bytes"
	"flag"
	"go/format"
	"io/ioutil"
	"log"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

var program = `package main
        import (
                // Import the entire framework
                _ "github.com/ReconfigureIO/sdaccel"
                aximemory "github.com/ReconfigureIO/sdaccel/axi/memory"
                axiprotocol "github.com/ReconfigureIO/sdaccel/axi/protocol"
                arbitrate "github.com/ReconfigureIO/sdaccel/axi/arbitrate"
        )

        func Top(
        		inputData uintptr,
        		outputData uintptr,
        		length uint32,
                // The second set of arguments will be the ports for interacting with memory
                memReadAddr chan<- axiprotocol.Addr,
                memReadData <-chan axiprotocol.ReadData,

                memWriteAddr chan<- axiprotocol.Addr,
                memWriteData chan<- axiprotocol.WriteData,
                memWriteResp <-chan axiprotocol.WriteResp) {

        // Read all of the input data into a channel
        inputChan := make(chan uint32, {{ .Processor.Replicate }})

        go aximemory.ReadBurstUInt32(
                memReadAddr, memReadData, true, inputData, length * ({{ .Processor.TypeWidth }} / 32), inputChan)


        // Read all of the input data into a channel
        elementChan := make(chan {{ .Processor.Type }}, 1)
        go {{ .Processor.Deserialize }}(inputChan, elementChan)

        // Read all of the input data into a channel
        // dataChan := make(chan [{{ .Processor.Replicate }}]{{ .Processor.Type }}, 1)

        {{ range $index, $spec := .Processors }}
        data{{ $spec.Index }} := make(chan {{ $.Processor.Type }}, 1)
        {{ end }}

        {{ range $index, $spec := .Processors }}
        interConnect{{ $spec.Index }} := make(chan {{ $.Processor.Type }}, 1)
        {{ end }}


                // Processor part

                {{ range $index, $spec := .Processors }}
            	go func() {
                    for {
        	    	{{ $.Processor.Function }}(data{{ $spec.Index }})
                    }
            	}()
                {{ end }}

                // Network part

                {{ range $index, $spec := .Processors }}
                go func() {
                    for {
                        {{ $.Router.Function }}(data{{ $spec.Index }}, interConnect{{ $spec.Index }}, interConnect{{ $spec.Index }})
                    }
                }()
                {{ end }}



        outputDataChan := make(chan uint32)
        go {{ .Router.Serialize }}(interConnect0, outputDataChan)


        // Write it back to the pointer the host requests
        aximemory.WriteBurstUInt32(
                memWriteAddr, memWriteData, memWriteResp, true, outputData, {{ .Processor.TypeWidth }} / 32, outputDataChan)
        }

        func Benchmark(n uint) {
        	ra := make(chan axiprotocol.Addr)
        	rd := make(chan axiprotocol.ReadData)

        	wa := make(chan axiprotocol.Addr)
        	wd := make(chan axiprotocol.WriteData)
        	wr := make(chan axiprotocol.WriteResp)

        	go func() {
        		for {
        			a := <-ra
        			for i := a.Len; i != 255; i-- {
        				rd <- axiprotocol.ReadData{Last: i == 0}
        			}
        		}
        	}()

        	go func() {
        		for {
        			<-wa
        			running := true
        			for running {
        				d := <-wd
        				running = !d.Last
        			}
        			wr <- axiprotocol.WriteResp{}
        		}
        	}()

        Top(0, 0, uint32(n), ra, rd, wa, wd, wr)
        }

`

func main() {
	var filename = flag.String("output", "network.go", "output file name")
	configFile, err := ioutil.ReadFile("reco.yml")
	if err != nil {
		log.Fatal("Error opening config file", err)
	}

	var funcs = template.FuncMap{}

	var buffer bytes.Buffer

	var d Data

	err = yaml.Unmarshal(configFile, &d)
	if err != nil {
		log.Fatal("Error reading config file", err)
	}

	// Generate main()
	t := template.Must(template.New("main").Funcs(funcs).Parse(program))
	if err := t.Execute(&buffer, d); err != nil {
		log.Fatal("template", err)
	}

	data, err := format.Source(buffer.Bytes())
	if err != nil {
		log.Fatal("format ", err, string(buffer.Bytes()))
	}

	if err := ioutil.WriteFile(*filename, data, 0644); err != nil {
		panic(err)
	}
}
