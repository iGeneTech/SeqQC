# SeqQC
Fast calculating QC20 and QC30 for 2nd-sequencing data


## Quick use

1. Download compiled version: [osx](https://github.com/iGeneTech/SeqQC/blob/master/osx-64/seqQC] or [linux https://github.com/iGeneTech/SeqQC/blob/master/linux-64/seqQC]
2. chmod +x seqQC
3. Run `./seqQC -c 2 *.gz`

## Usage

```bash
Usage: seqQC -c 8 *.gz

  -c int
    	Number of CPU/Core to use. (default 1)

```

## Compile

1. Install golang: https://golang.org/doc/install
2. `git clone https://github.com/iGeneTech/SeqQC.git`
3. `go build seqQC.go`

## Note

This code doesn't check the legality of the file format.