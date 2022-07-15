//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package common

import (
	"encoding/base64"
	"fmt"
	"os"
	"runtime"

	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

func Base64Encode(buffer []byte) string {
	return base64.StdEncoding.EncodeToString(buffer)
}

func Base64Decode(data string) E.Either[error, []byte] {
	return E.TryCatchError(func() ([]byte, error) {
		return base64.StdEncoding.DecodeString(data)
	})
}

func removeFile(file *os.File) error {
	fmt.Printf("Removing temp file [%s]...", file.Name())
	return os.Remove(file.Name())
}

func TempFile(data []byte) E.Either[error, *os.File] {
	result := F.Pipe1(
		E.TryCatchError(func() (*os.File, error) {
			return os.CreateTemp("", "*")
		}),
		E.Map[error](func(file *os.File) *os.File {
			runtime.SetFinalizer(file, removeFile)
			return file
		}),
	)
	return result
}
