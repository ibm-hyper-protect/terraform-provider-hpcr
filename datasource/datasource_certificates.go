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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	B "github.com/IBM/fp-go/bytes"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	R "github.com/IBM/fp-go/record"
	T "github.com/IBM/fp-go/tuple"
	"github.com/Masterminds/semver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
)

const (
	// default template used to download certificates
	defaultTemplate = "https://cloud.ibm.com/media/docs/downloads/hyper-protect-container-runtime/ibm-hyper-protect-container-runtime-{{.Major}}-{{.Minor}}-s390x-{{.Patch}}-encrypt.crt"

	// template key
	KeyMajor = "Major"
	KeyMinor = "Minor"
	KeyPatch = "Patch"
)

var (
	schemaTemplateIn = schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     defaultTemplate,
		Description: "Template used to download the encryption certificate, it may contain the placeholders for {{.Major}}, {{.Minor}} and {{.Patch}} as replacement tokens",
	}
	schemaVersionsIn = schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "List of strings, each denoting the version number of the certificate to download",
		Required:    true,
	}
	schemaCertificatesOut = schema.Schema{
		Type: schema.TypeMap,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "Map of certificates from version to certificate",
		Computed:    true,
	}
	getTemplateE = fp.ResourceDataGetE[string](common.KeyTemplate)
	getVersionsE = fp.ResourceDataGetE[[]any](common.KeyVersions)

	setCertificates = fp.ResourceDataSet[map[string]string](common.KeyCerts)
)

// parses a version string into a version object
var parseVersion = E.Eitherize1(semver.NewVersion)

// parseTemplate parses a template of the given name
func parseTemplate(name string) func(string) E.Either[error, *template.Template] {
	return E.Eitherize1(template.New(name).Parse)
}

// DatasourceEncryptionCertificates is a data source to download encryption certificates from an official location
func DatasourceEncryptionCertificates() *schema.Resource {
	return &schema.Resource{
		Read: handleDownload,
		Schema: map[string]*schema.Schema{
			// input parameters
			common.KeyTemplate: &schemaTemplateIn,
			common.KeyVersions: &schemaVersionsIn,
			// output parameters
			common.KeyCerts: &schemaCertificatesOut,
		},
		Description: "Downloads the encryption certificates for the given version numbers.",
	}
}

// resolveUrl computes the download URL from a version number
func resolveUrl(tmp *template.Template) func(version *semver.Version) E.Either[error, string] {
	return func(version *semver.Version) E.Either[error, string] {
		ctx := map[string]int64{
			KeyMajor: version.Major(),
			KeyMinor: version.Minor(),
			KeyPatch: version.Patch(),
		}
		var buffer bytes.Buffer
		err := tmp.Execute(&buffer, ctx)
		return E.TryCatchError(buffer.String(), err)
	}
}

func textFromResponse(url string) func(resp *http.Response) E.Either[error, string] {
	return func(resp *http.Response) E.Either[error, string] {
		defer resp.Body.Close()

		return F.Pipe1(
			E.TryCatchError(func() ([]byte, error) {
				data, err := io.ReadAll(resp.Body)
				if err != nil {
					return data, err
				}
				if resp.StatusCode != http.StatusOK {
					return data, fmt.Errorf("url: %s, status %d: cause: [%s]", url, resp.StatusCode, data)
				}
				return data, err
			}()),
			E.Map[error](B.ToString),
		)
	}
}

func downloadTextFromUrl(client *http.Client) func(url string) E.Either[error, string] {
	downloadE := E.Eitherize1(client.Get)
	return func(url string) E.Either[error, string] {
		return F.Pipe2(
			url,
			downloadE,
			E.Chain(textFromResponse(url)),
		)
	}
}

func downloadSingleVersion(client *http.Client) func(resolver func(version *semver.Version) E.Either[error, string]) func(version *semver.Version) E.Either[error, T.Tuple2[string, string]] {
	downloadE := downloadTextFromUrl(client)
	return func(resolver func(version *semver.Version) E.Either[error, string]) func(version *semver.Version) E.Either[error, T.Tuple2[string, string]] {
		return func(version *semver.Version) E.Either[error, T.Tuple2[string, string]] {
			return F.Pipe3(
				version,
				resolver,
				E.Chain(downloadE),
				E.Map[error](F.Bind1st(T.MakeTuple2[string, string], version.String())),
			)
		}
	}
}

func handleDownloadWithContext(ctx *Context) func(data fp.ResourceData) ResourceDataE {
	downloadE := downloadSingleVersion(ctx.client)
	return func(data fp.ResourceData) ResourceDataE {
		// some applicatives
		apE := fp.ResourceDataAp[fp.ResourceData](data)
		// final result
		resE := E.MapTo[error, []fp.ResourceData](data)
		// the list of version numbers
		versionsE := F.Pipe3(
			data,
			getVersionsE,
			E.Chain(E.TraverseArray(fp.ToTypeE[string])),
			E.Chain(E.TraverseArray(parseVersion)),
		)
		// resolve and download
		certificatesE := F.Pipe7(
			data,
			getTemplateE,
			E.Map[error](strings.TrimSpace),
			E.Chain(parseTemplate("downloadURL")),
			E.Map[error](F.Flow3(
				resolveUrl,
				downloadE,
				E.TraverseArray[error, *semver.Version, T.Tuple2[string, string]]),
			),
			E.Ap[E.Either[error, []T.Tuple2[string, string]]](versionsE),
			E.Flatten[error, []T.Tuple2[string, string]],
			E.Map[error](R.FromEntries[string, string]),
		)
		// output records
		certificatesMapE := F.Pipe2(
			certificatesE,
			E.Map[error](setCertificates),
			apE,
		)

		// combine all outputs
		return F.Pipe2(
			[]ResourceDataE{certificatesMapE},
			seqResourceData,
			resE,
		)

	}
}

// handleDownload is the data source callback, it downloads the versions
func handleDownload(data *schema.ResourceData, ctx any) error {
	// lift f into the context
	return F.Pipe5(
		ctx,
		toContextE,
		E.Map[error](handleDownloadWithContext),
		E.Ap[ResourceDataE](F.Pipe2(
			data,
			fp.CreateResourceDataProxy,
			setUniqueID,
		)),
		E.Flatten[error, fp.ResourceData],
		E.ToError[fp.ResourceData],
	)
}
