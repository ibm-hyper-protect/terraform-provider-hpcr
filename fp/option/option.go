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
package option

import (
	F "github.com/terraform-provider-hpcr/fp/function"
)

type Option[T any] interface {
	IsNone() bool
}

type none []struct {
}

type some[T any] struct {
	value T
}

func (none) IsNone() bool {
	return true
}

func (some[T]) IsNone() bool {
	return false
}

func IsNone[T any](val Option[T]) bool {
	return val.IsNone()
}

func Some[T any](value T) Option[T] {
	return some[T]{value: value}
}

func Of[T any](value T) Option[T] {
	return Some(value)
}

func None[T any]() Option[T] {
	return none{}
}

func fromPredicate[A any](a A, pred func(a A) bool) Option[A] {
	if pred(a) {
		return Some(a)
	}
	return None[A]()
}

func FromPredicate[A any](pred func(a A) bool) func(value A) Option[A] {
	return F.Bind2nd(fromPredicate[A], pred)
}

func fromValidation[A, B any](a A, f func(a A) (B, bool)) Option[B] {
	b, ok := f(a)
	if ok {
		return Some(b)
	}
	return None[B]()
}

func FromValidation[A, B any](f func(a A) (B, bool)) func(a A) Option[B] {
	return F.Bind2nd(fromValidation[A, B], f)
}

func MonadMap[A, B any](fa Option[A], f func(a A) B) Option[B] {
	return fold(fa, None[B], F.Flow2(f, Some[B]))
}

func Map[A, B any](f func(a A) B) func(fa Option[A]) Option[B] {
	return F.Bind2nd(MonadMap[A, B], f)
}

func fold[A, B any](ma Option[A], onNone func() B, onSome func(a A) B) B {
	if IsNone(ma) {
		return onNone()
	}
	return onSome(ma.(some[A]).value)
}

func Fold[A, B any](onNone func() B, onSome func(a A) B) func(ma Option[A]) B {
	return func(ma Option[A]) B {
		return fold(ma, onNone, onSome)
	}
}

func GetOrElse[A any](onNone func() A) func(Option[A]) A {
	return Fold(onNone, F.Identity[A])
}

func MonadChain[A, B any](fa Option[A], f func(a A) Option[B]) Option[B] {
	return fold(fa, None[B], f)
}

func Chain[A, B any](f func(a A) Option[B]) func(fa Option[A]) Option[B] {
	return F.Bind2nd(MonadChain[A, B], f)
}

func MonadSequence2[T1, T2, R any](o1 Option[T1], o2 Option[T2], f func(T1, T2) Option[R]) Option[R] {
	if IsNone(o1) {
		return None[R]()
	}
	if IsNone(o2) {
		return None[R]()
	}
	return f(o1.(some[T1]).value, o2.(some[T2]).value)
}

func Sequence2[T1, T2, R any](f func(T1, T2) Option[R]) func(Option[T1], Option[T2]) Option[R] {
	return func(o1 Option[T1], o2 Option[T2]) Option[R] {
		return MonadSequence2(o1, o2, f)
	}
}
