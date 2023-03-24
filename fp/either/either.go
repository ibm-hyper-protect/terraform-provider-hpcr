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
package either

import (
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	O "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/option"
)

type Either[E, A any] interface {
	IsLeft() bool
	IsRight() bool
}

type left[E any] struct {
	e E
}

type right[A any] struct {
	a A
}

func (left[E]) IsLeft() bool {
	return true
}

func (left[E]) IsRight() bool {
	return false
}

func (right[A]) IsLeft() bool {
	return false
}

func (right[A]) IsRight() bool {
	return true
}

func IsLeft[E, A any](val Either[E, A]) bool {
	return val.IsLeft()
}

func IsRight[E, A any](val Either[E, A]) bool {
	return val.IsRight()
}

func Left[E, A any](value E) Either[E, A] {
	return left[E]{e: value}
}

func Right[E, A any](value A) Either[E, A] {
	return right[A]{a: value}
}

func Of[E, A any](value A) Either[E, A] {
	return F.Pipe1(value, Right[E, A])
}

func MonadAp[E, A, B any](fab Either[E, func(a A) B], fa Either[E, A]) Either[E, B] {
	if IsLeft(fab) {
		return Left[E, B](fab.(left[E]).e)
	}
	if IsLeft(fa) {
		return Left[E, B](fa.(left[E]).e)
	}
	return Right[E](fab.(right[func(a A) B]).a(fa.(right[A]).a))
}

func Ap[E, A, B any](fa Either[E, A]) func(fab Either[E, func(a A) B]) Either[E, B] {
	return F.Bind2nd(MonadAp[E, A, B], fa)
}

func MonadMap[E, A, B any](fa Either[E, A], f func(a A) B) Either[E, B] {
	return MonadChain(fa, F.Flow2(f, Right[E, B]))
}

func MonadMapTo[E, A, B any](fa Either[E, A], b B) Either[E, B] {
	return MonadMap(fa, F.Constant1[A](b))
}

func MapTo[E, A, B any](b B) func(Either[E, A]) Either[E, B] {
	return F.Bind2nd(MonadMapTo[E, A, B], b)
}

func MonadMapLeft[E, A, B any](fa Either[E, A], f func(E) B) Either[B, A] {
	return fold(fa, F.Flow2(f, Left[B, A]), Right[B, A])
}

func Map[E, A, B any](f func(a A) B) func(fa Either[E, A]) Either[E, B] {
	return F.Bind2nd(MonadMap[E, A, B], f)
}

func MapLeft[E, A, B any](f func(E) B) func(fa Either[E, A]) Either[B, A] {
	return F.Bind2nd(MonadMapLeft[E, A, B], f)
}

func MonadChain[E, A, B any](fa Either[E, A], f func(a A) Either[E, B]) Either[E, B] {
	return fold(fa, Left[E, B], f)
}

func Chain[E, A, B any](f func(a A) Either[E, B]) func(Either[E, A]) Either[E, B] {
	return F.Bind2nd(MonadChain[E, A, B], f)
}

func Flatten[E, A any](mma Either[E, Either[E, A]]) Either[E, A] {
	return MonadChain(mma, F.Identity[Either[E, A]])
}

func TryCatch[E, A any](f func() (A, error), onThrow func(error) E) Either[E, A] {
	val, err := f()
	if err != nil {
		return F.Pipe2(err, onThrow, Left[E, A])
	}
	return F.Pipe1(val, Right[E, A])
}

func TryCatchError[A any](f func() (A, error)) Either[error, A] {
	return TryCatch(f, F.Identity[error])
}

func MonadSequence2[E, T1, T2, R any](e1 Either[E, T1], e2 Either[E, T2], f func(T1, T2) Either[E, R]) Either[E, R] {
	if IsLeft(e1) {
		return Left[E, R](e1.(left[E]).e)
	}
	if IsLeft(e2) {
		return Left[E, R](e2.(left[E]).e)
	}
	return f(e1.(right[T1]).a, e2.(right[T2]).a)
}

func MonadSequence3[E, T1, T2, T3, R any](e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3], f func(T1, T2, T3) Either[E, R]) Either[E, R] {
	if IsLeft(e1) {
		return Left[E, R](e1.(left[E]).e)
	}
	if IsLeft(e2) {
		return Left[E, R](e2.(left[E]).e)
	}
	if IsLeft(e3) {
		return Left[E, R](e3.(left[E]).e)
	}
	return f(e1.(right[T1]).a, e2.(right[T2]).a, e3.(right[T3]).a)
}

func Sequence2[E, T1, T2, R any](f func(T1, T2) Either[E, R]) func(Either[E, T1], Either[E, T2]) Either[E, R] {
	return func(e1 Either[E, T1], e2 Either[E, T2]) Either[E, R] {
		return MonadSequence2(e1, e2, f)
	}
}

func Sequence3[E, T1, T2, T3, R any](f func(T1, T2, T3) Either[E, R]) func(Either[E, T1], Either[E, T2], Either[E, T3]) Either[E, R] {
	return func(e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3]) Either[E, R] {
		return MonadSequence3(e1, e2, e3, f)
	}
}

func FromOption[E, A any](onNone func() E) func(O.Option[A]) Either[E, A] {
	return O.Fold(func() Either[E, A] { return Left[E, A](onNone()) }, Right[E, A])
}

func FromNillable[E, A any](e E) func(*A) Either[E, *A] {
	return FromPredicate(F.IsNonNil[A], F.Constant1[*A](e))
}

func Eitherize0[R any](f func() (R, error)) func() Either[error, R] {
	return func() Either[error, R] {
		return TryCatchError(f)
	}
}

func Eitherize1[T1, R any](f func(t1 T1) (R, error)) func(t1 T1) Either[error, R] {
	return func(t1 T1) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1)
		})
	}
}

func Eitherize2[T1, T2, R any](f func(t1 T1, t2 T2) (R, error)) func(t1 T1, t2 T2) Either[error, R] {
	return func(t1 T1, t2 T2) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, t2)
		})
	}
}

func fold[E, A, B any](ma Either[E, A], onLeft func(e E) B, onRight func(a A) B) B {
	if IsLeft(ma) {
		return onLeft(ma.(left[E]).e)
	}
	return onRight(ma.(right[A]).a)
}

func Fold[E, A, B any](onLeft func(e E) B, onRight func(a A) B) func(ma Either[E, A]) B {
	return func(ma Either[E, A]) B {
		return fold(ma, onLeft, onRight)
	}
}

func FromPredicate[E, A any](pred func(a A) bool, onFalse func(a A) E) func(value A) Either[E, A] {
	return func(a A) Either[E, A] {
		if pred(a) {
			return Right[E](a)
		}
		return Left[E, A](onFalse(a))
	}
}

func GetOrElse[E, A any](onLeft func(E) A) func(Either[E, A]) A {
	return Fold(onLeft, F.Identity[A])
}

func AltW[E, E1, A any](that func() Either[E1, A]) func(Either[E, A]) Either[E1, A] {
	return Fold(F.Ignore1[E](that), Right[E1, A])
}

func Alt[E, A any](that func() Either[E, A]) func(Either[E, A]) Either[E, A] {
	return AltW[E](that)
}

func ToError[A any](e Either[error, A]) error {
	return fold(e, F.Identity[error], F.Constant1[A, error](nil))
}

func MonadChainOptionK[E, A, B any](onNone func() E, ma Either[E, A], f func(A) O.Option[B]) Either[E, B] {
	return MonadChain(ma, F.Flow2(f, FromOption[E, B](onNone)))
}

func ChainOptionK[E, A, B any](onNone func() E) func(func(A) O.Option[B]) func(Either[E, A]) Either[E, B] {
	from := FromOption[E, B](onNone)
	return func(f func(A) O.Option[B]) func(Either[E, A]) Either[E, B] {
		return Chain(F.Flow2(f, from))
	}
}
