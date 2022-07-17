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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-provider-hpcr/data"
)

func TestDiagPath(t *testing.T) {
	assert.Len(t, DiagFolder(".", nil), 0)
}

func TestDiagNoPath(t *testing.T) {
	assert.Len(t, DiagFolder("../README.md", nil), 1)
}

func TestDiagCertificate(t *testing.T) {
	assert.Len(t, DiagCertificate(data.DefaultCertificate, nil), 0)
}
