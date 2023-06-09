package archiver

import (
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/unkn0wn3rr0r/zig-version-downloader/utils"
)

// that's the name of the file which resides in the 'test-file.zip' file
const tempTestFileName = "test-file.txt"

func TestNewArchiver(t *testing.T) {
	actual, err := NewArchiver()
	if err != nil {
		t.Errorf("Expected: %v, Actual: %v", nil, err)
	}

	concreteTypes := []string{"ZipArchiver", "TarArchiver"}
	if actual == nil {
		t.Errorf("Expected one of: %v, Actual: %v", concreteTypes, actual)
	}

	concreteType := reflect.TypeOf(actual).Elem()
	actualName := concreteType.Name()
	if actualName != "ZipArchiver" && actualName != "TarArchiver" {
		t.Errorf("Expected one of: %v, Actual: %v", concreteTypes, actualName)
	}

	supportedExtensions := []string{".zip", ".tar.xz"}
	extension, err := utils.GetOsFileExtension()
	if err != nil {
		t.Errorf("Expected one of: %v, Actual: %v", supportedExtensions, err)
	}

	var pathToZipFile string
	if extension == ".zip" {
		pathToZipFile, err = utils.CreateFilepath("test-files\\test-file.zip")
	} else {
		log.Println("no test setup found for 'tar.xz' files")
		return
	}

	if err != nil {
		t.Errorf("Expected: %v, Actual: %v", nil, err)
	}

	_, err = actual.Unzip(pathToZipFile)
	if err != nil {
		t.Errorf("Expected: %v, Actual: %v", nil, err)
	}

	pathToTempFile, err := utils.CreateFilepath(tempTestFileName)
	if err != nil {
		t.Errorf("Expected: %v, Actual: %v", nil, err)
	}

	time.Sleep(time.Second * 2) // that's just to see that the file is unzipped and then deleted by the cleanup function

	tempFileName, err := cleanupTempTestFile(pathToTempFile)
	if tempFileName != tempTestFileName || err != nil {
		t.Errorf("Expected: %v, Actual: %v", tempTestFileName, err)
	}
}

func cleanupTempTestFile(filepath string) (string, error) {
	err := os.Remove(filepath)
	if err != nil {
		return "", err
	}
	return tempTestFileName, nil
}
