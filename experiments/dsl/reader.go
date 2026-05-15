// reader.go — lexer + parser for the Genesis DSL.
//
// The DSL is a Lisp-like s-expression language: parenthesised lists of
// symbols, string literals, and nested lists. Comments run from ';' to EOL.
// This file turns raw source text into a slice of top-level *Node forms.
package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Kind tags what a Node holds.
type Kind int

const (
	List   Kind = iota // Items is meaningful
	Symbol             // Text is a bare atom (verb, noun, :keyword)
	Str                // Text is a "quoted string" literal
)

// Node is one element of the parsed program: an atom or a list.
type Node struct {
	Kind  Kind
	Text  string  // Symbol / Str payload
	Items []*Node // List children
}

// String renders a Node back to canonical s-expression text.
func (n *Node) String() string {
	switch n.Kind {
	case Symbol:
		return n.Text
	case Str:
		return fmt.Sprintf("%q", n.Text)
	default:
		parts := make([]string, len(n.Items))
		for i, it := range n.Items {
			parts[i] = it.String()
		}
		return "(" + strings.Join(parts, " ") + ")"
	}
}

// text returns a Node's displayable value (unquoted for strings).
func (n *Node) text() string {
	if n.Kind == List {
		return n.String()
	}
	return n.Text
}

// head returns the leading symbol of a list, or "" for anything else.
func (n *Node) head() string {
	if n.Kind == List && len(n.Items) > 0 && n.Items[0].Kind == Symbol {
		return n.Items[0].Text
	}
	return ""
}

type tokKind int

const (
	tLParen tokKind = iota
	tRParen
	tAtom
	tString
)

type token struct {
	kind tokKind
	text string
}

// tokenize splits source into tokens, discarding ';' comments.
func tokenize(src string) ([]token, error) {
	var toks []token
	rs := []rune(src)
	for i := 0; i < len(rs); {
		c := rs[i]
		switch {
		case c == ';': // comment to end of line
			for i < len(rs) && rs[i] != '\n' {
				i++
			}
		case unicode.IsSpace(c):
			i++
		case c == '(':
			toks = append(toks, token{tLParen, "("})
			i++
		case c == ')':
			toks = append(toks, token{tRParen, ")"})
			i++
		case c == '"':
			i++
			start := i
			for i < len(rs) && rs[i] != '"' {
				i++
			}
			if i >= len(rs) {
				return nil, fmt.Errorf("unterminated string literal")
			}
			toks = append(toks, token{tString, string(rs[start:i])})
			i++ // consume closing quote
		default:
			start := i
			for i < len(rs) && !unicode.IsSpace(rs[i]) && rs[i] != '(' && rs[i] != ')' && rs[i] != ';' {
				i++
			}
			toks = append(toks, token{tAtom, string(rs[start:i])})
		}
	}
	return toks, nil
}

// parse reads source into a slice of top-level forms.
func parse(src string) ([]*Node, error) {
	toks, err := tokenize(src)
	if err != nil {
		return nil, err
	}

	pos := 0
	var readForm func() (*Node, error)
	readForm = func() (*Node, error) {
		if pos >= len(toks) {
			return nil, fmt.Errorf("unexpected end of input")
		}
		t := toks[pos]
		switch t.kind {
		case tLParen:
			pos++
			list := &Node{Kind: List}
			for {
				if pos >= len(toks) {
					return nil, fmt.Errorf("unterminated list")
				}
				if toks[pos].kind == tRParen {
					pos++
					return list, nil
				}
				child, err := readForm()
				if err != nil {
					return nil, err
				}
				list.Items = append(list.Items, child)
			}
		case tRParen:
			return nil, fmt.Errorf("unexpected ')'")
		case tString:
			pos++
			return &Node{Kind: Str, Text: t.text}, nil
		default:
			pos++
			return &Node{Kind: Symbol, Text: t.text}, nil
		}
	}

	var forms []*Node
	for pos < len(toks) {
		f, err := readForm()
		if err != nil {
			return nil, err
		}
		forms = append(forms, f)
	}
	return forms, nil
}
