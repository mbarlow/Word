// Genesis DSL Interpreter
// Executes Hebrew creation commands and shows the universe being built
package main

import (
	"fmt"
	"strings"
	"time"
)

// ANSI colors for dramatic effect
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	BgBlack = "\033[40m"
)

// Universe state
type Universe struct {
	Exists     bool
	Light      bool
	Sky        bool
	Land       bool
	Seas       bool
	Plants     bool
	Sun        bool
	Moon       bool
	Stars      bool
	Fish       bool
	Birds      bool
	Animals    bool
	Humans     bool
	Day        int
	Status     string
	Evaluation string
}

func main() {
	fmt.Println(BgBlack + White + Bold)
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          GENESIS.DSL — The Creation Executable               ║")
	fmt.Println("║              בְּרֵאשִׁית בָּרָא אֱלֹהִים                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println(Reset)

	u := &Universe{}

	// Pre-creation state
	printState("INIT", "Before execution: TOHU VA-VOHU (formless and void)", u)
	fmt.Println(Dim + "  ├─ darkness: TRUE")
	fmt.Println("  ├─ deep: INFINITE")
	fmt.Println("  └─ spirit_hovering: TRUE" + Reset)
	pause()

	// Day 1: Light
	u.Day = 1
	execute("(yomer ELOHIM (yehi OR))", "וַיֹּאמֶר אֱלֹהִים יְהִי אוֹר")
	output("(va-yehi OR)", "Light instantiated")
	u.Light = true
	u.Evaluation = "TOV (Good)"
	printState("DAY 1", "LIGHT separated from DARKNESS", u)
	pause()

	// Day 2: Sky
	u.Day = 2
	execute("(yomer ELOHIM (yehi RAKIA))", "יְהִי רָקִיעַ בְּתוֹךְ הַמָּיִם")
	output("(va-yehi KEN)", "Firmament created, waters divided")
	u.Sky = true
	u.Evaluation = "TOV"
	printState("DAY 2", "SKY divides waters above/below", u)
	pause()

	// Day 3: Land & Plants
	u.Day = 3
	execute("(yomer ELOHIM (yikavu HA-MAYIM))", "יִקָּווּ הַמַּיִם...וְתֵרָאֶה הַיַּבָּשָׁה")
	output("(va-yehi KEN)", "Dry land appeared")
	u.Land = true
	u.Seas = true

	execute("(yomer ELOHIM (tadshe HA-ERETZ DESHE))", "תַּדְשֵׁא הָאָרֶץ דֶּשֶׁא")
	fmt.Println(Green + "  └─ COGNATE DETECTED: תַּדְשֵׁא...דֶּשֶׁא (vegetate→vegetation)" + Reset)
	fmt.Println(Green + "     └─ Recursive creation: plant.seed(seed) → self-replicating" + Reset)
	output("(va-yehi KEN)", "Vegetation, plants, trees instantiated")
	u.Plants = true
	u.Evaluation = "TOV"
	printState("DAY 3", "LAND + SEAS + PLANTS", u)
	pause()

	// Day 4: Celestial bodies
	u.Day = 4
	execute("(yomer ELOHIM (yehi MEOROT))", "יְהִי מְאֹרֹת בִּרְקִיעַ הַשָּׁמַיִם")
	output("(va-yaas ELOHIM)", "Lights created")
	u.Sun = true
	u.Moon = true
	u.Stars = true
	u.Evaluation = "TOV"
	printState("DAY 4", "SUN + MOON + STARS", u)
	fmt.Println(Yellow + "  ├─ sun:   { role: 'rule_day',   type: 'MAOR_GADOL' }")
	fmt.Println("  ├─ moon:  { role: 'rule_night', type: 'MAOR_KATAN' }")
	fmt.Println("  └─ stars: { count: 'uncountable' }" + Reset)
	pause()

	// Day 5: Sea & Sky creatures
	u.Day = 5
	execute("(yomer ELOHIM (yishretzu HA-MAYIM))", "יִשְׁרְצוּ הַמַּיִם שֶׁרֶץ נֶפֶשׁ חַיָּה")
	fmt.Println(Cyan + "  └─ NOTE: Using BARA (create ex nihilo) for NEPHESH (soul/life)" + Reset)
	output("(va-yivra ELOHIM)", "Sea creatures + birds created")
	u.Fish = true
	u.Birds = true
	u.Evaluation = "TOV"
	execute("(va-yevarekh ELOHIM)", "וַיְבָרֶךְ אֹתָם אֱלֹהִים")
	fmt.Println(Magenta + "  └─ BLESSING: peru u-revu (be fruitful, multiply)" + Reset)
	printState("DAY 5", "SEA CREATURES + BIRDS", u)
	pause()

	// Day 6: Animals & Humans
	u.Day = 6
	execute("(yomer ELOHIM (totze HA-ERETZ))", "תּוֹצֵא הָאָרֶץ נֶפֶשׁ חַיָּה")
	output("(va-yehi KEN)", "Land animals created")
	u.Animals = true

	fmt.Println()
	fmt.Println(Bold + Red + "  ╔════════════════════════════════════════════════════════╗")
	fmt.Println("  ║  DIVINE COUNCIL DELIBERATION                           ║")
	fmt.Println("  ║  נַעֲשֶׂה אָדָם בְּצַלְמֵנוּ כִּדְמוּתֵנוּ                            ║")
	fmt.Println("  ║  'Let US make human in OUR image'                      ║")
	fmt.Println("  ╚════════════════════════════════════════════════════════╝" + Reset)
	pause()

	execute("(va-yivra ELOHIM et HA-ADAM)", "וַיִּבְרָא אֱלֹהִים אֶת הָאָדָם בְּצַלְמוֹ")
	fmt.Println(Magenta + "  ├─ BARA (create): Used 3x for humans — maximum emphasis")
	fmt.Println("  ├─ be-tzelem: in_image(ELOHIM)")
	fmt.Println("  └─ zakhar u-nekevah: male AND female" + Reset)
	u.Humans = true

	execute("(va-yevarekh ELOHIM otam)", "וַיְבָרֶךְ אֹתָם אֱלֹהִים")
	fmt.Println(Magenta + "  └─ BLESSING: {")
	fmt.Println("       peru_u_revu: true,      // be fruitful")
	fmt.Println("       milu_et_ha_aretz: true, // fill earth")
	fmt.Println("       kivshuha: true,         // subdue it")
	fmt.Println("       redu: ['fish','birds','animals']  // rule over")
	fmt.Println("     }" + Reset)

	u.Evaluation = "TOV MEOD (Very Good)"
	printState("DAY 6", "ANIMALS + HUMANS (image of God)", u)
	fmt.Println(Bold + Green + "  └─ EVALUATION: " + u.Evaluation + " ← First 'VERY good'!" + Reset)
	pause()

	// Day 7: Rest
	u.Day = 7
	fmt.Println()
	fmt.Println(Bold + Blue + "╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        DAY 7: SHABBAT                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝" + Reset)

	execute("(va-yekhulu HA-SHAMAYIM ve-HA-ARETZ)", "וַיְכֻלּוּ הַשָּׁמַיִם וְהָאָרֶץ")
	fmt.Println(Dim + "  └─ Status: COMPLETE" + Reset)

	execute("(va-yishbot ba-YOM HA-SHEVII)", "וַיִּשְׁבֹּת בַּיּוֹם הַשְּׁבִיעִי")
	fmt.Println(Blue + "  └─ God.rest() — process IDLE" + Reset)

	execute("(va-yevarekh... va-yekadesh)", "וַיְבָרֶךְ אֱלֹהִים אֶת יוֹם הַשְּׁבִיעִי וַיְקַדֵּשׁ אֹתוֹ")
	fmt.Println(Magenta + "  └─ Day 7 blessed and sanctified (set apart)" + Reset)

	u.Status = "COMPLETE"
	u.Exists = true

	// Final output
	fmt.Println()
	fmt.Println(Bold + White + BgBlack)
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    EXECUTION COMPLETE                        ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║  (return UNIVERSE)                                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println(Reset)

	fmt.Println(Green + "UNIVERSE := {" + Reset)
	fmt.Printf("  exists:  %s%v%s\n", Green, u.Exists, Reset)
	fmt.Printf("  light:   %s%v%s\n", Yellow, u.Light, Reset)
	fmt.Printf("  sky:     %s%v%s\n", Cyan, u.Sky, Reset)
	fmt.Printf("  land:    %s%v%s\n", Green, u.Land, Reset)
	fmt.Printf("  seas:    %s%v%s\n", Blue, u.Seas, Reset)
	fmt.Printf("  plants:  %s%v%s\n", Green, u.Plants, Reset)
	fmt.Printf("  sun:     %s%v%s\n", Yellow, u.Sun, Reset)
	fmt.Printf("  moon:    %s%v%s\n", White, u.Moon, Reset)
	fmt.Printf("  stars:   %s%v%s\n", Cyan, u.Stars, Reset)
	fmt.Printf("  fish:    %s%v%s\n", Blue, u.Fish, Reset)
	fmt.Printf("  birds:   %s%v%s\n", Cyan, u.Birds, Reset)
	fmt.Printf("  animals: %s%v%s\n", Yellow, u.Animals, Reset)
	fmt.Printf("  humans:  %s%v%s  ← image_of(ELOHIM)\n", Magenta, u.Humans, Reset)
	fmt.Printf("  status:  %s%s%s\n", Green+Bold, u.Status, Reset)
	fmt.Printf("  eval:    %s%s%s\n", Green+Bold, u.Evaluation, Reset)
	fmt.Println(Green + "}" + Reset)

	fmt.Println()
	fmt.Println(Dim + "// Process exited. Universe running." + Reset)
	fmt.Println(Dim + "// To inspect: use conscience, prayer, or telescope." + Reset)
}

func execute(lisp, hebrew string) {
	fmt.Println()
	fmt.Print(Yellow + "► EXEC: " + Reset)
	fmt.Println(Dim + lisp + Reset)
	fmt.Println(Cyan + "  " + hebrew + Reset)
}

func output(lisp, desc string) {
	fmt.Println(Green + "  ✓ " + lisp + " → " + desc + Reset)
}

func printState(label, desc string, u *Universe) {
	fmt.Println()
	bar := strings.Repeat("─", 50)
	fmt.Println(Bold + bar)
	fmt.Printf("%s[%s]%s %s\n", Bold+Blue, label, Reset, desc)
	fmt.Println(bar + Reset)
}

func pause() {
	time.Sleep(100 * time.Millisecond)
}
