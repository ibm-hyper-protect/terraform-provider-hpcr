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
package either

import (
	RA "github.com/terraform-provider-hpcr/fp/array"
	F "github.com/terraform-provider-hpcr/fp/function"
)

func TraverseArray[E, A, B any](f func(A) Either[E, B]) func([]A) Either[E, []B] {
	return F.Pipe1(
		f,
		RA.Traverse[A, B, Either[E, A]](
			Of[E, []B],
			MonadMap[E, []B, func(B) []B],
			MonadAp[E, B, []B],
		),
	)
}

func SequenceArray[E, A any]() func([]Either[E, A]) Either[E, []A] {
	return TraverseArray(F.Identity[Either[E, A]])
}
