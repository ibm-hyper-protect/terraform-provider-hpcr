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

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
	"github.com/stretchr/testify/assert"
)

func TestDownloadCertificates(t *testing.T) {
	t.Skip()
	data := make(map[string]any)

	// prepare input data
	data[common.KeyVersions] = []any{"1.0.11", "1.0.10"}
	data[common.KeyTemplate] = defaultTemplate

	res := F.Pipe3(
		data,
		CreateResourceDataMock,
		handleDownloadWithContext(&defaultContext),
		E.ToError[fp.ResourceData],
	)

	assert.NoError(t, res)

	certificates, ok := data[common.KeyCerts].(map[string]string)
	assert.NotNil(t, certificates)
	assert.True(t, ok)

	assert.Contains(t, certificates, "1.0.11")
	assert.Contains(t, certificates, "1.0.10")
}
