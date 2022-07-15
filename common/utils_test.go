//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package common

import (
	"io/fs"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

var statE = E.Eitherize1(os.Stat)

func TestTempFile(t *testing.T) {
	resE := F.Pipe3(
		CreateTempE("", "*"),
		E.Map[error](func(f *os.File) string {
			return f.Name()
		}),
		E.Chain(statE),
		E.Map[error](fs.FileInfo.IsDir),
	)

	assert.Equal(t, resE, E.Of[error](false))
}
