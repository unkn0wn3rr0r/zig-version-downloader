package archiver

import (
	"reflect"
	"testing"
)

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
}
