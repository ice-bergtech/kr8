{
  _cluster+: {
    name: 'standard',
  },
  _components+: {
    external_dns: { path: "components/chart/external-dns" },
    sealed_secrets: { path: "components/chart/sealed-secrets" },
  },
}
