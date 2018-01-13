# go-pro

https://medium.com/the-recon/beyond-opencl-more-concurrent-than-parallel-78734d698462

Use Go for fast hardware prototyping,
this includes parameterized instantiations of memory blocks, processing units and the communication layout at SoC level.
We use Reconfigure.io's tooling to synthesize the system prototypes into FPGA based logic.

Since Go is a high level language explicit hierarchy, IO protocol, clocks or resets won't appear in the code. Instead these will be introduced by the HLS tool targetting hardware in a wiser way such that hardware level details are not reflected to the designer. Moreover by abstracting these low-level details we are guaranteeing faster modeling and simulation times.

## SSEM Model

Currently a high level model of a 4-stage in-order processor based on the Manchester Small-Scale Experimental Machine (SSEM) ISA is included. The SSEM has a 32-bit word length and a memory of 32 words. 



