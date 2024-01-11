#!/bin/sh

if [ -z "$GOPATH" ]; then
  GOPATH="${GOPATH:-$HOME/go}"
fi

if [ ! -d "$GOPATH" ]; then
  mkdir -p "$GOPATH/bin"
fi

cat >~/.terraformrc <<EOF
provider_installation {
  dev_overrides {
    "opsgenie/opsgenie" = "$GOPATH/bin"
  }
  direct {}
}
EOF
