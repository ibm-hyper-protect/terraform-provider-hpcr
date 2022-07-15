//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package array

import (
	F "github.com/terraform-provider-hpcr/fp/function"
)

/**
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
