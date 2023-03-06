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
package validation

import (
	"context"
	"log"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	D "github.com/ibm-hyper-protect/terraform-provider-hpcr/data"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
	A "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/array"
	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	J "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/json"
	O "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/option"
	S "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/string"
	Y "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/yaml"
	"github.com/qri-io/jsonschema"
)

type RawMap = map[string]any

func validate[A any](schema *jsonschema.Schema) func(A) []jsonschema.KeyError {
	return func(data A) []jsonschema.KeyError {
		return *schema.Validate(context.Background(), data).Errs
	}
}

var (
	// predicate to check if there exist errors
	hasErrorsO = O.FromPredicate(func(errs []jsonschema.KeyError) bool {
		return len(errs) > 0
	})

	// to convert from validation errors to a single error
	handleValidationErrorsO = F.Flow2(
		hasErrorsO,
		O.Map(func(errs []jsonschema.KeyError) error {
			log.Println(errs)
			return errs[0]
		}),
	)
)

// ValidateYAML validates a YAML file against the validator function by deserializing it, then validate the result
func ValidateYAML[A any](validator func(A) []jsonschema.KeyError) func(data string) E.Either[error, A] {
	return F.Flow4(
		S.ToBytes,
		Y.Parse[A],
		E.Map[error](F.Deref[A]),
		E.Chain(func(a A) E.Either[error, A] {
			return F.Pipe2(
				a,
				validator,
				F.Flow2(
					handleValidationErrorsO,
					O.Fold(F.Constant(E.Of[error](a)), E.Left[error, A]),
				),
			)
		}),
	)
}

func fromKeyError(err jsonschema.KeyError) []diag.Diagnostic {
	return diag.FromErr(err)
}

func schemaToDiagnostics(errs []jsonschema.KeyError) diag.Diagnostics {
	return F.Pipe2(
		errs,
		A.Map(fromKeyError),
		A.Flatten[diag.Diagnostic],
	)
}

func diagYAML[A any](validator func(A) []jsonschema.KeyError) func(data string) E.Either[error, diag.Diagnostics] {
	return F.Flow4(
		S.ToBytes,
		Y.Parse[A],
		E.Map[error](F.Deref[A]),
		E.Map[error](F.Flow2(
			validator,
			schemaToDiagnostics,
		),
		),
	)
}

// GetContractSchema reads the json schema from a string representation into a schema representation
func GetContractSchema() E.Either[error, *jsonschema.Schema] {
	return F.Pipe2(
		D.ContractSchema,
		S.ToBytes,
		J.Parse[jsonschema.Schema],
	)
}

// DiagContract validates that the given certificate is indeed a certificate
func DiagContract(data any, _ cty.Path) diag.Diagnostics {
	// convert the key
	dataE := F.Pipe1(
		data,
		fp.ToTypeE[string],
	)
	// combine
	return F.Pipe4(
		GetContractSchema(),
		E.Map[error](F.Flow2(
			validate[RawMap],
			diagYAML[RawMap],
		)),
		E.Ap[error, string, E.Either[error, diag.Diagnostics]](dataE),
		E.Flatten[error, diag.Diagnostics],
		E.GetOrElse(diag.FromErr),
	)
}
