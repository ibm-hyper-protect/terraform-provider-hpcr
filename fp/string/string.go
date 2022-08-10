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
package string

import (
	"strings"

	F "github.com/terraform-provider-hpcr/fp/function"
)

func ToBytes(a string) []byte {
	return []byte(a)
}

func Equals(a string, b string) bool {
	return a == b
}

// Includes returns a predicate that tests for the existence of the search string
func Includes(searchString string) func(s string) bool {
	return F.Bind2nd(strings.Contains, searchString)
}
