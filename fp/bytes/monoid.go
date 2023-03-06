// Copyright 2022 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package bytes

import (
	A "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/array"
)

// monoid for byte arrays
var Monoid = A.Monoid[byte]()

// ConcatAll concatenates all bytes
var ConcatAll = A.ConcatAll[byte]
