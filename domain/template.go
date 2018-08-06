package domain

// ConfigTemplate is a template of a configuration file.
const ConfigTemplate = `
---
# akoi - binary version control system
# https://github.com/suzuki-shunsuke/akoi
bin_dir: /usr/local/bin
link_dir: /usr/local/bin
bin_separator: "-"
# packages:
#   consul:
#     url: https://releases.hashicorp.com/{{.Name}}/{{.Version}}/{{.Name}}_{{.Version}}_darwin_amd64.zip
#     version: 1.2.1
#     files:
#     - name: consul
#       archive: consul
#   jq:
#     url: https://github.com/stedolan/{{.Name}}/releases/download/{{.Name}}-{{.Version}}/{{.Name}}-osx-amd64
#     version: 1.5
#     archive_type: unarchived
#     files:
#     - name: jq
`
