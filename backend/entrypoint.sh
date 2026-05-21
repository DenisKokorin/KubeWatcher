#!/bin/sh
set -e

if [ -n "$KUBECONFIG" ] && [ -f "$KUBECONFIG" ]; then
  cp "$KUBECONFIG" /tmp/kubeconfig.orig
  sed -E 's#([A-Za-z]):\\\\#\/host/\1/#g; s#\\\\#/#g' /tmp/kubeconfig.orig > /tmp/kubeconfig.fixed
  export KUBECONFIG=/tmp/kubeconfig.fixed
fi

exec ./app --kubeconfig "$KUBECONFIG"
