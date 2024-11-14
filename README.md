# SSO service repo

## Docs

* [CHANGELOG](docs/CHANGELOG.md)
* [CONTRIBUTING](docs/CONTRIBUTING.md)
* [RELEASING](docs/RELEASING.md)

## Environment

Create `.env` file in root dir of the project or assign in your OS env. Lookup [example](docs/.env).

## Tools

* [Taskfile](https://taskfile.dev/)
* Linter:

```bash
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
```

* Swagger:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Services

[Authentication](api/rest/v1/authentication/authentication.md)
