package archiver

import (
	"reflect"
	"testing"

	"github.com/unkn0wn3rr0r/zig-version-downloader/utils"
)

func TestNewArchiver(t *testing.T) {
	extension, err := utils.GetOsFileExtension()
	if err != nil {
		t.Errorf("Expected: %v, Actual: %v", nil, err)
	}

	if extension != ".zip" && extension != ".tar.xz" {
		t.Errorf("Expected one of: %v, Actual: %v", []string{".zip", ".tar.xz"}, extension)
	}

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
}
