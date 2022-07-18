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
	"fmt"
	"io/fs"
	"os"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/terraform-provider-hpcr/encrypt"
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	S "github.com/terraform-provider-hpcr/fp/string"
)

var (
	statE    = E.Eitherize1(os.Stat)
	Base64Re = regexp.MustCompile(`^((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)$`)
	TokenRe  = regexp.MustCompile(`^hyper-protect-basic\.((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)\.((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)$`)
)

func toDiagnostics[A any](value E.Either[error, A]) diag.Diagnostics {
	return F.Pipe1(
		value,
		E.Fold(diag.FromErr, F.Constant1[A, diag.Diagnostics](nil)),
	)
}

// validates that the given certificate is indeed a certificate
func DiagCertificate(data any, _ cty.Path) diag.Diagnostics {
	// convert the key
	return F.Pipe4(
		data,
		fp.ToTypeE[string],
		E.Map[error](S.ToBytes),
		E.Chain(encrypt.CertSerial),
		toDiagnostics[[]byte],
	)
}

// validates that the given path points to an existing folder
func DiagFolder(data any, _ cty.Path) diag.Diagnostics {
	return F.Pipe4(
		data,
		fp.ToTypeE[string],
		E.Chain(statE),
		E.Chain(E.FromPredicate(fs.FileInfo.IsDir, func(info fs.FileInfo) error {
			return fmt.Errorf("path %s is not a folder", info.Name())
		})),
		toDiagnostics[fs.FileInfo],
	)
}
