package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFilepath(t *testing.T) {
	dir, err := os.Getwd()
	printErr(err, t)

	table := []struct {
		filename string
	}{
		{"utils.go"},
		{"somepdf.pdf"},
		{"text.txt"},
		{"somefolder/text.txt"},
		{"somefolder/somenestedfolder/text.txt"},
		{"somefolder\\text.txt"},
		{"somefolder\\somenestedfolder\\text.txt"},
	}

	for _, data := range table {
		actual, err := CreateFilepath(data.filename)
		printErr(err, t)

		expected := filepath.FromSlash(filepath.Join(dir, data.filename))
		if expected != actual {
			t.Errorf("Expected: %v, Actual: %v", expected, actual)
		}
	}
}

func TestGetOs(t *testing.T) {
	actual, err := getOs()
	printErr(err, t)

	if windows != actual && linux != actual && macos != actual {
		t.Errorf("Expected one of: %v, Actual: %v", []string{windows, linux, macos}, actual)
	}

}

func printErr(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
