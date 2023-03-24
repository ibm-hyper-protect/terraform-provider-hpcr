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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	A "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/array"
	B "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/bytes"
	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	R "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/record"
)

type (
	FileList    = map[string][]byte
	ChecksumMap = map[string]string
)

var (
	// expression used to match one line of the attestation record
	reChecksum = regexp.MustCompile(`^\s*([0-9a-f]+)\s+([^\s]+)\s*$`)
)

func untarBase64(data string) E.Either[error, FileList] {
	// build the streams
	input := bytes.NewBuffer([]byte(data))
	base64Dec := base64.NewDecoder(base64.StdEncoding, input)
	gzip, err := gzip.NewReader(base64Dec)
	if err != nil {
		return E.Left[error, FileList](err)
	}
	untar := tar.NewReader(gzip)

	res := make(FileList)

	for {
		hdr, err := untar.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return E.Left[error, FileList](err)
		}
		if hdr.Typeflag == tar.TypeReg {
			// read the content
			data, err := io.ReadAll(untar)
			if err != nil {
				return E.Left[error, FileList](err)
			}
			// record this
			res[filepath.ToSlash(filepath.Clean(hdr.Name))] = data
		}
	}
	return E.Of[error](res)

}

// ParseAttestationRecord parses a text representation of a checksum record into a map
var ParseAttestationRecord = F.Flow5(
	F.Bind2nd(strings.Split, "\n"),
	A.Map(strings.TrimSpace),
	A.Map(reChecksum.FindStringSubmatch),
	A.Filter(A.IsNonEmpty[string]),
	A.Reduce(func(res ChecksumMap, entry []string) ChecksumMap { return R.UpsertAt(entry[2], entry[1])(res) }, make(ChecksumMap)),
)

// DecryptAttestation decrypts an attestation record into a mapping
func DecryptAttestation(decrypter func(privKey []byte) func(string) E.Either[error, []byte]) func(privKey []byte) func(string) E.Either[error, ChecksumMap] {
	return func(privKey []byte) func(string) E.Either[error, ChecksumMap] {
		return F.Flow2(
			decrypter(privKey),
			E.Map[error](F.Flow2(
				B.ToString,
				ParseAttestationRecord,
			)),
		)
	}
}
