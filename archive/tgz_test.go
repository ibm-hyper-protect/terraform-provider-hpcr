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
package archive

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"testing"

	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	"github.com/stretchr/testify/assert"
)

func TestTgz(t *testing.T) {
	var body bytes.Buffer

	base64 := base64.NewEncoder(base64.StdEncoding, &body)

	resE := F.Pipe3(
		base64,
		TarFolder[io.WriteCloser]("../samples/nginx-golang"),
		E.Chain(onClose[io.WriteCloser]),
		E.MapTo[error, any](true),
	)

	assert.Equal(t, E.Of[error](true), resE)

	fmt.Println(body.String())
}
