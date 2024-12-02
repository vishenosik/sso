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

### Tabnine

1. generate swagger comment for the `register` function including tags from github.com/swaggo/swag library.

## Services

[Authentication](api/rest/v1/authentication/authentication.md)

## Benchmarking

```bash
# benchmark running into file
go test -bench . -benchmem ./internal/lib/collections -cpuprofile=profile.out
# benchmark pdf version
go tool pprof --pdf profile.out > pprof.pdf 
```
