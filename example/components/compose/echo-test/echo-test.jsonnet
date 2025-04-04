local config = std.extVar('kr8');

local kompose_template = std.native('komposeFile')('compose/compose.yml', 'compose', {
    kr8_spec: config.kr8_spec,
    namespace: config.namespace,
    release_name: config.release_name,
    called_from: std.thisFile
});

[
    object
    for object in kompose_template
]