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
package identity

import (
	F "github.com/terraform-provider-hpcr/fp/function"
)

func MonadAp[A, B any](fab func(A) B, fa A) B {
	return fab(fa)
}

func Ap[A, B any](fa A) func(func(A) B) B {
	return F.Bind2nd(MonadAp[A, B], fa)
}

func MonadMap[A, B any](fa A, f func(A) B) B {
	return f(fa)
}

func Map[A, B any](f func(A) B) func(A) B {
	return f
}

func Of[A any](a A) A {
	return a
}

func MonadChain[A, B any](ma A, f func(A) B) B {
	return f(ma)
}

func Chain[A, B any](f func(A) B) func(A) B {
	return f
}
