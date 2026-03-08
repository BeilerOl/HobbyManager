# Integration / API contract tests

Place here tests that validate the running API against the OpenAPI contract:

- **Contract tests**: ensure responses match the schemas in [../../api/openapi.yaml](../../api/openapi.yaml). Tools such as Dredd or openapi-response-validator can be used.
- **Integration tests**: call the real backend (e.g. `http://localhost:8080`) and assert behaviour.

Backend and frontend unit tests remain in `backend/` and `frontend/` respectively.
