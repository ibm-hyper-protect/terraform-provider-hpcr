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
package array

import (
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	O "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/option"
)

func MakeBy[A any](n int, f func(int) A) []A {
	as := make([]A, n)
	for i := n - 1; i >= 0; i-- {
		as[i] = f(i)
	}
	return as
}

func Replicate[A any](n int, a A) []A {
	return MakeBy(n, F.Constant1[int](a))
}

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
		return append(bs, f(a)...)
	}, Empty[B]())
}

func Flatten[A any](mma [][]A) []A {
	return MonadChain(mma, F.Identity[[]A])
}

func filter[A any](fa []A, pred func(A) bool) []A {
	var result []A
	count := len(fa)
	for i := 0; i < count; i++ {
		a := fa[i]
		if pred(a) {
			result = append(result, a)
		}
	}
	return result
}

func Filter[A any](pred func(A) bool) func([]A) []A {
	return F.Bind2nd(filter[A], pred)
}

func Size[A any](fa []A) int {
	return len(fa)
}

func Head[A any](as []A) O.Option[A] {
	if len(as) <= 0 {
		return O.None[A]()
	}
	return O.Of(as[0])
}

func IsNonEmpty[A any](data []A) bool {
	return len(data) > 0
}

func Reduce[A, B any](f func(B, A) B, initial B) func([]A) B {
	return func(as []A) B {
		return reduce(as, f, initial)
	}
}
