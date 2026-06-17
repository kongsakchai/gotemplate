---
name: gotemplate-guideline
description: >
    Go + Echo v5 + sqlx REST API coding conventions for this project. Load this
    skill when writing or editing handlers, tests, middleware, routing, database
    queries, or domain modules — then read the relevant sub-skill below.
---

# Sub-skills

| When you are...                                 | Read this file                                       |
| ----------------------------------------------- | ---------------------------------------------------- |
| Creating a new domain or module, wiring, config | [`sub/architecture.md`](./sub/architecture.md)       |
| Writing handlers, services, or API responses    | [`sub/handler-service.md`](./sub/handler-service.md) |
| Writing any Go code (pointers, errors, logging) | [`sub/coding.md`](./sub/coding.md)                   |
| Writing or editing tests                        | [`sub/testing.md`](./sub/testing.md)                 |

## Makefile reference

### Commands

| Command          | Description                                    |
| ---------------- | ---------------------------------------------- |
| `make test`      | Run all tests with colorized output            |
| `make testcover` | Run all tests with coverage                    |
| `make coverage`  | Generate coverage report and open HTML         |
| `make init`      | Install mockery v3.7.3 + go-swagger            |
| `make gen-mock`  | Generate mocks via mockery                     |
| `make gendocs`   | Generate swagger spec to `./docs/swagger.yaml` |
| `make docs`      | Serve swagger UI from `./docs/swagger.yaml`    |

### Notes

- Test output is piped through `.script/colorize` for readability
- Init is only needed once per clone; tools are installed globally via `go install`
- Mockery config is read from `.mockery.yaml` at project root
- Swagger docs generated with `--tags=docs` — only endpoints tagged `docs` appear
