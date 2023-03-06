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
package predicate

import (
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
)

func Not[A any](predicate func(A) bool) func(A) bool {
	return func(a A) bool {
		return !predicate((a))
	}
}

// ContraMap creates a predicate from an existing predicate given a mapping functio
func ContraMap[A, B any](f func(B) A) func(func(A) bool) func(B) bool {
	return func(pred func(A) bool) func(B) bool {
		return F.Flow2(
			f,
			pred,
		)
	}
}
