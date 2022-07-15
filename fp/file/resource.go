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

	E "github.com/terraform-provider-hpcr/fp/either"
)

var (
	onCreate = func() E.Either[error, *os.File] {
		return E.TryCatchError(func() (*os.File, error) {
			return os.CreateTemp("", "*")
		})
	}
	onDelete = func(f *os.File) E.Either[error, any] {
		path := f.Name()
		f.Close() // #nosec
		return E.TryCatchError(func() (string, error) {
			return path, os.Remove(path)
		})
	}
)

func WithTempFile[A any]() func(func(*os.File) E.Either[error, A]) E.Either[error, A] {
	return E.WithResource[error, *os.File, A](onCreate, onDelete)
}
