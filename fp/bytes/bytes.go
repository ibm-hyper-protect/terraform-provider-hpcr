//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
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
