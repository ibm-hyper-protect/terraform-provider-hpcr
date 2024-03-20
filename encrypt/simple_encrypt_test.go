package encrypt

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/data"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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

func TestEncryptPasswordSuccess(t *testing.T) {
	passowrd, _, err := RandomPasswordGenerator()
	if err != nil {
		fmt.Println(err)
	}

	result, err := EncryptPassword(passowrd, data.DefaultCertificate)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("result - ", result)

	assert.NoError(t, err)
}

func TestMapToYaml(t *testing.T) {
	var contractMap map[string]interface{}

	file, err := os.Open("../samples/contracts/simple.yaml")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	contract, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(contract), &contractMap)
	if err != nil {
		fmt.Println(err)
	}

	_, err = MapToYaml(contractMap["env"].(map[string]interface{}))
	if err != nil {
		fmt.Println(err)
	}

	assert.NoError(t, err)
}

func TestEncryptContract(t *testing.T) {
	var contractMap map[string]interface{}

	file, err := os.Open("../samples/contracts/simple.yaml")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	contract, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(contract), &contractMap)
	if err != nil {
		fmt.Println(err)
	}

	passowrd, _, err := RandomPasswordGenerator()
	if err != nil {
		fmt.Println(err)
	}

	result, err := EncryptContract(passowrd, contractMap["workload"].(map[string]interface{}))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Result -", result)

	assert.NoError(t, err)
}

func TestEncryptFinalStr(t *testing.T) {
	var contractMap map[string]interface{}

	file, err := os.Open("../samples/contracts/simple.yaml")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	contract, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(contract), &contractMap)
	if err != nil {
		fmt.Println(err)
	}

	passowrd, encodedPassword, err := RandomPasswordGenerator()
	if err != nil {
		fmt.Println(err)
	}

	encryptedWorkload, err := EncryptContract(passowrd, contractMap["workload"].(map[string]interface{}))
	if err != nil {
		fmt.Println(err)
	}

	finalWorkload := EncryptFinalStr(encodedPassword, encryptedWorkload)

	fmt.Println("workload: ", finalWorkload)

	assert.NoError(t, err)
}
