// Copyright 2022 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package common

import (
	"encoding/base64"
	"fmt"
	"os"

	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	O "github.com/terraform-provider-hpcr/fp/option"
	S "github.com/terraform-provider-hpcr/fp/string"
)

var (
	Base64Encode  = base64.StdEncoding.EncodeToString
	Base64DecodeE = E.Eitherize1(base64.StdEncoding.DecodeString)

	Base64EncodeE = F.Flow2(
		Base64Encode,
		E.Of[error, string],
	)

	PlainTextEncodeE = F.Flow2(
		B.ToString,
		E.Of[error, string],
	)

	CreateTempE = E.Eitherize2(os.CreateTemp)

	MapBytesToStgE = E.Map[error](B.ToString)
	MapStgToBytesE = E.Map[error](S.ToBytes)
	MapRefAnyE     = E.Map[error](F.Ref[any])
)

func ToTypeO[A any](data any) O.Option[A] {
	value, ok := data.(A)
	if ok {
		return O.Some(value)
	}
	return O.None[A]()
}

func ToTypeE[A any](data any) E.Either[error, A] {
	return F.Pipe2(
		data,
		ToTypeO[A],
		E.FromOption[error, A](func() error {
			return fmt.Errorf("invalid type of input [%T]", data)
		}),
	)
}
