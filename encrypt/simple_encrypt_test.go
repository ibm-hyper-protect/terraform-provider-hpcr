package encrypt

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleExecCommand(t *testing.T) {
	_, err := SimpleExecCommand("openssl", "", "version")

	assert.NoError(t, err)
}

func TestCreateTempFile(t *testing.T) {
	text := "Testing"
	tmpfile, err := CreateTempFile(text)

	file, err1 := os.Open(tmpfile)
	if err1 != nil {
		fmt.Println(err1)
	}
	defer file.Close()

	content, err1 := io.ReadAll(file)
	if err1 != nil {
		fmt.Println(err1)
	}

	err1 = os.Remove(tmpfile)
	if err1 != nil {
		fmt.Println(err1)
	}

	assert.Equal(t, text, string(content))
	assert.NoError(t, err)
}

func TestEncodeToBase64(t *testing.T) {
	base64data := "c2FzaHdhdGs="
	result := EncodeToBase64("sashwatk")

	assert.Equal(t, result, base64data)
}

func TestOpensslCheck(t *testing.T) {
	err := OpensslCheck()

	assert.NoError(t, err)
}

func TestRandomPasswordGenerator(t *testing.T) {
	_, _, err := RandomPasswordGenerator()

	assert.NoError(t, err)
}
