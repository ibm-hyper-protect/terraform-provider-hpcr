# AGENTS.md — Agent (Coding) Mode

This file provides guidance to agents when writing or modifying code in this repository.

## Critical Coding Rules

- **Never use `terraform-plugin-sdk/v2`** — depguard will fail the build. All Terraform types and interfaces must come from `terraform-plugin-framework`.
- **All crypto/HPCR operations go through `contract-go/v2`** — never re-implement encryption, signing, TGZ, or image selection logic. The `contract-go/v2` library exports: `contract.HpcrTgz`, `contract.HpcrTextEncrypted`, `contract.HpcrContractSignedEncrypted`, `certificate.HpcrValidateEncryptionCertificate`, `certificate.HpcrGetEncryptionCertificateFromJson`, `image.HpcrSelectImage`.
- **`common.RefineContract()` is mandatory before passing YAML to `HpcrContractSignedEncrypted`** — it re-serializes `env` and `workload` keys as YAML block literals. Skipping this produces malformed contracts.
- **Private key generation requires a shell call to `openssl`** via `common.GeneratePrivateKey()`. This means `openssl` must be in PATH at test time. Tests that involve key generation will fail without it.
- **`Delete` on all resources must be a no-op** — these resources generate local computed values, not remote infrastructure.
- **`Read` must only reflect state back to Terraform** — no external calls, no recomputation.
- **Every schema attribute must have both `Description` and `MarkdownDescription`** — tests validate this.
- **Every resource/datasource file must have `var _ resource.Resource = &MyResource{}` or equivalent** — interface compliance is asserted at compile time.
- **IDs must use `stringplanmodifier.UseStateForUnknown()`** — prevents Terraform from showing a diff on existing resources.
- **Optional `types.String` fields must check `IsNull() && IsUnknown()` before calling `ValueString()`** to avoid operating on unknown/null values.

## Import Order Convention

```go
import (
    // 1. stdlib
    "context"
    "fmt"

    // 2. terraform-plugin-framework
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"

    // 3. contract-go/v2
    "github.com/ibm-hyper-protect/contract-go/v2/contract"

    // 4. local common
    "github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)
```

## Adding a New Resource

1. Create `internal/provider/resources/resource_<name>.go` and `resource_<name>_test.go`
2. Add `var _ resource.Resource = &MyResource{}` interface check
3. Register in `internal/provider/provider.go` `Resources()` slice
4. Update `expectedCount` in `provider_test.go` (`TestHPCRProvider_Resources` checks exact count = 8 currently)
5. Add `examples/resources/hpcr_<name>/` with an example `.tf` file
6. Run `make generate` to regenerate docs

## Adding a New Datasource

Same pattern as above but in `datasources/` package; register in `DataSources()` in `provider.go`; update `expectedCount` in `TestHPCRProvider_DataSources` (currently 4).

## License Headers

All new `.go` files must have the Apache-2.0 header with `Copyright 2026 IBM Corp.`. Run `make generate` to auto-add via `copywrite`.
