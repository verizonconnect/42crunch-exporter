# Changelog

## [1.2.6](https://github.com/verizonconnect/42crunch-exporter/compare/v1.2.5...v1.2.6) (2025-09-25)


### Miscellaneous

* add GitHub Actions to Dependabot configuration ([d671e1e](https://github.com/verizonconnect/42crunch-exporter/commit/d671e1ee2ef412bc7decfcc713b35c4899220329))
* **dependabot:** update GitHub Actions directory path in configuration ([fda2d0f](https://github.com/verizonconnect/42crunch-exporter/commit/fda2d0f826607d689a5d760546415970662920b9))

## [1.2.5](https://github.com/verizonconnect/42crunch-exporter/compare/v1.2.4...v1.2.5) (2025-09-25)


### Bug Fixes

* correct import path for yaml package in go.mod ([342db3d](https://github.com/verizonconnect/42crunch-exporter/commit/342db3d6bf7d6f12cceaa1fd46b08f201210cadd))


### Miscellaneous

* change CollectionInclRegex from pointer to string and improve logging messages ([00d9e02](https://github.com/verizonconnect/42crunch-exporter/commit/00d9e024ab8279192ef2b72b119e2c305aa29b5e))
* fix spelling mistake ([2326688](https://github.com/verizonconnect/42crunch-exporter/commit/232668838de4f30e2e249057b31023f674e87b33))
* improve logging for regex matching and ensure context is used during server shutdown ([bb326f1](https://github.com/verizonconnect/42crunch-exporter/commit/bb326f120f7500a18af6d5bd288ae1bb801dd88b))
* package import ([5f69411](https://github.com/verizonconnect/42crunch-exporter/commit/5f6941111a56b95cf67eaf371bb20bd7b48046c0))
* remove redundant null check ([f7b9f2d](https://github.com/verizonconnect/42crunch-exporter/commit/f7b9f2d0ad0e9a396a976cdd0662df89010466d0))
* remove unnecessary check ([1909e8a](https://github.com/verizonconnect/42crunch-exporter/commit/1909e8aea2f83e5190320d22eaf566b65da5f517))
* update logging to use slog package and improve error handling ([cbe78d0](https://github.com/verizonconnect/42crunch-exporter/commit/cbe78d01ac565dfb83359d333091fb38dd17ae00))
* use context.Background for server shutdown instead of context.TODO ([80d3b11](https://github.com/verizonconnect/42crunch-exporter/commit/80d3b113daf82f61ae7362803a49a968498e2d35))

## [1.2.3](https://github.com/verizonconnect/42crunch-exporter/compare/v1.2.2...v1.2.3) (2023-07-12)


### Miscellaneous

* bump deps ([#12](https://github.com/verizonconnect/42crunch-exporter/issues/12)) ([bb789ce](https://github.com/verizonconnect/42crunch-exporter/commit/bb789cea0601e58e3e1eee29f44977ddef6cf9cc))

## [1.2.2](https://github.com/verizonconnect/42crunch-exporter/compare/v1.2.1...v1.2.2) (2023-07-10)


### Bug Fixes

* docker entrypoint ([#10](https://github.com/verizonconnect/42crunch-exporter/issues/10)) ([0d99bd9](https://github.com/verizonconnect/42crunch-exporter/commit/0d99bd990aee0b753bfbeba73af1ae09f945350a))

## [1.2.1](https://github.com/verizonconnect/42crunch-exporter/compare/v1.2.0...v1.2.1) (2023-07-06)


### Bug Fixes

* multi-architecture builds ([fa85e46](https://github.com/verizonconnect/42crunch-exporter/commit/fa85e46d42fa7e54aaff0e2bd80cbe8956b9f047))

## [1.2.0](https://github.com/verizonconnect/42crunch-exporter/compare/v1.1.2...v1.2.0) (2023-07-06)


### Features

* add docker image ([#7](https://github.com/verizonconnect/42crunch-exporter/issues/7)) ([72ace06](https://github.com/verizonconnect/42crunch-exporter/commit/72ace06dc7296122d37ef9aa7aac5f32f0ab124e))

## [1.1.2](https://github.com/verizonconnect/42crunch-exporter/compare/v1.1.1...v1.1.2) (2023-06-27)


### Miscellaneous

* transfer ownership ([e2f4040](https://github.com/verizonconnect/42crunch-exporter/commit/e2f404091cf0e33696f1d12191c3ee22770b6c1c))

## [1.1.1](https://github.com/verizonconnect/42crunch-exporter/compare/v1.1.0...v1.1.1) (2023-06-26)


### Miscellaneous

* cleanup go.mod ([b8f9388](https://github.com/verizonconnect/42crunch-exporter/commit/b8f93880d468d52dd34fd0d51a9191c270252c4b))

## [1.1.0](https://github.com/verizonconnect/42crunch-exporter/compare/v1.0.0...v1.1.0) (2023-06-23)


### Features

* get assessment report state ([efe1cd1](https://github.com/verizonconnect/42crunch-exporter/commit/efe1cd1ecc97673ebe3be5e9eb894ad623efe4d1))


### Miscellaneous

* fix lint warning ([7f0d546](https://github.com/verizonconnect/42crunch-exporter/commit/7f0d5461216a035fde70800aac1c89b240be2da6))
* update help labels ([19a9d2e](https://github.com/verizonconnect/42crunch-exporter/commit/19a9d2e1e15905eabd4c2f3f2027f35ea4d86bc3))

## 1.0.0 (2023-06-18)


### Features

* report api audit metrics ([#1](https://github.com/verizonconnect/42crunch-exporter/issues/1)) ([832a00d](https://github.com/verizonconnect/42crunch-exporter/commit/832a00da1802707c292b852b212a057395fd352b))


### Miscellaneous

* improve exporter ([3926cb0](https://github.com/verizonconnect/42crunch-exporter/commit/3926cb08f9980dfc8add9a607531441c6c5aa3a1))
* improve exporter ([#2](https://github.com/verizonconnect/42crunch-exporter/issues/2)) ([3926cb0](https://github.com/verizonconnect/42crunch-exporter/commit/3926cb08f9980dfc8add9a607531441c6c5aa3a1))
