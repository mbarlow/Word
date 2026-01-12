package translation

import (
	"fmt"
	"strings"
)

// TranslatorPrompt generates the prompt for the draft translation pass.
func TranslatorPrompt(profile *Profile, sourceVID string, sourceText string, kjvRef string, isHotVerse bool) string {
	var sb strings.Builder

	sb.WriteString("You are a biblical translator. Translate the following source text according to the profile.\n\n")

	// Profile context
	sb.WriteString("## Translation Profile\n")
	sb.WriteString(fmt.Sprintf("- Profile: %s\n", profile.Profile))
	sb.WriteString(fmt.Sprintf("- Target Language: %s\n", profile.TargetLang))
	sb.WriteString(fmt.Sprintf("- Audience: %s\n", profile.Audience))
	sb.WriteString(fmt.Sprintf("- Reading Level: %s\n", profile.ReadingLevel))
	sb.WriteString(fmt.Sprintf("- Tone: %s\n", profile.Tone))
	sb.WriteString("\n### Style Rules\n")
	for _, rule := range profile.StyleRules {
		sb.WriteString(fmt.Sprintf("- %s\n", rule))
	}

	if len(profile.TerminologyPolicy) > 0 {
		sb.WriteString("\n### Terminology Policy\n")
		for term, rendering := range profile.TerminologyPolicy {
			sb.WriteString(fmt.Sprintf("- %s → %s\n", term, rendering))
		}
	}

	// Hot verse warning
	if isHotVerse {
		sb.WriteString("\n## ⚠️ HOT VERSE WARNING\n")
		sb.WriteString("This verse is theologically significant. Be extra conservative in translation choices.\n")
		sb.WriteString("Include translator notes explaining any ambiguity or key decisions.\n")
	}

	// Source text
	sb.WriteString("\n## Source Text\n")
	sb.WriteString(fmt.Sprintf("VID: %s\n", sourceVID))
	sb.WriteString(fmt.Sprintf("```\n%s\n```\n", sourceText))

	// KJV reference (if available)
	if kjvRef != "" {
		sb.WriteString("\n## KJV Reference (for context only, do not copy)\n")
		sb.WriteString(fmt.Sprintf("```\n%s\n```\n", kjvRef))
	}

	// Instructions
	sb.WriteString("\n## Instructions\n")
	sb.WriteString("1. Translate the source text into the target language\n")
	sb.WriteString("2. Follow all style rules and terminology policy\n")
	sb.WriteString("3. If there is ambiguity, choose the most conservative rendering and add a note\n")
	sb.WriteString("4. Do not add information not in the source\n")
	sb.WriteString("5. Do not harmonize with other passages\n")

	sb.WriteString("\n## Output Format\n")
	sb.WriteString("Respond with ONLY a JSON object:\n")
	sb.WriteString("```json\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"text\": \"Your translated verse text here\",\n")
	sb.WriteString("  \"notes\": [\"Optional translator notes\"]\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")

	return sb.String()
}

// CriticPrompt generates the prompt for the QA/critic pass.
func CriticPrompt(profile *Profile, sourceVID string, sourceText string, draftText string, kjvRef string, isHotVerse bool) string {
	var sb strings.Builder

	sb.WriteString("You are a biblical translation critic. Review the draft translation for issues.\n\n")
	sb.WriteString("DO NOT rewrite the translation. Only identify potential issues.\n\n")

	// Profile context
	sb.WriteString("## Translation Profile\n")
	sb.WriteString(fmt.Sprintf("- Profile: %s\n", profile.Profile))
	sb.WriteString(fmt.Sprintf("- Audience: %s\n", profile.Audience))
	sb.WriteString(fmt.Sprintf("- Tone: %s\n", profile.Tone))

	if len(profile.Consistency.TrackKeyTerms) > 0 {
		sb.WriteString("\n### Key Terms to Track\n")
		for _, term := range profile.Consistency.TrackKeyTerms {
			sb.WriteString(fmt.Sprintf("- %s\n", term))
		}
	}

	// Hot verse context
	if isHotVerse {
		sb.WriteString("\n## ⚠️ HOT VERSE\n")
		sb.WriteString("This verse is theologically significant. Apply extra scrutiny.\n")
	}

	// Source and draft
	sb.WriteString("\n## Source Text\n")
	sb.WriteString(fmt.Sprintf("VID: %s\n", sourceVID))
	sb.WriteString(fmt.Sprintf("```\n%s\n```\n", sourceText))

	if kjvRef != "" {
		sb.WriteString("\n## KJV Reference\n")
		sb.WriteString(fmt.Sprintf("```\n%s\n```\n", kjvRef))
	}

	sb.WriteString("\n## Draft Translation\n")
	sb.WriteString(fmt.Sprintf("```\n%s\n```\n", draftText))

	// Flag types
	sb.WriteString("\n## Flag Types to Check\n")
	sb.WriteString("- `semantic_drift`: Meaning has shifted from source\n")
	sb.WriteString("- `missing_subject`: Subject or object unclear or omitted\n")
	sb.WriteString("- `tense_aspect`: Verb tense/aspect doesn't match source\n")
	sb.WriteString("- `theological_bias`: Translation implies doctrine not in source\n")
	sb.WriteString("- `consistency`: Key term rendered differently than expected\n")
	sb.WriteString("- `ambiguity`: Source allows multiple readings, needs note\n")
	sb.WriteString("- `hot_verse`: Extra caution needed (only if this is a hot verse)\n")

	sb.WriteString("\n## Output Format\n")
	sb.WriteString("Respond with ONLY a JSON object. If no issues, return empty flags array:\n")
	sb.WriteString("```json\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"flags\": [\n")
	sb.WriteString("    {\"type\": \"flag_type\", \"detail\": \"explanation\"}\n")
	sb.WriteString("  ]\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")

	return sb.String()
}
