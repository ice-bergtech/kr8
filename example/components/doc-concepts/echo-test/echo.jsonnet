# a jsonnet external variable from kr8 that gets contains cluster-level configuration
#local kr8_cluster = std.extVar('kr8_cluster');

local deployment = std.parseYaml( std.extVar("echoFile") );

deployment