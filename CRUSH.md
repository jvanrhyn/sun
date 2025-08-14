# CRUSH.md

Repository quick-reference for agentic coding

Build/lint/test
- Build: make build (outputs to ./bin). Alt: go build -v ./...
- Run: go run main.go
- Test all: make test or go test -race -v ./...
- Test single package: go test -race -v ./cmd/sun
- Test single file: go test -race -v ./cmd/sun -run TestName -count=1
- Test single test: go test -race -v ./cmd/sun -run '^TestExactName$' -count=1
- Benchmarks: make bench or go test -benchmem -count=3 -bench ./...
- Coverage: make coverage (produces cover.out and cover.html)
- Lint: make lint (golangci-lint). Auto-fix: make lint-fix
- Dev watch (tests/bench): make watch-test / make watch-bench (requires reflex; install via make tools)

Environment and tooling
- Go 1.22.x. Modules enabled. Run go mod tidy after dependency changes.
- Required env vars (from .env): WEATHER_ACCESS_TOKEN, DEFAULT_LOCATION, NO_OF_DAYS. Build copies .env to ./bin.
- Install helpers: make tools (reflex, gotest, go-mod-outdated, goweight, golangci-lint, cover, nancy)

Code style guidelines
- Imports: standard lib first, then external, grouped and alphabetized. Use goimports formatting.
- Formatting: gofmt/goimports default; run golangci-lint locally before commits.
- Types and naming: exported identifiers use PascalCase with doc comments; unexported use camelCase; constants in PascalCase; keep var names clear and short.
- Errors: do not panic in libraries; return errors with context using fmt.Errorf("...: %w", err); only main may exit or log.Fatal. Prefer sentinel errors or errors.Is/As.
- HTTP/resources: always defer Close and handle close errors; check status codes explicitly; set timeouts where applicable.
- Concurrency: prefer context.Context for cancelation; avoid shared state; use race detector in tests (-race).
- CLI/flags: validate inputs; provide sane defaults from env; avoid global mutable state when possible.

CI
- GitHub Actions: .github/workflows/go.yml runs go build and go test on Go 1.22.

Notes
- No Cursor or Copilot rules detected. If added later (.cursor/rules/, .cursorrules, or .github/copilot-instructions.md), mirror key constraints here.
