package array

import (
	F "github.com/terraform-provider-hpcr/fp/function"
	O "github.com/terraform-provider-hpcr/fp/option"
)

func MakeBy[A any](n int, f func(int) A) []A {
	as := make([]A, n)
	for i := n - 1; i >= 0; i-- {
		as[i] = f(i)
	}
	return as
}

func MonadMap[A, B any](as []A, f func(a A) B) []B {
	count := len(as)
	bs := make([]B, count)
	for i := count - 1; i >= 0; i-- {
		bs[i] = f(as[i])
	}
	return bs
}

func MonadMapRef[A, B any](as []A, f func(a *A) B) []B {
	count := len(as)
	bs := make([]B, count)
	for i := count - 1; i >= 0; i-- {
		bs[i] = f(&as[i])
	}
	return bs
}

func Map[A, B any](f func(a A) B) func([]A) []B {
	return F.Bind2nd(MonadMap[A, B], f)
}

func MapRef[A, B any](f func(a *A) B) func([]A) []B {
	return F.Bind2nd(MonadMapRef[A, B], f)
}

func filter[A any](fa []A, pred func(a A) bool) []A {
	var result []A
	count := len(fa)
	for i := 0; i < count; i++ {
		a := fa[i]
		if pred(a) {
			result = append(result, a)
		}
	}
	return result
}

func filterRef[A any](fa []A, pred func(a *A) bool) []A {
	var result []A
	count := len(fa)
	for i := 0; i < count; i++ {
		a := fa[i]
		if pred(&a) {
			result = append(result, a)
		}
	}
	return result
}

func filterMap[A, B any](fa []A, pred func(a A) bool, f func(a A) B) []B {
	var result []B
	count := len(fa)
	for i := 0; i < count; i++ {
		a := fa[i]
		if pred(a) {
			result = append(result, f(a))
		}
	}
	return result
}

func filterMapRef[A, B any](fa []A, pred func(a *A) bool, f func(a *A) B) []B {
	var result []B
	count := len(fa)
	for i := 0; i < count; i++ {
		a := fa[i]
		if pred(&a) {
			result = append(result, f(&a))
		}
	}
	return result
}

func Filter[A any](pred func(a A) bool) func([]A) []A {
	return F.Bind2nd(filter[A], pred)
}

func FilterRef[A any](pred func(a *A) bool) func([]A) []A {
	return F.Bind2nd(filterRef[A], pred)
}

func FilterMap[A, B any](pred func(a A) bool, f func(a A) B) func([]A) []B {
	return func(fa []A) []B {
		return filterMap(fa, pred, f)
	}
}

func FilterMapRef[A, B any](pred func(a *A) bool, f func(a *A) B) func([]A) []B {
	return func(fa []A) []B {
		return filterMapRef(fa, pred, f)
	}
}

func reduce[A, B any](fa []A, f func(B, A) B, initial B) B {
	current := initial
	count := len(fa)
	for i := 0; i < count; i++ {
		current = f(current, fa[i])
	}
	return current
}

func reduceRef[A, B any](fa []A, f func(B, *A) B, initial B) B {
	current := initial
	count := len(fa)
	for i := 0; i < count; i++ {
		current = f(current, &fa[i])
	}
	return current
}

func Reduce[A, B any](f func(B, A) B, initial B) func([]A) B {
	return func(as []A) B {
		return reduce(as, f, initial)
	}
}

func ReduceRef[A, B any](f func(B, *A) B, initial B) func([]A) B {
	return func(as []A) B {
		return reduceRef(as, f, initial)
	}
}

func Append[A any](as []A, a A) []A {
	return append(as, a)
}

func IsEmpty[A any](as []A) bool {
	return len(as) == 0
}

func IsNonEmpty[A any](as []A) bool {
	return len(as) > 0
}

func Empty[A any]() []A {
	return make([]A, 0)
}

func Zero[A any]() []A {
	return Empty[A]()
}

func Of[A any](a A) []A {
	return []A{a}
}

func MonadChain[A, B any](fa []A, f func(a A) []B) []B {
	return reduce(fa, func(bs []B, a A) []B {
		return append(bs, f(a)...)
	}, Zero[B]())
}

func Chain[A, B any](f func(a A) []B) func([]A) []B {
	return F.Bind2nd(MonadChain[A, B], f)
}

func MonadAp[A, B any](fab []func(A) B, fa []A) []B {
	return MonadChain(fab, F.Bind1st(MonadMap[A, B], fa))
}

func Ap[A, B any](fa []A) func([]func(A) B) []B {
	return F.Bind2nd(MonadAp[A, B], fa)
}

func Match[A, B any](onEmpty func() B, onNonEmpty func([]A) B) func([]A) B {
	return func(as []A) B {
		if IsEmpty(as) {
			return onEmpty()
		}
		return onNonEmpty(as)
	}
}

func Tail[A any](as []A) O.Option[[]A] {
	if IsEmpty(as) {
		return O.None[[]A]()
	}
	return O.Some(as[1:])
}

func Head[A any](as []A) O.Option[A] {
	if IsEmpty(as) {
		return O.None[A]()
	}
	return O.Some(as[0])
}

func PrependAll[A any](middle A) func([]A) []A {
	return func(as []A) []A {
		count := len(as)
		dst := count * 2
		result := make([]A, dst)
		for i := count - 1; i >= 0; i-- {
			dst--
			result[dst] = as[i]
			dst--
			result[dst] = middle
		}
		return result
	}
}

func Intersperse[A any](middle A) func([]A) []A {
	prepend := PrependAll(middle)
	return func(as []A) []A {
		if IsEmpty(as) {
			return as
		}
		return prepend(as)[1:]
	}
}

func Flatten[A any](mma [][]A) []A {
	return MonadChain(mma, F.Identity[[]A])
}
