package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/unkn0wn3rr0r/zig-version-downloader/archiver"
	"github.com/unkn0wn3rr0r/zig-version-downloader/utils"
)

const (
	zigDownloadUrl   = "https://ziglang.org/download/index.json"
	zigArchivePrefix = "https://ziglang.org/builds"
)

var startTime = time.Now()

func main() {
	client := &http.Client{}

	res := utils.MakeRequest(http.MethodGet, zigDownloadUrl, client)
	defer res.Body.Close()

	zigSuffix := utils.GetUrlVersionSuffix(res)

	dlUrl := fmt.Sprintf("%s/%s", zigArchivePrefix, zigSuffix)
	res = utils.MakeRequest(http.MethodGet, dlUrl, client)
	defer res.Body.Close()

	pathToFile := utils.CreateFilepath(zigSuffix)

	newArchiver := archiver.NewArchiver()
	switch a := newArchiver.(type) {
	case *archiver.ZipArchiver:
		a.CreateArchive(pathToFile, res)
	case *archiver.TarArchiver:
		a.CreateArchive(pathToFile, res)
	default:
		panic(fmt.Sprintf("No such archiver type %T!\n", a))
	}

	shouldReturn := readUserInput()
	if shouldReturn {
		return
	}

	written := newArchiver.Unzip(pathToFile)
	pathToFile = strings.TrimSuffix(pathToFile, utils.GetOsFileExtension())
	log.Printf("succesfully downloaded and extracted total of %fmbs at: %s", float64(written)/1_048_576, pathToFile)
	log.Printf("time took: %f seconds", time.Since(startTime).Seconds())
}

func readUserInput() bool {
	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read user input err: %s", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(answer), "y") {
		log.Printf("time took: %f seconds", time.Since(startTime).Seconds())
		return true
	}
	return false
}

// go build -o myapp -ldflags="-s -w" -tags netgo -installsuffix netgo --ldflags="-extldflags=-static" -ldflags "-linkmode external -extldflags -static" -v main.go
// GOOS=darwin GOARCH=amd64 go build -o app-amd64-darwin main.go
// GOOS=linux GOARCH=amd64 go build -o app-amd64-linux main.go
