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
	E "github.com/terraform-provider-hpcr/fp/either"
	I "github.com/terraform-provider-hpcr/fp/identity"
)

func ToType[A any](data any) E.Either[error, A] {
	value, ok := data.(A)
	if ok {
		return E.Of[error](value)
	}
	return E.Left[error, A](fmt.Errorf("invalid type"))
}

func ResourceDataGet[A any](key string) func(*schema.ResourceData) E.Either[error, A] {
	return func(d *schema.ResourceData) E.Either[error, A] {
		data, ok := d.GetOk(key)
		if ok {
			return ToType[A](data)
		}
		return E.Left[error, A](fmt.Errorf("key [%s] has not been declared", key))
	}
}

func ResourceDataSet[A any](key string) func(A) func(*schema.ResourceData) E.Either[error, *schema.ResourceData] {
	return func(value A) func(*schema.ResourceData) E.Either[error, *schema.ResourceData] {
		return func(d *schema.ResourceData) E.Either[error, *schema.ResourceData] {
			if err := d.Set(key, value); err != nil {
				return E.Left[error, *schema.ResourceData](err)
			}
			return E.Of[error](d)
		}
	}
}

func ResourceDataAp[A any](d *schema.ResourceData) func(E.Either[error, func(*schema.ResourceData) E.Either[error, A]]) E.Either[error, A] {
	return E.Chain(I.Ap[*schema.ResourceData, E.Either[error, A]](d))
}
