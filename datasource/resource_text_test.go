//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package datasource

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-provider-hpcr/common"
	D "github.com/terraform-provider-hpcr/data"
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	"github.com/terraform-provider-hpcr/validation"
)

func TestUnencryptedText(t *testing.T) {
	data := make(map[string]any)

	// prepare input data
	data[common.KeyText] = "sample text"

	res := F.Pipe3(
		data,
		CreateResourceDataMock,
		resourceText,
		E.ToError[fp.ResourceData],
	)

	assert.NoError(t, res)
	assert.Regexp(t, validation.Base64Re, data[common.KeyRendered])
}

func TestEncryptedText(t *testing.T) {
	data := make(map[string]any)

	// prepare input data
	data[common.KeyText] = "sample text"
	data[common.KeyCert] = D.DefaultCertificate

	res := F.Pipe4(
		data,
		CreateResourceDataMock,
		resourceEncText,
		E.Chain(resourceEncText),
		E.ToError[fp.ResourceData],
	)

	assert.NoError(t, res)
	assert.Regexp(t, validation.TokenRe, data[common.KeyRendered])
}
