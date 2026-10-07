// Package mailschema carries the Mail Action Protocol 0.3 artifacts byte for byte as
// mailschema.org publishes them: the profile record, the JSON-LD context, and the core, type
// contract and implementation record schemas. The profile record binds the context and the
// first two schemas by SHA-256.
package mailschema

import (
	"embed"
	"fmt"
)

// Profile is the MAP profile these artifacts define.
const Profile = "https://mailschema.org/profiles/map/0.3"

// Context is the JSON-LD context every MAP 0.3 description names.
const Context = "https://mailschema.org/contexts/map-0.3.jsonld"

// Artifact names a bundled file.
type Artifact string

const (
	ProfileRecord        Artifact = "profile.json"
	ContextDocument      Artifact = "context.jsonld"
	CoreSchema           Artifact = "core.schema.json"
	ContractSchema       Artifact = "contract.schema.json"
	ImplementationSchema Artifact = "implementation.schema.json"
)

//go:embed artifacts/*
var files embed.FS

// Bytes returns an independent copy of a bundled artifact's exact bytes.
func Bytes(name Artifact) ([]byte, error) {
	contents, err := files.ReadFile("artifacts/" + string(name))
	if err != nil {
		return nil, fmt.Errorf("mailschema: unknown artifact %q", name)
	}
	return append([]byte(nil), contents...), nil
}
