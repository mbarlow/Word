package translation

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Draft represents a translated verse output.
type Draft struct {
	ID         string     `json:"id"`
	SourceIDs  []string   `json:"source_ids"`
	Lang       string     `json:"lang"`
	Profile    string     `json:"profile"`
	Text       string     `json:"text"`
	Notes      []string   `json:"notes,omitempty"`
	Provenance Provenance `json:"provenance"`
}

// Provenance tracks the origin and reproducibility of a translation.
type Provenance struct {
	Model        string    `json:"model"`
	PromptHash   string    `json:"prompt_hash"`
	Timestamp    time.Time `json:"timestamp"`
	Temperature  float64   `json:"temperature"`
	InputsDigest string    `json:"inputs_digest"`
}

// CriticResult contains QA flags for a translated verse.
type CriticResult struct {
	ID    string       `json:"id"`
	Flags []CriticFlag `json:"flags"`
}

// CriticFlag represents a single QA issue.
type CriticFlag struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

// Flag types for critic pass.
const (
	FlagSemanticDrift     = "semantic_drift"
	FlagMissingSubject    = "missing_subject"
	FlagTenseAspect       = "tense_aspect"
	FlagTheologicalBias   = "theological_bias"
	FlagConsistency       = "consistency"
	FlagHotVerse          = "hot_verse"
	FlagAmbiguity         = "ambiguity"
)

// ComputeInputsDigest creates a SHA256 hash of source texts.
func ComputeInputsDigest(sourceTexts ...string) string {
	h := sha256.New()
	for _, t := range sourceTexts {
		h.Write([]byte(t))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// ComputePromptHash creates a SHA256 hash of the prompt template.
func ComputePromptHash(prompt string) string {
	h := sha256.Sum256([]byte(prompt))
	return hex.EncodeToString(h[:])[:16]
}
