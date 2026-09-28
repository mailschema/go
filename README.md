# MailSchema for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/mailschema/go.svg)](https://pkg.go.dev/github.com/mailschema/go)
[![CI](https://github.com/mailschema/go/actions/workflows/test.yml/badge.svg)](https://github.com/mailschema/go/actions/workflows/test.yml)

The Mail Action Protocol 0.2 core artifacts and the MailSchema Registry contribution schema, embedded in Go.

[Specification](https://mailschema.org/specification) · [Registry](https://mailschema.org/registry) · [Tools](https://mailschema.org/tools) · [Source](https://github.com/mailschema/go)

## Install

```sh
go get github.com/mailschema/go@v0.1.2
```

## Use an artifact

```go
schema, err := mailschema.Schema(mailschema.MAPSchema)
if err != nil {
	panic(err)
}
```

`MAPSchema`, `MAPContext`, `ContractFormatSchema` and `FormsSchema` return independent copies of the files the [profile record](https://mailschema.org/profiles/map/0.2.json) binds by SHA-256, byte for byte. `ContributionSchema` describes Registry contributions. Use the schemas with a Draft 2020-12 validator. Type contracts are not bundled: a client obtains them from the [Registry catalogue](https://mailschema.org/registry/catalog.json) by digest. Processing MAP messages is outside this module; see the [profile](https://mailschema.org/specification/profile).

## Trust boundary

A valid document is structured input. Schema validation does not authenticate a service, grant authority or establish product conformance. The module performs no network requests. Module versions and MAP profile versions advance independently. MIT licensed.
