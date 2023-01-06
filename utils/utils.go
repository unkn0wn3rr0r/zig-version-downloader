package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	windows = "windows"
	linux   = "linux"
	macos   = "macos"
	darwin  = "darwin"
)

func MakeRequest(method, url string, client *http.Client) *http.Response {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf("request to %s failed with: %s", url, err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("response from %s failed with: %s", url, err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("response from %s failed with: %s", url, res.Status)
	}
	return res
}

func CreateFilepath(filename string) string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %s", err)
	}
	path := filepath.Join(dir, filename)
	return filepath.FromSlash(path)
}

func GetUrlVersionSuffix(res *http.Response) string {
	return fmt.Sprintf("zig-%s-%s-%s%s", getOs(), getArch(), getZigLatestVersion(res), GetOsFileExtension())
}

func GetOsFileExtension() string {
	os := getOs()
	switch os {
	case windows:
		return ".zip"
	case linux:
		return ".tar.xz"
	case macos:
		return ".tar.xz"
	default:
		panic("failed to get file extension for os: " + os)
	}
}

func getArch() string {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	default:
		panic("unsupported architecture: " + arch)
	}
}

func getOs() string {
	goos := runtime.GOOS
	switch goos {
	case windows:
		return windows
	case linux:
		return linux
	case darwin:
		return macos
	default:
		panic("unsupported operation system: " + goos)
	}
}

func getZigLatestVersion(res *http.Response) string {
	var payload struct {
		Master struct {
			Version string `json:"version"`
		} `json:"master"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		log.Fatalf("parsing response body failed with: %s", err)
	}
	return payload.Master.Version
}
