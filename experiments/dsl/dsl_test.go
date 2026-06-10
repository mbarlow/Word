package main

import (
	"reflect"
	"strings"
	"testing"
)

// --- reader.go ------------------------------------------------------------

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    []token
		wantErr bool
	}{
		{
			name: "atoms and parens",
			src:  "(bara ELOHIM)",
			want: []token{{tLParen, "("}, {tAtom, "bara"}, {tAtom, "ELOHIM"}, {tRParen, ")"}},
		},
		{
			name: "string literal",
			src:  `(program "bereshit")`,
			want: []token{{tLParen, "("}, {tAtom, "program"}, {tString, "bereshit"}, {tRParen, ")"}},
		},
		{
			name: "comment discarded",
			src:  "yehi ; let there be\nOR",
			want: []token{{tAtom, "yehi"}, {tAtom, "OR"}},
		},
		{
			name:    "unterminated string",
			src:     `(yomer "unclosed`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenize(tt.src)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("tokenize: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokens = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		wantForms int
		wantErr   bool
	}{
		{"single form", "(bara ELOHIM)", 1, false},
		{"two top-level forms", "(yehi OR) (va-yehi OR)", 2, false},
		{"nested lists", "(yomer ELOHIM (yehi OR))", 1, false},
		{"unterminated list", "(yehi OR", 0, true},
		{"stray rparen", ")", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forms, err := parse(tt.src)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if len(forms) != tt.wantForms {
				t.Errorf("got %d forms, want %d", len(forms), tt.wantForms)
			}
		})
	}
}

func TestParseRoundTrip(t *testing.T) {
	src := `(yomer ELOHIM (yehi OR :day 1) "and there was light")`
	forms, err := parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(forms) != 1 {
		t.Fatalf("got %d forms, want 1", len(forms))
	}
	if got := forms[0].String(); got != src {
		t.Errorf("round trip = %q, want %q", got, src)
	}
}

func TestNodeHead(t *testing.T) {
	forms, err := parse(`(bara ELOHIM) ("not-a-symbol" head)`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := forms[0].head(); got != "bara" {
		t.Errorf("head = %q, want bara", got)
	}
	if got := forms[1].head(); got != "" {
		t.Errorf("head of string-led list = %q, want empty", got)
	}
	sym := &Node{Kind: Symbol, Text: "OR"}
	if got := sym.head(); got != "" {
		t.Errorf("head of symbol = %q, want empty", got)
	}
}

// --- interp.go ------------------------------------------------------------

func TestStripPrefix(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"HA-ERETZ", "ERETZ"},
		{"u-ve-OF", "OF"},
		{"ET-HA-ADAM", "ADAM"},
		{":day", "DAY"},
		{"OR", "OR"},
		{"adam", "ADAM"},
		{"HA-", "HA-"}, // bare prefix is not stripped to empty
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := stripPrefix(tt.in); got != tt.want {
				t.Errorf("stripPrefix(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestKeywordArgs(t *testing.T) {
	forms, err := parse(`(va-yehi :day 6 :eval TOV)`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	kw := keywordArgs(forms[0])
	if len(kw) != 2 {
		t.Fatalf("got %d keyword args, want 2", len(kw))
	}
	if d, ok := kw[":day"]; !ok || d.text() != "6" {
		t.Errorf(":day = %v, want 6", d)
	}
	if e, ok := kw[":eval"]; !ok || e.text() != "TOV" {
		t.Errorf(":eval = %v, want TOV", e)
	}
}

func TestCollectNouns(t *testing.T) {
	forms, err := parse(`(va-yaas ET-HA-ADAM ve-ET-HA-BEHEMAH (nested HA-OR))`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var ordered []string
	collectNouns(forms[0], map[string]bool{}, &ordered)
	want := []string{"ADAM", "BEHEMAH", "OR"}
	if !reflect.DeepEqual(ordered, want) {
		t.Errorf("collectNouns = %v, want %v", ordered, want)
	}
}

func TestBaraCount(t *testing.T) {
	forms, err := parse(`(va-yivra ELOHIM bara bara)`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := baraCount(forms[0]); got != 3 {
		t.Errorf("baraCount = %d, want 3", got)
	}
}

func TestConsonants(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"DESHE", "dsh"},
		{"tadshe", "tdsh"},
		{"va-yehi", "vyh"},
		{"OR", "r"},
	}
	for _, tt := range tests {
		if got := consonants(tt.in); got != tt.want {
			t.Errorf("consonants(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestDetectCognates(t *testing.T) {
	forms, err := parse(`(tadshe HA-ERETZ DESHE)`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := detectCognates(forms[0])
	if len(got) != 1 {
		t.Fatalf("got %d cognate pairs, want 1: %v", len(got), got)
	}
	if !strings.Contains(got[0], "tadshe") || !strings.Contains(got[0], "DESHE") {
		t.Errorf("cognate pair = %q, want tadshe ↔ DESHE", got[0])
	}

	// No cognates between unrelated symbols.
	forms, err = parse(`(yehi OR ba-RAKIA)`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := detectCognates(forms[0]); len(got) != 0 {
		t.Errorf("got cognates %v, want none", got)
	}
}

func TestBlessingText(t *testing.T) {
	forms, err := parse(`(va-yevarekh peru u-revu u-milu) (va-yevarekh nothing-known)`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := blessingText(forms[0]); got != "be fruitful, multiply, fill" {
		t.Errorf("blessingText = %q", got)
	}
	if got := blessingText(forms[1]); got != "blessed" {
		t.Errorf("blessingText fallback = %q, want \"blessed\"", got)
	}
}

func TestInterpMissingAndExitCode(t *testing.T) {
	in := newInterp()
	if in.exitCode() != 1 {
		t.Error("fresh interp should exit 1 (universe incomplete)")
	}
	if missing := in.missing(); len(missing) != 13 {
		t.Errorf("fresh interp missing %d fields, want 13", len(missing))
	}

	// Set every creation field and complete the run.
	u := in.u
	u.Exists, u.Light, u.Sky, u.Land, u.Seas, u.Plants = true, true, true, true, true, true
	u.Sun, u.Moon, u.Stars, u.Fish, u.Birds, u.Animals, u.Humans = true, true, true, true, true, true, true
	u.Status = "COMPLETE"
	if missing := in.missing(); len(missing) != 0 {
		t.Errorf("complete universe still missing %v", missing)
	}
	if in.exitCode() != 0 {
		t.Error("complete universe should exit 0")
	}
}
