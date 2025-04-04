local config = std.extVar('kr8');

local helm_template = std.native('helmTemplate')(config.release_name, './vendor/'+"sealed-secrets-"+config.chart_version, {
    calledFrom: std.thisFile,
    namespace: config.namespace,
    values: if 'helm_values' in config then config.helm_values else {},
});

[
    object
    for object in std.objectValues(helm_template)
    if 'kind' in object && object.kind != 'Secret'
]