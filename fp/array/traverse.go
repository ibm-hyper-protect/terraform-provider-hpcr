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
)

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<[]B>
HKTA = HKT<A>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func traverse[A, B, HKTA, HKTB, HKTAB, HKTRB any](
	_of func([]B) HKTRB,
	_map func(HKTRB, func([]B) func(B) []B) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,
) func([]A, func(A) HKTB) HKTRB {
	cb := F.Curry2(Append[B])

	return func(ta []A, f func(A) HKTB) HKTRB {
		return reduce(ta, func(r HKTRB, a A) HKTRB {
			return _ap(
				_map(r, cb),
				f(a),
			)
		}, _of(Empty[B]()))
	}
}

func Traverse[A, B, HKTA, HKTB, HKTAB, HKTRB any](
	_of func([]B) HKTRB,
	_map func(HKTRB, func([]B) func(B) []B) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,
) func(func(A) HKTB) func([]A) HKTRB {
	delegate := traverse[A, B, HKTA](_of, _map, _ap)
	return func(f func(A) HKTB) func([]A) HKTRB {
		return F.Bind2nd(delegate, f)
	}
}
