package mailschema

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBundledSchemas(t *testing.T) {
	for _, name := range []SchemaName{ContributionSchema, MAP01Schema, ContentReview01Schema} {
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
	}
}

func TestStrictDecodeAndValidation(t *testing.T) {
	request := `{
		"kind":"MapRequest",
		"profile":"https://mailschema.org/profiles/map/0.1",
		"requestId":"urn:uuid:018f47a2-b4d3-7c02-b491-7bdf2eaac67c",
		"interactionId":"urn:uuid:018f47a2-5d7c-7b11-9a3d-4d2160b85b10",
		"requestedAt":"2026-09-23T01:06:00Z",
		"type":{"id":"https://mailschema.org/types/content-review","version":"0.1","recordDigest":"sha-256:0e6365df1bf904f2475f972ec66a5561edc0bfaae69fd7a1b41f53b40bf8ac1c"},
		"operation":"approve",
		"target":{"id":"https://reviews.example/content/campaign-42","revision":"4","digest":"sha-256:9096ce5e9226696d3bdc80f7d62a8e1ca4995c4fd8c4b9f0155287529aefe8fd"},
		"input":{}
	}`
	value, err := Decode[Request](strings.NewReader(request))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateRequest(value); err != nil {
		t.Fatal(err)
	}

	_, err = Decode[Request](strings.NewReader(strings.Replace(request, `"input":{}`, `"input":{},"token":"secret"`, 1)))
	if err == nil {
		t.Fatal("unknown credential field was accepted")
	}
}
