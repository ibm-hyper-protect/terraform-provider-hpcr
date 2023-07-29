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
	"io/fs"
	"testing"

	"os"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

var statE = E.Eitherize1(os.Stat)

func TestTempFile(t *testing.T) {
	resE := F.Pipe3(
		CreateTempE("", "*"),
		E.Map[error](func(f *os.File) string {
			return f.Name()
		}),
		E.Chain(statE),
		E.Map[error](fs.FileInfo.IsDir),
	)

	assert.Equal(t, resE, E.Of[error](false))
}
