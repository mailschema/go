package mailschema

import (
	"encoding/json"
	"testing"
)

func TestBundledArtifacts(t *testing.T) {
	ids := map[SchemaName]string{
		MAPSchema:            "https://mailschema.org/schemas/map-0.2.schema.json",
		ContractFormatSchema: "https://mailschema.org/schemas/type-contract-0.2.schema.json",
		FormsSchema:          "https://mailschema.org/schemas/forms-0.1.schema.json",
		ContributionSchema:   "",
	}
	for name, id := range ids {
		contents, err := Schema(name)
		if err != nil {
			t.Fatal(err)
		}
		var document map[string]any
		if err := json.Unmarshal(contents, &document); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if document["$schema"] != "https://json-schema.org/draft/2020-12/schema" {
			t.Fatalf("%s: unexpected dialect", name)
		}
		if id != "" && document["$id"] != id {
			t.Fatalf("%s: $id = %q", name, document["$id"])
		}
	}
	contents, err := Schema(MAPContext)
	if err != nil {
		t.Fatal(err)
	}
	var context struct {
		Context map[string]any `json:"@context"`
	}
	if err := json.Unmarshal(contents, &context); err != nil {
		t.Fatal(err)
	}
	if context.Context["MailAction"] != "map:MailAction" {
		t.Fatalf("context MailAction = %v", context.Context["MailAction"])
	}
	if Profile != "https://mailschema.org/profiles/map/0.2" {
		t.Fatalf("Profile = %q", Profile)
	}
}

func TestCopiesAreIndependent(t *testing.T) {
	first, _ := Schema(MAPSchema)
	first[0] = 'x'
	second, _ := Schema(MAPSchema)
	if second[0] == 'x' {
		t.Fatal("Schema returned shared bytes")
	}
	if _, err := Schema("unknown"); err == nil {
		t.Fatal("unknown schema accepted")
	}
}
