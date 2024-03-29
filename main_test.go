package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllSampleFiles(t *testing.T) {
	dirEntry, err := os.ReadDir("./samples")
	assert.Nil(t, err, "error should be nil")
	for _, entry := range dirEntry {
		if !entry.IsDir() {
			stringSplit := strings.Split(entry.Name(), ".")
			assert.Len(t, stringSplit, 2)
			fileName := stringSplit[0]

			inputFile, err := os.Open("./samples/" + fileName + ".txt")
			assert.Nil(t, err, "error should be nil")
			defer inputFile.Close()
			var b bytes.Buffer
			outputWriter := bufio.NewWriter(&b)
			processData(inputFile, outputWriter)
			assert.Nil(t, err, "error should be nil")
			// Ensure writer's contents are flushed
			err = outputWriter.Flush()
			assert.Nil(t, err, "error should be nil")
			outFile, err := os.Open("./samples/" + fileName + ".out")
			assert.Nil(t, err, "error should be nil")
			defer outFile.Close()
			equal, err := contentsEqual(outFile, b)
			assert.Nil(t, err, "should be nil")
			assert.Equal(t, true, equal, "should be true filename: %v", fileName)
		}
	}
}

func Test100MFile(t *testing.T) {
	file, err := os.Open("./measurements_10000000.txt")
	assert.Nil(t, err, "error should be nil")
	var b bytes.Buffer
	outputWriter := bufio.NewWriter(&b)
	processData(file, outputWriter)
}

func contentsEqual(file *os.File, buffer bytes.Buffer) (bool, error) {
	// Reset both files' read pointers to the start
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}

	// Create readers for both the file and the buffer
	fileReader := bufio.NewReader(file)
	// This assumes you have a way to get the bytes written by the writer,
	// like writing to a bytes.Buffer as its underlying writer.
	// This part needs to be adjusted based on your actual implementation.
	bufferReader := bufio.NewReader(&buffer)

	// Compare contents
	for {
		b1, err1 := fileReader.ReadByte()
		b2, err2 := bufferReader.ReadByte()
		// fmt.Printf("%v = %v\n", b1, b2)
		if err1 != err2 {
			return false, nil // Different errors or EOF states
		}

		if err1 == io.EOF {
			break // Both have reached EOF together
		}

		if err1 != nil {
			return false, err1 // Some other error occurred
		}

		if b1 != b2 {
			return false, nil // Bytes differ
		}
	}

	return true, nil // Contents are the same
}
