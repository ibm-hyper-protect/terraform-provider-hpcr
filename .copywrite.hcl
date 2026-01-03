schema_version = 1

project {
  license        = "Apache-2.0"
  copyright_year = 2022

  # (OPTIONAL) A list of globs that should not have copyright/license headers.
  # Supports doublestar glob patterns for more flexibility in defining which
  # files or folders should be ignored
  header_ignore = [
    # examples used within documentation
    "examples/**",
    # GitHub issue config
    ".github/ISSUE_TEMPLATE/*.yml",
    # golangci-lint tooling configuration
    ".golangci.yml",
    # GoReleaser tooling configuration
    ".goreleaser.yml",
    # Semantic release configuration
    ".releaserc",
    # Dependency management
    "renovate.json",
    # Security scanning
    ".whitesource",
    # Test fixtures
    "**/testdata/**",
    # Generated files
    "**/*_generated.go",
    # Vendor directory
    "vendor/**",
  ]
}
