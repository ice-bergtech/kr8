local config = std.extVar('kr8');

local kompose_template = std.native('komposeFile')('compose', 'compose', {
    kr8_spec: config.kr8_spec,
    namespace: config.namespace,
    release_name: config.release_name,
});

[
    object
    for object in std.objectValues(kompose_template)
]