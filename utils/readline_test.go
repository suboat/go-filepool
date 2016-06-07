package utils

import (
	"testing"
)

var testPath = "readline_test.go"

func TestReadOneLine(t *testing.T) {
	rd, _ := NewReader(testPath)
	var err error
	var line string
	for ; err == nil; line, err = rd.ReadOneline() {
		println(line)
	}
	if err != rd.EOF {
		t.Log(err.Error())
	}
}
