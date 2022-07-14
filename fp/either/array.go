package either

import (
	RA "github.com/terraform-provider-hpcr/fp/array"
	F "github.com/terraform-provider-hpcr/fp/function"
)

func TraverseArray[E, A, B any](f func(A) Either[E, B]) func([]A) Either[E, []B] {
	return F.Pipe1(
		f,
		RA.Traverse[A, B, Either[E, A]](
			Of[E, []B],
			MonadMap[E, []B, func(B) []B],
			MonadAp[E, B, []B],
		),
	)
}

func SequenceArray[E, A any]() func([]Either[E, A]) Either[E, []A] {
	return TraverseArray(F.Identity[Either[E, A]])
}
