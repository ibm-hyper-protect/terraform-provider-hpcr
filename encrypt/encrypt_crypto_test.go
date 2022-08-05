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
package encrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

func TestCryptRandomPassword(t *testing.T) {
	n := keylen
	pwd := CryptoRandomPassword(n)

	lenE := F.Pipe1(
		pwd(),
		E.Map[error](B.Len),
	)

	assert.Equal(t, E.Of[error](n), lenE)
}
