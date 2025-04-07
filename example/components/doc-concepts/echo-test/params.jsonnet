{
  namespace: 'kube-system',
  release_name: 'echo-test',
  kr8_spec: {
    includes: [ "echo.jsonnet" ],
    extfiles: [
      {"echoFile": "./vendor/" + self.version + "/echo.yml"}
    ],
  },
  version: "v1.0.0"
}