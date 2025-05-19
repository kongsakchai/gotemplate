# Go Template

> inspired by [gotemplate](https://github.com/pallat/gotemplate)

## Installation

### Install `gonew`

```bash
go install golang.org/x/tools/cmd/gonew@latest
```

### New Project

```bash
gonew github.com/kongsakchai/gotemplate github.com/yourname/projectname
```

## Structure

```painttext
├── app
│   └── todo
├── cache
├── config
├── database
├── error
├── example
│   ├── echo
│   └── gin
├── logger
├── middleware
├── pkg
│   └── generate
└── test
    └── appmock
```

- `app`: Application layer and business logic.
- `cache`: Cache layer.
- `config`: Configuration files.
- `database`: Database connection.
- `error`: Custom error package.
- `example`: Example usage of the framework.
- `logger`: Logger package.
- `middleware`: Middleware package.
- `pkg`: Shared packages such as utils, helpers, etc.
- `test`: Test files and mock data for unit test.
