// Copyright 2023 IBM Corp.
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
	"fmt"
	"sort"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	O "github.com/IBM/fp-go/option"
	R "github.com/IBM/fp-go/record"
	T "github.com/IBM/fp-go/tuple"
	"github.com/Masterminds/semver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
)

var (
	schemaCertCertificatesIn = schema.Schema{
		Type: schema.TypeMap,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "Map of certificates from version to certificate",
		Required:    true,
	}
	schemaCertSpecIn = schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Default:          "*",
		Description:      "Semantic version range defining the HPCR certificate",
		ValidateDiagFunc: validateSpecFunc,
	}
	schemaCertVersionOut = schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Version number of the selected certificate",
	}
	schemaCertCertificateOut = schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Selected certificate",
	}
	// getters
	getCertSpecE         = fp.ResourceDataGetE[string](common.KeySpec)
	getCertCertificatesE = fp.ResourceDataGetE[map[string]any](common.KeyCerts)

	setCertCertificate = fp.ResourceDataSet[string](common.KeyCert)
	setCertVersion     = fp.ResourceDataSet[string](common.KeyVersion)
)

// toCertificateMap performs a type conversion of the certificates
func toCertificateMap(data map[string]any) E.Either[error, map[string]string] {
	return E.TryCatchError(func() (map[string]string, error) {
		var dst = make(map[string]string, len(data))
		for k, v := range data {
			s, ok := v.(string)
			if !ok {
				return dst, fmt.Errorf("invalid type")
			}
			dst[k] = s
		}
		return dst, nil
	}())
}

// parses a constraint string into a constraint object
var parseConstraint = E.Eitherize1(semver.NewConstraint)

func checkCertContraintPredicate(cstr *semver.Constraints) func(T.Tuple2[*semver.Version, string]) bool {
	return func(t T.Tuple2[*semver.Version, string]) bool {
		return cstr.Check(t.F1)
	}
}

// sortCertByVersion sorts the structs by version in inverse order, i.e. the first item will be the latest version
func sortCertByVersion(imgs []T.Tuple2[*semver.Version, string]) []T.Tuple2[*semver.Version, string] {
	sort.SliceStable(imgs, func(left, right int) bool {
		return imgs[left].F1.Compare(imgs[right].F1) > 0
	})
	return imgs
}

// selectCertBySpec selects the latest version that matches the specification
func selectCertBySpec(spec *semver.Constraints) func(img []T.Tuple2[*semver.Version, string]) O.Option[T.Tuple2[*semver.Version, string]] {
	return F.Flow3(
		A.Filter(checkCertContraintPredicate(spec)),
		I.Map(sortCertByVersion),
		A.Head[T.Tuple2[*semver.Version, string]],
	)
}

// parseEntry parses the version part of an entry
func parseEntry(entry T.Tuple2[string, string]) E.Either[error, T.Tuple2[*semver.Version, string]] {
	return F.Pipe2(
		entry.F1,
		parseVersion,
		E.Map[error](F.Bind2nd(T.MakeTuple2[*semver.Version, string], entry.F2)),
	)
}

func versionToString(version *semver.Version) string {
	return version.String()
}

func handleCertificateWithContext(ctx *Context) func(data fp.ResourceData) ResourceDataE {
	return func(data fp.ResourceData) ResourceDataE {
		// some applicatives
		apE := fp.ResourceDataAp[fp.ResourceData](data)
		// final result
		resE := E.MapTo[error, []fp.ResourceData](data)
		// access the certificates and convert to tuple
		certsE := F.Pipe4(
			data,
			getCertCertificatesE,
			E.Chain(toCertificateMap),
			E.Map[error](R.ToEntries[string, string]),
			E.Chain(E.TraverseArray(parseEntry)),
		)
		// access the spec and parse it
		selectedE := F.Pipe5(
			data,
			getCertSpecE,
			E.Chain(parseConstraint),
			E.Map[error](selectCertBySpec),
			E.Ap[O.Option[T.Tuple2[*semver.Version, string]]](certsE),
			E.Chain(E.FromOption[T.Tuple2[*semver.Version, string]](func() error { return fmt.Errorf("unable to select a version") })),
		)
		// output records
		certificateMapE := F.Pipe2(
			selectedE,
			E.Map[error](F.Flow2(
				T.Second[*semver.Version, string],
				setCertCertificate,
			)),
			apE,
		)
		versionMapE := F.Pipe2(
			selectedE,
			E.Map[error](F.Flow3(
				T.First[*semver.Version, string],
				versionToString,
				setCertVersion,
			)),
			apE,
		)

		// combine all outputs
		return F.Pipe2(
			[]ResourceDataE{certificateMapE, versionMapE},
			seqResourceData,
			resE,
		)

	}
}

// handleCertificate is the data source callback, it selects the version
func handleCertificate(data *schema.ResourceData, ctx any) error {
	// lift f into the context
	return F.Pipe5(
		ctx,
		toContextE,
		E.Map[error](handleCertificateWithContext),
		E.Ap[ResourceDataE](F.Pipe2(
			data,
			fp.CreateResourceDataProxy,
			setUniqueID,
		)),
		E.Flatten[error, fp.ResourceData],
		E.ToError[fp.ResourceData],
	)
}

// DatasourceEncryptionCertificate is a data source to select a certificate from a map, where the key is the version and the value the certificate
func DatasourceEncryptionCertificate() *schema.Resource {
	return &schema.Resource{
		Read: handleCertificate,
		Schema: map[string]*schema.Schema{
			// input parameters
			common.KeyCerts: &schemaCertCertificatesIn,
			common.KeySpec:  &schemaCertSpecIn,
			// output parameters
			common.KeyCert:    &schemaCertCertificateOut,
			common.KeyVersion: &schemaCertVersionOut,
		},
		Description: "Selects the best matching certificate based on the semantic version.",
	}
}
