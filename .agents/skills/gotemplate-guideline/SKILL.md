---
name: gotemplate-guideline
description: "Conventions and architecture for the gotemplate Go project. Use when writing, reviewing, or refactoring any Go code in this codebase. Triggers on: internal/{domain}, app/{domain}app, consumer/{domain}consumer, service_{core}.go, adapter_{3party}.go, serror, handleError, app.Error, app.ErrorMiddleware, app.Request, EchoApp.RegisterRoute, mockery, go-sqlmock, echotest, go-swagger/docs annotations, pkg/config env tags, request/response envelope {code, success, message, data}. Use ONLY for this repository's Go code, not for general-purpose Go tips."
license: MIT
metadata:
  version: '1.0.0'
---

# Gotemplate Guidelines

Format that goes with the `gotemplate` project layout (`internal/`, `app/`, `consumer/`, `pkg/`). Follow these conventions when adding or changing Go code in this repo.

## Quick Reference

| Topic | When to Use | Reference |
|-------|-------------|-----------|
| **Architecture** | Project layout, domain naming, file prefixes, service file splitting | [architecture.md](references/architecture.md) |
| **App / Handler** | Building handlers, route registration, request validation, responses | [app-handler.md](references/app-handler.md) |
| **Middleware** | Global middleware stack, traceID/tag, JWT claims in context | [middleware.md](references/middleware.md) |
| **Errors** | serror, app.Error, error codes, per-module handleError | [errors.md](references/errors.md) |
| **Testing** | mockery mocks, go-sqlmock, echotest handler tests, testify | [testing.md](references/testing.md) |
| **Code Review** | Checking existing code against these conventions | [review.md](references/review.md) |
| **Tooling** | pkg/ shared infra, config env tags, migrate, makefile, swagger | [tooling.md](references/tooling.md) |

## Rules That Never Change

- Return concrete types, never interfaces. Check existence by returning `bool`, not pointers.
- Wrap 3rd-party/external errors with `serror` so they can be traced; define business codes with `serror.NewCoded`.
- Register error handling via `app.ErrorMiddleware` on the group, never per handler.
- Common code goes in `app/`; business code goes in each module; reusable infra goes in `pkg/`.
- Log through the echo context logger (`ctx.Logger()`), which already carries traceID/tag.