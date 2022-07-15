//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package File

import (
	"os"

	"github.com/terraform-provider-hpcr/common"
	E "github.com/terraform-provider-hpcr/fp/either"
)

var (
	onCreate = func() E.Either[error, *os.File] {
		return common.CreateTempE("", "*")
	}
	onDelete = func(f *os.File) E.Either[error, any] {
		f.Close() // #nosec
		return E.TryCatchError(func() (any, error) {
			return nil, os.Remove(f.Name())
		})
	}
)

func WithTempFile[A any]() func(func(*os.File) E.Either[error, A]) E.Either[error, A] {
	return E.WithResource[error, *os.File, A](onCreate, onDelete)
}
