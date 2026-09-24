package mailschema

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestBundledSchemas(t *testing.T) {
	for _, name := range []SchemaName{ContributionSchema, MAP01Schema, ContentReview01Schema, ContentReview02Schema} {
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
	var contract map[string]any
	if err := json.Unmarshal(ContentReview01Contract(), &contract); err != nil {
		t.Fatal(err)
	}
	if contract["id"] != ContentReviewType {
		t.Fatalf("contract id = %q", contract["id"])
	}
	if err := json.Unmarshal(ContentReview02Contract(), &contract); err != nil {
		t.Fatal(err)
	}
	if contract["version"] != "0.2" {
		t.Fatalf("current contract version = %q", contract["version"])
	}
}

func TestResultCarriesExactTypeReference(t *testing.T) {
	fixture, err := os.Open("testdata/map-result.json")
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()

	result, err := Decode[Result](fixture)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateResult(result); err != nil {
		t.Fatal(err)
	}
	if result.Type.ID != ContentReviewType {
		t.Fatalf("result type = %q", result.Type.ID)
	}
}

func TestProblemCorrelation(t *testing.T) {
	problem := Problem{
		Type:          "https://mailschema.org/problems/authentication-required",
		Title:         "Authentication required",
		Status:        401,
		Detail:        "Authenticate before executing this action.",
		Instance:      "https://reviews.example/map/results/018f47a2-b4d3-7c02-b491-7bdf2eaac67c",
		Profile:       Profile01,
		RequestID:     "urn:uuid:018f47a2-b4d3-7c02-b491-7bdf2eaac67c",
		InteractionID: "urn:uuid:018f47a2-5d7c-7b11-9a3d-4d2160b85b10",
		Code:          "authentication-required",
	}
	if err := ValidateProblem(problem); err != nil {
		t.Fatal(err)
	}
	problem.Status = 403
	if err := ValidateProblem(problem); err == nil {
		t.Fatal("contradictory problem status was accepted")
	}
}

func TestStrictDecodeAndValidation(t *testing.T) {
	request := `{
		"kind":"MapRequest",
		"profile":"https://mailschema.org/profiles/map/0.1",
		"requestId":"urn:uuid:018f47a2-b4d3-7c02-b491-7bdf2eaac67c",
		"interactionId":"urn:uuid:018f47a2-5d7c-7b11-9a3d-4d2160b85b10",
		"type":{"id":"https://mailschema.org/types/content-review","version":"0.2","contractDigest":"sha-256:6ee5b086138799f70f2489e851e1db6e11ff91cd1141bd36fe202a47d8c1e637"},
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

func TestResultNotFoundDoesNotInventInteraction(t *testing.T) {
	problem := Problem{
		Type:      "https://mailschema.org/problems/result-not-found",
		Title:     "Result not found",
		Status:    404,
		Detail:    "No retained result exists for this request identifier.",
		Instance:  "https://reviews.example/map/results/urn%3Auuid%3A018f47a2-b4d3-7c02-b491-7bdf2eaac699",
		Profile:   Profile01,
		RequestID: "urn:uuid:018f47a2-b4d3-7c02-b491-7bdf2eaac699",
		Code:      "result-not-found",
	}
	if err := ValidateProblem(problem); err != nil {
		t.Fatal(err)
	}
}
