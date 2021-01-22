local c = import 'github.com/rancher/local-path-provisioner/deploy/chart/Chart.yaml';
local k = import 'ksonnet-util/kausal.libsonnet';

local defaultValues = {
  nodeSelector: {
    'node-role.kubernetes.io/master': 'true',
  },
  tolerations: [
    {
      key: 'node-role.kubernetes.io/master',
      operator: 'Exists',
      effect: 'NoSchedule',
    },
  ],
  nodePathMap: [
    {
      node: 'DEFAULT_PATH_FOR_NON_LISTED_NODES',
      paths: [
        '/data/local-path-provisioner',
      ],
    },
  ],
  configmap: {
    helperPod: std.manifestYamlDoc({
      apiVersion: 'v1',
      kind: 'Pod',
      metadata: {
        name: 'helper-pod',
      },
      spec: {
        containers: [
          {
            name: 'helper-pod',
            image: 'docker.io/library/busybox',
          },
        ],
      },
    }),
  },
};

{
  helm: function(values={}) std.native('helmTemplate')(
    c.name,
    c.__dirname,
    defaultValues + values + { calledFrom: '/' },
  ),
}
