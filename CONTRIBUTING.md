# Contributing to terraform-provider-hpcr

Thank you for your interest in contributing to the IBM Hyper Protect Container Runtime Terraform Provider! This document provides guidelines for contributing to the project.

## Code of Conduct

This project adheres to our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the maintainers listed in [MAINTAINERS.md](MAINTAINERS.md).

## How Can I Contribute?

### Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates. When creating a bug report, use the issue template and include as much detail as possible:

- **Clear title**: Describe the problem concisely
- **Steps to reproduce**: Provide detailed steps to reproduce the issue
- **Expected behavior**: What you expected to happen
- **Actual behavior**: What actually happened
- **Environment**: Terraform version, provider version, OS, etc.
- **Configuration**: Relevant Terraform configuration (sanitize sensitive data)
- **Logs**: Error messages and relevant log output

**Security vulnerabilities** should **never** be reported as public issues. Use [GitHub Security Advisories](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/security/advisories) instead.

### Suggesting Features

We welcome feature suggestions! Before creating a feature request:

1. Check existing issues and discussions to see if it's already been suggested
2. Consider if the feature aligns with the project's goals
3. Use the feature request template and include:
   - **Problem statement**: What problem does this solve?
   - **Proposed solution**: How would you like to see it implemented?
   - **Alternatives**: What alternatives have you considered?
   - **Use case**: Describe your specific use case

### Asking Questions

Before asking a question:

1. Check the [documentation](README.md) and existing examples
2. Search existing issues and discussions

For general questions:
- Use [GitHub Discussions](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/discussions)
- Tag your discussion appropriately

### Contributing Code

We appreciate code contributions! To ensure a smooth process:

1. **Open an issue first**: Discuss your approach before investing time in implementation
2. **Avoid duplicate work**: Check if someone else is already working on it
3. **Follow conventions**: Adhere to our coding standards and commit message format
4. **Include tests**: All code changes should include appropriate tests
5. **Update documentation**: Update relevant documentation for user-facing changes

## Getting Started


### Prerequisites

- **Go**: Version 1.24 or later
- **Terraform**: Version 1.0 or later (for testing)
- **Make**: For running build tasks
- **Git**: For version control

### Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR-USERNAME/terraform-provider-hpcr.git
   cd terraform-provider-hpcr
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/ibm-hyper-protect/terraform-provider-hpcr.git
   ```

4. **Install dependencies**:
   ```bash
   make install-deps
   ```

5. **Verify your setup**:
   ```bash
   # Build the provider
   make build

   # Run tests
   make test

   # Run linting
   make lint
   ```

## Development Workflow

1. **Create a feature branch** from `main`:
   ```bash
   git checkout main
   git pull upstream main
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our [coding standards](#coding-standards)

3. **Run tests** frequently during development:
   ```bash
   make test
   ```

4. **Tidy dependencies**:
   ```bash
   make tidy
   ```

5. **Format your code** as per standards:
   ```bash
   make lint
   ```

6. **Commit your changes** with [proper commit messages](#commit-messages)

7. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

8. **Open a Pull Request** from your fork to the main repository

## Coding Standards

This project follows standard Go coding conventions:

- Follow [Effective Go](https://golang.org/doc/effective_go.html) principles
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting (enforced by CI)
- Use `golangci-lint` for linting

### Key Practices

1. **Small, focused functions**: Each function should do one thing well
2. **Meaningful names**: Use clear, descriptive names for variables and functions
3. **Error handling**: Always handle errors explicitly, never ignore them
4. **Documentation**: Add godoc comments for exported functions and types
5. **Minimal exports**: Only export what needs to be public
6. **Testing**: Write table-driven tests, test edge cases and error conditions

### Terraform Provider Best Practices

- Follow [Terraform Plugin Development](https://developer.hashicorp.com/terraform/plugin) guidelines
- Use the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- Validate inputs early and provide helpful error messages
- Support both simple and complex use cases
- Maintain backwards compatibility whenever possible

## Commit Message Convention

This project follows [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Common Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation changes only
- **refactor**: Code changes that neither fix bugs nor add features
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks, dependency updates
- **ci**: Changes to CI/CD configuration

### Examples

```
feat(contract): add support for contract expiry validation

fix(encryption): handle malformed certificates gracefully

docs: update README with new resource examples

test(contract): add edge case tests for empty workload
```

## Pull Request Process

1. **Fork and create a branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Write clear, focused commits
   - Follow coding standards
   - Add or update tests
   - Update documentation

3. **Ensure quality**:
   ```bash
   make fmt      # Format code
   make lint     # Run linters
   make test     # Run tests
   make tidy     # Tidy dependencies
   ```

4. **Push and create PR**:
   - Push to your fork
   - Create a pull request against `main`
   - Link to related issues
   - Fill out the PR template completely

5. **Code review**:
   - Address review comments promptly
   - Keep the discussion focused and professional
   - Be open to feedback and suggestions

6. **Merging**:
   - All CI checks must pass
   - At least one maintainer approval required
   - All review comments must be resolved

### PR Checklist

Before submitting your PR, ensure:

- [ ] Code follows project conventions and style
- [ ] Tests pass locally (`make test`)
- [ ] New tests added for new functionality
- [ ] Documentation updated (README, examples, etc.)
- [ ] Commit messages follow Conventional Commits
- [ ] PR is linked to related issue(s)
- [ ] No sensitive data (credentials, keys) in code or commits

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test
go test -v -run TestResourceContract ./...

# Run with coverage
make test-coverage
```

### Writing Tests

- Use **table-driven tests** for multiple test cases
- Test both **success and failure** paths
- Test **edge cases** and **boundary conditions**
- Use **descriptive test names** that explain what's being tested
- Mock external dependencies appropriately

### Acceptance Tests

For Terraform provider acceptance tests:

```bash
# Set required environment variables
export TF_ACC=1
export IBM_CLOUD_API_KEY=your_api_key

# Run acceptance tests
make testacc
```

**Note**: Acceptance tests interact with real infrastructure and may incur costs.

### Test Provider Locally

To test the provider locally with Terraform without publishing to the registry:

```bash
# Build and install the provider to your Go bin directory
make build
make install

# Configure Terraform to use the local provider
# The provider will be installed to $GOPATH/bin (typically ~/go/bin)
cat <<EOF > ~/.terraformrc
provider_installation {
  dev_overrides {
    "ibm-hyper-protect/hpcr" = "$HOME/go/bin"
  }
  direct {}
}
EOF

# Enable debug logging (optional)
export TF_LOG=DEBUG

# Test with any example
cd examples/resources/hpcr_contract_encrypted

# Run the complete Terraform lifecycle
terraform init
terraform plan
terraform apply --auto-approve
terraform destroy --auto-approve
```

**Tips:**
- Remove or comment out the `dev_overrides` section in `~/.terraformrc` when done testing
- The provider binary must be rebuilt (`make build && make install`) after code changes
- Use `TF_LOG=TRACE` for even more detailed logging if needed

## Documentation

Keep documentation up to date:

- **README.md**: General overview and quick start
- **docs/**: Detailed resource and data source documentation
- **examples/**: Working example configurations
- **Inline comments**: For complex logic

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.

All contributions are subject to:
- [Developer Certificate of Origin (DCO) Version 1.1](https://developercertificate.org/)
- [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0.txt)

## Getting Help

- **General questions**: Use [GitHub Discussions](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/discussions)
- **Bugs**: Create an [issue](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/issues)
- **Security**: See [SECURITY.md](SECURITY.md)
- **Maintainers**: See [MAINTAINERS.md](MAINTAINERS.md)

## Recognition

We value all contributions and contributors will be:
- Acknowledged in release notes
- Listed as contributors on GitHub
- Mentioned in relevant documentation

Thank you for contributing to terraform-provider-hpcr!
