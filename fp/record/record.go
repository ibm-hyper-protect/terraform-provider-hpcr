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
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	O "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/option"
)

func reduceWithIndex[K comparable, V, R any](r map[K]V, f func(K, R, V) R, initial R) R {
	current := initial
	for k, v := range r {
		current = f(k, current, v)
	}
	return current
}

func MonadMap[K comparable, V, R any](r map[K]V, f func(V) R) map[K]R {
	return MonadMapWithIndex(r, F.Ignore1st[K](f))
}

func MonadMapWithIndex[K comparable, V, R any](r map[K]V, f func(K, V) R) map[K]R {
	return reduceWithIndex(r, func(k K, dst map[K]R, v V) map[K]R {
		return upsertAtReadWrite(dst, k, f(k, v))
	}, make(map[K]R, len(r)))
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
