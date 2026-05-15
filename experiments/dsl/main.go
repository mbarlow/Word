// Genesis DSL interpreter.
//
// Reads a Lisp-like creation program (experiments/genesis.dsl), parses it
// into s-expressions, and evaluates it form by form — building a Universe
// from the Hebrew verbs and nouns it actually encounters.
//
//	go run ./experiments/dsl            # auto-locates genesis.dsl
//	go run ./experiments/dsl path.dsl   # run a specific program
//	go run ./experiments/dsl -fast      # skip the dramatic pauses
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// ANSI colors for dramatic effect.
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	BgBlack = "\033[40m"
)

var fast bool

func main() {
	var path string
	for _, a := range os.Args[1:] {
		if a == "-fast" || a == "--fast" {
			fast = true
			continue
		}
		path = a
	}

	src, resolved, err := loadSource(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
	}

	forms, err := parse(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error: "+err.Error())
		os.Exit(1)
	}

	banner()
	fmt.Println(Dim + "// source: " + resolved + Reset)
	fmt.Printf(Dim+"// %d top-level forms parsed\n"+Reset, len(forms))
	pause()

	in := newInterp()
	for _, f := range forms {
		in.eval(f)
	}

	os.Exit(in.exitCode())
}

// loadSource resolves the program path, trying sensible fallbacks so the
// interpreter runs from any working directory.
func loadSource(arg string) (string, string, error) {
	var candidates []string
	if arg != "" {
		candidates = append(candidates, arg,
			filepath.Join("experiments", arg),
			filepath.Join("experiments", "dsl", arg))
	}
	if _, file, _, ok := runtime.Caller(0); ok {
		// main.go lives in experiments/dsl/ — genesis.dsl is one up.
		candidates = append(candidates, filepath.Join(filepath.Dir(file), "..", "genesis.dsl"))
	}
	candidates = append(candidates, "experiments/genesis.dsl", "genesis.dsl", "../genesis.dsl")

	for _, c := range candidates {
		if b, err := os.ReadFile(c); err == nil {
			abs, _ := filepath.Abs(c)
			return string(b), abs, nil
		}
	}
	return "", "", fmt.Errorf("could not locate a .dsl program (tried %d paths)", len(candidates))
}

func banner() {
	fmt.Println(BgBlack + White + Bold)
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          GENESIS.DSL — The Creation Interpreter              ║")
	fmt.Println("║          read · parse · evaluate · instantiate              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println(Reset)
}

func pause() {
	if !fast {
		time.Sleep(120 * time.Millisecond)
	}
}
