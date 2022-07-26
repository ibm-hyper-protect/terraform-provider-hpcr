// Copyright 2022 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package function

// what a mess, golang does not have proper types ...

func Pipe1[A, R any](a A, f1 func(a A) R) R {
	return f1(a)
}

func Pipe2[A, T1, R any](a A, f1 func(a A) T1, f2 func(t1 T1) R) R {
	return Pipe1(Pipe1(a, f1), f2)
}

func Pipe3[A, T1, T2, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) R) R {
	return Pipe1(Pipe2(a, f1, f2), f3)
}

func Pipe4[A, T1, T2, T3, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) R) R {
	return Pipe1(Pipe3(a, f1, f2, f3), f4)
}

func Pipe5[A, T1, T2, T3, T4, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) R) R {
	return Pipe1(Pipe4(a, f1, f2, f3, f4), f5)
}

func Pipe6[A, T1, T2, T3, T4, T5, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) R) R {
	return Pipe1(Pipe5(a, f1, f2, f3, f4, f5), f6)
}

func Pipe7[A, T1, T2, T3, T4, T5, T6, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) R) R {
	return Pipe1(Pipe6(a, f1, f2, f3, f4, f5, f6), f7)
}

func Pipe8[A, T1, T2, T3, T4, T5, T6, T7, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) T7, f8 func(t7 T7) R) R {
	return Pipe1(Pipe7(a, f1, f2, f3, f4, f5, f6, f7), f8)
}

func Pipe9[A, T1, T2, T3, T4, T5, T6, T7, T8, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) T7, f8 func(t7 T7) T8, f9 func(t8 T8) R) R {
	return Pipe1(Pipe8(a, f1, f2, f3, f4, f5, f6, f7, f8), f9)
}

func Pipe10[A, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) T7, f8 func(t7 T7) T8, f9 func(t8 T8) T9, f10 func(t9 T9) R) R {
	return Pipe1(Pipe9(a, f1, f2, f3, f4, f5, f6, f7, f8, f9), f10)
}

func Flow2[A, T1, R any](f1 func(a A) T1, f2 func(t1 T1) R) func(a A) R {
	return func(a A) R {
		return Pipe2(a, f1, f2)
	}
}

func Flow3[A, T1, T2, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) R) func(a A) R {
	return func(a A) R {
		return Pipe3(a, f1, f2, f3)
	}
}

func Flow4[A, T1, T2, T3, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) R) func(a A) R {
	return func(a A) R {
		return Pipe4(a, f1, f2, f3, f4)
	}
}

func Flow5[A, T1, T2, T3, T4, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) R) func(a A) R {
	return func(a A) R {
		return Pipe5(a, f1, f2, f3, f4, f5)
	}
}

func Flow6[A, T1, T2, T3, T4, T5, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) R) func(a A) R {
	return func(a A) R {
		return Pipe6(a, f1, f2, f3, f4, f5, f6)
	}
}

func Flow7[A, T1, T2, T3, T4, T5, T6, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) R) func(a A) R {
	return func(a A) R {
		return Pipe7(a, f1, f2, f3, f4, f5, f6, f7)
	}
}

func Flow8[A, T1, T2, T3, T4, T5, T6, T7, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) T7, f8 func(t7 T7) R) func(a A) R {
	return func(a A) R {
		return Pipe8(a, f1, f2, f3, f4, f5, f6, f7, f8)
	}
}

func Flow9[A, T1, T2, T3, T4, T5, T6, T7, T8, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) T7, f8 func(t7 T7) T8, f9 func(t8 T8) R) func(a A) R {
	return func(a A) R {
		return Pipe9(a, f1, f2, f3, f4, f5, f6, f7, f8, f9)
	}
}

func Identity[A any](a A) A {
	return a
}

func Constant[A any](a A) func() A {
	return func() A {
		return a
	}
}

func Constant1[B, A any](a A) func(B) A {
	return func(_ B) A {
		return a
	}
}
