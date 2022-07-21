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
	O "github.com/terraform-provider-hpcr/fp/option"
	T "github.com/terraform-provider-hpcr/fp/tuple"
)

func IsEmpty[K comparable, V any](r map[K]V) bool {
	return len(r) == 0
}

func IsNonEmpty[K comparable, V any](r map[K]V) bool {
	return len(r) > 0
}

func Keys[K comparable, V any](r map[K]V) []K {
	return collect(r, F.First[K, V])
}

func Values[K comparable, V any](r map[K]V) []V {
	return collect(r, F.Second[K, V])
}

func collect[K comparable, V, R any](r map[K]V, f func(K, V) R) []R {
	count := len(r)
	result := make([]R, count)
	idx := 0
	for k, v := range r {
		result[idx] = f(k, v)
		idx++
	}
	return result
}

func Collect[K comparable, V, R any](f func(K, V) R) func(map[K]V) []R {
	return F.Bind2nd(collect[K, V, R], f)
}

func reduce[K comparable, V, R any](r map[K]V, f func(R, V) R, initial R) R {
	current := initial
	for _, v := range r {
		current = f(current, v)
	}
	return current
}

func reduceWithIndex[K comparable, V, R any](r map[K]V, f func(K, R, V) R, initial R) R {
	current := initial
	for k, v := range r {
		current = f(k, current, v)
	}
	return current
}

func reduceRef[K comparable, V, R any](r map[K]V, f func(R, *V) R, initial R) R {
	current := initial
	for _, v := range r {
		current = f(current, &v) // #nosec G601
	}
	return current
}

func reduceRefWithIndex[K comparable, V, R any](r map[K]V, f func(K, R, *V) R, initial R) R {
	current := initial
	for k, v := range r {
		current = f(k, current, &v) // #nosec G601
	}
	return current
}

func Reduce[K comparable, V, R any](f func(R, V) R, initial R) func(map[K]V) R {
	return func(r map[K]V) R {
		return reduce(r, f, initial)
	}
}

func ReduceWithIndex[K comparable, V, R any](f func(K, R, V) R, initial R) func(map[K]V) R {
	return func(r map[K]V) R {
		return reduceWithIndex(r, f, initial)
	}
}

func ReduceRef[K comparable, V, R any](f func(R, *V) R, initial R) func(map[K]V) R {
	return func(r map[K]V) R {
		return reduceRef(r, f, initial)
	}
}

func ReduceRefWithIndex[K comparable, V, R any](f func(K, R, *V) R, initial R) func(map[K]V) R {
	return func(r map[K]V) R {
		return reduceRefWithIndex(r, f, initial)
	}
}

func MonadMap[K comparable, V, R any](r map[K]V, f func(V) R) map[K]R {
	return MonadMapWithIndex(r, F.Ignore1st[K](f))
}

func MonadMapWithIndex[K comparable, V, R any](r map[K]V, f func(K, V) R) map[K]R {
	return reduceWithIndex(r, func(k K, dst map[K]R, v V) map[K]R {
		return upsertAtReadWrite(dst, k, f(k, v))
	}, make(map[K]R, len(r)))
}

func MonadMapRefWithIndex[K comparable, V, R any](r map[K]V, f func(K, *V) R) map[K]R {
	return reduceRefWithIndex(r, func(k K, dst map[K]R, v *V) map[K]R {
		return upsertAtReadWrite(dst, k, f(k, v))
	}, make(map[K]R, len(r)))
}

func MonadMapRef[K comparable, V, R any](r map[K]V, f func(*V) R) map[K]R {
	return MonadMapRefWithIndex(r, F.Ignore1st[K](f))
}

func Map[K comparable, V, R any](f func(V) R) func(map[K]V) map[K]R {
	return F.Bind2nd(MonadMap[K, V, R], f)
}

func MapRef[K comparable, V, R any](f func(*V) R) func(map[K]V) map[K]R {
	return F.Bind2nd(MonadMapRef[K, V, R], f)
}

func MapWithIndex[K comparable, V, R any](f func(K, V) R) func(map[K]V) map[K]R {
	return F.Bind2nd(MonadMapWithIndex[K, V, R], f)
}

func MapRefWithIndex[K comparable, V, R any](f func(K, *V) R) func(map[K]V) map[K]R {
	return F.Bind2nd(MonadMapRefWithIndex[K, V, R], f)
}

func lookup[K comparable, V any](r map[K]V, k K) O.Option[V] {
	if val, ok := r[k]; ok {
		return O.Some(val)
	}
	return O.None[V]()
}

func Lookup[K comparable, V any](k K) func(map[K]V) O.Option[V] {
	return F.Bind2nd(lookup[K, V], k)
}

func Has[K comparable, V any](k K, r map[K]V) bool {
	_, ok := r[k]
	return ok
}

func Empty[K comparable, V any]() map[K]V {
	return make(map[K]V)
}

func Size[K comparable, V any](r map[K]V) int {
	return len(r)
}

func ToArray[K comparable, V any](r map[K]V) []T.Tuple2[K, V] {
	return collect(r, T.MakeTuple2[K, V])
}

func duplicate[K comparable, V any](r map[K]V) map[K]V {
	return MonadMap(r, F.Identity[V])
}

func upsertAt[K comparable, V any](r map[K]V, k K, v V) map[K]V {
	dup := duplicate(r)
	dup[k] = v
	return dup
}

func upsertAtReadWrite[K comparable, V any](r map[K]V, k K, v V) map[K]V {
	r[k] = v
	return r
}

func UpsertAt[K comparable, V any](k K, v V) func(map[K]V) map[K]V {
	return func(ma map[K]V) map[K]V {
		return upsertAt(ma, k, v)
	}
}

func Singleton[K comparable, V any](k K, v V) map[K]V {
	return map[K]V{k: v}
}
