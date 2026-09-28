// Package mailschema carries the Mail Action Protocol 0.2 core artifacts and the
// MailSchema Registry contribution schema, byte for byte as mailschema.org publishes them.
//
// Type contracts are not bundled: a client obtains them from the Registry catalogue
// by digest. Processing MAP messages is outside this module.
package mailschema

import (
	"embed"
	"fmt"
)

// Profile is the MAP profile whose core artifacts this module carries.
const Profile = "https://mailschema.org/profiles/map/0.2"

// SchemaName identifies an artifact bundled with this module.
type SchemaName string

const (
	// MAPSchema is the MAP 0.2 core schema: descriptions, requests, results and problems.
	MAPSchema SchemaName = "map-0.2"
	// MAPContext is the MAP 0.2 JSON-LD context.
	MAPContext SchemaName = "map-0.2-context"
	// ContractFormatSchema is the type contract format every MAP 0.2 contract follows.
	ContractFormatSchema SchemaName = "type-contract-0.2"
	// FormsSchema is the form fields block contracts pin.
	FormsSchema SchemaName = "forms-0.1"
	// ContributionSchema describes Registry contributions and expanded type records.
	ContributionSchema SchemaName = "contribution"
)

//go:embed schemas/*.json contexts/*.jsonld
var files embed.FS

var paths = map[SchemaName]string{
	MAPSchema:            "schemas/map-0.2.schema.json",
	MAPContext:           "contexts/map-0.2.jsonld",
	ContractFormatSchema: "schemas/type-contract-0.2.schema.json",
	FormsSchema:          "schemas/forms-0.1.schema.json",
	ContributionSchema:   "schemas/contribution.schema.json",
}

// Schema returns an independent copy of a bundled artifact's exact bytes.
func Schema(name SchemaName) ([]byte, error) {
	path, ok := paths[name]
	if !ok {
		return nil, fmt.Errorf("mailschema: unknown schema %q", name)
	}
	contents, err := files.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("mailschema: read %s: %w", name, err)
	}
	return append([]byte(nil), contents...), nil
}
