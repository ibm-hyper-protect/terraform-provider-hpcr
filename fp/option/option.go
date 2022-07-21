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
package Option

import (
	"fmt"

	F "github.com/terraform-provider-hpcr/fp/function"
)

type Option[T any] interface {
	IsNone() bool
	IsSome() bool
}

type none []struct {
}

type some[T any] struct {
	value T
}

func (none) IsNone() bool {
	return true
}

func (none) IsSome() bool {
	return false
}

func (none) String() string {
	return "None"
}

func (some[T]) IsNone() bool {
	return false
}

func (some[T]) IsSome() bool {
	return true
}

func (s some[T]) String() string {
	return fmt.Sprintf("Some[%v]", s.value)
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

func IsSome[T any](val Option[T]) bool {
	return val.IsSome()
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

func FromNillable[A any](a *A) Option[*A] {
	return fromPredicate(a, F.IsNonNil[A])
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

func MonadAp[A, B any](fab Option[func(A) B], fa Option[A]) Option[B] {
	if IsNone(fab) || IsNone(fa) {
		return None[B]()
	}
	// apply
	return Some(fab.(some[func(A) B]).value(fa.(some[A]).value))
}

func Ap[A, B any](fa Option[A]) func(Option[func(A) B]) Option[B] {
	return F.Bind2nd(MonadAp[A, B], fa)
}

func MonadMap[A, B any](fa Option[A], f func(a A) B) Option[B] {
	return fold(fa, None[B], F.Flow2(f, Some[B]))
}

func Map[A, B any](f func(a A) B) func(fa Option[A]) Option[B] {
	return F.Bind2nd(MonadMap[A, B], f)
}

func TryCatch[A any](f func() (A, error)) Option[A] {
	val, err := f()
	if err != nil {
		return None[A]()
	}
	return Some(val)
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

func MonadChainTo[A, B any](ma Option[A], mb Option[B]) Option[B] {
	return mb
}

func ChainTo[E, A, B any](mb Option[B]) func(Option[A]) Option[B] {
	return F.Bind2nd(MonadChainTo[A, B], mb)
}

func MonadChainFirst[E, A, B any](ma Option[A], f func(a A) Option[B]) Option[A] {
	return MonadChain(ma, func(a A) Option[A] {
		return MonadMap(f(a), F.Constant1[B](a))
	})
}

func ChainFirst[E, A, B any](f func(a A) Option[B]) func(fa Option[A]) Option[A] {
	return F.Bind2nd(MonadChainFirst[E, A, B], f)
}

func Flatten[A any](mma Option[Option[A]]) Option[A] {
	return MonadChain(mma, F.Identity[Option[A]])
}

func Alt[A any](that func() Option[A]) func(Option[A]) Option[A] {
	return Fold(that, Of[A])
}
