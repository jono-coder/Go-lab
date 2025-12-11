package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	file *os.File
)

func TestCsvParser_Parse(t *testing.T) {
	beforeEach()
	defer afterEach()

	f, err := os.Open(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			t.Errorf("error closing file: %v", err)
		}
	}()

	csvParser := NewCsvParser()

	numLines := 1
	err = csvParser.Parse(f, func(_ []string) {
		numLines++
	})
	if err != nil {
		t.Errorf("parsing error: %s", err)
	}

	req := require.New(t)
	req.Equal(100_001, numLines)
}

func beforeEach() {
	f, err := os.CreateTemp("", "*.csv")
	if err != nil {
		log.Fatalf("error creating tempfile: %s", err)
	}
	w := bufio.NewWriterSize(f, 4<<20)

	for i := 1; i <= 100_000; i++ {
		_, err := w.WriteString(fmt.Sprintf("%d,testing_%d\n", i, i))
		if err != nil {
			return
		}
	}

	if err = w.Flush(); err != nil {
		fmt.Println("error flushing writer:", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		fmt.Println("error seeking to beginning of file:", err)
	}

	if err := f.Close(); err != nil {
		fmt.Printf("error closing file: %s\n", err)
	}

	file = f
}

func afterEach() {
	if err := os.Remove(file.Name()); err != nil {
		fmt.Printf("error removing file: %s\n", err)
	}
}
