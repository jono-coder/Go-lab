package utils

import (
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
	if file == nil {
		return NPE("file")
	}
	if onParsed == nil {
		return NPE("func onParsed")
	}

	r := bufio.NewReaderSize(file, 1<<22) // 4MB
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
