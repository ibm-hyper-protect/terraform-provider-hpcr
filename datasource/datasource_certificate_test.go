// Copyright 2023 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package datasource

package datasource

import (
	"testing"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	"github.com/stretchr/testify/assert"
)

func TestSelectCertificate(t *testing.T) {
	data := make(map[string]any)

	// prepare input data
	data[common.KeyCerts] = map[string]any{
		"1.0.10": "cert 1.0.10",
		"1.0.11": "cert 1.0.11",
	}
	data[common.KeySpec] = "^1.0.0"

	res := F.Pipe3(
		data,
		CreateResourceDataMock,
		handleCertificateWithContext(&defaultContext),
		E.ToError[fp.ResourceData],
	)

	assert.NoError(t, res)

	version, ok := data[common.KeyVersion].(string)
	assert.NotNil(t, version)
	assert.True(t, ok)

	assert.Equal(t, "1.0.11", version)
}
