package translation

import (
	"fmt"
	"strings"
)

// TranslatorPrompt generates the prompt for the draft translation pass.
func TranslatorPrompt(profile *Profile, sourceVID string, sourceText string, kjvRef string, isHotVerse bool) string {
	var sb strings.Builder

	// Select appropriate system prompt based on profile type
	if profile.IsStructural() {
		return structuralTranslatorPrompt(profile, sourceVID, sourceText, kjvRef, isHotVerse)
	}
	if profile.IsCodeFormat() {
		return codeFormatTranslatorPrompt(profile, sourceVID, sourceText, kjvRef, isHotVerse)
	}

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

// structuralTranslatorPrompt generates the prompt for structural/scholarly translations.
func structuralTranslatorPrompt(profile *Profile, sourceVID string, sourceText string, kjvRef string, isHotVerse bool) string {
	var sb strings.Builder

	sb.WriteString("You are a Hebrew linguistic analyst and biblical translator specializing in preserving literary structures.\n\n")

	sb.WriteString("## Profile: structural_en\n")
	sb.WriteString("Target: Scholars, students, programmers interested in Hebrew literary structure\n")
	sb.WriteString("Goal: Preserve Hebrew literary devices that are typically lost in conventional translations\n\n")

	sb.WriteString("## Hebrew Literary Devices to Detect\n")
	sb.WriteString("1. **Cognate Accusative (Figura Etymologica)**: Verb + noun from SAME ROOT\n")
	sb.WriteString("   - Example: תַּדְשֵׁא...דֶּשֶׁא (tadshé...déshe) = 'vegetate vegetation' (root: דשא)\n")
	sb.WriteString("   - Example: מַזְרִיעַ זֶרַע (mazría zéra) = 'seeding seed' (root: זרע)\n")
	sb.WriteString("2. **Repetition (Epanalepsis)**: Same word repeated for emphasis\n")
	sb.WriteString("3. **Merism**: Two extremes representing totality (heavens+earth = everything)\n")
	sb.WriteString("4. **Chiasm**: ABBA inverted parallelism\n\n")

	// Hot verse warning
	if isHotVerse {
		sb.WriteString("## ⚠️ HOT VERSE - Apply extra scholarly rigor\n\n")
	}

	sb.WriteString("## Source Text\n")
	sb.WriteString(fmt.Sprintf("VID: %s\n", sourceVID))
	sb.WriteString(fmt.Sprintf("```hebrew\n%s\n```\n", sourceText))

	if kjvRef != "" {
		sb.WriteString("\n## KJV Reference (for semantic context only)\n")
		sb.WriteString(fmt.Sprintf("```\n%s\n```\n", kjvRef))
	}

	sb.WriteString("\n## Required Output Layers\n")
	sb.WriteString("Analyze and produce ALL of the following layers:\n\n")

	sb.WriteString("### Layer 1: Root Analysis\n")
	sb.WriteString("Identify Hebrew roots and detect cognate pairs (verb+noun sharing same root).\n")
	sb.WriteString("Format: `[n] ROOT (transliteration): word1 (form) ↔ word2 (form) — pattern_type`\n\n")

	sb.WriteString("### Layer 2: Structural Representation\n")
	sb.WriteString("Express the verse as pseudo-code showing the functional/declarative structure.\n")
	sb.WriteString("Use format like: `god.say(earth.vegetate<דשא>(vegetation<דשא>, plant.seed<זרע>(seed<זרע>))) => TRUE`\n\n")

	sb.WriteString("### Layer 3: Readable Translation with Markers\n")
	sb.WriteString("Translate to readable English with superscript markers (¹²³) linking cognate pairs.\n")
	sb.WriteString("Preserve awkward English if needed to show Hebrew structure.\n")
	sb.WriteString("Example: 'Let the earth **vegetate¹ vegetation¹**—plants **seeding² seed²**'\n\n")

	sb.WriteString("### Layer 4: Cognate Glossary\n")
	sb.WriteString("Table of cognate pairs found: marker, root, Hebrew forms, glosses, pattern type.\n\n")

	sb.WriteString("### Layer 5: Polysemy Ranges\n")
	sb.WriteString("For key Hebrew words, show semantic range: `word: meaning1|meaning2|meaning3`\n\n")

	sb.WriteString("### Layer 6: Literary Devices\n")
	sb.WriteString("List any detected devices: cognate_accusative, repetition, merism, chiasm, inclusio.\n\n")

	sb.WriteString("## Output Format\n")
	sb.WriteString("Respond with ONLY a JSON object (no markdown code fences):\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"root_analysis\": \"[1] דשא (d-sh-a): תַּדְשֵׁא (hiphil) ↔ דֶּשֶׁא (noun) — cognate_accusative\\n[2] ...\",\n")
	sb.WriteString("  \"structural\": \"god.say(earth.vegetate<דשא>(vegetation<דשא>)) => TRUE\",\n")
	sb.WriteString("  \"text\": \"And God said, 'Let the earth **vegetate¹ vegetation¹**...'\",\n")
	sb.WriteString("  \"cognates\": [\n")
	sb.WriteString("    {\"marker\": 1, \"root\": \"דשא\", \"verb\": \"תַּדְשֵׁא\", \"noun\": \"דֶּשֶׁא\", \"gloss\": \"vegetate/vegetation\", \"pattern\": \"cognate_accusative\"}\n")
	sb.WriteString("  ],\n")
	sb.WriteString("  \"polysemy\": {\"אֱלֹהִים\": \"God|gods|divine-council\", \"אֶרֶץ\": \"earth|land|ground\"},\n")
	sb.WriteString("  \"literary_devices\": [\"cognate_accusative\", \"merism\"],\n")
	sb.WriteString("  \"notes\": [\"Optional scholarly notes\"]\n")
	sb.WriteString("}\n")

	return sb.String()
}

// codeFormatTranslatorPrompt generates the prompt for Lisp/YAML code-style translations.
func codeFormatTranslatorPrompt(profile *Profile, sourceVID string, sourceText string, kjvRef string, isHotVerse bool) string {
	var sb strings.Builder

	format := profile.OutputFormat.Format // "lisp" or "yaml"

	sb.WriteString("You are a biblical translator who represents scripture as executable code structures.\n\n")

	sb.WriteString(fmt.Sprintf("## Profile: %s\n", profile.Profile))
	sb.WriteString(fmt.Sprintf("Output Format: %s\n", strings.ToUpper(format)))
	sb.WriteString(fmt.Sprintf("Audience: %s\n\n", profile.Audience))

	sb.WriteString("## Source Text\n")
	sb.WriteString(fmt.Sprintf("VID: %s\n", sourceVID))
	sb.WriteString(fmt.Sprintf("```hebrew\n%s\n```\n", sourceText))

	if kjvRef != "" {
		sb.WriteString("\n## English Reference (for semantic guidance)\n")
		sb.WriteString(fmt.Sprintf("```\n%s\n```\n", kjvRef))
	}

	sb.WriteString("\n## Translation Instructions\n")

	if format == "lisp" {
		sb.WriteString("Express this verse as a Lisp S-expression that captures the Hebrew structure.\n\n")
		sb.WriteString("### Lisp Conventions:\n")
		sb.WriteString("- Use Hebrew verb roots as function names (transliterated): `(yomer ...)`, `(bara ...)`, `(yehi ...)`\n")
		sb.WriteString("- Use CAPS for divine/proper names: `ELOHIM`, `YHWH`, `HA-SHAMAYIM`\n")
		sb.WriteString("- Use keywords for attributes: `:state`, `:result`, `:kind`\n")
		sb.WriteString("- Show cognate structures: `(tadshe HA-ERETZ DESHE)` — vegetate(vegetation)\n")
		sb.WriteString("- Use `va-yehi KEN` for 'and it was so' → returns TRUE\n")
		sb.WriteString("- Include Hebrew comments: `;; תַּדְשֵׁא...דֶּשֶׁא (cognate accusative)`\n\n")
		sb.WriteString("### Example:\n")
		sb.WriteString("```lisp\n")
		sb.WriteString("(yomer ELOHIM                         ;; God said\n")
		sb.WriteString("  (yehi OR))                          ;; 'Let there be light'\n")
		sb.WriteString("(va-yehi OR)                          ;; And there was light → OUTPUT\n")
		sb.WriteString("```\n\n")
	} else if format == "yaml" {
		sb.WriteString("Express this verse as YAML that captures the Hebrew declarative structure.\n\n")
		sb.WriteString("### YAML Conventions:\n")
		sb.WriteString("- Use `speaker:` for who is speaking\n")
		sb.WriteString("- Use `command:` for the divine command\n")
		sb.WriteString("- Use `result:` for the outcome\n")
		sb.WriteString("- Use `cognates:` to note same-root verb/noun pairs\n")
		sb.WriteString("- Use `|` for multi-line strings\n")
		sb.WriteString("- Include Hebrew in comments or `hebrew:` field\n\n")
		sb.WriteString("### Example:\n")
		sb.WriteString("```yaml\n")
		sb.WriteString("verse: GEN.1.3\n")
		sb.WriteString("speaker: ELOHIM\n")
		sb.WriteString("command:\n")
		sb.WriteString("  action: yehi  # let there be\n")
		sb.WriteString("  object: OR    # light\n")
		sb.WriteString("result: va-yehi OR  # and there was light\n")
		sb.WriteString("evaluation: null  # no 'good' evaluation this verse\n")
		sb.WriteString("```\n\n")
	}

	sb.WriteString("## Output Format\n")
	sb.WriteString("Respond with ONLY a JSON object (no markdown code fences):\n")
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("  \"text\": \"Your %s code here as a single string (use \\\\n for newlines)\",\n", format))
	sb.WriteString("  \"notes\": [\"Optional notes about structure or Hebrew features\"]\n")
	sb.WriteString("}\n")

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
