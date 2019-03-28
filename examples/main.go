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
	// p.InputResult = bb.Styler(bb.FGCyan)
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
	// c.InputResult = bb.Styler(bb.FGCyan)
	checkError(c.Run())
	checkError(bb.Confirm("easy confirm", "Y", false))

	s := bb.Select{
		Label:   "test select",
		Items:   []string{"one", "two", "three", "1one", "1two", "1three", "2one", "2two", "2three"},
		Default: 3,
	}
	_, res, err := s.Run()
	checkError(res, err)
	checkError(bb.PromptAfterSelect("prompt after select", []string{"feat", "fix", "doc", "other"}))

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
				// fmt.Println("Validate multiline")
				if s == "Err" {
					// fmt.Println("input must not be equal to 'err'")
					return bb.NewValidationError("input must not be equal to 'err'")
				}
				if s == "" {
					// fmt.Println("input must not be empty string")
					return bb.NewValidationError("input must not be empty string")
				}
				return nil
			},
		},
	}
	// ml.InputResult = bb.Styler(bb.FGCyan)
	checkError(ml.Run())

	mle := bb.MultilinePrompt{
		BasicPrompt: bb.BasicPrompt{
			Label:   "test multiline onError",
			Default: "Request",
			Validate: func(s string) error {
				// fmt.Println("Validate multiline")
				if s == "Err" {
					// fmt.Println("input must not be equal to 'err'")
					return bb.NewValidationError("input must not be equal to 'err'")
				}
				if s == "" {
					// fmt.Println("input must not be empty string")
					return bb.NewValidationError("input must not be empty string")
				}
				return nil
			},
		},
		OnError: func(s string) (string, error) {
			// numlines := 0
			defer func() {
				fmt.Print(bb.ClearUpLines(1))
			}()
			for {
				yn, oerr := bb.Confirm("Edit test multiline", "", true)
				// numlines++
				if oerr != nil {
					return s, oerr
				}
				if yn == "Y" {
					s, oerr = bb.Editor("", s)
					if oerr != nil {
						// bb.ClearUpLines(numlines)
						return s, oerr
					}
					fmt.Print(bb.ClearUpLines(1))
					continue
				} else {
					// fmt.Print(bb.ClearUpLines(1))
					break
				}
			}
			return s, nil
		},
	}
	// mle.InputResult = bb.Styler(bb.FGCyan)
	checkError(mle.Run())

	checkError(bb.MultiLine("easy multi", "Ready to go"))

}

func checkError(res string, err error) {
	if err != nil {
		if err == bb.ErrInterrupt {
			fmt.Println()
			os.Exit(1)
		}
		fmt.Println(err)
		// os.Exit(1)
	} else {
		fmt.Println("User input:", res)
	}
}
