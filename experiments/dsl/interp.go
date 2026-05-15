// interp.go — evaluator for the parsed Genesis DSL.
//
// The interpreter walks the s-expression forms and dispatches on the head
// verb. Creative verbs scan their form for known Hebrew nouns and set the
// matching field on the Universe. State is real: the final UNIVERSE struct
// reflects what the program actually instantiated, not a hardcoded script.
package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Universe is the program's mutable state — the thing being created.
type Universe struct {
	Exists, Light, Sky, Land, Seas, Plants         bool
	Sun, Moon, Stars, Fish, Birds, Animals, Humans bool
	Day                                            int
	Status                                         string
	Evaluation                                     string
}

// Interp holds evaluation state across all forms.
type Interp struct {
	u    *Universe
	env  map[string]*Node  // define bindings
	meta map[string]string // program metadata
}

func newInterp() *Interp {
	return &Interp{
		u:    &Universe{Status: "RUNNING", Evaluation: "—"},
		env:  map[string]*Node{},
		meta: map[string]string{},
	}
}

// effect maps a Hebrew noun to a creation it triggers.
type effect struct {
	label string
	apply func(*Universe)
}

// nouns is the instantiation table: when one of these appears inside a
// creative form, the corresponding Universe field is set true.
var nouns = map[string]effect{
	"OR":           {"light", func(u *Universe) { u.Light = true }},
	"RAKIA":        {"sky/firmament", func(u *Universe) { u.Sky = true }},
	"YABASHAH":     {"dry land", func(u *Universe) { u.Land = true }},
	"YAMIM":        {"seas", func(u *Universe) { u.Seas = true }},
	"DESHE":        {"vegetation", func(u *Universe) { u.Plants = true }},
	"ESEV":         {"plants", func(u *Universe) { u.Plants = true }},
	"ETZ":          {"trees", func(u *Universe) { u.Plants = true }},
	"GADOL":        {"sun (greater light)", func(u *Universe) { u.Sun = true }},
	"KATAN":        {"moon (lesser light)", func(u *Universe) { u.Moon = true }},
	"KOKHAVIM":     {"stars", func(u *Universe) { u.Stars = true }},
	"SHERETZ":      {"sea swarmers", func(u *Universe) { u.Fish = true }},
	"TANINIM":      {"great sea creatures", func(u *Universe) { u.Fish = true }},
	"OF":           {"birds", func(u *Universe) { u.Birds = true }},
	"BEHEMAH":      {"livestock", func(u *Universe) { u.Animals = true }},
	"REMES":        {"creeping things", func(u *Universe) { u.Animals = true }},
	"CHAYTO-ERETZ": {"wild animals", func(u *Universe) { u.Animals = true }},
	"ADAM":         {"humans", func(u *Universe) { u.Humans = true }},
}

// stripable are hyphen-anchored Hebrew proclitic prefixes peeled before a
// noun lookup (e.g. HA-ERETZ -> ERETZ, u-ve-OF -> OF).
var stripable = []string{"HA-", "VE-", "VA-", "U-", "BE-", "BI-", "LE-", "LA-", "MI-", "ME-", "KI-", "AL-", "ET-"}

// stripPrefix peels leading proclitics and upper-cases for table lookup.
func stripPrefix(s string) string {
	up := strings.ToUpper(strings.TrimPrefix(s, ":"))
	for again := true; again; {
		again = false
		for _, p := range stripable {
			if len(up) > len(p) && strings.HasPrefix(up, p) {
				up = up[len(p):]
				again = true
				break
			}
		}
	}
	return up
}

// eval dispatches one form on its head verb.
func (in *Interp) eval(form *Node) {
	if form == nil || form.Kind != List || len(form.Items) == 0 {
		return
	}
	switch form.head() {
	case "program":
		in.doProgram(form)
	case "bara":
		in.doBara(form)
	case "define":
		in.doDefine(form)
	case "yomer":
		in.doSpeak(form)
	case "yehi", "vihi", "yikavu", "ve-teraeh", "tadshe", "yishretzu", "totze":
		in.doCommand(form)
	case "hineh":
		in.doProvision(form)
	case "naaseh":
		in.doCouncil(form)
	case "va-yehi":
		in.doVayehi(form)
	case "va-yar":
		in.doSee(form)
	case "va-yavdel":
		in.doSeparate(form)
	case "va-yikra":
		in.doName(form)
	case "va-yaas", "va-totze", "va-yivra":
		in.doMake(form)
	case "va-yevarekh":
		in.doBless(form)
	case "va-yekhulu", "va-yekhal":
		in.doFinish(form)
	case "va-yishbot":
		in.doRest(form)
	case "return":
		in.doReturn(form)
	default:
		// No registered verb — still a clause to instantiate (e.g. a
		// noun-led "ve-OF yeofeif" / "and birds fly").
		in.exec(form)
		in.fire(form)
	}
}

// exec prints the EXEC trace line for a form (truncated for long forms).
func (in *Interp) exec(form *Node) {
	line := []rune(form.String())
	if len(line) > 72 {
		line = append(line[:71], '…')
	}
	fmt.Println()
	fmt.Println(Yellow + "► EXEC  " + Reset + Dim + string(line) + Reset)
}

// fire applies every noun the form mentions and reports instantiations,
// cognate echoes, and bara-root emphasis.
func (in *Interp) fire(form *Node) {
	var ordered []string
	seen := map[string]bool{}
	collectNouns(form, seen, &ordered)
	if len(ordered) > 0 {
		var labels []string
		for _, k := range ordered {
			nouns[k].apply(in.u)
			labels = append(labels, nouns[k].label)
		}
		fmt.Println(Green + "  ✓ instantiated: " + strings.Join(labels, ", ") + Reset)
	}
	for _, c := range detectCognates(form) {
		fmt.Println(Dim + "  ≈ cognate: " + c + " (verb/noun share a root)" + Reset)
	}
	if n := baraCount(form); n >= 2 {
		fmt.Printf(Magenta+"  ⚑ root ב־ר־א (create) used %d× — emphatic creation\n"+Reset, n)
	}
}

// --- verb handlers --------------------------------------------------------

func (in *Interp) doProgram(form *Node) {
	in.exec(form)
	for _, it := range form.Items[1:] {
		if it.Kind == Str {
			in.meta["name"] = it.Text
			break
		}
	}
	for k, v := range keywordArgs(form) {
		in.meta[strings.TrimPrefix(k, ":")] = v.text()
	}
	fmt.Printf(Cyan+"  program %q  v%s  · encoding %s · direction %s\n"+Reset,
		in.meta["name"], in.meta["version"], in.meta["encoding"], in.meta["direction"])
}

func (in *Interp) doBara(form *Node) {
	in.exec(form)
	in.u.Exists = true
	fmt.Println(Green + "  ✓ bara — heavens and earth called into being (still unformed)" + Reset)
	pause()
}

func (in *Interp) doDefine(form *Node) {
	in.exec(form)
	if len(form.Items) > 1 && form.Items[1].Kind == Symbol {
		name := form.Items[1].Text
		in.env[name] = form
		fmt.Println(Blue + "  defined " + name + Reset)
	}
	kw := keywordArgs(form)
	if st, ok := kw[":state"]; ok {
		bar := strings.Repeat("─", 58)
		fmt.Println()
		fmt.Println(Bold + bar + Reset)
		fmt.Println(Bold + Blue + "[INIT]" + Reset + " pre-creation state: " + Cyan + st.String() + Reset)
		fmt.Println(Bold + bar + Reset)
		if s, ok := kw[":surface"]; ok {
			fmt.Println(Dim + "  ├─ surface:  " + s.String() + Reset)
		}
		if h, ok := kw[":hovering"]; ok {
			fmt.Println(Dim + "  └─ hovering: " + h.String() + Reset)
		}
	}
	pause()
}

func (in *Interp) doSpeak(form *Node) {
	in.exec(form)
	speaker := "ELOHIM"
	if len(form.Items) > 1 && form.Items[1].Kind == Symbol {
		speaker = form.Items[1].Text
	}
	fmt.Println(Magenta + "  " + speaker + " speaks — evaluating command(s):" + Reset)
	for _, it := range form.Items[2:] {
		if it.Kind == List {
			in.eval(it)
		}
	}
}

func (in *Interp) doCommand(form *Node) {
	in.exec(form)
	in.fire(form)
}

func (in *Interp) doProvision(form *Node) {
	in.exec(form)
	fmt.Println(Green + "  ✓ provision — seed-bearing plants and fruit trees given for food" + Reset)
}

func (in *Interp) doCouncil(form *Node) {
	in.exec(form)
	fmt.Println()
	fmt.Println(Bold + Red + "  ╔════════════════════════════════════════════════════════╗")
	fmt.Println("  ║  DIVINE COUNCIL — plural deliberation                    ║")
	fmt.Println("  ║  naaseh ADAM be-tzalmenu ki-dmutenu                      ║")
	fmt.Println("  ║  'Let US make ADAM in OUR image, by OUR likeness'        ║")
	fmt.Println("  ╚════════════════════════════════════════════════════════╝" + Reset)
	fmt.Println(Dim + "  · deliberation only — ADAM not yet instantiated" + Reset)
	pause()
}

func (in *Interp) doVayehi(form *Node) {
	in.exec(form)
	if d, ok := keywordArgs(form)[":day"]; ok {
		in.u.Day = atoi(d.text())
		fmt.Println(Blue + "  ↻ evening and morning — day boundary committed" + Reset)
		in.printDay()
		pause()
		return
	}
	arg := ""
	if len(form.Items) > 1 && form.Items[1].Kind == Symbol {
		arg = form.Items[1].Text
	}
	if strings.EqualFold(arg, "KEN") {
		fmt.Println(Green + "  ✓ (va-yehi KEN) → TRUE — and it was so" + Reset)
		return
	}
	in.fire(form)
	fmt.Println(Green + "  ✓ (va-yehi " + arg + ") → confirmed" + Reset)
}

func (in *Interp) doSee(form *Node) {
	in.exec(form)
	in.fire(form)
	eval := "TOV (good)"
	for _, s := range symbols(form) {
		if strings.EqualFold(s, "MEOD") {
			eval = "TOV MEOD (very good)"
		}
	}
	in.u.Evaluation = eval
	fmt.Println(Green + "  ✓ ELOHIM evaluates: " + eval + Reset)
}

func (in *Interp) doSeparate(form *Node) {
	in.exec(form)
	fmt.Println(Cyan + "  ⟂ va-yavdel — domains divided" + Reset)
	in.fire(form)
}

func (in *Interp) doName(form *Node) {
	in.exec(form)
	fmt.Println(Cyan + "  ✎ va-yikra — naming" + Reset)
	in.fire(form)
}

func (in *Interp) doMake(form *Node) {
	in.exec(form)
	in.fire(form)
}

func (in *Interp) doBless(form *Node) {
	in.exec(form)
	txt := blessingText(form)
	for _, s := range symbols(form) {
		if strings.EqualFold(s, "va-yekadesh") {
			txt += "; sanctified (set apart)"
		}
	}
	fmt.Println(Magenta + "  ✦ blessing — " + txt + Reset)
}

func (in *Interp) doFinish(form *Node) {
	in.exec(form)
	in.u.Status = "COMPLETE"
	if in.u.Day < 7 {
		in.u.Day = 7
	}
	fmt.Println(Dim + "  └─ heavens and earth finished — status: COMPLETE" + Reset)
}

func (in *Interp) doRest(form *Node) {
	in.exec(form)
	fmt.Println(Blue + "  ⏸  va-yishbot — God.rest() — process IDLE on day 7" + Reset)
	pause()
}

func (in *Interp) doReturn(form *Node) {
	in.exec(form)
	in.renderUniverse()
}

// --- helpers --------------------------------------------------------------

// keywordArgs collects :keyword value pairs from a form's items.
func keywordArgs(form *Node) map[string]*Node {
	m := map[string]*Node{}
	for i := 0; i+1 < len(form.Items); i++ {
		if it := form.Items[i]; it.Kind == Symbol && strings.HasPrefix(it.Text, ":") {
			m[it.Text] = form.Items[i+1]
		}
	}
	return m
}

// symbols returns every symbol atom in a form, recursively.
func symbols(n *Node) []string {
	var out []string
	var walk func(*Node)
	walk = func(x *Node) {
		if x.Kind == Symbol {
			out = append(out, x.Text)
		}
		for _, c := range x.Items {
			walk(c)
		}
	}
	walk(n)
	return out
}

// collectNouns gathers, in first-seen order, the noun-table keys a form uses.
func collectNouns(n *Node, seen map[string]bool, ordered *[]string) {
	if n.Kind == Symbol {
		if k := stripPrefix(n.Text); !seen[k] {
			if _, ok := nouns[k]; ok {
				seen[k] = true
				*ordered = append(*ordered, k)
			}
		}
	}
	for _, c := range n.Items {
		collectNouns(c, seen, ordered)
	}
}

// baraCount counts occurrences of the create-root (bara / yivra) in a form.
func baraCount(form *Node) int {
	n := 0
	for _, s := range symbols(form) {
		l := strings.ToLower(s)
		if strings.Contains(l, "bara") || strings.Contains(l, "yivra") {
			n++
		}
	}
	return n
}

// consonants strips vowels and joiners — a crude Hebrew root skeleton.
func consonants(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch r {
		case 'a', 'e', 'i', 'o', 'u', '-', ':', '\'':
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// detectCognates flags verb/noun pairs whose consonant skeletons nest —
// the figura etymologica the structural translation profile cares about
// (e.g. tadshe ↔ DESHE, yishretzu ↔ SHERETZ).
func detectCognates(form *Node) []string {
	syms := symbols(form)
	var out []string
	seenPair := map[string]bool{}
	for i := 0; i < len(syms); i++ {
		for j := i + 1; j < len(syms); j++ {
			a, b := syms[i], syms[j]
			if strings.EqualFold(a, b) {
				continue
			}
			ca, cb := consonants(stripPrefix(a)), consonants(stripPrefix(b))
			if len(ca) < 3 || len(cb) < 3 {
				continue
			}
			if strings.Contains(ca, cb) || strings.Contains(cb, ca) {
				key := a + "|" + b
				if !seenPair[key] {
					seenPair[key] = true
					out = append(out, a+" ↔ "+b)
				}
			}
		}
	}
	return out
}

var blessGloss = map[string]string{
	"peru": "be fruitful", "revu": "multiply", "u-revu": "multiply",
	"milu": "fill the earth", "u-milu": "fill", "kivshuha": "subdue it",
	"ve-kivshuha": "subdue it", "redu": "rule over", "u-redu": "rule over",
	"yirev": "increase",
}

// blessingText glosses the imperatives inside a blessing form.
func blessingText(form *Node) string {
	var parts []string
	seen := map[string]bool{}
	for _, s := range symbols(form) {
		if g, ok := blessGloss[strings.ToLower(s)]; ok && !seen[g] {
			seen[g] = true
			parts = append(parts, g)
		}
	}
	if len(parts) == 0 {
		return "blessed"
	}
	return strings.Join(parts, ", ")
}

func atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

// printDay renders the running universe at a day boundary.
func (in *Interp) printDay() {
	u := in.u
	bar := strings.Repeat("─", 58)
	fmt.Println()
	fmt.Println(Bold + bar + Reset)
	fmt.Printf(Bold+Blue+"[DAY %d]"+Reset+"  evaluation: "+Green+"%s"+Reset+"\n", u.Day, u.Evaluation)
	fmt.Println(Bold + bar + Reset)
	rows := []struct {
		k string
		v bool
	}{
		{"light", u.Light}, {"sky", u.Sky}, {"land", u.Land}, {"seas", u.Seas},
		{"plants", u.Plants}, {"sun", u.Sun}, {"moon", u.Moon}, {"stars", u.Stars},
		{"fish", u.Fish}, {"birds", u.Birds}, {"animals", u.Animals}, {"humans", u.Humans},
	}
	for _, r := range rows {
		mark, val := Dim+"·"+Reset, Dim+"false"+Reset
		if r.v {
			mark, val = Green+"✓"+Reset, Green+"true"+Reset
		}
		fmt.Printf("  %s %-8s %s\n", mark, r.k, val)
	}
}

// renderUniverse prints the final returned struct and sets exit status.
func (in *Interp) renderUniverse() {
	u := in.u
	fmt.Println()
	fmt.Println(Bold + White + BgBlack)
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    EXECUTION COMPLETE                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println(Reset)

	fmt.Println(Green + "UNIVERSE := {" + Reset)
	row := func(k string, v bool, col string) {
		fmt.Printf("  %-9s %s%v%s\n", k+":", col, v, Reset)
	}
	row("exists", u.Exists, Green)
	row("light", u.Light, Yellow)
	row("sky", u.Sky, Cyan)
	row("land", u.Land, Green)
	row("seas", u.Seas, Blue)
	row("plants", u.Plants, Green)
	row("sun", u.Sun, Yellow)
	row("moon", u.Moon, White)
	row("stars", u.Stars, Cyan)
	row("fish", u.Fish, Blue)
	row("birds", u.Birds, Cyan)
	row("animals", u.Animals, Yellow)
	fmt.Printf("  %-9s %s%v%s  ← image_of(ELOHIM)\n", "humans:", Magenta, u.Humans, Reset)
	fmt.Printf("  %-9s %s%d%s\n", "days:", Bold, u.Day, Reset)
	fmt.Printf("  %-9s %s%s%s\n", "status:", Green+Bold, u.Status, Reset)
	fmt.Printf("  %-9s %s%s%s\n", "eval:", Green+Bold, u.Evaluation, Reset)
	fmt.Println(Green + "}" + Reset)
	fmt.Println()

	if missing := in.missing(); len(missing) == 0 && u.Status == "COMPLETE" {
		fmt.Println(Dim + "// process exited 0 — universe running." + Reset)
		fmt.Println(Dim + "// to inspect: use conscience, prayer, or telescope." + Reset)
	} else {
		fmt.Println(Red + "// process exited 1 — universe incomplete: " +
			strings.Join(missing, ", ") + Reset)
	}
}

// missing lists creation fields that never got set.
func (in *Interp) missing() []string {
	u := in.u
	checks := []struct {
		k string
		v bool
	}{
		{"exists", u.Exists}, {"light", u.Light}, {"sky", u.Sky}, {"land", u.Land},
		{"seas", u.Seas}, {"plants", u.Plants}, {"sun", u.Sun}, {"moon", u.Moon},
		{"stars", u.Stars}, {"fish", u.Fish}, {"birds", u.Birds},
		{"animals", u.Animals}, {"humans", u.Humans},
	}
	var out []string
	for _, c := range checks {
		if !c.v {
			out = append(out, c.k)
		}
	}
	return out
}

// exitCode is 0 when the program built a complete universe.
func (in *Interp) exitCode() int {
	if len(in.missing()) == 0 && in.u.Status == "COMPLETE" {
		return 0
	}
	return 1
}
