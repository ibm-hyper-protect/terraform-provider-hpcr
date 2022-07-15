//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package tuple

type Tuple2[T1, T2 any] struct {
	F1 T1
	F2 T2
}

type Tuple3[T1, T2, T3 any] struct {
	F1 T1
	F2 T2
	F3 T3
}

func FirstOf2[T1, T2 any](t Tuple2[T1, T2]) T1 {
	return t.F1
}

func SecondOf2[T1, T2 any](t Tuple2[T1, T2]) T2 {
	return t.F2
}

func MakeTuple2[T1, T2 any](t1 T1, t2 T2) Tuple2[T1, T2] {
	return Tuple2[T1, T2]{F1: t1, F2: t2}
}

func MakeTuple3[T1, T2, T3 any](t1 T1, t2 T2, t3 T3) Tuple3[T1, T2, T3] {
	return Tuple3[T1, T2, T3]{F1: t1, F2: t2, F3: t3}
}
