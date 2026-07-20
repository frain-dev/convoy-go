#!/usr/bin/env bash
set -euo pipefail

# Regenerate the API client from Convoy's OpenAPI spec with oapi-codegen into
# client/client.gen.go. The hand-written client (repo root package) and
# webhook verify (webhook.go) are not generated and never touched here.
#
# Requires: go 1.22+, curl. Run from the repo root.

SPEC_URL="${SPEC_URL:-https://raw.githubusercontent.com/frain-dev/convoy/main/docs/v3/openapi3.yaml}"
# Pin so regeneration output is reproducible; bump deliberately.
GENERATOR_VERSION="v2.8.0"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

curl -fsSL "$SPEC_URL" -o "$tmp/openapi3.yaml"

go run "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@${GENERATOR_VERSION}" \
  -config .oapi-codegen.yaml "$tmp/openapi3.yaml"

gofmt -w client/client.gen.go

echo "Generated client synced into client/client.gen.go"
