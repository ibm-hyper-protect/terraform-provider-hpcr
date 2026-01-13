# Security Policy

## Supported Versions

We release security patches for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.x.x   | :white_check_mark: |

Versions older than those listed above are no longer supported and will not receive security patches.

We recommend always using the latest release to ensure you have the most up-to-date security patches.

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

If you discover a security vulnerability in this project, please report it responsibly:

### Preferred Method: GitHub Security Advisories

1. Go to the [Security](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/security) tab
2. Click on "Report a vulnerability"
3. Fill out the advisory form with details

### Alternative: Direct Contact

You can also email the maintainers directly. See [MAINTAINERS.md](MAINTAINERS.md) for contact information.

### What to Include

When reporting a vulnerability, please include:

- **Description**: A clear description of the vulnerability
- **Impact**: What could an attacker achieve by exploiting this vulnerability
- **Steps to reproduce**: Detailed steps to reproduce the issue
- **Affected versions**: Which versions of the provider are affected
- **Suggested fix**: If you have ideas on how to fix it (optional)
- **Your contact information**: So we can follow up with you

## Response Timeline

We are committed to responding to security reports promptly:

- **Acknowledgment**: Within 3 business days
- **Initial assessment**: Within 5 business days
- **Fix for high-severity issues**: Within 30 days
- **Fix for medium/low-severity issues**: Within 90 days

We will keep you informed about the progress of fixing the vulnerability.

## Security Best Practices

### For Users

When using this Terraform provider, follow these security best practices:

1. **Keep dependencies updated**: Regularly update to the latest version of the provider
2. **Protect sensitive data**:
   - Never commit encryption keys, private keys, or CA certificates to version control
   - Use Terraform variables and secure secret management solutions (HashiCorp Vault, AWS Secrets Manager, etc.)
   - Mark sensitive outputs with `sensitive = true`
3. **Validate contracts**: Always validate contract schemas before deployment
4. **Use HTTPS**: Ensure certificate downloads use HTTPS (default behavior)
5. **Secure state files**: Terraform state files may contain sensitive data - store them securely with encryption at rest
6. **Review generated contracts**: Verify contract checksums (sha256_in, sha256_out) for integrity

### For Contributors

If you're contributing to this project:

1. **Code review**: All code changes must be reviewed before merging
2. **Dependency auditing**: Regularly audit dependencies for known vulnerabilities
3. **Security-focused testing**: Write tests that verify security properties
4. **Avoid hardcoded secrets**: Never commit real credentials, even in test code
5. **Follow secure coding practices**: Input validation, proper error handling, secure defaults

## Security Considerations

This Terraform provider:

- Generates and handles cryptographic keys (RSA 4096-bit)
- Performs encryption operations using HPVS encryption certificates
- Creates signed contracts for HPCR deployments
- Handles sensitive workload and environment data
- Relies on OpenSSL for cryptographic operations

Users should be aware that:

- Generated contracts contain encrypted sensitive data
- Private signing keys must be protected and stored securely
- Terraform state may contain sensitive information
- The provider delegates cryptographic operations to the underlying contract-go library

## Disclosure Policy

When we receive a security report:

1. We will confirm the vulnerability and determine its severity
2. We will develop and test a fix
3. We will prepare a security advisory
4. We will release a patched version
5. We will publish the security advisory with appropriate credits

We follow a coordinated disclosure model. We request that security researchers:

- Allow reasonable time for us to fix the vulnerability before public disclosure
- Provide us with a reasonable amount of detail to reproduce the issue
- Do not exploit the vulnerability beyond what's necessary to demonstrate it

## Non-Retaliation

We will not take legal action against security researchers who:

- Follow this responsible disclosure policy
- Act in good faith
- Do not cause harm to the project or its users
- Do not access or modify data beyond what's necessary to demonstrate the vulnerability

## Attribution

We appreciate the work of security researchers and will acknowledge your contribution in:

- The security advisory (unless you prefer to remain anonymous)
- Release notes for the patched version
- This document's hall of fame (coming soon)

## Questions?

If you have questions about this security policy, please open a discussion in the [GitHub Discussions](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/discussions) or contact the maintainers listed in [MAINTAINERS.md](MAINTAINERS.md).
