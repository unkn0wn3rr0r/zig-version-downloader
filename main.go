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

	res, err := utils.MakeRequest(http.MethodGet, zigDownloadUrl, client)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	zigSuffix, err := utils.GetUrlVersionSuffix(res)
	if err != nil {
		log.Println(err)
		return
	}

	dlUrl := fmt.Sprintf("%s/%s", zigArchivePrefix, zigSuffix)
	res, err = utils.MakeRequest(http.MethodGet, dlUrl, client)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	pathToFile, err := utils.CreateFilepath(zigSuffix)
	if err != nil {
		log.Println(err)
		return
	}

	newArchiver, err := archiver.NewArchiver()
	if err != nil {
		log.Println(err)
		return
	}

	switch a := newArchiver.(type) {
	case *archiver.ZipArchiver:
		if err := a.CreateArchive(pathToFile, res); err != nil {
			log.Println(err)
			return
		}
	case *archiver.TarArchiver:
		if err := a.CreateArchive(pathToFile, res); err != nil {
			log.Println(err)
			return
		}
	}

	log.Printf("successfully downloaded archive at: %s", pathToFile)
	log.Println("do you want to unzip it? - [y]/[n]")
	shouldReturn, err := readUserInput()
	if err != nil {
		log.Println(err)
		return
	}
	if shouldReturn {
		return
	}

	written, err := newArchiver.Unzip(pathToFile)
	if err != nil {
		log.Println(err)
		return
	}
	extension, err := utils.GetOsFileExtension()
	if err != nil {
		log.Println(err)
		return
	}

	pathToFile = strings.TrimSuffix(pathToFile, extension)
	if written > 0 {
		log.Printf("succesfully downloaded and extracted total of %fmbs at: %s", float64(written)/1_048_576, pathToFile)
	} else {
		log.Printf("succesfully downloaded and extracted files at: %s", pathToFile)
	}
	log.Printf("time took: %f seconds", time.Since(startTime).Seconds())
}

func readUserInput() (shouldReturn bool, err error) {
	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read user input err: %s", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(answer), "y") {
		log.Printf("time took: %f seconds", time.Since(startTime).Seconds())
		return true, nil
	}
	return false, nil
}
