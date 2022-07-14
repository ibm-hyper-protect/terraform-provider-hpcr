package common

import (
	"fmt"
	"testing"
)

func TestTempFile(t *testing.T) {
	newFile := TempFile([]byte("Carsten"))
	fmt.Println(newFile)
}
