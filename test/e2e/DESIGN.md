# E2E Test Project — Design

## Context

After exploratory testing of RustFS, critical functional flows were identified and formalized into this test suite.
AWS CLI was used to stub endpoints and verify S3-compatible API behavior during exploration.

The goal is to validate core S3 API operations against a live RustFS instance — not exhaustive coverage, but the critical CRUD data flows and bucket administration operations that matter most.
User and group management is intentionally offloaded to the Playwright GUI test project.

## Coverage

- Bucket administration: create, delete, list, lifecycle
- Object operations: upload, download, delete, list, lifecycle
- Negative paths: non-existent resources, duplicate creation, constraint violations

### Limitations

The items below represent areas to extend coverage. They were identified through exploratory testing and general knowledge of S3-compatible storage solutions — not from deep review of RustFS documentation.

- No multipart upload testing (files >5MB use a different S3 code path)
- No presigned URL generation and expiry validation
- No bucket versioning flows
- No access control validation — users are created via GUI tests but policy enforcement is not verified
- No concurrent write testing
- No edge cases for special characters or long names in bucket/object keys

## Project structure

```
test/e2e/
├── cases/           # test case functions, one file per domain
├── helpers/         # AWS SDK client setup and utility functions
├── configuration/   # environment loading and config struct
└── main_test.go     # test suite orchestration
```

**`cases/`** — each file groups test case functions by domain (bucket, object). Functions are self-contained and can be composed freely into suites. This allows granular targeting (`-run TestBucketHappyFlowSuite`) and clear separation of concerns.

**`helpers/`** — wraps the AWS SDK into reusable Go functions. Cases call helpers rather than the SDK directly, keeping test logic readable and SDK details isolated.

**`configuration/`** — handles `.env` loading, environment variable resolution with defaults, and exposes a single `Config` struct. Defaults target a standard local RustFS setup — no configuration required for the common case.

**`main_test.go`** — orchestrates which case functions run in which suite. Adding a new test case is: write a function in `cases/`, register it with `t.Run(...)` in the relevant suite. No framework boilerplate.

This structure allows flexibility in grouping by domain, component, business area, or functional purpose — suites can be reorganized without touching the case functions themselves.
