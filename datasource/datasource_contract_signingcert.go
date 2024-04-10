package datasource

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/encrypt"
	"gopkg.in/yaml.v3"
)

func ResourceContractEncryptedSigningCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceContractEncryptedSigningCertCreate,
		Read:   resourceContractEncryptedSigningCertRead,
		Delete: resourceContractEncryptedSigningCertDelete,
		Schema: map[string]*schema.Schema{
			// input parameters
			common.KeyContract:   &schemaContractIn,
			common.KeyCert:       &schemaCertIn,
			common.KeyPrivKey:    &schemaPrivKeyIn,
			common.KeyExpiryDays: &schemaExpiryDaysIn,
			common.KeyCaCert:     &schemaCaCertIn,
			common.KeyCaKey:      &schemaCaKeyIn,
			common.KeyCsrParams:  &schemaCsrParamsIn,
			common.KeyCsrfile:    &schemaCsrFileIn,
			// output parameters
			common.KeyRendered: &schemaRenderedOut,
		},
		Description: "Generates an encrypted and signed user data field with contract expiry enabled",
	}
}

func resourceContractEncryptedSigningCertCreate(d *schema.ResourceData, meta interface{}) error {
	// OpenSSL check
	err := encrypt.OpensslCheck()
	if err != nil {
		return fmt.Errorf("OpenSSL not installed correctly %s", err.Error())
	}

	contract := d.Get(common.KeyContract).(string)
	encryptCertificate := d.Get(common.KeyCert).(string)
	privateKey := d.Get(common.KeyPrivKey).(string)
	expiryDays := d.Get(common.KeyExpiryDays).(int)
	caCert := d.Get(common.KeyCaCert).(string)
	caKey := d.Get(common.KeyCaKey).(string)
	csrParams := d.Get(common.KeyCsrParams).(map[string]interface{})
	csrData := d.Get(common.KeyCsrfile).(string)

	if privateKey == "" {
		return fmt.Errorf("private key missing")
	}

	if len(csrParams) != 0 && csrData != "" {
		return fmt.Errorf("CSR parameters and csr.pem has been parsed together")
	}

	csrJsonStr, err := json.Marshal(csrParams)
	if err != nil {
		return fmt.Errorf("error working on CSR data %s", err.Error())
	}
	finalContract, err := EncryptAndSign(contract, encryptCertificate, privateKey, caCert, caKey, string(csrJsonStr), csrData, expiryDays)
	if err != nil {
		return fmt.Errorf("error generating contract %s", err.Error())
	}
	err = d.Set(common.KeyRendered, finalContract)
	if err != nil {
		return fmt.Errorf("error saving contract %s", err.Error())
	}

	newUUID := uuid.New()
	d.SetId(newUUID.String())

	return resourceContractEncryptedSigningCertRead(d, meta)
}

func resourceContractEncryptedSigningCertRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceContractEncryptedSigningCertDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func EncryptAndSign(contract, encryptCert, privateKey, cacert, caKey, csrDataStr, csrPemData string, expiryDays int) (string, error) {
	var contractMap map[string]interface{}

	err := yaml.Unmarshal([]byte(contract), &contractMap)
	if err != nil {
		return "", err
	}

	randomPassword, err := encrypt.RandomPasswordGenerator()
	if err != nil {
		return "", err
	}

	encryptedRandomPassword, err := encrypt.EncryptPassword(randomPassword, encryptCert)
	if err != nil {
		return "", err
	}

	encryptedWorkload, err := encrypt.EncryptContract(randomPassword, contractMap["workload"].(map[string]interface{}))
	if err != nil {
		return "", err
	}

	finalWorkload := encrypt.EncryptFinalStr(encryptedRandomPassword, encryptedWorkload)

	signingCert, err := encrypt.CreateSigningCert(privateKey, cacert, caKey, csrDataStr, csrPemData, expiryDays)
	if err != nil {
		return "", err
	}

	signingKeyInjectedEnv, err := encrypt.KeyValueInjector(contractMap["env"].(map[string]interface{}), "signingKey", signingCert)
	if err != nil {
		return "", err
	}

	var envMap map[string]interface{}

	err = yaml.Unmarshal([]byte(signingKeyInjectedEnv), &envMap)
	if err != nil {
		return "", err
	}

	encryptedEnv, err := encrypt.EncryptContract(randomPassword, envMap)
	if err != nil {
		return "", err
	}

	finalEnv := encrypt.EncryptFinalStr(encryptedRandomPassword, encryptedEnv)

	workloadEnvSignature, err := encrypt.SignContract(finalWorkload, finalEnv, privateKey)
	if err != nil {
		return "", err
	}

	finalContract, err := encrypt.GenFinalSignedContract(finalWorkload, finalEnv, workloadEnvSignature)
	if err != nil {
		return "", err
	}

	return finalContract, nil
}
