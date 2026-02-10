# Go Template

## Structure

```sh
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

### app

Application Layer and business logic. It contains the application logic and is organized into subdirectories for business domains or features.

#### Example

```sh
└── app
    ├── user
    └── admin
        ├── admin.go
        └── handler.go
```

### cache

Cache connector, e.g., Redis

### config

Configuration files and environment variables

### database

Database connector, e.g., MySQL, PostgreSQL

### errs

Custom error types and error handling

### httpclient

HTTP Client for call external service

### logger

Logger setup and configuration

### middleware

Middleware for request processing, e.g., authentication, logging

### validator

Validation logic for request data, e.g., using `go-playground/validator`

### migrations

Migration files `.sql` for database schema changes, e.g., using `kongsakchai/simple-migrate`

---

## Usage

### Install `gonew`

```sh
go install golang.org/x/tools/cmd/gonew@latest
```

### Create a new project

```sh
gonew github.com/kongsakchai/gotemplate github.com/yourname/projectname
```
