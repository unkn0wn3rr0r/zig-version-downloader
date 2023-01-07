package archiver

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/unkn0wn3rr0r/zig-version-downloader/utils"
)

type archiver interface {
	Unzip(string) int64
}

type CreateArchive func(pathToFile string, res *http.Response)

type ZipArchiver struct {
	CreateArchive
}
type TarArchiver struct {
	CreateArchive
}

func NewArchiver() archiver {
	if utils.GetOsFileExtension() == ".zip" {
		return &ZipArchiver{CreateArchive: createArchive()}
	}
	return &TarArchiver{CreateArchive: createArchive()}
}

func (a *ZipArchiver) Unzip(pathToFile string) (written int64) {
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

func (a *TarArchiver) Unzip(pathToFile string) (written int64) {
	cmd := exec.Command("tar", "-J", "-xf", pathToFile)
	if err := cmd.Run(); err != nil {
		log.Fatalf("error while executing command on file: %s", pathToFile)
	}
	return
}

func createArchive() CreateArchive {
	return func(pathToFile string, res *http.Response) {
		archiveDestination, err := os.Create(pathToFile)
		if err != nil {
			log.Fatalf("failed to create destination dir from file %s error: %s", pathToFile, err)
		}
		defer archiveDestination.Close()

		_, err = io.Copy(archiveDestination, res.Body)
		if err != nil {
			log.Fatalf("failed to write file into destination %s error: %s", pathToFile, err)
		}
	}
}
