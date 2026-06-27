# Unit Testing

## Standards

Every Titan service must have unit test coverage for all public functions and methods. Unit tests are deterministic, fast (sub-millisecond per test), and isolated from external dependencies.

## Expectations

Coverage minimum is 80% line coverage with business logic achieving 90%+ branch coverage. No network calls are permitted — external services are mocked or stubbed. Go tests use table-driven patterns while Python tests use parametrize. Test naming follows `Test_functionName_scenario` or `test__function_name__scenario`. Every test follows Arrange-Act-Assert structure.

## Tools by Language

| Language | Framework | Mocking | Assertions |
|----------|-----------|---------|------------|
| Go | `testing` standard library | `gomock` / moq | `testify/require` |
| Rust | `#[cfg(test)]` + `cargo test` | `mockall` | `assert_eq!`, `assert!(…)` |
| Python | `pytest` | `unittest.mock` / `pytest-mock` | `assert` + `pytest.approx` |
| TypeScript | `vitest` | `vitest.mock` | `expect` matchers |

## What to Test

Business logic and validation rules, edge cases including empty inputs, nil pointers, boundary values, and overflow, error paths including every error return value, and concurrency safety with race condition detection via `-race` flag (Go) or `loom` (Rust).

## What Not to Test

Generated code (protobuf stubs, OpenAPI clients), framework internals (ORM, router, middleware unless customized), and configuration loading (tested once in a shared utility).

## Running Tests Locally

Go tests run with `go test ./... -count=1 -race -cover`. Rust tests run with `cargo test --all-features`. Python tests run with `pytest --cov --cov-fail-under=80`. TypeScript tests run with `vitest run --coverage`.