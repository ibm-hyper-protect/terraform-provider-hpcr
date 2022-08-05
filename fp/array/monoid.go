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
	M "github.com/terraform-provider-hpcr/fp/monoid"
)

func concat[A any](left, right []A) []A {
	buf := make([]A, len(left)+len(right))
	copy(buf[copy(buf, left):], right)
	return buf
}

// Monoid for arrays
func Monoid[A any]() M.Monoid[[]A] {
	return M.MakeMonoid(concat[A], Empty[A]())
}

func addLen[A any](count int, data []A) int {
	return count + len(data)
}

// ConcatAll efficiently concatenates the input arrays into a final array
func ConcatAll[A any](data ...[]A) []A {
	// get the full size
	count := reduce(data, addLen[A], 0)
	buf := make([]A, count)
	// copy
	reduce(data, func(idx int, seg []A) int {
		return idx + copy(buf[idx:], seg)
	}, 0)
	// returns the final array
	return buf
}
