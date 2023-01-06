package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/unkn0wn3rr0r/zig-version-downloader/utils"
)

const (
	zigDownloadUrl   = "https://ziglang.org/download/index.json"
	zigArchivePrefix = "https://ziglang.org/builds"
)

func main() {
	startTime := time.Now()

	client := &http.Client{}

	res := request(http.MethodGet, zigDownloadUrl, client)
	defer res.Body.Close()

	zigSuffix := utils.GetUrlVersionSuffix(res)

	dlUrl := fmt.Sprintf("%s/%s", zigArchivePrefix, zigSuffix)
	res = request(http.MethodGet, dlUrl, client)
	defer res.Body.Close()

	pathToFile := createFilepath(zigSuffix)
	destination, err := os.Create(pathToFile)
	if err != nil {
		log.Fatalf("failed to create destination dir from file %s error: %s", pathToFile, err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, res.Body)
	if err != nil {
		log.Fatalf("failed to write file into destination %s error: %s", pathToFile, err)
	}

	log.Printf("successfully downloaded archive at: %s", pathToFile)
	log.Println("do you want to unzip it? - [y]/[n]")
	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read user input err: %s", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(answer), "y") {
		log.Printf("time took: %f seconds", time.Since(startTime).Seconds())
		return
	}

	written := unzipFile(pathToFile)
	pathToFile = strings.TrimSuffix(pathToFile, utils.GetOsFileExtension())
	log.Printf("succesfully downloaded and extracted total of %fmbs at: %s", float64(written)/1_048_576, pathToFile)
	log.Printf("time took: %f seconds", time.Since(startTime).Seconds())
}

func unzipFile(pathToFile string) (written int64) {
	reader, err := zip.OpenReader(pathToFile)
	if err != nil {
		log.Fatalf("failed to get a reader for the zip file err: %s", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			log.Fatalf("failed to open the current file %s err: %s", file.Name, err)
		}

		targetDirectory := filepath.Join(".", file.Name)
		log.Printf("writing to current target directory: %s", targetDirectory)

		if file.FileInfo().IsDir() {
			log.Printf("creating a target directory: %s", targetDirectory)
			if err = os.MkdirAll(targetDirectory, file.Mode()); err != nil {
				log.Fatalf("creating target directory %s failed with err: %s", targetDirectory, err)
			}
		} else {
			openedFile, err := os.OpenFile(targetDirectory, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				log.Fatalf("failed to open the current file %s err: %s", targetDirectory, err)
			}
			w, err := io.Copy(openedFile, zippedFile)
			if err != nil {
				log.Fatalf("failed to copy contents from %s to  %s err: %s", file.Name, openedFile.Name(), err)
			}
			written += w
			openedFile.Close()
		}
		zippedFile.Close()
	}
	return
}

func createFilepath(filename string) string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %s", err)
	}
	path := filepath.Join(dir, filename)
	return filepath.FromSlash(path)
}

func request(method, url string, client *http.Client) *http.Response {
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

// go build -o myapp -ldflags="-s -w" -tags netgo -installsuffix netgo --ldflags="-extldflags=-static" -ldflags "-linkmode external -extldflags -static" -v main.go
// GOOS=darwin GOARCH=amd64 go build -o app-amd64-darwin main.go
// GOOS=linux GOARCH=amd64 go build -o app-amd64-linux main.go
