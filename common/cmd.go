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
	"bytes"
	"os/exec"

	E "github.com/terraform-provider-hpcr/fp/either"
	T "github.com/terraform-provider-hpcr/fp/tuple"
)

type CommandOutput = T.Tuple2[[]byte, []byte]

var (
	GetStdOut = T.FirstOf2[[]byte, []byte]
	GetStdErr = T.SecondOf2[[]byte, []byte]
)

func ExecCommand(name string, arg ...string) func([]byte) E.Either[error, CommandOutput] {
	return func(dataIn []byte) E.Either[error, CommandOutput] {
		return E.TryCatchError(func() (CommandOutput, error) {
			// command input
			cmd := exec.Command(name, arg...)
			cmd.Stdin = bytes.NewReader(dataIn)
			// command result
			var stdOut bytes.Buffer
			var stdErr bytes.Buffer

			cmd.Stdout = &stdOut
			cmd.Stderr = &stdErr

			err := cmd.Run()
			return T.MakeTuple2(stdOut.Bytes(), stdErr.Bytes()), err
		})
	}
}
