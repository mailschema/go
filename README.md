# MailSchema for Go

Typed Mail Action Protocol 0.1 documents and canonical JSON Schemas for Go.

```sh
go get github.com/mailschema/go@v0.1.0
```

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

`Schema(mailschema.MAP01Schema)`, `Schema(mailschema.ContentReview01Schema)` and `Schema(mailschema.ContributionSchema)` return independent copies of the bundled schemas. The validation helpers check the fixed MAP identifiers and core references; use the bundled Draft 2020-12 schemas when complete schema validation is required.

The package performs no network requests. It does not send email, establish endpoint trust or grant service authorization.

MIT licensed.
