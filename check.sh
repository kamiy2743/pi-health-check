#!/usr/bin/env bash

set -euo pipefail

CONF="/home/kamiy2743/workspace/health-check/check.conf"

if [ ! -f "$CONF" ]; then
    echo "config not found: $CONF" >&2
    exit 1
fi

fail_count=0
failed_urls=""

while IFS= read -r url || [ -n "$url" ]; do
    case "$url" in
      "" | \#*) continue ;;  # 空行 or #で始まる行はスキップ
    esac

    echo "checking: $url"
    if curl -fsS --max-time 5 "$url" >/dev/null; then
        :
    else
        code=$?
        fail_count=$((fail_count + 1))
        failed_urls="${failed_urls}${url} (code=${code})"$'\n'
    fi
done < "$CONF"

if [ "$fail_count" -gt 0 ]; then
    echo "failed url list:" >&2
    printf '%s' "$failed_urls" >&2
    echo "failed checks: $fail_count" >&2
    exit 1
fi

echo "all checks passed"
