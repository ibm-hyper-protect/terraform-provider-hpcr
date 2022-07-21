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
// limitations under the License.
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
