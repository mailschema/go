package mailschema

import (
	"embed"
	"fmt"
)

// SchemaName identifies a schema bundled with this module.
type SchemaName string

const (
	ContributionSchema    SchemaName = "contribution"
	MAP01Schema           SchemaName = "map-0.1"
	ContentReview01Schema SchemaName = "content-review-0.1"
)

//go:embed schemas/*.json
var schemaFiles embed.FS

//go:embed contracts/content-review-0.1.json
var contentReviewContract []byte

var schemaPaths = map[SchemaName]string{
	ContributionSchema:    "schemas/contribution.schema.json",
	MAP01Schema:           "schemas/map-0.1.schema.json",
	ContentReview01Schema: "schemas/content-review-0.1.schema.json",
}

// Schema returns an independent copy of a bundled JSON Schema document.
func Schema(name SchemaName) ([]byte, error) {
	path, ok := schemaPaths[name]
	if !ok {
		return nil, fmt.Errorf("mailschema: unknown schema %q", name)
	}
	contents, err := schemaFiles.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("mailschema: read %s: %w", name, err)
	}
	return append([]byte(nil), contents...), nil
}

// ContentReview01Contract returns an independent copy of the canonical type contract.
func ContentReview01Contract() []byte {
	return append([]byte(nil), contentReviewContract...)
}
