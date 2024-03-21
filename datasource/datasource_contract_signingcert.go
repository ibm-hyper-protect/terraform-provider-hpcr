package datasource

import (
	"fmt"

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
			common.KeyCsrParams:  &schemaCsrParams,
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
	// contract := d.Get(common.KeyContract).(string)
	// encryptCertificate := d.Get(common.KeyCert).(string)
	// privateKey := d.Get(common.KeyPrivKey).(string)
	// expiryDays := d.Get(common.KeyExpiryDays).(int)
	// caCert := d.Get(common.KeyCaCert).(string)
	// caKey := d.Get(common.KeyCaKey).(string)
	// csrParams := d.Get(common.KeyCsrParams).(map[string]interface{})

	// fmt.Println("contract - ", contract)
	// fmt.Println("Encrypt Certificate - ", encryptCertificate)
	// fmt.Println("Private Key - ", privateKey)
	// fmt.Println("expiryDays - ", expiryDays)
	// fmt.Println("CA Cert - ", caCert)
	// fmt.Println("CA Key - ", caKey)

	// for key, val := range csrParams {
	// 	fmt.Println(key, val)
	// }
	// encryptWorkload(contractMap)

	return resourceContractEncryptedSigningCertRead(d, meta)
}

func resourceContractEncryptedSigningCertRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceContractEncryptedSigningCertDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func EncryptAndSign(contract, encryptCert string) (string, error) {
	var contractMap map[string]interface{}

	err := yaml.Unmarshal([]byte(contract), &contractMap)
	if err != nil {
		return "", err
	}

	randomPassword, encodedRandomPassword, err := encrypt.RandomPasswordGenerator()
	if err != nil {
		return "", err
	}

	encryptedRandomPassword, err := encrypt.EncryptPassword(randomPassword, encryptCert)
	if err != nil {
		return "", err
	}

}
