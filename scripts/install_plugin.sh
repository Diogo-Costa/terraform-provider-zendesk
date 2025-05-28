#!/usr/bin/env bash
#
# install_plugins.sh
#
# This script installs the plugin in ~/.terraform.d/plugins

set -e

oss=( darwin )
archs=( arm64 )
plugins_dir="${HOME}/.terraform.d/plugins"

install_plugin() {
  plugin=$1
  version=1.0.0
  plugin_name=terraform-provider-$(basename "${plugin}")
  echo "Installing Terraform plugin ${plugin}..."
  for os in "${oss[@]}"
  do
    for arch in "${archs[@]}"
    do
      file="${plugin_name}_v${version}"
      plugin_dst="${plugins_dir}/${plugin}/terraform.local/local/${plugin_name}/${version}/${file}"
      mkdir -p "$(dirname "${plugin_dst}")"
      go build -ldflags="-X main.version=${version} -X main.commit=n/a"
      mv terraform-provider-$(basename "${plugin}") ${plugin_dst}
      echo "Copied to ${plugin_dst}"
    done
  done
}

install_plugin "terraform.local/local/zendesk"

# Release to registry
# /opt/homebrew/Cellar/goreleaser/1.18.2/bin/goreleaser release --clean