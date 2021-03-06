{
  module: 'github.com/x/c',
  jpath: './vendor',
  replace: {
    'github.com/x/b': '../b',
  },
  require: {
    'github.com/grafana/jsonnet-libs':: 'v0.0.0-20210226150703-815de90bad0e',
    'github.com/jsonnet-libs/docsonnet':: 'v0.0.0-20200817072722-3e1757637edf',
    'github.com/jsonnet-libs/k8s-alpha':: 'v0.0.0-20210118111845-5e0d0738721f',
    'github.com/rancher/local-path-provisioner':: 'v0.0.19',
    'github.com/x/a':: 'v0.0.0',
    'github.com/x/b': 'v0.0.0',
  },
}
