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
package yaml

import (
	E "github.com/IBM/fp-go/either"
	"gopkg.in/yaml.v3"
)

func Parse[A any](data []byte) E.Either[error, A] {
	var result A
	err := yaml.Unmarshal(data, &result)
	return E.TryCatchError(result, err)
}

func Stringify[A any](a A) E.Either[error, []byte] {
	return E.TryCatchError(yaml.Marshal(a))
}
