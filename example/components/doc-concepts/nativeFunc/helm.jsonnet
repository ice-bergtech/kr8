local config = std.extVar('kr8');

local helm_template = std.native('helmTemplate')(config.release_name, './vendor/'+"sealed-secrets-"+config.chart_version, {
    calledFrom: std.thisFile,
    namespace: config.namespace,
    values: if 'helm_values' in config then config.helm_values else {},
});

local standardAnnotations = {
    "app.kubernetes.io/name": "sealed-secrets",
    "app.kubernetes.io/part-of": config.release_name,
    "app.kubernetes.io/version": config.chart_version,
};

[
    object {
        metadata+: if 'annotations' in object.metadata then
        if object.metadata.annotations != null then {
            annotations+: standardAnnotations
        } else {
            annotations: standardAnnotations
        } else {}
    }
    for object in std.objectValues(helm_template)
    if 'kind' in object && object.kind != 'Secret'
]