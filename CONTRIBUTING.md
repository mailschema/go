# Contributing

This repository publishes the Go implementation for MailSchema. The protocol schemas are maintained in [`mailschema/mailschema`](https://github.com/mailschema/mailschema); Go-specific types, decoding, validation and packaging changes belong here and must remain compatible with those canonical files.

Run the package checks before opening a pull request:

```sh
go test ./...
go vet ./...
```

Changes to MAP behavior, shared schemas or Registry records should begin in the [main project repository](https://github.com/mailschema/mailschema/blob/main/CONTRIBUTING.md).
