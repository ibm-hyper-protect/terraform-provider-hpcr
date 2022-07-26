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

func filter[A any](fa []A, pred func(a A) bool) []A {
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

func filterRef[A any](fa []A, pred func(a *A) bool) []A {
	var result []A
	count := len(fa)
	for i := 0; i < count; i++ {
		a := fa[i]
		if pred(&a) {
			result = append(result, a)
		}
	}
	return result
}

func filterMap[A, B any](fa []A, pred func(a A) bool, f func(a A) B) []B {
	var result []B
	count := len(fa)
	for i := 0; i < count; i++ {
		a := fa[i]
		if pred(a) {
			result = append(result, f(a))
		}
	}
	return result
}

func filterMapRef[A, B any](fa []A, pred func(a *A) bool, f func(a *A) B) []B {
	var result []B
	count := len(fa)
	for i := 0; i < count; i++ {
		a := fa[i]
		if pred(&a) {
			result = append(result, f(&a))
		}
	}
	return result
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
