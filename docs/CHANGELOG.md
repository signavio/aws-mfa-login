## [0.1.10](https://github.com/signavio/aws-mfa-login/compare/v0.1.9...v0.1.10) (2022-03-31)


### Bug Fixes

* support ARM64 architecture for linux and mac new assets ([05182ba](https://github.com/signavio/aws-mfa-login/commit/05182ba50424ca643cb6cf6f6f5caae3fa7ca2f3)), closes [#46](https://github.com/signavio/aws-mfa-login/issues/46)

## [0.1.9](https://github.com/signavio/aws-mfa-login/compare/v0.1.8...v0.1.9) (2022-03-31)


### Bug Fixes

* support ARM64 architecture for linux and mac ([3583c73](https://github.com/signavio/aws-mfa-login/commit/3583c73a5123a7998821c17bc3154c9e2701590d)), closes [#46](https://github.com/signavio/aws-mfa-login/issues/46)

## [0.1.8](https://github.com/signavio/aws-mfa-login/compare/v0.1.7...v0.1.8) (2022-02-26)


### Bug Fixes

* **deps:** update golang dependencies ([6cc8219](https://github.com/signavio/aws-mfa-login/commit/6cc8219b2cc10c7918f6868518dd7dfa25fc0fbd))

## [0.1.7](https://github.com/signavio/aws-mfa-login/compare/v0.1.6...v0.1.7) (2021-06-05)


### Bug Fixes

* add default region for login ([3b6cf98](https://github.com/signavio/aws-mfa-login/commit/3b6cf980f6ee1aaa5008fc0dd379eeb8ee850998)), closes [#16](https://github.com/signavio/aws-mfa-login/issues/16)

## [0.1.6](https://github.com/signavio/aws-mfa-login/compare/v0.1.5...v0.1.6) (2021-06-01)


### Bug Fixes

* **deps:** update golang dependencies ([70ad52a](https://github.com/signavio/aws-mfa-login/commit/70ad52a70549e8369e4c51f49f6b655117a0da5e))

## [0.1.5](https://github.com/signavio/aws-mfa-login/compare/v0.1.4...v0.1.5) (2021-05-16)

### BREAKING CHANGE
sorry I used wrong commit message so no major version increase
* replace aws cli with golang sdk2 api calls to update kubeconfig ([1bb1895](https://github.com/signavio/aws-mfa-login/commit/1bb18953a76b3ae1956285559f432422f8d4d17e))

### Features
* upgrade to golang 1.14 ([1bb1895](https://github.com/signavio/aws-mfa-login/commit/1bb18953a76b3ae1956285559f432422f8d4d17e))
* assume more than one role per cluster is possible ([1bb1895](https://github.com/signavio/aws-mfa-login/commit/1bb18953a76b3ae1956285559f432422f8d4d17e))
* colorized output ([1bb1895](https://github.com/signavio/aws-mfa-login/commit/1bb18953a76b3ae1956285559f432422f8d4d17e))
* migrate to aws sdk 2 ([1bb1895](https://github.com/signavio/aws-mfa-login/commit/1bb18953a76b3ae1956285559f432422f8d4d17e))
* improved performance for kubeconfig update ([1bb1895](https://github.com/signavio/aws-mfa-login/commit/1bb18953a76b3ae1956285559f432422f8d4d17e))


### Bug Fixes

* add command line flag to disable colorized output ([acfa42e](https://github.com/signavio/aws-mfa-login/commit/acfa42e02588c8e7c6bc8619464cfd9136049b4e))
* fix auto-complete and support additionally zsh, fish and powershell ([0c6a2fb](https://github.com/signavio/aws-mfa-login/commit/0c6a2fb154efe562775b66c5ebb04c7bd1f9ea94))
* update README ([6925138](https://github.com/signavio/aws-mfa-login/commit/692513829af3cc1a5f85fa3acbdfcbc483fe8ec0))
* add test for write kubeconfig ([1bb1895](https://github.com/signavio/aws-mfa-login/commit/1bb18953a76b3ae1956285559f432422f8d4d17e))

## [0.1.4](https://github.com/signavio/aws-mfa-login/compare/v0.1.3...v0.1.4) (2021-05-08)


### Bug Fixes

* add semantic release and update pipeline ([b9671a8](https://github.com/signavio/aws-mfa-login/commit/b9671a8a035bc7bda18d09bf9669a9b69468bfe7))
