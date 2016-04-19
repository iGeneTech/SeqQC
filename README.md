# SeqQC
Fast calculating QC20 and QC30 for 2nd-sequencing data


## Compile

1. Install golang: https://golang.org/doc/install
2. `git clone https://github.com/iGeneTech/SeqQC.git`
3. `go build seqQC.go`


## Usage

```bash
Usage: seqQC -c 8 *.gz

  -c int
    	Number of CPU/Core to use. (default 1)

```



## Note

This code doesn't check the legality of the file format.