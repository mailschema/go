# mailschema

The [Mail Action Protocol](https://mailschema.org) 0.3 artifacts for Go, embedded byte for byte as mailschema.org publishes them: the profile record, the JSON-LD context, and the core, type contract and implementation record schemas. The profile record binds the context and the first two schemas by SHA-256.

```sh
go get github.com/mailschema/go@v0.3.0
```

```go
schema, err := mailschema.Bytes(mailschema.CoreSchema)
```

`Bytes` returns an independent copy of `ProfileRecord`, `ContextDocument`, `CoreSchema`, `ContractSchema` or `ImplementationSchema`. Schema validation alone does not establish a valid description; see the [specification](https://mailschema.org/specification/core).

Source: [mailschema/go](https://github.com/mailschema/go). License: MIT.
