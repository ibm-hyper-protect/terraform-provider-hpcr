package identity

import (
	F "github.com/terraform-provider-hpcr/fp/function"
)

func MonadAp[A, B any](fab func(A) B, fa A) B {
	return fab(fa)
}

func Ap[A, B any](fa A) func(func(A) B) B {
	return F.Bind2nd(MonadAp[A, B], fa)
}

func MonadMap[A, B any](fa A, f func(A) B) B {
	return f(fa)
}

func Map[A, B any](f func(A) B) func(A) B {
	return f
}

func Of[A any](a A) A {
	return a
}

func MonadChain[A, B any](ma A, f func(A) B) B {
	return f(ma)
}

func Chain[A, B any](f func(A) B) func(A) B {
	return f
}
