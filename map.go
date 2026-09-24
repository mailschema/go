package mailschema

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	Profile01         = "https://mailschema.org/profiles/map/0.1"
	Context01         = "https://mailschema.org/contexts/map-0.1.jsonld"
	ContentReviewType = "https://mailschema.org/types/content-review"
)

var (
	uuidURN = regexp.MustCompile(`^urn:uuid:[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	digest  = regexp.MustCompile(`^sha-256:[a-f0-9]{64}$`)
)

type TypeReference struct {
	ID             string `json:"id"`
	Version        string `json:"version"`
	ContractDigest string `json:"contractDigest"`
}

type Target struct {
	ID       string `json:"id"`
	Revision string `json:"revision"`
	Title    string `json:"title,omitempty"`
	Digest   string `json:"digest"`
}

type Operation struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema string `json:"inputSchema"`
}

type Authorization struct {
	Kind     string   `json:"kind"`
	Schemes  []string `json:"schemes"`
	Audience string   `json:"audience,omitempty"`
}

type Execution struct {
	URL                    string `json:"url"`
	Method                 string `json:"method"`
	RequestMediaType       string `json:"requestMediaType"`
	ResultMediaType        string `json:"resultMediaType"`
	ResultURLTemplate      string `json:"resultUrlTemplate"`
	ResultRetentionSeconds int    `json:"resultRetentionSeconds"`
}

type Service struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Execution     Execution     `json:"execution"`
	HumanURL      string        `json:"humanUrl"`
	Authorization Authorization `json:"authorization"`
}

type Description struct {
	Context     string        `json:"@context"`
	Kind        string        `json:"@type"`
	ID          string        `json:"@id"`
	Profile     string        `json:"profile"`
	Type        TypeReference `json:"type"`
	DescribedAt time.Time     `json:"describedAt"`
	ExpiresAt   time.Time     `json:"expiresAt"`
	Service     Service       `json:"service"`
	Target      Target        `json:"target"`
	Operations  []Operation   `json:"operations"`
}

type Request struct {
	Kind          string         `json:"kind"`
	Profile       string         `json:"profile"`
	RequestID     string         `json:"requestId"`
	InteractionID string         `json:"interactionId"`
	Type          TypeReference  `json:"type"`
	Operation     string         `json:"operation"`
	Target        Target         `json:"target"`
	Input         map[string]any `json:"input"`
}

type Result struct {
	Kind          string         `json:"kind"`
	Profile       string         `json:"profile"`
	RequestID     string         `json:"requestId"`
	InteractionID string         `json:"interactionId"`
	Type          TypeReference  `json:"type"`
	Operation     string         `json:"operation"`
	State         string         `json:"state"`
	Target        Target         `json:"target"`
	RecordedAt    time.Time      `json:"recordedAt"`
	ResultURL     string         `json:"resultUrl"`
	Output        map[string]any `json:"output"`
}

type Problem struct {
	Type          string  `json:"type"`
	Title         string  `json:"title"`
	Status        int     `json:"status"`
	Detail        string  `json:"detail"`
	Instance      string  `json:"instance,omitempty"`
	Profile       string  `json:"profile,omitempty"`
	RequestID     string  `json:"requestId,omitempty"`
	InteractionID string  `json:"interactionId,omitempty"`
	Code          string  `json:"code,omitempty"`
	Target        *Target `json:"target,omitempty"`
}

// Decode reads one strict JSON document. Unknown fields and trailing values are rejected.
func Decode[T any](reader io.Reader) (T, error) {
	var value T
	decoder := json.NewDecoder(io.LimitReader(reader, 1<<20+1))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		return value, fmt.Errorf("mailschema: decode: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return value, fmt.Errorf("mailschema: expected one JSON document")
	}
	return value, nil
}

func validHTTPS(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil && parsed.Scheme == "https" && parsed.Host != ""
}

func validateType(value TypeReference) error {
	if !validHTTPS(value.ID) || value.Version == "" || !digest.MatchString(value.ContractDigest) {
		return fmt.Errorf("invalid type reference")
	}
	return nil
}

func validateTarget(value Target) error {
	if !validHTTPS(value.ID) || value.Revision == "" || !digest.MatchString(value.Digest) {
		return fmt.Errorf("invalid target reference")
	}
	return nil
}

// ValidateDescription checks the fixed MAP 0.1 identifiers and core trust-relevant fields.
// Use the bundled JSON Schema for complete structural validation.
func ValidateDescription(value Description) error {
	if value.Context != Context01 || value.Kind != "MailAction" || value.Profile != Profile01 || !uuidURN.MatchString(value.ID) {
		return fmt.Errorf("mailschema: unsupported or invalid MAP description identity")
	}
	if err := validateType(value.Type); err != nil {
		return fmt.Errorf("mailschema: %w", err)
	}
	if err := validateTarget(value.Target); err != nil {
		return fmt.Errorf("mailschema: %w", err)
	}
	if !value.ExpiresAt.After(value.DescribedAt) {
		return fmt.Errorf("mailschema: expiry must be later than description time")
	}
	if !validHTTPS(value.Service.ID) || !validHTTPS(value.Service.Execution.URL) || !validHTTPS(value.Service.HumanURL) {
		return fmt.Errorf("mailschema: invalid service URL")
	}
	if value.Service.Execution.Method != "POST" || value.Service.Authorization.Kind != "service-configured" {
		return fmt.Errorf("mailschema: unsupported execution or authorization mode")
	}
	if len(value.Operations) == 0 {
		return fmt.Errorf("mailschema: no operations offered")
	}
	seen := make(map[string]bool, len(value.Operations))
	for _, operation := range value.Operations {
		if operation.ID == "" || seen[operation.ID] || strings.TrimSpace(operation.Name) == "" {
			return fmt.Errorf("mailschema: invalid or duplicate operation")
		}
		seen[operation.ID] = true
	}
	return nil
}

// ValidateRequest checks a MAP 0.1 request's fixed identifiers and references.
// Type-specific input constraints remain in the type schema.
func ValidateRequest(value Request) error {
	if value.Kind != "MapRequest" || value.Profile != Profile01 || !uuidURN.MatchString(value.RequestID) || !uuidURN.MatchString(value.InteractionID) {
		return fmt.Errorf("mailschema: unsupported or invalid MAP request identity")
	}
	if err := validateType(value.Type); err != nil {
		return fmt.Errorf("mailschema: %w", err)
	}
	if err := validateTarget(value.Target); err != nil {
		return fmt.Errorf("mailschema: %w", err)
	}
	if value.Operation == "" || value.Input == nil {
		return fmt.Errorf("mailschema: operation and input are required")
	}
	return nil
}

// ValidateResult checks a MAP 0.1 result's fixed identifiers and references.
// Type-specific output constraints remain in the type schema.
func ValidateResult(value Result) error {
	if value.Kind != "MapResult" || value.Profile != Profile01 || !uuidURN.MatchString(value.RequestID) || !uuidURN.MatchString(value.InteractionID) {
		return fmt.Errorf("mailschema: unsupported or invalid MAP result identity")
	}
	if err := validateType(value.Type); err != nil {
		return fmt.Errorf("mailschema: %w", err)
	}
	if err := validateTarget(value.Target); err != nil {
		return fmt.Errorf("mailschema: %w", err)
	}
	states := map[string]bool{
		"accepted": true, "completed": true, "failed": true, "pending": true, "approval-required": true,
	}
	if value.Operation == "" || !states[value.State] || value.RecordedAt.IsZero() || !validHTTPS(value.ResultURL) || value.Output == nil {
		return fmt.Errorf("mailschema: invalid MAP result state or fields")
	}
	return nil
}

// ValidateProblem checks a MAP 0.1 problem's fixed identifiers and the
// required relationship between its code, type URI and HTTP status.
func ValidateProblem(value Problem) error {
	if strings.TrimSpace(value.Type) == "" || strings.TrimSpace(value.Title) == "" || strings.TrimSpace(value.Detail) == "" || value.Status < 400 || value.Status > 599 {
		return fmt.Errorf("mailschema: unsupported or invalid MAP problem identity")
	}
	if value.Code == "" {
		if value.Instance != "" || value.Profile != "" || value.RequestID != "" || value.InteractionID != "" || value.Target != nil {
			return fmt.Errorf("mailschema: incomplete MAP problem correlation")
		}
		return nil
	}
	if value.Profile != Profile01 || !validHTTPS(value.Instance) || !uuidURN.MatchString(value.RequestID) {
		return fmt.Errorf("mailschema: unsupported or invalid MAP problem correlation")
	}
	if value.Code != "result-not-found" && !uuidURN.MatchString(value.InteractionID) {
		return fmt.Errorf("mailschema: MAP problem requires an interaction identifier")
	}
	if value.InteractionID != "" && !uuidURN.MatchString(value.InteractionID) {
		return fmt.Errorf("mailschema: invalid MAP problem interaction identifier")
	}
	expected := map[string]struct {
		typeURI string
		status  int
	}{
		"invalid-request":         {"https://mailschema.org/problems/invalid-request", 400},
		"authentication-required": {"https://mailschema.org/problems/authentication-required", 401},
		"refused":                 {"https://mailschema.org/problems/refused", 403},
		"result-not-found":        {"https://mailschema.org/problems/result-not-found", 404},
		"stale-target":            {"https://mailschema.org/problems/stale-target", 409},
		"idempotency-conflict":    {"https://mailschema.org/problems/idempotency-conflict", 409},
		"request-in-progress":     {"https://mailschema.org/problems/request-in-progress", 409},
		"expired-interaction":     {"https://mailschema.org/problems/expired-interaction", 410},
		"unsupported-profile":     {"https://mailschema.org/problems/unsupported-profile", 422},
		"unsupported-type":        {"https://mailschema.org/problems/unsupported-type", 422},
		"unsupported-operation":   {"https://mailschema.org/problems/unsupported-operation", 422},
	}
	relation, ok := expected[value.Code]
	if !ok || value.Type != relation.typeURI || value.Status != relation.status {
		return fmt.Errorf("mailschema: contradictory MAP problem code, type or status")
	}
	if value.Code == "stale-target" && value.Target == nil {
		return fmt.Errorf("mailschema: stale-target problem requires a target")
	}
	return nil
}
