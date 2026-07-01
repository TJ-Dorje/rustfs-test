# RustFS Test Framework — Demo

> **This is a demo test framework.** It is not production-ready and is intended to demonstrate an approach to testing a RustFS deployment — covering the S3 API and the admin console UI. Use it as a starting point and extend to fit your needs.

[RustFS](https://rustfs.com) is a Rust-based S3-compatible object storage. This framework exercises it across two test modules:

| Module | Location | What it tests | Design |
|--------|----------|---------------|--------|
| **e2e** | `test/e2e/` | S3 API via AWS SDK (buckets, objects) | [DESIGN.md](test/e2e/DESIGN.md) |
| **playwright** | `test/playwright-test/` | Admin console GUI (login, buckets, users, groups) | [DESIGN.md](test/playwright-test/DESIGN.md) |

## Design & architecture decisions

- [E2E test project design](test/e2e/DESIGN.md)
- [Playwright test project design](test/playwright-test/DESIGN.md)

---

# Prerequisites

| Tool | Purpose | Install |
|------|---------|---------|
| Go 1.23+ | Build and run tests | https://go.dev/dl |
| Task | Task runner | `brew install go-task` or https://taskfile.dev |
| gotestsum | Test reporting + JUnit XML | `task setup` |
| Playwright browsers | Chromium for GUI tests | `task setup` |
| Docker & Docker Compose | Run RustFS + Docker delivery | https://docs.docker.com/get-docker |

> Run `task setup` once after cloning — installs `gotestsum` and Playwright browser binaries.

---

# Option 1: Raw Go commands

Requires RustFS running locally (`docker compose up`).

```shell
# e2e — run all suites
cd test/e2e && go test -count=1 . -v

# e2e — run specific suite
cd test/e2e && go test -count=1 . -v -run TestBucketHappyFlowSuite

# e2e — run specific test case
cd test/e2e && go test -count=1 . -v -run TestBucketHappyFlowSuite/BucketCreation

# playwright — run all suites
cd test/playwright-test && go test -count=1 . -v

# playwright — run specific suite
cd test/playwright-test && go test -count=1 . -v -run TestUserManagementHappyFlowSuite

# playwright — run specific test case
cd test/playwright-test && go test -count=1 . -v -run TestUserManagementHappyFlowSuite/UserCreation
```

Default configuration targets `localhost:9000` / `localhost:9001` with credentials `rustfsadmin`/`rustfsadmin` — no `.env` needed if running RustFS with defaults. To override, copy `.env.example` to `.env` in the module directory and adjust.

---

# Option 2: Taskfile

Requires RustFS running locally (`docker compose up`).

### Install Task (if not already installed):

```shell
brew install go-task/tap/go-task
```

### One-time setup:

```shell
task setup
```

### Run tests:

```shell
task e2e-test              # all e2e collection
task pw-test               # all playwright collection
task test:report           # both collections + JUnit XML in created reports/ directory
task clean                 # cleans reports dir
```

### Run specific suites:

```shell
task e2e-test:bucket       # bucket suites only
task e2e-test:object       # object suites only
task pw-test:login         # login suites only
task pw-test:bucket        # bucket management suites
task pw-test:user          # user management suites
```

### Run without teardown (inspect created resources after tests):

```shell
task e2e-test:no-teardown
task pw-test:no-teardown
```

List all available tasks:

```shell
task
```

---

# Option 3: Docker Compose

Self-contained — starts RustFS, runs all tests, stops everything. No local Go or RustFS required.

```shell
task docker-test
```

Starts RustFS, runs all tests. Containers remain stopped after for inspection via Docker Desktop or logs. Reports written to `reports/e2e.xml` and `reports/pw.xml` on the host.

To clean up after (must use `-f` to target the test compose):

```shell
docker compose -f docker-compose.test.yml down -v
```

---

# Environment variables

Both modules read configuration from a `.env` file in their directory. Values can also be set as environment variables (env vars take precedence over `.env`).

### `test/e2e/.env`

| Variable | Default | Description |
|----------|---------|-------------|
| `RUSTFS_ENDPOINT` | `http://localhost:9000` | S3 API endpoint |
| `AWS_ACCESS_KEY_ID` | `rustfsadmin` | Access key |
| `AWS_SECRET_ACCESS_KEY` | `rustfsadmin` | Secret key |
| `AWS_DEFAULT_REGION` | `us-east-1` | Region |
| `TEARDOWN_ENABLED` | `true` | Delete created resources after tests |

### `test/playwright-test/.env`

| Variable | Default | Description |
|----------|---------|-------------|
| `RUSTFS_CONSOLE_URL` | `http://localhost:9001` | Admin console URL |
| `RUSTFS_USERNAME` | `rustfsadmin` | Admin username |
| `RUSTFS_PASSWORD` | `rustfsadmin` | Admin password |
| `HEADLESS` | `true` | Run browser headless (`false` shows browser) |
| `NO_SANDBOX` | `false` | Disable Chromium sandbox (required in Docker/CI) |
| `TEARDOWN_ENABLED` | `true` | Delete created resources after tests |

---

# Test suites

### E2E (`test/e2e/`)

| Suite | `-run` flag | Cases |
|-------|-------------|-------|
| Bucket happy flow | `TestBucketHappyFlowSuite` | Create, delete, lifecycle, list |
| Bucket negative flow | `TestBucketNegativeFlowSuite` | Duplicate creation, delete non-empty, non-existent |
| Object happy flow | `TestObjectHappyFlowSuite` | Upload, download, delete, list, lifecycle |
| Object negative flow | `TestObjectNegativeFlowSuite` | Non-existent bucket/object, invalid operations |

### Playwright (`test/playwright-test/`)

| Suite | `-run` flag | Cases |
|-------|-------------|-------|
| Login happy flow | `TestLoginHappyFlowSuite` | Valid credentials |
| Login negative flow | `TestLoginNegativeFlowSuite` | Invalid credentials |
| Bucket management | `TestBucketManagementHappyFlowSuite` | Create, delete, lifecycle via UI |
| User management | `TestUserManagementHappyFlowSuite` | Create user, create group, add user to group |

---

# Adding new tests

1. Add a case function in `test/e2e/cases/` or `test/playwright-test/cases/`
2. Register it in the relevant `TestXxxSuite` in `main_test.go`
3. Add a Taskfile task in `Taskfile.yml` if a dedicated run target is useful

Helper packages:
- `test/e2e/helpers/` — S3 client setup, random name generators
- `test/playwright-test/helpers/actions/` — page-level actions (login, bucket, user, group)
- `test/playwright-test/helpers/components/` — reusable UI component interactions (dialogs, tables, toasts)

---

# CI

GitHub Actions workflow runs on every push and pull request to `main`.
RustFS runs as a Docker service; tests run natively on the runner.
JUnit XML reports are uploaded as workflow artifacts after each run.

See `.github/workflows/tests.yml`.
