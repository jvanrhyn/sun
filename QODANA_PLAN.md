# Qodana Remediation Plan

This plan prioritizes addressing likely Qodana (GoLand/Go inspections) findings observed in this repository. It focuses on safe, incremental improvements with minimal behavior changes.

Scope analyzed:
- main.go (Bubble Tea TUI)
- cmd/sun (client.go, models.go, sun.go)
- internal/ui (adapt.go, theme.go)
- Tests and configuration (qodana.yaml)

## Priorities and Actions

### 1) High priority – Correctness and robustness

1.1 Replace panic-based error handling in cmd/sun/sun.go ✓
- Problem: Multiple `panic` calls for network and JSON errors; fragile and flagged by Qodana.
- Action:
  - Refactor Run to return an error (`func Run() error`) or handle errors gracefully (print to stderr, exit with non-zero) if kept as an entrypoint.
  - Replace `panic(err)` and `panic("Weather Api not available")` with formatted errors.
  - Add error context for easier troubleshooting.
- Acceptance criteria: No `panic` usages remain in non-test code; Qodana no longer reports error-handling issues for this file.

1.2 Use context-aware HTTP with timeouts in cmd/sun/sun.go ✓
- Problem: Uses `http.Get` without context; no timeout; magic number `200`.
- Action:
  - Create `http.Client{Timeout: 10 * time.Second}`.
  - Use `context.WithTimeout` and `http.NewRequestWithContext`.
  - Compare with `http.StatusOK` instead of `200`.
- Acceptance criteria: Network calls are context-aware with explicit timeout; constants used for status checks.

1.3 Avoid variable shadowing in main.go ✓
- Problem: `key := os.Getenv("WEATHER_ACCESS_TOKEN")` shadows imported package `key` from bubbles.
- Action: Rename to `apiKey` (or similar) and update usage.
- Acceptance criteria: No shadowed identifiers reported by Qodana in main.go.

1.4 Explicit decision on .env loading errors in main.go ✓
- Problem: `_ = godotenv.Load()` ignores error silently.
- Action: Add a comment explaining intentional ignore (e.g., `.env` optional), or log a warning to stderr once.
- Acceptance criteria: Qodana either sees the rationale comment or a handled error path.

### 2) Medium priority – Maintainability and duplication

2.1 Deduplicate column width logic ✓
- Problem: main.go defines `defaultColumns()` and `setAdaptiveDimensions()` duplicating logic found in internal/ui/adapt.go.
- Action:
  - Replace in main.go with a call to `ui.AdaptiveColumns(width)`.
  - Optionally remove `defaultColumns()` if not needed; or keep as initial columns but align widths with AdaptiveColumns.
- Acceptance criteria: Single source of truth for table columns (internal/ui).

2.2 Idiomatic formatting and readability ✓
- Problem: Several single-line `if` statements with semicolons and inline returns; harder to read and sometimes flagged.
- Action: Expand to idiomatic multi-line statements.
- Acceptance criteria: gofmt-compliant, readable control flow; Qodana style warnings resolved.

2.3 Replace magic numbers with constants where meaningful ✓
- Problem: Repeated layout values (e.g., min widths) and status codes.
- Action: Introduce small, well-named constants for recurring values (only where it improves clarity).
- Acceptance criteria: Reduced magic numbers in business logic; clearer intent.

### 3) Low priority – Documentation

3.1 Add GoDoc comments for exported identifiers ✓
- Targets:
  - `cmd/sun.Client` and method `Forecast`.
  - `cmd/sun.Weather` struct.
  - `internal/ui.Theme` and `DefaultTheme`.
- Action: Add concise comments starting with the identifier name.
- Acceptance criteria: Qodana no longer flags missing documentation for exported names.

## File-by-file checklist

- cmd/sun/sun.go
  - [ ] Replace `panic` with error handling or return values.
  - [ ] Use `http.Client` with timeout and `context.WithTimeout`.
  - [ ] Use `http.StatusOK` instead of `200`.
  - [ ] Prefer `defer res.Body.Close()` with error check/report; no panics.

- main.go
  - [ ] Rename `key` variable to `apiKey`.
  - [ ] Add rationale for ignoring `.env` load error or log a warning once.
  - [ ] Consider using `ui.AdaptiveColumns` in `setAdaptiveDimensions` to remove duplication.
  - [ ] Expand one-line statements with semicolons to multiline for readability.

- cmd/sun/client.go
  - [ ] Add GoDoc for `Client` and `Forecast`.
  - [ ] Consider configurable user-agent and base URL (optional).

- cmd/sun/models.go
  - [ ] Add GoDoc for `Weather`.

- internal/ui/theme.go
  - [ ] Add GoDoc for `Theme` and `DefaultTheme`.

## Validation steps

1) Unit tests
- Run: `go test ./...` (must remain green)

2) go vet
- Run: `go vet ./...` (must be clean)

3) Qodana
- Local or CI: ensure qodana-go runs with `qodana.yaml` starter profile and produces no high-severity problems; address any residuals iteratively.

## Effort estimates
- High-priority items: ~1–2 hours
- Medium-priority items: ~1 hour
- Low-priority docs: ~30–45 minutes

## Risks and rollbacks
- The proposed changes are primarily refactors with no behavior changes in the TUI path; keep commits small and focused.
- If cmd/sun/sun.go is used externally, converting Run to return error may be breaking; alternatively, keep signature and handle errors internally (log + os.Exit(1)).

## Next steps
- Review this plan, confirm priorities and acceptance criteria.
- Apply high-priority changes first and re-run Qodana.
- Tackle medium then low priority items; iterate until the Problems tab is clean or acceptable.
