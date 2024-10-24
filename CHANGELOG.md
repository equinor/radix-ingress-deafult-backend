# Changelog

## [1.0.2](https://github.com/equinor/radix-ingress-default-backend/compare/v1.0.1...v1.0.2) (2024-10-24)


### Bug Fixes

* rename helm chart ([aaf6f50](https://github.com/equinor/radix-ingress-default-backend/commit/aaf6f5086eb437b029b88372b336e47bba93978d))

## 1.0.1 (2024-10-24)

**Full Changelog**: https://github.com/equinor/radix-ingress-default-backend/compare/v1.0.0...v1.0.1

## [1.0.0](https://github.com/equinor/radix-ingress-default-backend/compare/v0.2.0...v1.0.0) (2024-10-24)


### ⚠ BREAKING CHANGES

* First release

### Features

* embed files to improve security, only serve html content ([eb85e7d](https://github.com/equinor/radix-ingress-default-backend/commit/eb85e7d5f21c35460f4c259ec4bb21ab34ee55b7))
* First release ([ccff4db](https://github.com/equinor/radix-ingress-default-backend/commit/ccff4db3732c189771f260d446a0f0d1c82bf4d0))
* Log all headers (except authorization) ([86805a7](https://github.com/equinor/radix-ingress-default-backend/commit/86805a7a391b06e69f9c8fa7c15d84719d89d4b5))


### Bug Fixes

* Cleanup main and use httptest.NewServer ([c7a0e51](https://github.com/equinor/radix-ingress-default-backend/commit/c7a0e512da2168df6c5e2c5aa972d92eb8232917))
* container port ([b614bea](https://github.com/equinor/radix-ingress-default-backend/commit/b614bea48494ae1762e741d3f5ccb5063d535136))
* default backend service should use port 80 ([0dffae2](https://github.com/equinor/radix-ingress-default-backend/commit/0dffae243f797c1d9c13453fb66d5ca860491973))
* dont log cookies ([7a54334](https://github.com/equinor/radix-ingress-default-backend/commit/7a5433452c029c58c6eac24a5cb87e0596fc568d))
* fix add default appName if missing ([0d397ef](https://github.com/equinor/radix-ingress-default-backend/commit/0d397efe4988140fd89dc2f9aeb3d467bd65ceae))
* mark removed authorization headers as removed ([9349640](https://github.com/equinor/radix-ingress-default-backend/commit/9349640fc67b217bcdb5cd1e16ed7115895e3de2))
* run tests on pr ([f69ceec](https://github.com/equinor/radix-ingress-default-backend/commit/f69ceecf753c8ec95a6e69fcb236479d72cc789d))
* use context logger for logging request logger ([7d6f7bf](https://github.com/equinor/radix-ingress-default-backend/commit/7d6f7bf93813b9971be88a304f35595def0b5db9))

## [0.2.0](https://github.com/equinor/radix-ingress-default-backend/compare/v0.1.0...v0.2.0) (2024-10-22)


### Features

* Add Equinor.com support ([859c72d](https://github.com/equinor/radix-ingress-default-backend/commit/859c72d40e0e2214a54195d02dc5172ab26d9b83))
* update readme ([b0fd9a7](https://github.com/equinor/radix-ingress-default-backend/commit/b0fd9a7a909ab04716173306128f6ba51c88b320))
