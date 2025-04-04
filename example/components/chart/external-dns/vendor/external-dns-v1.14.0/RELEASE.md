### Added

- Added the option to explicitly enable or disable service account token automounting. ([#3983](https://github.com/kubernetes-sigs/external-dns/pull/3983)) [@gilles-gosuin](https://github.com/gilles-gosuin)
- Added the option to configure revisionHistoryLimit on the K8s Deployment resource. ([#4008](https://github.com/kubernetes-sigs/external-dns/pull/4008)) [@arnisoph](https://github.com/arnisoph)
- Added support for webhook providers, as a sidecar. ([#4032](https://github.com/kubernetes-sigs/external-dns/pull/4032) [@mloiseleur](https://github.com/mloiseleur)
- Added the option to configure ipFamilyPolicy and ipFamilies of external-dns Service.  ([#4153](https://github.com/kubernetes-sigs/external-dns/pull/4153)) [@dongjiang1989](https://github.com/dongjiang1989)

### Changed

- Avoid unnecessary pod restart on each helm chart version. ([#4103](https://github.com/kubernetes-sigs/external-dns/pull/4103)) [@jkroepke](https://github.com/jkroepke)
- Updated _ExternalDNS_ OCI image version to [v0.14.0](https://github.com/kubernetes-sigs/external-dns/releases/tag/v0.14.0). ([#4073](https://github.com/kubernetes-sigs/external-dns/pull/4073)) [@appkins](https://github.com/appkins)

### Deprecated

- The `secretConfiguration` value has been deprecated in favour of creating secrets external to the Helm chart and configuring their use via the `extraVolumes` & `extraVolumeMounts` values. ([#4161](https://github.com/kubernetes-sigs/external-dns/pull/4161)) [@stevehipwell](https://github.com/stevehipwell)
