#!/usr/bin/env python3
"""Keep published delivery-attempt APIs while allowing generated pagination APIs.

Fail closed when the expected generator shape changes. Run after oapi-codegen,
never directly patch client.gen.go: regeneration must preserve these contracts.
"""
from pathlib import Path
import re

path = Path("client/client.gen.go")
source = path.read_text()
# Rename identifiers, not substrings (GetDeliveryAttemptsParams stays stable).
renames = {
    "GetDeliveryAttempts": "GetDeliveryAttemptsWithParams",
    "NewGetDeliveryAttemptsRequest": "NewGetDeliveryAttemptsWithParamsRequest",
    "GetDeliveryAttemptsWithResponse": "GetDeliveryAttemptsWithParamsWithResponse",
    "GetDeliveryAttemptsResponse": "GetDeliveryAttemptsWithParamsResponse",
    "ParseGetDeliveryAttemptsResponse": "ParseGetDeliveryAttemptsWithParamsResponse",
}
for old, new in renames.items():
    if not re.search(r"\b" + old + r"\b", source):
        raise SystemExit("Expected generated identifier missing: " + old)
    source = re.sub(r"\b" + old + r"\b", new, source)
# Do not expand existing interfaces: downstream mocks must still implement them.
for new, old, response in [
    ("GetDeliveryAttemptsWithParams", "GetDeliveryAttempts", "http.Response"),
    ("GetDeliveryAttemptsWithParamsWithResponse", "GetDeliveryAttemptsWithResponse", "GetDeliveryAttemptsResponse"),
]:
    pattern = r"(?m)^\t" + new + r"\(ctx context.Context, projectID string, eventDeliveryID string, params \*GetDeliveryAttemptsParams, reqEditors \.\.\.RequestEditorFn\) \([^\n]+\)$"
    replacement = "\t" + old + "(ctx context.Context, projectID string, eventDeliveryID string, reqEditors ...RequestEditorFn) (*" + response + ", error)"
    source, count = re.subn(pattern, replacement, source)
    if count != 1:
        raise SystemExit("Expected exactly one interface method: " + new)
source = source.replace("(the `GetDeliveryAttemptsWithParams` operationId)", "(the `GetDeliveryAttempts` operationId)")
source = source.replace("\t// GetDeliveryAttemptsWithParams List", "\t// GetDeliveryAttempts List")
source = source.replace("\t// GetDeliveryAttemptsWithParamsWithResponse List", "\t// GetDeliveryAttemptsWithResponse List")
path.write_text(source)
