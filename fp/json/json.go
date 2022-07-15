//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package Json

import (
	"encoding/json"

	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

func Parse[A any](data []byte) E.Either[error, *A] {
	return E.TryCatch(func() (*A, error) {
		var result A
		err := json.Unmarshal(data, &result)
		return &result, err
	}, F.Identity[error])
}

func Stringify[A any](a *A) E.Either[error, []byte] {
	return E.TryCatch(func() ([]byte, error) {
		b, err := json.Marshal(a)
		return b, err
	}, F.Identity[error])

}
