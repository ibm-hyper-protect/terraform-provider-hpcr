## Description

<!-- Provide a brief description of the changes in this pull request -->

## Related Issue

<!-- Link to the issue this PR addresses -->
Fixes #(issue_number)

## Type of Change

Please select the relevant option(s):

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring (no functional changes)
- [ ] Performance improvement
- [ ] Test coverage improvement
- [ ] Build/CI configuration change

## What Changed

<!-- Describe what changed and why. Include any context that would help reviewers understand your changes. -->

### Resources/Data Sources Affected

<!-- List any resources or data sources that were modified or added -->

- [ ] `hpcr_tgz`
- [ ] `hpcr_tgz_encrypted`
- [ ] `hpcr_text`
- [ ] `hpcr_text_encrypted`
- [ ] `hpcr_json`
- [ ] `hpcr_json_encrypted`
- [ ] `hpcr_contract_encrypted`
- [ ] `hpcr_contract_encrypted_contract_expiry`
- [ ] `hpcr_image` (data source)
- [ ] `hpcr_attestation` (data source)
- [ ] `hpcr_encryption_certs` (data source)
- [ ] `hpcr_encryption_cert` (data source)
- [ ] Provider configuration
- [ ] Other: <!-- specify -->

## Testing

<!-- Describe the tests you ran and how to reproduce them -->

### Test Configuration

```hcl
# Provide a minimal Terraform configuration that demonstrates the changes
terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = ">= 1.2.0"
    }
  }
}

# Your test configuration here
```

### Test Results

<!-- Paste the output of `make test` or `make testacc` -->

```
# Test output here
```

### Acceptance Tests

<!-- Only check this if your PR changes data source logic, encryption functionality, or IBM Cloud integrations -->

- [ ] **This PR requires acceptance tests** (changes to data sources, encryption, or IBM Cloud interactions)
  - [ ] I have added the `run-acceptance-tests` label to trigger CI acceptance tests
  - [ ] OR I have run `make testacc` locally and all tests pass

**Note**: Acceptance tests interact with real IBM Cloud infrastructure and cost money. Only use the `run-acceptance-tests` label when necessary.

## Checklist

Before submitting this PR, please make sure:

- [ ] My code follows the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) style guidelines
- [ ] I have run `make fmt` and `make lint` with no errors
- [ ] I have run `make tidy` to update go.mod and go.sum
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
  - [ ] Updated resource/data source documentation (if applicable)
  - [ ] Updated README.md (if applicable)
  - [ ] Added/updated examples in `examples/` directory (if applicable)
- [ ] My changes generate no new warnings or errors
- [ ] I have added tests that prove my fix is effective or that my feature works
  - [ ] Unit tests added/updated
  - [ ] Acceptance tests added/updated (if applicable)
- [ ] All new and existing unit tests pass locally (`make test`)
- [ ] All new and existing acceptance tests pass (if run with `make testacc`)
- [ ] I have verified my changes work with:
  - [ ] Latest Terraform version
  - [ ] Minimum supported Terraform version (0.13)
- [ ] My commit messages follow [Conventional Commits](https://www.conventionalcommits.org/) format
  - Examples: `feat: add new resource`, `fix: resolve encryption issue`, `docs: update README`

## Breaking Changes

<!-- If this is a breaking change, describe the impact and migration path for users -->

- [ ] This PR introduces breaking changes
- [ ] Migration guide included (if breaking changes)

## Screenshots (if applicable)

<!-- Add screenshots to help explain your changes -->

## Additional Context

<!-- Add any other context about the pull request here -->

### Dependencies

<!-- List any dependencies that are required for this change -->

- [ ] Requires update to contract-go library: <!-- specify version -->
- [ ] Requires changes to related projects: <!-- specify -->

### Documentation

<!-- Link to any external documentation or references -->

- Related documentation:
- API documentation:

## Reviewer Notes

<!-- Any specific areas you want reviewers to focus on -->

---

**For Maintainers:**

- [ ] Reviewed by at least one maintainer
- [ ] All CI checks pass
- [ ] No unresolved review comments
- [ ] Commit messages are clean and follow conventions
- [ ] Documentation is updated
- [ ] Ready to merge
