# PromptUI

[![GoDoc](https://godoc.org/github.com/karantin2020/promptui?status.svg)](github.com/karantin2020/promptui) [![Go Report Card](https://goreportcard.com/badge/github.com/karantin2020/promptui)](https://goreportcard.com/report/github.com/karantin2020/promptui) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A library for building interactive prompts.


![demo](https://github.com/karantin2020/promptui/raw/master/examples/screen.gif)

See `examples/main.go`.  

```go
package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	bb "github.com/karantin2020/promptui"
)

func main() {
	p := bb.Prompt{
		BasicPrompt: bb.BasicPrompt{
			Label: "test masked",
			Validate: func(s string) error {
				if s == "err" {
					return bb.NewValidationError("input must not be equal to 'err'")
				}
				if s == "" {
					return bb.NewValidationError("input must not be empty string")
				}
				return nil
			},
			Default: "hi there",
		},
		Mask: '*',
		// IsVimMode: false,
	}
	checkError(p.Run())
	checkError(bb.Ask("easy prompt", "Yepp"))
	checkError(bb.AskMasked("easy masked prompt", "Yepp"))
	c := bb.ConfirmPrompt{
		BasicPrompt: bb.BasicPrompt{
			Label:   "test confirm",
			Default: "y",
		},
		ConfirmOpt: "?",
	}
	checkError(c.Run())
	checkError(bb.Confirm("easy confirm", "Y", false))

	ml := bb.MultilinePrompt{
		BasicPrompt: bb.BasicPrompt{
			Label:   "test multiline",
			Default: "Request",
			Formatter: func(s string) string {
				var upl = func(sl string) string {
					rs := []rune(sl)
					if len(rs) > 0 {
						rs[0] = unicode.ToUpper(rs[0])
					}
					return string(rs)
				}
				out := []string{}
				ins := strings.Split(s, "\n")
				for i := range ins {
					out = append(out, upl(ins[i]))
				}
				return strings.Join(out, "\n")
			},
			Validate: func(s string) error {
				if s == "Err" {
					return bb.NewValidationError("input must not be equal to 'err'")
				}
				if s == "" {
					return bb.NewValidationError("input must not be empty string")
				}
				return nil
			},
		},
	}
	checkError(ml.Run())
	checkError(bb.MultiLine("easy multi", "Ready to go"))

	s := bb.Select{
		Label:   "test select",
		Items:   []string{"one", "two", "three", "1one", "1two", "1three", "2one", "2two", "2three"},
		Default: 3,
	}
	_, res, err := s.Run()
	checkError(res, err)
	checkError(bb.PromptAfterSelect("prompt after select", []string{"feat", "fix", "doc", "other"}))
}
```