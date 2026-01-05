package utils

import (
	"Go-lab/internal/utils/validate"
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

type CsvParser struct {
}

func NewCsvParser() *CsvParser {
	return &CsvParser{}
}

func (p *CsvParser) Parse(file *os.File, onParsed func(row []string)) error {
	err := validate.Required("file", file)
	if err != nil {
		return err
	}
	err = validate.Required("onParsed", onParsed)
	if err != nil {
		return err
	}

	r := bufio.NewReaderSize(file, 4<<20) // 4MB
	reader := csv.NewReader(r)
	reader.ReuseRecord = true // be kind to the garbage collector

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		onParsed(row)
	}

	return nil
}

// Rows uses an iterator (better imo)
func (p *CsvParser) Rows(r io.Reader) Iter[[]string] {
	return func(yield func([]string) bool) {
		cr := csv.NewReader(r)
		cr.ReuseRecord = true // be kind to the garbage collector

		for {
			row, err := cr.Read()
			if err == io.EOF {
				return
			}
			if err != nil {
				// For now, panic. Could also store error externally
				panic(err)
			}

			if !yield(row) {
				// early termination requested
				return
			}
		}
	}
}
