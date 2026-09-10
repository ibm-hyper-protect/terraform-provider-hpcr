# AGENTS.md — Plan Mode

This file provides guidance for planning changes to this repository.

## Architectural Constraints

- **All resources are stateless and compute-only** — they generate encrypted/signed artifacts locally. There is no remote state reconciliation in `Read`. Planning changes that add remote API calls to `Read` would break the design pattern.
- **`contract-go/v2` is the single source of truth for HPCR operations** — the Terraform provider is a thin wrapper. Any new encryption, signing, or image-selection logic must be added to `contract-go/v2` first, then consumed here.
- **`tools/` is a standalone Go module** — it cannot import from the root module. Any shared generation logic must stay self-contained in `tools/`.
- **The `common/` package is deliberately minimal** — it exists only for: UUID generation, RSA key generation (via openssl subprocess), file reading, checksum map parsing, and the YAML contract refinement transform. Do not add business logic here.
- **Adding or removing resources/datasources requires updating hardcoded expected counts** in `provider_test.go` (`TestHPCRProvider_Resources` expects exactly 8, `TestHPCRProvider_DataSources` expects exactly 4). This is an intentional guard against accidental omissions.
- **Release pipeline is sequential and gated**: unit tests must pass before release; goreleaser only runs if semantic-release creates a new tag. Breaking the release pipeline requires fixing `.releaserc`, the workflow in `.github/workflows/build.yml`, and the GPG signing setup in goreleaser.
- **CGO is disabled** in goreleaser builds — any dependency that requires CGO cannot be added to this provider.
- **Commit message format gates releases** — CI enforces Conventional Commits on PRs. A malformed commit will block the PR. Breaking changes require `!` suffix or `BREAKING CHANGE:` footer to trigger a major version bump.
- **`depguard` bans `terraform-plugin-sdk/v2`** — any plan to migrate or add SDK v2 code will fail linting.
