# Go Template

## Usage

- Install `gonew`

```sh
go install golang.org/x/tools/cmd/gonew@latest
```

- Create a new project

```sh
gonew github.com/kongsakchai/gotemplate github.com/yourname/projectname
```

## Project structure

```sh
./
├── app
├── cache
├── config
├── database
├── errs
├── httpclient
├── logger
├── middleware
├── migrations
└── validator
```

- **app** Application layer and business logic. Contains core application logic, each subdirectory represents a business domain or feature.

```sh
./
└── app
    ├── user
    └── admin
        ├── admin.go
        └── handler.go
```

- **cache** Cache layer connectors and helpers, such as Redis.
  Used for caching data, session storage, or performance optimization.

- **config** Application configuration management.
  Includes environment variables, configuration structs, and config loaders.

- **database** Database connectors and setup, e.g., MySQL or PostgreSQL.

- **errs** Custom error types and centralized error handling for error tracking.

- **httpclient** HTTP client utilities for calling external services or APIs.

- **logger** Logging configuration and shared logger instances.
  Ensures consistent logging format and levels across the application.

- **middleware** HTTP middleware components for request processing, such as authentication, authorization, logging, and request tracing.

- **validator** Request validation logic. Typically integrates with libraries such as go-playground/validator to validate incoming payloads.

- **migrations** Database migration files (.sql) for schema versioning and changes.
  Compatible with tools such as kongsakchai/simple-migrate.

---

## Guidline for this template

### package `app/`

business logic และ application layer ต่าง ๆ ควรอยู่ใน package นี้, และควรแบ่งแยกแต่ละ module ให้ชัดเจน example:

- app/register
- app/booking
- app/product

> [!CAUTION]
> ไม่ควรเขียน business logic ใน package อื่น ๆ ควรอยู่ใน package `app/` เท่านั้น
