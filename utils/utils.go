package utils

import (
	"encoding/json"
	"fmt"
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

func MakeRequest(method, url string, client *http.Client) (res *http.Response, err error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed with: %s", url, err)
	}
	res, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response from %s failed with: %s", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response from %s failed with status: %s", url, res.Status)
	}
	return res, nil
}

func CreateFilepath(filename string) (path string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %s", err)
	}
	return filepath.FromSlash(filepath.Join(dir, filename)), nil
}

func GetUrlVersionSuffix(res *http.Response) (suffix string, err error) {
	version, err := getZigLatestVersion(res)
	if err != nil {
		return "", fmt.Errorf("getting latest version failed: %s", err)
	}
	extension, err := GetOsFileExtension()
	if err != nil {
		return "", fmt.Errorf("getting file extension failed: %s", err)
	}
	os, err := getOs()
	if err != nil {
		return "", fmt.Errorf("getting os failed: %s", err)
	}
	arch, err := getArch()
	if err != nil {
		return "", fmt.Errorf("getting architecture failed: %s", err)
	}
	return fmt.Sprintf("zig-%s-%s-%s%s", os, arch, version, extension), nil
}

func GetOsFileExtension() (extension string, err error) {
	os, err := getOs()
	if err != nil {
		return "", err
	}
	switch os {
	case windows:
		return ".zip", nil
	case linux:
		return ".tar.xz", nil
	case macos:
		return ".tar.xz", nil
	default:
		return "", fmt.Errorf("failed to get file extension for os: %s", os)
	}
}

func getArch() (ar string, err error) {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		return "x86_64", nil
	case "arm64":
		return "aarch64", nil
	default:
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}
}

func getOs() (opSystem string, err error) {
	goos := runtime.GOOS
	switch goos {
	case windows:
		return windows, nil
	case linux:
		return linux, nil
	case darwin:
		return macos, nil
	default:
		return "", fmt.Errorf("unsupported operation system: %s", goos)
	}
}

func getZigLatestVersion(res *http.Response) (version string, err error) {
	var payload struct {
		Master struct {
			Version string `json:"version"`
		} `json:"master"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("parsing response body failed with: %s", err)
	}
	return payload.Master.Version, nil
}
