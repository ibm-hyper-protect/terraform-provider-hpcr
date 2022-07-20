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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandOk(t *testing.T) {

	cmdE := ExecCommand("openssl", "help")(make([]byte, 0))

	assert.True(t, cmdE.IsRight())
}

func TestCommandFail(t *testing.T) {

	cmdE := ExecCommand("openssl", "help1")(make([]byte, 0))

	fmt.Println(cmdE)

	assert.True(t, cmdE.IsLeft())
}
