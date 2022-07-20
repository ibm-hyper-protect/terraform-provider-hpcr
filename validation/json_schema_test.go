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
	"fmt"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
)

//go:embed samples/simple.yml
var TrivialContract string

func TestJsonSchema(t *testing.T) {
	schemaE := GetContractSchema()
	assert.True(t, schemaE.IsRight())
}

func TestValidate(t *testing.T) {
	// validator function
	validatorE := F.Pipe1(
		GetContractSchema(),
		E.Map[error](F.Flow2(
			validate[RawMap],
			ValidateYAML[RawMap],
		)),
	)
	// validate the data
	resE := F.Pipe1(
		validatorE,
		E.Chain(I.Ap[string, E.Either[error, RawMap]](TrivialContract)),
	)

	fmt.Println(resE)
}
