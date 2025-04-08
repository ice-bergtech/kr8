{
  local comp = self,
  namespace: 'kube-system',
  release_name: 'echo-test',
  kr8_spec: {
    includes: [ "echo.jsonnet" ],
    extfiles: {
      echoFile: "./vendor/" + comp.version + "/echo.yml"
    },
  },
  version: "v1.0.0",
}