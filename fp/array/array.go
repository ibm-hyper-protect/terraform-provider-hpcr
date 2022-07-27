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
package array

import (
	F "github.com/terraform-provider-hpcr/fp/function"
)

func MonadMap[A, B any](as []A, f func(a A) B) []B {
	count := len(as)
	bs := make([]B, count)
	for i := count - 1; i >= 0; i-- {
		bs[i] = f(as[i])
	}
	return bs
}

func Map[A, B any](f func(a A) B) func([]A) []B {
	return F.Bind2nd(MonadMap[A, B], f)
}

func reduce[A, B any](fa []A, f func(B, A) B, initial B) B {
	current := initial
	count := len(fa)
	for i := 0; i < count; i++ {
		current = f(current, fa[i])
	}
	return current
}

func Append[A any](as []A, a A) []A {
	return append(as, a)
}

func Empty[A any]() []A {
	return make([]A, 0)
}

func MonadChain[A, B any](fa []A, f func(a A) []B) []B {
	return reduce(fa, func(bs []B, a A) []B {
		return concat(bs, f(a))
	}, Empty[B]())
}

func Flatten[A any](mma [][]A) []A {
	return MonadChain(mma, F.Identity[[]A])
}
