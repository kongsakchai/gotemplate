# Makefile reference

## Commands

| Command | Description |
|---------|-------------|
| `make test` | Run all tests with colorized output |
| `make testcover` | Run all tests with coverage |
| `make coverage` | Generate coverage report and open HTML |
| `make init` | Install mockery v3.7.3 + go-swagger |
| `make gen-mock` | Generate mocks via mockery |
| `make gendocs` | Generate swagger spec to `./docs/swagger.yaml` |
| `make docs` | Serve swagger UI from `./docs/swagger.yaml` |

## Notes

- Test output is piped through `.script/colorize` for readability
- Init is only needed once per clone; tools are installed globally via `go install`
- Mockery config is read from `.mockery.yaml` at project root
- Swagger docs generated with `--tags=docs` — only endpoints tagged `docs` appear
