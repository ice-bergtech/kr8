{
  _cluster+: {
    name: 'docs-kube-1',
    cluster_type: "aws",
  },
  _components+: {
    jsonnetStorageClasses: { path: "components/doc-concepts/jsonnetStorageClasses" },
    nativeFuncs: { path: "components/doc-concepts/nativeFunc" },
    echo_test: { path: "components/doc-concepts/echo-test" },
    echo_test_compose: { path: "components/doc-concepts/echo-test-compose" },
  },
}
