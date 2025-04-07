# a jsonnet external variable from kr8 that gets contains cluster-level configuration
local kr8_cluster = std.extVar('kr8_cluster');

# a jsonnet function for creating a storageclass
local StorageClass(name, type, default=false) = {
  apiVersion: 'storage.k8s.io/v1',
  kind: 'StorageClass',
  metadata: {
    name: name,
    annotations: {
      'storageclass.kubernetes.io/is-default-class': if default then 'true' else 'false',
    },
  },
  parameters: {
    type: type,
  },
};

# check the cluster configuration for a type, if it's AWS make a gp2 type storageclass
if kr8_cluster.cluster_type == 'aws' then std.objectValues(
  {
    // default gp2 storage class, not tied to a zone
    ebs_gp2: StorageClass('gp2', 'gp2', true) {},
  }
) else [] # don't make a storageclass if it's not AWS
