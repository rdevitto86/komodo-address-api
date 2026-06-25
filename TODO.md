# TODO

> **Current Version:** V1

## V1 (Current)

> Status: Implemented but provider calls are stubs. Stateless — no DB needed.

### OpenAPI

- **[M]** Complete `openapi.yaml` — required by order-api and cart-api for address validation codegen

### Open Items

- **[H]** Wire real address validation provider (SmartyStreets, Google Address Validation, or similar) — replace all 3 stub `TODO` bodies in `internal/provider/address.go`
- **[M]** Add provider API key secret (`ADDRESS_PROVIDER_API_KEY`) to LocalStack init seed
- **[L]** Add unit tests for provider error handling and response mapping
- **[L]** Wire `GET /health/ready` in `cmd/public/main.go` once address provider is wired — checkers: `HTTPChecker("address-provider", os.Getenv("ADDRESS_PROVIDER_URL")+"/health")`; stateless service, no DB checker needed; blocked on forge SDK `api/handlers/health` release and provider selection
- **[L]** Live address lookup/autocomplete on keypress — add a typeahead endpoint (e.g. `GET /addresses/suggest`) backed by the provider's autocomplete API (e.g. Google Places Autocomplete, SmartyStreets Autocomplete) so client address fields can show suggestions as the user types

## Testing

- **[M]** **Implement CI test stack** — add `github.com/stretchr/testify` and `go.uber.org/mock` to `go.mod`; generate mocks from the provider interface via `mockgen -source`; convert stub `*_test.go` files to real unit tests (table-driven, `t.Run` subtests) with `net/http/httptest` for handler layer; add `testutil.Component(t)` / `testutil.Integration(t)` tier decorators from the SDK (`github.com/rdevitto86/komodo-forge-sdk-go/testing/testutil`, `TEST_TIER`-gated; default tier is `unit`); add `testcontainers-go` for integration tests against LocalStack if any AWS calls are added; apply section banners. Reference auth-api as the canonical pattern.
