// +build ignore
package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func main() {
	in, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	out, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(in)
	r.Comma = '\t'
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	w := csv.NewWriter(out)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record[7] != "ADMD" {
			continue
		}
		if err := w.Write(record); err != nil {
			log.Fatal(err)
		}
	}

	w.Flush()
	in.Close()
	out.Close()
}
