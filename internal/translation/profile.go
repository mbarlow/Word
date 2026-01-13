package translation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Profile defines a translation target audience and style constraints.
type Profile struct {
	Profile           string             `json:"profile"`
	TargetLang        string             `json:"target_lang"`
	Audience          string             `json:"audience"`
	ReadingLevel      string             `json:"reading_level"`
	Tone              string             `json:"tone"`
	StyleRules        []string           `json:"style_rules"`
	TerminologyPolicy map[string]string  `json:"terminology_policy"`
	Consistency       ConsistencyConfig  `json:"consistency"`
	OutputFormat      OutputFormatConfig `json:"output_format,omitempty"`
	LiteraryDevices   LiteraryDevices    `json:"literary_devices,omitempty"`
}

// ConsistencyConfig defines term tracking rules.
type ConsistencyConfig struct {
	PreferConsistentRenderings bool     `json:"prefer_consistent_renderings"`
	ShowAllOptions             bool     `json:"show_all_options,omitempty"`
	TrackKeyTerms              []string `json:"track_key_terms"`
}

// OutputFormatConfig defines what layers to include in output.
type OutputFormatConfig struct {
	IncludeHebrew          bool   `json:"include_hebrew,omitempty"`
	IncludeStructural      bool   `json:"include_structural,omitempty"`
	IncludeReadable        bool   `json:"include_readable,omitempty"`
	IncludeCognates        bool   `json:"include_cognates,omitempty"`
	IncludePolysemy        bool   `json:"include_polysemy,omitempty"`
	IncludeLiteraryDevices bool   `json:"include_literary_devices,omitempty"`
	Format                 string `json:"format,omitempty"` // "prose", "lisp", "yaml", "structural"
}

// LiteraryDevices configures detection of Hebrew literary patterns.
type LiteraryDevices struct {
	DetectCognateAccusative bool `json:"detect_cognate_accusative,omitempty"`
	DetectRepetition        bool `json:"detect_repetition,omitempty"`
	DetectChiasm            bool `json:"detect_chiasm,omitempty"`
	DetectMerism            bool `json:"detect_merism,omitempty"`
	DetectInclusio          bool `json:"detect_inclusio,omitempty"`
}

// IsStructural returns true if this profile uses structural output format.
func (p *Profile) IsStructural() bool {
	return p.OutputFormat.IncludeStructural || p.OutputFormat.Format == "structural"
}

// IsCodeFormat returns true if output should be in a programming language format.
func (p *Profile) IsCodeFormat() bool {
	return p.OutputFormat.Format == "lisp" || p.OutputFormat.Format == "yaml"
}

// LoadProfile reads a profile from a JSON file.
func LoadProfile(path string) (*Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile: %w", err)
	}

	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse profile: %w", err)
	}

	return &p, nil
}

// LoadProfileByName loads a profile from the profiles directory.
func LoadProfileByName(profilesDir, name string) (*Profile, error) {
	path := filepath.Join(profilesDir, name+".json")
	return LoadProfile(path)
}

// ListProfiles returns all available profile names in a directory.
func ListProfiles(profilesDir string) ([]string, error) {
	entries, err := os.ReadDir(profilesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read profiles directory: %w", err)
	}

	var profiles []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			name := e.Name()[:len(e.Name())-5] // strip .json
			profiles = append(profiles, name)
		}
	}
	return profiles, nil
}
