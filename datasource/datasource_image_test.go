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
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	"github.com/stretchr/testify/assert"
)

//go:embed samples/images.json
var Images string

func TestParseImages(t *testing.T) {
	res := F.Pipe1(
		Images,
		parseImages,
	)
	assert.True(t, E.IsRight(res))
}

func TestReHyperProtectOS(t *testing.T) {
	assert.True(t, reHyperProtectOS.MatchString("hyper-protect-1-0-s390x"))
}

func TestReHyperProtectName(t *testing.T) {
	assert.True(t, reHyperProtectName.MatchString("ibm-hyper-protect-container-runtime-1-0-s390x-8"))
}

func TestFilterImages(t *testing.T) {
	res := F.Pipe2(
		Images,
		parseImages,
		E.Map[error](F.Flow2(
			A.Filter(isCandidateImage),
			A.Size[Image],
		)),
	)
	assert.Equal(t, E.Of[error](2), res)
}

func TestFilterAndSortImages(t *testing.T) {
	res := F.Pipe2(
		Images,
		parseImages,
		E.Map[error](F.Flow3(
			A.Filter(isCandidateImage),
			A.Map(imageVersionFomImage),
			I.Map(sortByVersion),
		)),
	)
	fmt.Println(res)
}

func TestSelectBySpec(t *testing.T) {
	res := F.Pipe4(
		Images,
		parseImages,
		E.Map[error](F.Flow2(
			A.Filter(isCandidateImage),
			A.Map(imageVersionFomImage),
		)),
		E.ChainOptionK[error, []ImageVersion, ImageVersion](noMatchingVersionFound)(selectBySpec("*")),
		E.Map[error](func(version ImageVersion) string {
			return version.Version.String()
		}),
	)
	assert.Equal(t, E.Of[error]("1.0.8"), res)
}
