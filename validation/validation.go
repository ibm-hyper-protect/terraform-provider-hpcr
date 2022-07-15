//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package validation

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/terraform-provider-hpcr/encrypt"
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	S "github.com/terraform-provider-hpcr/fp/string"
)

func ToDiagnostics[A any](value E.Either[error, A]) diag.Diagnostics {
	return F.Pipe1(
		value,
		E.Fold(diag.FromErr, F.Constant1[A, diag.Diagnostics](nil)),
	)
}

func DiagCertificate(data any, _ cty.Path) diag.Diagnostics {
	// convert the key
	return F.Pipe4(
		data,
		fp.ToType[string],
		E.Map[error](S.ToBytes),
		E.Chain(encrypt.CertSerial),
		ToDiagnostics[[]byte],
	)
}
