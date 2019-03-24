package promptui

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/karantin2020/readline"
)

var (
	// ErrorIncorrect describes incorrect default confirm value
	ErrorIncorrect = errors.New("Incorrect default values")
)

// ConfirmPrompt represents a single line text field input.
type ConfirmPrompt struct {
	BasicPrompt

	// ConfirmOpt is the 3d answer option
	ConfirmOpt string

	confirmDefault string
}

// Run func implements confirn prompt
func (cp *ConfirmPrompt) Run() (string, error) {
	switch cp.Default {
	case "Y", "N", "n", "y":
	case "":
		cp.Default = "N"
	default:
		return "", ErrorIncorrect
	}
	err := cp.Init()
	if err != nil {
		return "", err
	}

	cp.confirmDefault = strings.ToUpper(cp.Default)

	cp.c.Stdin = ioutil.NopCloser(io.MultiReader(bytes.NewBuffer([]byte(cp.out)), os.Stdin))

	cp.rl, err = readline.NewEx(cp.c)
	if err != nil {
		return "", err
	}

	cp.punctuation = "?"
	answers := "y/N"
	if strings.ToLower(cp.Default) == "y" {
		answers = "Y/n"
	}
	if cp.ConfirmOpt != "" {
		answers = answers + "/" + cp.ConfirmOpt
	}
	cp.suggestedAnswer = " " + faint("["+answers+"]")
	// cp.confirmDefault = strings.ToUpper(cp.Default)
	// cp.Default = ""
	cp.prompt = cp.LabelInitial(cp.Label) + cp.punctuation + cp.suggestedAnswer + " "
	// cp.out = cp.Default
	// cp.c.Stdin = ioutil.NopCloser(io.MultiReader(bytes.NewBuffer([]byte(cp.out)), os.Stdin))

	setupConfirm(cp.c, cp.prompt, cp, cp.rl)
	cp.out, err = cp.rl.Readline()
	if cp.out == "" {
		cp.out = cp.confirmDefault
	}
	if err != nil {
		if err.Error() == "Interrupt" {
			err = ErrInterrupt
		}
		cp.rl.Write([]byte("\n"))
		return "", err
	}
	cp.out = strings.ToUpper(cp.out)
	cp.state = cp.IconGood
	cp.out = cp.Formatter(cp.out)
	cp.rl.Write([]byte(cp.Indent + cp.state + " " + cp.prompt + cp.InputResult(cp.out) + "\n"))
	return cp.out, err
}

func setupConfirm(c *readline.Config, prompt string,
	cp *ConfirmPrompt, rl *readline.Instance) {
	filterInput := func(r rune) (rune, bool) {
		switch r {
		case readline.CharCtrlZ,
			readline.CharInterrupt,
			readline.CharEnter,
			readline.CharBackward,
			readline.CharForward,
			readline.CharNext,
			readline.CharPrev,
			readline.CharBackspace:
			break
		case 'Y', 'y', 'N', 'n':
			break
		default:
			if len(cp.ConfirmOpt) > 0 {
				if r == []rune(cp.ConfirmOpt)[0] ||
					r == []rune(strings.ToUpper(cp.ConfirmOpt))[0] {
					break
				}
			}
			return r, false
		}
		return r, true
	}
	c.FuncFilterInputRune = filterInput
	var confirmReader = func(line []rune, pos int, key rune) ([]rune, int, bool) {
		if key == readline.CharEnter {
			return nil, 0, false
		}
		if key == readline.CharBackward ||
			key == readline.CharForward ||
			key == readline.CharNext ||
			key == readline.CharPrev {
			return nil, 0, false
		}

		var state string
		var separator = " "
		if cp.NoIcons {
			separator = ""
		}

		switch key {
		case 'Y', 'y', 'N', 'n':
			if len(line) > 1 {
				return line[0:1], 1, true
			}
		default:
			if cp.ConfirmOpt != "" {
				if len(line) > 1 {
					return line[0:1], 1, true
				}
				if key == []rune(cp.ConfirmOpt)[0] ||
					key == []rune(strings.ToUpper(cp.ConfirmOpt))[0] {
					break
				}
			}
			state = cp.IconInitial
			rl.SetPrompt(cp.Indent + state + separator + cp.PromptInitial(prompt))
			rl.Refresh()
			return []rune(""), 0, true
		}

		if string(line) == "" {
			state = cp.IconInitial
		} else {
			state = cp.IconGood
		}

		rl.SetPrompt(cp.Indent + state + separator + cp.PromptInitial(prompt))
		rl.Refresh()

		return nil, 0, false
	}
	c.SetListener(confirmReader)
}

// Confirm func is predefined easy confirm promp
var Confirm = func(label, answer string, noIcons bool) (string, error) {
	cp := ConfirmPrompt{
		BasicPrompt: BasicPrompt{
			Label:   label,
			Default: answer,
			NoIcons: noIcons,
		},
	}
	return cp.Run()
}
