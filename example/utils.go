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
package example

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-provider-hpcr/common"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
	O "github.com/terraform-provider-hpcr/fp/option"
	S "github.com/terraform-provider-hpcr/fp/string"
	Y "github.com/terraform-provider-hpcr/fp/yaml"
	"github.com/terraform-provider-hpcr/provider"
)

var (
	providerName      = "hpcr"
	providerFactories = map[string]func() (*schema.Provider, error){
		providerName: func() (*schema.Provider, error) { return provider.Provider(), nil },
	}
)

func getOutputO(s *terraform.State) func(string) O.Option[string] {
	return F.Flow3(
		O.FromValidation(func(name string) (*terraform.OutputState, bool) {
			ms := s.RootModule()
			rs, ok := ms.Outputs[name]
			return rs, ok
		}),
		O.Map(func(os *terraform.OutputState) any {
			return os.Value
		}),
		O.Chain(common.ToTypeO[string]),
	)
}

func TestCheckOutput(name string, check func(value string) O.Option[error]) resource.TestCheckFunc {
	return F.Flow5(
		getOutputO,
		I.Ap[string, O.Option[string]](name),
		E.FromOption[error, string](func() error {
			return fmt.Errorf("output [%s] not found", name)
		}),
		E.Fold(O.Of[error], check),
		O.GetOrElse(F.Constant[error](nil)),
	)
}

var validateUserData = F.Flow4(
	S.ToBytes,
	Y.Parse[map[string]string],
	E.Map[error](F.Deref[map[string]string]),
	E.Fold(O.Of[error], F.Constant1[map[string]string](O.None[error]())),
)
