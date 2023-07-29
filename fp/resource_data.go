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
package fp

import (
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	O "github.com/IBM/fp-go/option"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

type ResourceData interface {
	GetOk(string) (any, bool)
	SetID(string)
	Set(key string, value any) error
}

type resourceDataProxy struct {
	delegate *schema.ResourceData
}

func (proxy resourceDataProxy) GetOk(key string) (any, bool) {
	return proxy.delegate.GetOk(key)
}

func (proxy resourceDataProxy) SetID(value string) {
	proxy.delegate.SetId(value)
}

func (proxy resourceDataProxy) Set(key string, value any) error {
	return proxy.delegate.Set(key, value)
}

func CreateResourceDataProxy(delegate *schema.ResourceData) ResourceData {
	return resourceDataProxy{delegate: delegate}
}

func typeError() error {
	return fmt.Errorf("invalid type")
}

func ToTypeE[A any](data any) E.Either[error, A] {
	return F.Pipe2(
		data,
		common.ToTypeO[A],
		E.FromOption[error, A](typeError),
	)
}

func ResourceDataGetO[A any](key string) func(ResourceData) O.Option[A] {
	return func(d ResourceData) O.Option[A] {
		return F.Pipe2(
			key,
			O.FromValidation(d.GetOk),
			O.Chain(common.ToTypeO[A]),
		)
	}
}

func ResourceDataGetE[A any](key string) func(ResourceData) E.Either[error, A] {
	return func(d ResourceData) E.Either[error, A] {
		data, ok := d.GetOk(key)
		if ok {
			return ToTypeE[A](data)
		}
		return E.Left[A](fmt.Errorf("key [%s] has not been declared", key))
	}
}

func ResourceDataSet[A any](key string) func(A) func(ResourceData) E.Either[error, ResourceData] {
	return func(value A) func(ResourceData) E.Either[error, ResourceData] {
		return func(d ResourceData) E.Either[error, ResourceData] {
			if err := d.Set(key, value); err != nil {
				return E.Left[ResourceData](err)
			}
			return E.Of[error](d)
		}
	}
}

func ResourceDataAp[A any](d ResourceData) func(E.Either[error, func(ResourceData) E.Either[error, A]]) E.Either[error, A] {
	return E.Chain(I.Ap[E.Either[error, A]](d))
}
