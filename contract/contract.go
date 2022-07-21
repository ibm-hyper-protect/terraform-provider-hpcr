// Copyright 2022 IBM Corp.
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

package contract

import (
	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/encrypt"
	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	R "github.com/terraform-provider-hpcr/fp/record"
)

var (
	KeyEnv        = "env"
	KeyWorkload   = "workload"
	KeySigningKey = "signingKey"

	getWorkload   = R.Lookup[string, any](KeyWorkload)
	getEnv        = R.Lookup[string, any](KeyEnv)
	getSigningKey = R.Lookup[string, any](KeySigningKey)
)

type RawMap = map[string]any

func toAny[A any](a A) any {
	return a
}

func addSigningKey(key []byte) func(any) E.Either[error, any] {
	// function to add the pkey into a map
	pemE := F.Pipe4(
		key,
		encrypt.PublicKey,
		E.Map[error](B.ToString),
		E.Map[error](toAny[string]),
		E.Map[error](F.Bind1st(R.UpsertAt[string, any], KeySigningKey)),
	)

	return func(data any) E.Either[error, any] {
		// insert into the map
		return F.Pipe2(
			pemE,
			E.Ap[error, RawMap, RawMap](common.ToTypeE[RawMap](data)),
			E.Map[error](toAny[RawMap]),
		)
	}
}
