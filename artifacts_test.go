package mailschema

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
)

func TestProfileRecordBindsTheBundledBytes(t *testing.T) {
	contents, err := Bytes(ProfileRecord)
	if err != nil {
		t.Fatal(err)
	}
	var record struct {
		ID        string `json:"id"`
		Artifacts map[string]struct {
			SHA256 string `json:"sha256"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal(contents, &record); err != nil {
		t.Fatal(err)
	}
	if record.ID != Profile {
		t.Fatalf("profile record names %q", record.ID)
	}
	for key, artifact := range map[string]Artifact{
		"context":        ContextDocument,
		"schema":         CoreSchema,
		"contractFormat": ContractSchema,
	} {
		bytes, err := Bytes(artifact)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(bytes)
		if hex.EncodeToString(sum[:]) != record.Artifacts[key].SHA256 {
			t.Fatalf("%s differs from the bytes the profile record binds", artifact)
		}
	}
	if _, err := Bytes(ImplementationSchema); err != nil {
		t.Fatal(err)
	}
	if _, err := Bytes("other.json"); err == nil {
		t.Fatal("an unknown artifact was returned")
	}
}
