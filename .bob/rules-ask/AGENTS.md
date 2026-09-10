# AGENTS.md — Ask Mode

This file provides guidance for answering questions about this repository.

## Non-Obvious Context

- **`tools/` is a separate Go module** (`tools/go.mod`) with its own dependencies. It's only used via `go generate` with the `generate` build tag. Running `go build ./...` from the root does NOT include it.
- **`common/` is the only local utility package** — there's no `pkg/`, `lib/`, or `util/` directory. All shared logic lives in `common/common.go`.
- **The provider has no API key or cloud credentials** — it requires no Terraform provider configuration block. All resources work locally by calling `contract-go/v2` library functions.
- **`CHANGELOG.md` is auto-generated** — never manually edited; `semantic-release` writes it on merge to `main` based on Conventional Commits.
- **`examples/` is excluded from copyright headers and linting** — see `.copywrite.hcl` and `.golangci.yml`.
- **Acceptance tests (`make testacc`) call real IBM Cloud APIs** and require `IBM_CLOUD_API_KEY`. They take ~90 minutes and may incur costs. Unit tests (`make test`) do NOT require any credentials.
- **`openssl` binary is a runtime dependency for unit tests** — `TestGeneratePrivateKey` and any test that creates a contract with a generated key will shell out to `openssl genrsa 4096`. Override path with `OPENSSL_BIN`.
- **The provider is published at `registry.terraform.io/ibm-hyper-protect/hpcr`** — this is the address hardcoded in `main.go`.
- **`docs/` is generated from resource schemas + `examples/`** using `tfplugindocs`. Do not hand-edit it.
- **Releases are fully automated**: push to `main` → semantic-release analyzes commits → creates tag → GoReleaser builds + GPG-signs + publishes to Terraform Registry.
