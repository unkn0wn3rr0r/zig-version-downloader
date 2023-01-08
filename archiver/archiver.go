package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/unkn0wn3rr0r/zig-version-downloader/utils"
)

type archiver interface {
	Unzip(string) (written int64, err error)
}

type CreateArchive func(pathToFile string, res *http.Response) error

type ZipArchiver struct {
	CreateArchive
}
type TarArchiver struct {
	CreateArchive
}

func NewArchiver() (a archiver, err error) {
	extension, err := utils.GetOsFileExtension()
	if err != nil {
		return nil, err
	}
	if extension == ".zip" {
		return &ZipArchiver{CreateArchive: createArchive()}, nil
	}
	return &TarArchiver{CreateArchive: createArchive()}, nil
}

func (a *ZipArchiver) Unzip(pathToFile string) (written int64, err error) {
	reader, err := zip.OpenReader(pathToFile)
	if err != nil {
		return 0, fmt.Errorf("failed to get a reader for the zip file err: %s", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return 0, fmt.Errorf("failed to open the current file %s err: %s", file.Name, err)
		}

		targetDirectory := filepath.Join(".", file.Name)
		log.Printf("writing to current target directory: %s", targetDirectory)

		if file.FileInfo().IsDir() {
			log.Printf("creating a target directory: %s", targetDirectory)
			if err = os.MkdirAll(targetDirectory, file.Mode()); err != nil {
				return 0, fmt.Errorf("creating target directory %s failed with err: %s", targetDirectory, err)
			}
		} else {
			copyContents := func() (int64, error) {
				defer zippedFile.Close()
				openedFile, err := os.OpenFile(targetDirectory, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
				if err != nil {
					return 0, fmt.Errorf("failed to open the current file %s err: %s", targetDirectory, err)
				}
				defer openedFile.Close()
				w, err := io.Copy(openedFile, zippedFile)
				if err != nil {
					return 0, fmt.Errorf("failed to copy contents from file: %s to %s err: %s", file.Name, openedFile.Name(), err)
				}
				return w, nil
			}
			if w, err := copyContents(); err != nil {
				return 0, err
			} else {
				written += w
			}
		}
	}
	return
}

func (a *TarArchiver) Unzip(pathToFile string) (written int64, err error) {
	cmd := exec.Command("tar", "-J", "-xf", pathToFile)
	if err = cmd.Run(); err != nil {
		return 0, fmt.Errorf("error %s while executing command on file: %s", err, pathToFile)
	}
	return
}

// i might want to try a different and more efficient approach here
func createArchive() CreateArchive {
	createArch := func(pathToFile string, res *http.Response) error {
		log.Printf("downloading archive at: %s ...", pathToFile)
		archiveDestination, err := os.Create(pathToFile)
		if err != nil {
			return fmt.Errorf("failed to create destination dir from file: %s error: %s", pathToFile, err)
		}
		defer archiveDestination.Close()
		if _, err = io.Copy(archiveDestination, res.Body); err != nil {
			return fmt.Errorf("failed to write file into destination: %s error: %s", pathToFile, err)
		}
		return nil
	}
	return createArch
}
