# MailSchema for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/mailschema/go.svg)](https://pkg.go.dev/github.com/mailschema/go)
[![CI](https://github.com/mailschema/go/actions/workflows/test.yml/badge.svg)](https://github.com/mailschema/go/actions/workflows/test.yml)

Decode typed Mail Action Protocol documents and use the canonical MailSchema JSON Schemas in Go.

[Specification](https://mailschema.org/specification/) · [Registry](https://mailschema.org/registry/) · [Tools](https://mailschema.org/tools/) · [Source](https://github.com/mailschema/go)

## Install

```sh
go get github.com/mailschema/go@v0.1.1
```

## Decode and check a request

```go
package main

import (
	"os"

	mailschema "github.com/mailschema/go"
)

func main() {
	request, err := mailschema.Decode[mailschema.Request](os.Stdin)
	if err != nil {
		panic(err)
	}
	if err := mailschema.ValidateRequest(request); err != nil {
		panic(err)
	}
}
```

`Description`, `Request`, `Result` and `Problem` model MAP 0.1 documents. `ValidateDescription`, `ValidateRequest`, `ValidateResult` and `ValidateProblem` check the fixed profile identifiers and core references while decoding rejects unknown fields and trailing JSON.

## Use the schemas

```go
schema, err := mailschema.Schema(mailschema.MAP01Schema)
if err != nil {
	panic(err)
}
```

`MAP01Schema`, `ContentReview01Schema`, `ContentReview02Schema` and `ContributionSchema` expose independent copies of the bundled Draft 2020-12 schemas. Use them with a Draft 2020-12 validator when complete schema validation is required.

`ContentReview01Contract()` and `ContentReview02Contract()` return immutable canonical type contracts. New integrations should use Content Review 0.2.

## Trust boundary

A valid document is structured input. Decoding and validation do not authenticate a service, grant authority, approve an action or establish product conformance. Implementations must apply their own endpoint trust, credentials, permissions and policy before executing a request.

The package performs no network requests. Module versions and MAP profile versions advance independently. MIT licensed.
