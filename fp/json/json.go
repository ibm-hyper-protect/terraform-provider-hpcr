package Json

import (
	"encoding/json"

	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

func Parse[A any](data []byte) E.Either[error, *A] {
	return E.TryCatch(func() (*A, error) {
		var result A
		err := json.Unmarshal(data, &result)
		return &result, err
	}, F.Identity[error])
}

func Stringify[A any](a *A) E.Either[error, []byte] {
	return E.TryCatch(func() ([]byte, error) {
		b, err := json.Marshal(a)
		return b, err
	}, F.Identity[error])

}
