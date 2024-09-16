SSO service contracts repo.

ProtoBuf install:

[windows]
1. Download release from https://github.com/protocolbuffers/protobuf/releases
2. Unzip into any `folder`
3. Add `folder` to `$PATH`

[linux]
```bash
apt install -y protobuf-compiler
```

ProtoBuf plugins 4 Golang:
```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```