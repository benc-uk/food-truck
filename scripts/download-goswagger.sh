#!/bin/bash

scriptDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
  jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url')
curl -o "$scriptDir"/../bin/swagger -L'#' "$download_url"
chmod +x "$scriptDir"/../bin/swagger