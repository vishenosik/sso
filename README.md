# SSO service repo.


## Tools to install

### Linter
```bash
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
```

### ProtoBuf:

#### Windows
1. Download release from https://github.com/protocolbuffers/protobuf/releases
2. Unzip into any `folder`
3. Add `folder` to `$PATH`

#### Linux
```bash
apt install -y protobuf-compiler
```

#### Protobuf golang plugins:
```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

### [CHANGELOG](CHANGELOG.md)
## Services

[Authentication](api/rest/v1/authentication/authentication.md)