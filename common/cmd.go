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
	"fmt"
	"os/exec"

	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	T "github.com/terraform-provider-hpcr/fp/tuple"
)

type CommandOutput = T.Tuple2[[]byte, []byte]

var (
	GetStdOut = T.FirstOf2[[]byte, []byte]
	GetStdErr = T.SecondOf2[[]byte, []byte]
)

func ExecCommand(name string, arg ...string) func([]byte) E.Either[error, CommandOutput] {
	return func(dataIn []byte) E.Either[error, CommandOutput] {
		// command result
		var stdOut bytes.Buffer
		var stdErr bytes.Buffer
		// execute the command
		return F.Pipe1(
			// run the command
			E.TryCatchError(func() (CommandOutput, error) {
				// command input
				cmd := exec.Command(name, arg...)
				cmd.Stdin = bytes.NewReader(dataIn)

				cmd.Stdout = &stdOut
				cmd.Stderr = &stdErr

				err := cmd.Run()
				return T.MakeTuple2(stdOut.Bytes(), stdErr.Bytes()), err
			}),
			// enrich the error
			E.MapLeft[error, CommandOutput](func(cause error) error {
				return fmt.Errorf("command execution of [%s][%s] failed, stdout [%s], stderr [%s], cause [%w]", name, arg, stdOut.String(), stdErr.String(), cause)
			}),
		)
	}
}
