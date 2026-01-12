package translation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// LLMClient interface for translation generation.
type LLMClient interface {
	Generate(ctx context.Context, prompt string, temperature float64) (string, error)
	ModelName() string
}

// Translator coordinates the translation process.
type Translator struct {
	client      LLMClient
	profile     *Profile
	temperature float64
}

// NewTranslator creates a new translator with the given LLM client and profile.
func NewTranslator(client LLMClient, profile *Profile, temperature float64) *Translator {
	return &Translator{
		client:      client,
		profile:     profile,
		temperature: temperature,
	}
}

// SourceVerse represents a verse to be translated.
type SourceVerse struct {
	VID    string
	Text   string
	KJVRef string // Optional KJV reference text
}

// TranslateVerse performs the draft translation pass.
func (t *Translator) TranslateVerse(ctx context.Context, source SourceVerse) (*Draft, error) {
	isHot := IsHotVerse(source.VID)
	prompt := TranslatorPrompt(t.profile, source.VID, source.Text, source.KJVRef, isHot)

	response, err := t.client.Generate(ctx, prompt, t.temperature)
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	// Parse response
	var result struct {
		Text  string   `json:"text"`
		Notes []string `json:"notes"`
	}

	// Extract JSON from response (handle markdown code blocks)
	jsonStr := extractJSON(response)
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w (response: %s)", err, response)
	}

	// Build target VID: t_<profile>/<book>/<chapter>/<verse>
	targetVID := buildTargetVID(t.profile.Profile, source.VID)

	draft := &Draft{
		ID:        targetVID,
		SourceIDs: []string{source.VID},
		Lang:      t.profile.TargetLang,
		Profile:   t.profile.Profile,
		Text:      result.Text,
		Notes:     result.Notes,
		Provenance: Provenance{
			Model:        t.client.ModelName(),
			PromptHash:   ComputePromptHash(prompt),
			Timestamp:    time.Now().UTC(),
			Temperature:  t.temperature,
			InputsDigest: ComputeInputsDigest(source.Text),
		},
	}

	return draft, nil
}

// CriticVerse performs the QA/critic pass on a draft.
func (t *Translator) CriticVerse(ctx context.Context, source SourceVerse, draft *Draft) (*CriticResult, error) {
	isHot := IsHotVerse(source.VID)
	prompt := CriticPrompt(t.profile, source.VID, source.Text, draft.Text, source.KJVRef, isHot)

	response, err := t.client.Generate(ctx, prompt, 0.0) // Use temperature 0 for critic
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	// Parse response
	var result struct {
		Flags []CriticFlag `json:"flags"`
	}

	jsonStr := extractJSON(response)
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse critic response: %w", err)
	}

	// Add hot verse flag if applicable
	if isHot {
		hasHotFlag := false
		for _, f := range result.Flags {
			if f.Type == FlagHotVerse {
				hasHotFlag = true
				break
			}
		}
		if !hasHotFlag {
			result.Flags = append(result.Flags, CriticFlag{
				Type:   FlagHotVerse,
				Detail: "Theologically significant verse - manual review recommended",
			})
		}
	}

	return &CriticResult{
		ID:    draft.ID,
		Flags: result.Flags,
	}, nil
}

// buildTargetVID creates target VID from source VID.
// "heb-wlc/GEN/1/1" -> "t_techdoc_en/GEN/1/1"
func buildTargetVID(profile, sourceVID string) string {
	parts := strings.SplitN(sourceVID, "/", 2)
	if len(parts) < 2 {
		return fmt.Sprintf("t_%s/%s", profile, sourceVID)
	}
	return fmt.Sprintf("t_%s/%s", profile, parts[1])
}

// extractJSON extracts JSON from a response that may contain markdown code blocks.
func extractJSON(response string) string {
	response = strings.TrimSpace(response)

	// Try to extract from ```json ... ``` blocks
	if idx := strings.Index(response, "```json"); idx != -1 {
		start := idx + 7
		if end := strings.Index(response[start:], "```"); end != -1 {
			return strings.TrimSpace(response[start : start+end])
		}
	}

	// Try to extract from ``` ... ``` blocks
	if idx := strings.Index(response, "```"); idx != -1 {
		start := idx + 3
		// Skip language identifier if present
		if nl := strings.Index(response[start:], "\n"); nl != -1 {
			start += nl + 1
		}
		if end := strings.Index(response[start:], "```"); end != -1 {
			return strings.TrimSpace(response[start : start+end])
		}
	}

	// Return as-is if no code blocks
	return response
}
