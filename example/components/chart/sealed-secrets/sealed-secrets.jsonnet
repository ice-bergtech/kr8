local config = std.extVar('kr8');

local helm_template = std.native('helmTemplate')(config.release_name, './vendor/'+"sealed-secrets-"+config.chart_version, {
    calledFrom: std.thisFile,
    namespace: config.namespace,
    values: config.helm_values,
});

[
    object
    for object in std.objectValues(helm_template)
    if 'kind' in object && object.kind != 'Secret'
]