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
package file

import (
	"os"

	E "github.com/IBM/fp-go/either"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

var (
	onCreate = func() E.Either[error, *os.File] {
		return common.CreateTempE("", "*")
	}
	onDelete = func(f *os.File) E.Either[error, any] {
		f.Close() // #nosec
		return E.TryCatchError[any](nil, os.Remove(f.Name()))
	}
)

func WithTempFile[A any]() func(func(*os.File) E.Either[error, A]) E.Either[error, A] {
	return E.WithResource[error, *os.File, A](onCreate, onDelete)
}
