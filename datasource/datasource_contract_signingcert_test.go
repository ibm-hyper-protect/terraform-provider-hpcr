package datasource

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/data"
	"github.com/stretchr/testify/assert"
)

func Commoner() (string, string, string, string, string, error) {
	contractPath, err := os.Open("../samples/contracts/simple.yaml")
	if err != nil {
		fmt.Println("Error parsing contract - ", err)
	}
	defer contractPath.Close()

	contract, err := io.ReadAll(contractPath)
	if err != nil {
		fmt.Println(err)
		return "", "", "", "", "", err
	}

	encryptCert := data.DefaultCertificate

	privateKeyPath, err := os.Open("../samples/contract-expiry/private.pem")
	if err != nil {
		fmt.Println("Error parsing Private Key - ", err)
		return "", "", "", "", "", err
	}
	defer privateKeyPath.Close()

	privateKey, err := io.ReadAll(privateKeyPath)
	if err != nil {
		fmt.Println(err)
		return "", "", "", "", "", err
	}

	caCertPath, err := os.Open("../samples/contract-expiry/personal_ca.crt")
	if err != nil {
		fmt.Println("Error parsing CA certificate - ", err)
		return "", "", "", "", "", err
	}
	defer caCertPath.Close()

	caCert, err := io.ReadAll(caCertPath)
	if err != nil {
		fmt.Println(err)
		return "", "", "", "", "", err
	}

	caKeyPath, err := os.Open("../samples/contract-expiry/personal_ca.key")
	if err != nil {
		fmt.Println("Error parsing CA certificate - ", err)
		return "", "", "", "", "", err
	}
	defer caCertPath.Close()

	caKey, err := io.ReadAll(caKeyPath)
	if err != nil {
		fmt.Println(err)
		return "", "", "", "", "", err
	}

	return string(contract), encryptCert, string(privateKey), string(caCert), string(caKey), nil
}

func TestEncryptAndSign(t *testing.T) {

	contract, encryptCert, privateKey, caCert, caKey, err := Commoner()
	if err != nil {
		fmt.Println(err)
	}

	csrDataMap := map[string]interface{}{
		"country":  "IN",
		"state":    "Karnataka",
		"location": "Bangalore",
		"org":      "IBM",
		"unit":     "ISDL",
		"domain":   "HPVS",
		"mail":     "sashwat.k@ibm.com",
	}
	csrDataStr, err := json.Marshal(csrDataMap)
	if err != nil {
		fmt.Println(err)
	}

	expiryDays := 365

	_, err = EncryptAndSign(contract, encryptCert, privateKey, caCert, caKey, string(csrDataStr), "", expiryDays)

	assert.NoError(t, err)
}

func TestEncryptAndSignCsrFile(t *testing.T) {

	contract, encryptCert, privateKey, caCert, caKey, err := Commoner()
	if err != nil {
		fmt.Println(err)
	}

	csrPath, err := os.Open("../samples/contract-expiry/csr.pem")
	if err != nil {
		fmt.Println("Error parsing CSR - ", err)
	}
	defer csrPath.Close()

	csr, err := io.ReadAll(csrPath)
	if err != nil {
		fmt.Println(err)
	}

	expiryDays := 365

	_, err = EncryptAndSign(contract, encryptCert, privateKey, caCert, caKey, "", string(csr), expiryDays)

	assert.NoError(t, err)
}
