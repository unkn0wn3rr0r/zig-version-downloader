package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/unkn0wn3rr0r/zig-version-downloader/utils"
)

func TestCreateFilePath(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

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
		actual, err := utils.CreateFilepath(data.filename)
		if err != nil {
			t.Error(err)
		}
		expected := filepath.FromSlash(filepath.Join(dir, data.filename))
		if expected != actual {
			t.Errorf("Expected: %v, Actual: %v", expected, actual)
		}
	}
}
