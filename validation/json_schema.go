//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package validation

import (
	"context"
	"log"

	"github.com/qri-io/jsonschema"
	D "github.com/terraform-provider-hpcr/data"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	J "github.com/terraform-provider-hpcr/fp/json"
	O "github.com/terraform-provider-hpcr/fp/option"
	S "github.com/terraform-provider-hpcr/fp/string"
	Y "github.com/terraform-provider-hpcr/fp/yaml"
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

// validates a YAML file against the validator function by deserializing it, then validate the result
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

// reads the json schema from a string representation into a schema representation
func GetContractSchema() E.Either[error, *jsonschema.Schema] {
	return F.Pipe2(
		D.ContractSchema,
		S.ToBytes,
		J.Parse[jsonschema.Schema],
	)
}
