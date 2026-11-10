#!/usr/bin/env bash
set -euo pipefail
# Fail closed if the test file or any mandatory regression test disappears.
for test in TestDeliveryAttemptsLegacyCalls TestDeliveryAttemptsPaginationRequest TestDeliveryAttemptsPaginationResponse TestDeliveryAttemptsLegacyInterfaces; do
  go test ./client -list "^${test}$" | grep -Fx "$test" >/dev/null
done
go test ./client -run '^TestDeliveryAttempts' -count=1
