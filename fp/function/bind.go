//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package function

func Bind1st[T1, T2, R any](f func(T1, T2) R, t1 T1) func(T2) R {
	return func(t2 T2) R {
		return f(t1, t2)
	}
}
func Bind2nd[T1, T2, R any](f func(T1, T2) R, t2 T2) func(T1) R {
	return func(t1 T1) R {
		return f(t1, t2)
	}
}

func Bind1[T1, R any](f func(T1) R, t1 T1) func() R {
	return func() R {
		return f(t1)
	}
}

func SK[T1, T2 any](_ T1, t2 T2) T2 {
	return t2
}

func Ignore1[T, R any](f func() R) func(T) R {
	return func(_ T) R {
		return f()
	}
}

func Ignore1st[T1, T2, R any](f func(T2) R) func(T1, T2) R {
	return func(_ T1, t2 T2) R {
		return f(t2)
	}
}

func Ignore2nd[T1, T2, R any](f func(T1) R) func(T1, T2) R {
	return func(t1 T1, _ T2) R {
		return f(t1)
	}
}
