package translation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Profile defines a translation target audience and style constraints.
type Profile struct {
	Profile           string            `json:"profile"`
	TargetLang        string            `json:"target_lang"`
	Audience          string            `json:"audience"`
	ReadingLevel      string            `json:"reading_level"`
	Tone              string            `json:"tone"`
	StyleRules        []string          `json:"style_rules"`
	TerminologyPolicy map[string]string `json:"terminology_policy"`
	Consistency       ConsistencyConfig `json:"consistency"`
}

// ConsistencyConfig defines term tracking rules.
type ConsistencyConfig struct {
	PreferConsistentRenderings bool     `json:"prefer_consistent_renderings"`
	TrackKeyTerms              []string `json:"track_key_terms"`
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
