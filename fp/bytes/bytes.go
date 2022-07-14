package bytes

import "fmt"

func ToString(a []byte) string {
	return string(a)
}

func ToHexString(a []byte) string {
	return fmt.Sprintf("%x", a)
}

func Slice(start int, end int) func([]byte) []byte {
	return func(a []byte) []byte {
		return a[start:end]
	}
}
