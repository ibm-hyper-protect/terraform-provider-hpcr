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

package semigroup

import (
	M "github.com/terraform-provider-hpcr/fp/magma"
)

type Semigroup[A any] interface {
	M.Magma[A]
}

type semigroup[A any] struct {
	c func(A, A) A
}

func (self semigroup[A]) Concat(x A, y A) A {
	return self.c(x, y)
}

func MakeSemigroup[A any](c func(A, A) A) Semigroup[A] {
	return semigroup[A]{c: c}
}

func Reverse[A any](m Semigroup[A]) Semigroup[A] {
	return MakeSemigroup(M.Reverse[A](m).Concat)
}
