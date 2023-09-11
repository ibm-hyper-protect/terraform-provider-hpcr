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

package attestation

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/encrypt"
	"github.com/stretchr/testify/assert"

	B "github.com/IBM/fp-go/bytes"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	R "github.com/IBM/fp-go/record"
)

//go:embed samples/encrypted/attestation.base64
var encBase64 string

//go:embed samples/encrypted/attestation.pub
var encPub string

//go:embed samples/encrypted/attestation
var encPriv []byte

var (
	getChecksums     = R.Lookup[[]byte]("se-checksums.txt.enc")
	getContract      = R.Lookup[string]("cidata/user-data")
	defaultDecryptor = encrypt.DefaultDecryption()
)

func TestParseAndDecryptAttestation(t *testing.T) {
	// decode the sample
	token := F.Pipe3(
		encBase64,
		untarBase64,
		E.ChainOptionK[FileList, []byte](func() error { return fmt.Errorf("unable to read checksums") })(getChecksums),
		E.Map[error](B.ToString),
	)

	// get the decryptor
	dec := DecryptAttestation(defaultDecryptor.DecryptBasic)(encPriv)

	// decrypt the record
	checksum := F.Pipe2(
		token,
		E.Chain(dec),
		E.ChainOptionK[ChecksumMap, string](func() error { return fmt.Errorf("unable to read checksum for contracz") })(getContract),
	)

	assert.Equal(t, E.Of[error]("a6f6228bbf820e766ebe43c51e97332dda92e9744e719a646f611fe0681d2458"), checksum)
}
