# Playwright Test Project — Design

## Why Playwright

Playwright is widely adopted, straightforward to integrate into any tech stack, and supports multiple browsers.
A key advantage for this project: Playwright's MCP integration allows an AI agent to crawl the UI,
identify CSS selectors and user flows, and produce test code rapidly — significantly reducing the time to write reliable UI tests against an unfamiliar interface.

## Coverage

- Login: valid credentials, invalid credentials
- Bucket management: create, delete, lifecycle via UI
- User management: create user, create group, add user to group

S3 API operations are covered by the e2e project. This project focuses on admin console journeys only.

### Limitations

The items below represent areas to extend coverage. They were identified through exploratory testing and general knowledge of S3-compatible storage admin consoles — not from deep review of RustFS documentation.

- No screenshot capture on test failure
- No access control validation — user is created but login as that user is not verified
- No user removal from group
- No policy assignment to users or groups
- No CORS or bucket policy configuration via UI

## Approach

The structure does not implement the Screenplay pattern directly, but follows an **action-based approach** as a practical approximation.

RustFS is typically operated by a small set of admin users — DevOps engineers, on-call support — rather than a broad user base with diverse roles.
Testing from an action perspective maps naturally to this: validate the journeys an admin actually performs, e.g. managing users, creating buckets, configuring groups.

## Project structure

```
test/playwright-test/
├── cases/                    # test case functions, one file per user journey
├── helpers/
│   ├── actions/              # user-level workflows (login, bucket, user, group)
│   └── components/           # UI component interactions (dialogs, tables, toasts)
├── configuration/            # browser setup, environment loading, config struct
└── main_test.go              # test suite orchestration
```

**`cases/`** — each file covers a user journey (login, bucket management, user management). Functions are composed into suites in `main_test.go`.

**`helpers/actions/`** — page-level workflows a user can perform: navigate to a section, create a resource, delete a resource. Actions are the building blocks cases are assembled from.

**`helpers/components/`** — reusable interactions with specific UI components: confirming dialogs, waiting for table rows, reading toast notifications. Components are called by actions, keeping selector knowledge in one place.

**`configuration/`** — initializes and manages the Playwright browser instance, loads environment variables, exposes a `Config` struct. `HEADLESS`, `NO_SANDBOX`, and `TEARDOWN_ENABLED`
flags allow the same code to run locally with a visible browser, in Docker, and in CI without modification.
