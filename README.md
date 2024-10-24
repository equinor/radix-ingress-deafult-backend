[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![Release](https://github.com/equinor/radix-ingress-default-backend/actions/workflows/release-please.yaml/badge.svg?branch=main&event=push)](https://github.com/equinor/radix-ingress-default-backend/actions/workflows/release-please.yaml)
[![SCM Compliance](https://scm-compliance-api.radix.equinor.com/repos/equinor/radix-ingress-default-backend/badge)](https://developer.equinor.com/governance/scm-policy/)

# Radix Ingress Defailt Backend

Responds with a nice error page for Radix (and Equinor) based on the requested domain

## How we work

Commits to the main branch must follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) and uses Release Please to create new versions.


## Deployment

Radix Ingress Default Backend follows the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) and uses Release Please to create new versions. 
All commits to main should either start with `feat: ` or with `fix: ` (a `!` will suggest a major version bump).

Release-Please will create a PR with for a new Release with all unreleased commits, bumping the version according to the commit message, linking the PRs and writing a nice [CHANGELOG](CHANGELOG). 

After a new version is released, we need to update radix-flux to start using it.

Each commit to `main` will also build a _nightly_ build that can be used if required.

## Pull request checking

Radix Ingress Default Backend makes use of [GitHub Actions](https://github.com/features/actions) for build checking in every pull request to the `main` branch. Refer to the [configuration file](.github/workflows/pr.yaml) of the workflow for more details.

## Contributing

Read our [contributing guidelines](./CONTRIBUTING.md)

------------------

[Security notification](./SECURITY.md)
