package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
	expected := []string{windows, linux, macos}

	actual, err := getOs()
	if err != nil {
		t.Errorf("Expected one of: %v, Actual: %v", expected, err)
	}

	if actual != windows && actual != linux && actual != macos {
		t.Errorf("Expected one of: %v, Actual: %v", expected, actual)
	}
}

func TestGetArch(t *testing.T) {
	expected := []string{"x86_64", "aarch64"}

	actual, err := getArch()
	if err != nil {
		t.Errorf("Expected one of: %v, Actual: %v", expected, err)
	}

	if actual != "x86_64" && actual != "aarch64" {
		t.Errorf("Expected one of: %v, Actual: %v", expected, actual)
	}
}

func TestGetOsFileExtension(t *testing.T) {
	expected := []string{".zip", ".tar.xz"}

	actual, err := GetOsFileExtension()
	if err != nil {
		t.Errorf("Expected one of: %v, Actual: %v", expected, err)
	}

	if actual != ".zip" && actual != ".tar.xz" {
		t.Errorf("Expected one of: %v, Actual: %v", expected, actual)
	}
}

func TestGetZigLatestVersion(t *testing.T) {
	expected := "some-version"
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
	if err != nil {
		t.Errorf("Expected: %v, Actual: %v", expected, err)
	}

	if actual != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, actual)
	}

	recorder.Write(out)

	actual, err = getZigLatestVersion(res)
	if err == nil {
		t.Errorf("Expected: %v, Actual: %v", err, actual)
	}
}

func TestGetUrlVersionSuffix(t *testing.T) {
	var payload struct {
		Master struct {
			Version string `json:"version"`
		} `json:"master"`
	}
	payload.Master.Version = "some-version"

	out, err := json.Marshal(payload)
	printErr(t, err)

	recorder := httptest.NewRecorder()
	recorder.Write(out)
	res := recorder.Result()

	expected := []string{".zip", ".tar.xz"}

	actual, err := GetUrlVersionSuffix(res)
	if err != nil {
		t.Errorf("Expected one of: %v, Actual: %v", expected, err)
	}

	if !strings.HasSuffix(actual, "some-version.zip") && !strings.HasSuffix(actual, "some-version.tar.xz") {
		t.Errorf("Expected one of: %v, Actual: %v", expected, actual)
	}
}

func TestMakeRequest(t *testing.T) {
	expected := "some-response"

	roundTripFn := RoundTripFn(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(expected)),
		}
	})
	client := &http.Client{
		Transport: roundTripFn,
	}

	res, err := MakeRequest(http.MethodGet, "some-url.com", client)
	if err != nil {
		t.Errorf("Expected Status: %v, Actual Status: %v", http.StatusOK, res.StatusCode)
	}

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected: %v, Actual: %v", expected, err)
	}

	actual := string(bs)
	if actual != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, actual)
	}
}

type RoundTripFn func(req *http.Request) *http.Response

func (fn RoundTripFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req), nil
}

func printErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
