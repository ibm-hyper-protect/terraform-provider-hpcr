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
	_ "embed"
	"testing"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/attestation"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
	"github.com/stretchr/testify/assert"
)

//go:embed samples/attestation/unencrypted.txt
var UnencryptedAttestation string

//go:embed samples/attestation/se-checksums.txt.enc
var EncryptedAttestation string

//go:embed samples/attestation/attestation
var EncPrivKey string

func TestUnencryptedAttestation(t *testing.T) {
	data := make(map[string]any)

	// prepare input data
	data[common.KeyAttestation] = UnencryptedAttestation

	res := F.Pipe3(
		data,
		CreateResourceDataMock,
		handleAttestationWithContext(&defaultContext),
		E.ToError[fp.ResourceData],
	)

	assert.NoError(t, res)

	checksums, ok := data[common.KeyChecksums].(attestation.ChecksumMap)
	assert.NotNil(t, checksums)
	assert.True(t, ok)

	assert.Equal(t, "a6f6228bbf820e766ebe43c51e97332dda92e9744e719a646f611fe0681d2458", checksums["cidata/user-data"])
}

func TestEncryptedAttestation(t *testing.T) {
	data := make(map[string]any)

	// prepare input data
	data[common.KeyAttestation] = EncryptedAttestation
	data[common.KeyPrivKey] = EncPrivKey

	res := F.Pipe3(
		data,
		CreateResourceDataMock,
		handleAttestationWithContext(&defaultContext),
		E.ToError[fp.ResourceData],
	)

	assert.NoError(t, res)

	checksums, ok := data[common.KeyChecksums].(attestation.ChecksumMap)
	assert.NotNil(t, checksums)
	assert.True(t, ok)

	assert.Equal(t, "a6f6228bbf820e766ebe43c51e97332dda92e9744e719a646f611fe0681d2458", checksums["cidata/user-data"])
}
