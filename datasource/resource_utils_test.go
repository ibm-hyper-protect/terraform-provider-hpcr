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
	"fmt"
	"testing"

	"github.com/terraform-provider-hpcr/common"
	D "github.com/terraform-provider-hpcr/data"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
)

func TestHashWithCert(t *testing.T) {

	data := make(map[string]any)

	// prepare input data
	data[common.KeyCert] = D.DefaultCertificate

	test := []byte("Carsten")

	hashE := F.Pipe3(
		data,
		CreateResourceDataMock,
		createHashWithCert,
		I.Ap[[]byte, E.Either[error, string]](test),
	)

	fmt.Println(hashE)
}
