package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

const VERSION = 1.0

// Flag and Args
var cpu = flag.Int("c", 1, "Number of CPU/Core to use.")

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type ResultData struct {
	FileName  string
	QC20      float64
	QC30      float64
	BaseCount int64
	ReadCount int64
}

func Worker(jobs <-chan string, results chan<- ResultData) {
	for inFile := range jobs {
		// log.Printf("Processing... %s", inFile)

		fh, err := os.Open(inFile)
		CheckErr(err)
		defer fh.Close()

		var r *bufio.Reader
		tr, err := gzip.NewReader(fh)
		CheckErr(err)

		r = bufio.NewReaderSize(tr, 100000)

		qc20 := float64(0)
		qc30 := float64(0)
		readCount := int64(0)
		baseCount := int64(0)

		for {
			line, err := r.ReadBytes('\n')

			if len(line) == 0 && err == nil {
				continue
			}

			if len(line) == 0 && err != nil {
				break
			}

			if line[0] != 43 {
				continue
			}

			line, err = r.ReadBytes('\n')

			if len(line) == 0 && err == nil {
				continue
			}

			if len(line) == 0 && err != nil {
				break
			}

			readCount++

			lineSize := len(line)

			for i := 0; i < lineSize; i++ {
				qc := int(line[i]) - 33
				if qc >= 20 {
					qc20++
				}

				if qc >= 30 {
					qc30++
				}
			}

			baseCount += int64(lineSize)

			if err != nil {
				break
			}

		}

		results <- ResultData{QC20: qc20, QC30: qc30, BaseCount: baseCount, ReadCount: readCount, FileName: inFile}
	}
}

func main() {
	startTime := time.Now()

	self, err := os.Stat(os.Args[0])
	CheckErr(err)

	flag.Usage = func() {
		fmt.Printf("SeqQC [v%.1f]: fast calculating QC20 and QC30 for 2nd-sequencing data\n\nUsage: %s -c 8 *.gz\n\n", VERSION, self.Name())
		flag.PrintDefaults()
		fmt.Println("\n")
		fmt.Println("Copyright 2015, iGeneTech Biotech Co., Ltd.")
	}

	flag.Parse()

	runtime.GOMAXPROCS(*cpu)

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	jobs := make(chan string, *cpu)
	results := make(chan ResultData)

	go func() {
		for _, inFile := range flag.Args() {
			jobs <- inFile
		}
		close(jobs)
	}()

	var wg sync.WaitGroup
	wg.Add(*cpu)

	for i := 0; i < *cpu; i++ {
		go func() {
			Worker(jobs, results)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("FileName\tReadsCount\tBaseCount\tQC20\tQC30")
	for r := range results {
		fmt.Printf("%s\t%d\t%d\t%.2f\t%.2f\n", r.FileName, r.ReadCount, r.BaseCount, 100*r.QC20/float64(r.BaseCount), 100*r.QC30/float64(r.BaseCount))
		// log.Printf("%s time used: %s.\n", r.FileName, time.Since(startTime).String())
	}

	log.Printf("Total time used: %s.\n", time.Since(startTime).String())
}
