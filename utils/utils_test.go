package utils

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFilepath(t *testing.T) {
	dir, err := os.Getwd()
	printErr(t, err)

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
		printErr(t, err)

		expected := filepath.FromSlash(filepath.Join(dir, data.filename))
		if expected != actual {
			t.Errorf("Expected: %v, Actual: %v", expected, actual)
		}
	}
}

func TestGetOs(t *testing.T) {
	actual, err := getOs()
	printErr(t, err)

	if actual != windows && actual != linux && actual != macos {
		t.Errorf("Expected one of: %v, Actual: %v", []string{windows, linux, macos}, actual)
	}

}

func TestGetArch(t *testing.T) {
	actual, err := getArch()
	printErr(t, err)

	if actual != "x86_64" && actual != "aarch64" {
		t.Errorf("Expected one of: %v, Actual: %v", []string{"x86_64", "aarch64"}, actual)
	}

}

func TestGetOsFileExtension(t *testing.T) {
	actual, err := GetOsFileExtension()
	printErr(t, err)

	if actual != ".zip" && actual != ".tar.xz" {
		t.Errorf("Expected one of: %v, Actual: %v", []string{".zip", ".tar.xz"}, actual)
	}

}

func TestGetZigLatestVersion(t *testing.T) {
	expected := "someversion"
	var payload struct {
		Master struct {
			Version string `json:"version"`
		} `json:"master"`
	}
	payload.Master.Version = expected

	out, err := json.Marshal(payload)
	printErr(t, err)

	recorder := httptest.NewRecorder()
	recorder.Write(out)
	res := recorder.Result()

	actual, err := getZigLatestVersion(res)
	printErr(t, err)

	if actual != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, actual)
	}
}

func printErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
