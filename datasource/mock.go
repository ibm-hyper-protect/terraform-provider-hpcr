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
// limitations under the License.package datasource

package datasource

import (
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
)

type resourceDataMock struct {
	data map[string]any
}

func (mock resourceDataMock) GetOk(key string) (any, bool) {
	value, exists := mock.data[key]
	return value, exists
}

func (mock resourceDataMock) SetID(value string) {
	// noop
}

func (mock resourceDataMock) Set(key string, value any) error {
	mock.data[key] = value
	return nil
}

func CreateResourceDataMock(data map[string]any) fp.ResourceData {
	return resourceDataMock{data: data}
}
