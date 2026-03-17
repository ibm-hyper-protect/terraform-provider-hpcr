# Maintainers

This document lists the maintainers of the terraform-provider-hpcr project.

## Current Maintainers

| Name | GitHub Handle | Role | Email | Focus Areas |
|------|---------------|------|-------|-------------|
| Sashwat K | [@Sashwat-K](https://github.com/Sashwat-K) | Lead Maintainer | Sashwat.K@ibm.com | Overall project direction, releases, core features |
| Vikas Sharma | [@vikas-sharma24](https://github.com/vikas-sharma24) | Maintainer | Vikas.Sharma24@ibm.com | Releases and core features |
| Lokesh Puthalapattu | [@Lokesh-Puthalapattu](https://github.com/Lokesh-Puthalapattu) | Security Lead | lokesh.puthalapattu@ibm.com | Overall Security |

## Responsibilities

### Lead Maintainer

The Lead Maintainer has the following responsibilities:

- Set the technical direction for the project
- Manage releases and versioning
- Review and merge pull requests
- Manage community interactions
- Coordinate security responses
- Make final decisions when consensus cannot be reached

### All Maintainers

All maintainers are expected to:

- Be responsive to issues and pull requests
  - Acknowledge new issues within 3 business days
  - Provide initial feedback on pull requests within 3 business days
  - Respond to security reports within 3 business days
- Maintain high code quality standards
- Support contributors and provide constructive feedback
- Uphold the [Code of Conduct](CODE_OF_CONDUCT.md)
- Stay informed about developments in the Terraform ecosystem and IBM Hyper Protect technologies

## Becoming a Maintainer

We welcome new maintainers who have demonstrated:

- Consistent, high-quality contributions to the project
- Active engagement with the community (issues, discussions, reviews)
- Understanding of the project's goals and architecture
- Alignment with the project's values and Code of Conduct

### Criteria for Maintainer Nomination

A maintainer candidate typically demonstrates:

**Technical Excellence:**
- 10+ merged pull requests with high code quality
- Contributions across multiple areas (features, tests, docs)
- Understanding of Terraform provider architecture
- Knowledge of IBM Hyper Protect technologies

**Community Engagement:**
- Active participation in issue discussions (20+ meaningful comments)
- Helpful code reviews for other contributors (5+ reviews)
- Responsive to questions and supportive of newcomers
- Alignment with project values and Code of Conduct

**Commitment:**
- 6+ months of consistent contribution history
- Demonstrated reliability and follow-through
- Willingness to take on maintainer responsibilities
- Available for regular engagement (issues, PRs, discussions)

**Note:** These are guidelines, not strict requirements. Quality and impact matter more than quantity.

### Process

1. An existing maintainer nominates a candidate
2. Current maintainers discuss the nomination
3. If consensus is reached, the nominee is invited to become a maintainer
4. Upon acceptance, the new maintainer is added to this document and granted appropriate permissions

## Decision-Making Process

### Regular Decisions

For most decisions (bug fixes, minor features, documentation):

- Decisions are made through normal GitHub issues and pull requests
- Maintainers review and provide feedback
- Changes can be merged with approval from at least one maintainer
- If the code is written by a maintainer, it requires one additional maintainer's approval

**Examples of Regular Decisions:**
- Bug fixes that don't change public APIs
- Documentation improvements
- Test additions or improvements
- Dependency updates (non-breaking)
- Example code additions
- Performance optimizations (no API changes)

### Major Decisions

For significant changes (new features, architectural changes, breaking changes):

1. Create an RFC (Request for Comments) as a GitHub issue
2. Allow at least one week for discussion
3. Maintainers work towards consensus
4. If consensus cannot be reached, the Lead Maintainer makes the final decision

**Examples of Major Decisions:**
- New resources or data sources
- Breaking changes to existing APIs
- Significant architectural changes
- Changes to release process or versioning
- Addition of new dependencies with large footprint
- Changes to security or cryptographic implementations

## Code Review and Merging

### Review Guidelines

When reviewing pull requests, maintainers should:

- Be **constructive** - Focus on helping the contributor improve their work
- Be **timely** - Provide feedback within 3 business days when possible
- Be **thorough** - Check code quality, tests, documentation, and adherence to standards
- Be **respectful** - Disagreement is fine, but keep discourse civil and professional

### Merging Criteria

Pull requests should only be merged when:

- All CI checks pass
- At least one maintainer has approved (two for maintainer PRs)
- All review comments are resolved
- The commit messages follow [Conventional Commits](https://www.conventionalcommits.org/)
- Documentation is updated (if applicable)
- Tests are included for new functionality

## Communication Channels

- **General questions and discussions**: [GitHub Discussions](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/discussions)
- **Bug reports and feature requests**: [GitHub Issues](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/issues)
- **Security issues**: [GitHub Security Advisories](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/security/advisories) (never public issues)
- **Code of Conduct violations**: Direct contact with maintainers via email

## Stepping Down

Maintainers may step down at any time by:

1. Notifying other maintainers
2. Submitting a pull request to remove themselves from this document
3. Coordinating the transition of any ongoing responsibilities

We appreciate all contributions from past maintainers and thank them for their service to the project.

## Emeritus Maintainers

We recognize and thank former maintainers for their contributions:

<!-- This section will list maintainers who have stepped down -->
<!-- Currently empty -->

---
