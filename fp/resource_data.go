//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package fp

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-hpcr/common"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
	O "github.com/terraform-provider-hpcr/fp/option"
)

type ResourceData interface {
	GetOk(string) (any, bool)
	SetId(string)
	Set(key string, value any) error
}

type resourceDataProxy struct {
	delegate *schema.ResourceData
}

func (proxy resourceDataProxy) GetOk(key string) (any, bool) {
	return proxy.delegate.GetOk(key)
}

func (proxy resourceDataProxy) SetId(value string) {
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
		return E.Left[error, A](fmt.Errorf("key [%s] has not been declared", key))
	}
}

func ResourceDataSet[A any](key string) func(A) func(ResourceData) E.Either[error, ResourceData] {
	return func(value A) func(ResourceData) E.Either[error, ResourceData] {
		return func(d ResourceData) E.Either[error, ResourceData] {
			if err := d.Set(key, value); err != nil {
				return E.Left[error, ResourceData](err)
			}
			return E.Of[error](d)
		}
	}
}

func ResourceDataAp[A any](d ResourceData) func(E.Either[error, func(ResourceData) E.Either[error, A]]) E.Either[error, A] {
	return E.Chain(I.Ap[ResourceData, E.Either[error, A]](d))
}
