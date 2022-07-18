package predicate

import (
	F "github.com/terraform-provider-hpcr/fp/function"
)

func ContraMap[A, B any](f func(B) A) func(func(A) bool) func(B) bool {
	return func(pred func(A) bool) func(B) bool {
		return F.Flow2(f, pred)
	}
}

func Not[A any](predicate func(A) bool) func(A) bool {
	return func(a A) bool {
		return !predicate((a))
	}
}
