package datasource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

func ResourceContractEncryptedSigningCert() *schema.Resource {
	return &schema.Resource{
		Create: nil,
		Read:   nil,
		Delete: nil,
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
