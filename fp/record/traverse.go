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
// limitations under the License.package datasource

package record

import (
	F "github.com/terraform-provider-hpcr/fp/function"
)

func upsertTraverse[K comparable, A, B any]() func(K) func(map[K]B) func(B) map[K]B {
	return func(k K) func(map[K]B) func(B) map[K]B {
		return func(r map[K]B) func(B) map[K]B {
			return func(b B) map[K]B {
				return upsertAt(r, k, b)
			}
		}
	}
}

/**
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<map[K]B>
HKTA = HKT<A>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func traverseWithIndex[K comparable, A, B, HKTA, HKTB, HKTAB, HKTRB any](
	_of func(map[K]B) HKTRB,
	_map func(HKTRB, func(map[K]B) func(B) map[K]B) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,
) func(map[K]A, func(K, A) HKTB) HKTRB {
	cb := upsertTraverse[K, A, B]()

	return func(ta map[K]A, f func(K, A) HKTB) HKTRB {
		return reduceWithIndex(ta, func(k K, r HKTRB, a A) HKTRB {
			return _ap(
				_map(r, cb(k)),
				f(k, a),
			)
		}, _of(Empty[K, B]()))
	}
}

func traverse[K comparable, A, B, HKTA, HKTB, HKTAB, HKTRB any](
	_of func(map[K]B) HKTRB,
	_map func(HKTRB, func(map[K]B) func(B) map[K]B) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,
) func(map[K]A, func(A) HKTB) HKTRB {
	delegate := traverseWithIndex[K, A, B, HKTA](_of, _map, _ap)
	return func(r map[K]A, f func(A) HKTB) HKTRB {
		return delegate(r, F.Ignore1st[K](f))
	}
}

func TraverseWithIndex[K comparable, A, B, HKTA, HKTB, HKTAB, HKTRB any](
	_of func(map[K]B) HKTRB,
	_map func(HKTRB, func(map[K]B) func(B) map[K]B) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,
) func(func(K, A) HKTB) func(map[K]A) HKTRB {
	delegate := traverseWithIndex[K, A, B, HKTA](_of, _map, _ap)
	return func(f func(K, A) HKTB) func(map[K]A) HKTRB {
		return F.Bind2nd(delegate, f)
	}
}

func Traverse[K comparable, A, B, HKTA, HKTB, HKTAB, HKTRB any](
	_of func(map[K]B) HKTRB,
	_map func(HKTRB, func(map[K]B) func(B) map[K]B) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,
) func(func(A) HKTB) func(map[K]A) HKTRB {
	delegate := traverse[K, A, B, HKTA](_of, _map, _ap)
	return func(f func(A) HKTB) func(map[K]A) HKTRB {
		return F.Bind2nd(delegate, f)
	}
}
