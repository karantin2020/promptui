package promptui

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/karantin2020/readline"
)

// IconSet struct holds all inner icons to render
type IconSet struct {
	// IconInitial icon to render
	IconInitial string
	// IconGood icon to render
	IconGood string
	// IconQuest icon to render
	IconQuest string
	// IconWarn icon to render
	IconWarn string
	// IconBad icon to render
	IconBad string
}

// LabelStyle contains Label styles
type LabelStyle struct {
	// LabelInitial styler for initial Label text
	LabelInitial StyleFn
	// LabelResult styler for result Label text
	LabelResult StyleFn
}

// PromptStyle contains Prompt styles for including label, punctuation, suggested
type PromptStyle struct {
	// PromptInitial styler for initial Prompt text
	PromptInitial StyleFn
	// PromptResult styler for result Prompt text
	PromptResult StyleFn
}

// InputStyle contains Input styles
type InputStyle struct {
	// InputInitial styler for initial Input text
	InputInitial StyleFn
	// InputResult styler for result Input text
	InputResult StyleFn
}

// BasicPrompt represents a single line text field input.
type BasicPrompt struct {
	Label string // Label is the value displayed on the command line prompt

	Default string // Default is the initial value to populate in the prompt

	// Validate is optional. If set, this function is used to validate the input
	// after each character entry.
	Validate ValidateFunc

	// Indent will be placed before the prompt's state symbol
	Indent string

	// InterruptPrompt to send to readline
	InterruptPrompt string

	// NoIcons flag to set empty string icons
	NoIcons bool
	// IconSet contains prompt icons
	IconSet
	// LabelStyle contains Label styles
	LabelStyle
	// InputStyle contains Input styles
	InputStyle
	// PromptStyle contains Prompt styles
	PromptStyle

	// IsVimMode option
	IsVimMode bool
	// Preamble option
	Preamble *string

	// Formatter formats input result
	Formatter StyleFn

	stdin           io.Reader
	stdout          io.Writer
	c               *readline.Config
	rl              *readline.Instance
	suggestedAnswer string
	punctuation     string
	state           string
	prompt          string
	validFn         func(string) error
	out             string
}

// Init func to setup BasicPrompt
func (bp *BasicPrompt) Init() error {
	bp.c = &readline.Config{}

	err := bp.c.Init()
	if err != nil {
		return err
	}

	if bp.stdin != nil {
		bp.c.Stdin = ioutil.NopCloser(bp.stdin)
	}

	if bp.stdout != nil {
		bp.c.Stdout = bp.stdout
	}

	if bp.IsVimMode {
		bp.c.VimMode = true
	}

	if bp.Preamble != nil {
		fmt.Println(*bp.Preamble)
	}

	if bp.IconInitial == "" && !bp.NoIcons {
		bp.IconInitial = bold(IconInitial)
	}
	if bp.IconGood == "" && !bp.NoIcons {
		bp.IconGood = bold(IconGood)
	}
	if bp.IconQuest == "" && !bp.NoIcons {
		bp.IconQuest = bold(IconQuest)
	}
	if bp.IconWarn == "" && !bp.NoIcons {
		bp.IconWarn = bold(IconWarn)
	}
	if bp.IconBad == "" && !bp.NoIcons {
		bp.IconBad = bold(IconBad)
	}
	if bp.LabelInitial == nil {
		bp.LabelInitial = func(s string) string { return s }
	}
	if bp.LabelResult == nil {
		bp.LabelResult = func(s string) string { return s }
	}
	if bp.PromptInitial == nil {
		bp.PromptInitial = func(s string) string { return bold(s) }
	}
	if bp.PromptResult == nil {
		bp.PromptResult = func(s string) string { return s }
	}
	if bp.InputInitial == nil {
		bp.InputInitial = func(s string) string { return s }
	}
	if bp.InputResult == nil {
		bp.InputResult = func(s string) string { return faint(s) }
	}
	if bp.Formatter == nil {
		bp.Formatter = func(s string) string { return s }
	}
	bp.c.Painter = &defaultPainter{style: bp.InputInitial}

	bp.suggestedAnswer = ""
	bp.punctuation = ":"
	bp.c.UniqueEditLine = true

	bp.state = bp.IconInitial
	bp.prompt = bp.LabelInitial(bp.Label) + bp.punctuation + bp.suggestedAnswer + " "

	bp.c.Prompt = bp.Indent + bp.state + " " + bp.PromptInitial(bp.prompt)
	bp.c.HistoryLimit = -1

	bp.c.InterruptPrompt = bp.InterruptPrompt

	bp.validFn = func(x string) error {
		return nil
	}

	if bp.Validate != nil {
		bp.validFn = bp.Validate
	}

	return nil
}

// Prompt represents a single line text field input.
type Prompt struct {
	BasicPrompt

	// If mask is set, this value is displayed instead of the actual input
	// characters.
	Mask rune

	// Handlers catch key input events, user defined
	Handlers []func(line []rune, pos int, key rune) ([]rune, int, bool)
}

// Run runs the prompt, returning the validated input.
func (p *Prompt) Run() (string, error) {
	err := p.Init()
	if err != nil {
		return "", err
	}

	p.c.Stdin = ioutil.NopCloser(io.MultiReader(bytes.NewBuffer([]byte(p.Default)), os.Stdin))

	p.rl, err = readline.NewEx(p.c)
	if err != nil {
		return "", err
	}

	var (
		firstListen = true
		wroteErr    = false
		caughtup    = true
	)

	if p.Default != "" {
		caughtup = false
	}

	if p.Mask != 0 {
		p.c.EnableMask = true
		p.c.MaskRune = p.Mask
	}

	var onelineReader = func(line []rune, pos int, key rune) ([]rune, int, bool) {
		if key == readline.CharEnter {
			return nil, 0, false
		}
		if p.Mask != 0 && key == 42 {
			p.c.EnableMask = !p.c.EnableMask
			p.rl.Refresh()
			return append(line[:pos-1], line[pos:]...), pos - 1, true
		}

		if firstListen {
			firstListen = false
			return nil, 0, false
		}

		if !caughtup && p.out != "" {
			if string(line) == p.out {
				caughtup = true
			}
			if wroteErr {
				return nil, 0, false
			}
		}

		err := p.validFn(string(line))
		if err != nil {
			if _, ok := err.(*ValidationError); ok {
				p.state = p.IconBad
			} else {
				p.rl.Close()
				return nil, 0, false
			}
		} else {
			if string(line) == "" {
				p.state = p.IconInitial
			} else {
				p.state = p.IconGood
			}
		}

		p.rl.SetPrompt(p.Indent + p.state + " " + p.PromptInitial(p.prompt))
		p.rl.Refresh()
		wroteErr = false

		return nil, 0, false
	}

	p.c.SetListener(onelineReader)

	for {
		p.out, err = p.rl.Readline()

		var msg string
		valid := true
		oerr := p.validFn(p.out)
		if oerr != nil {
			if verr, ok := oerr.(*ValidationError); ok {
				msg = verr.msg
				valid = false
				p.state = p.IconBad
			} else {
				p.rl.Close()
				return "", oerr
			}
		}

		if valid {
			p.state = p.IconGood
			break
		}

		if err != nil {
			switch err {
			case readline.ErrInterrupt:
				err = ErrInterrupt
			case io.EOF:
				err = ErrEOF
			}

			break
		}

		caughtup = false

		p.c.Stdin = ioutil.NopCloser(io.MultiReader(bytes.NewBuffer([]byte(p.out)), os.Stdin))
		p.rl, _ = readline.NewEx(p.c)

		firstListen = true
		wroteErr = true
		p.rl.SetPrompt("\n" + red("Error: ") + msg + upLine(1) + "\r" + p.Indent + p.state + " " + p.PromptInitial(p.prompt))
		p.rl.Refresh()
	}

	// if wroteErr {
	// 	rl.Write([]byte(downLine(1) + clearLine + upLine(1) + "\r"))
	// }

	if err != nil {
		if err.Error() == "Interrupt" {
			err = ErrInterrupt
		}
		p.rl.Write([]byte("\n"))
		return "", err
	}

	echo := p.out
	if p.Mask != 0 {
		echo = strings.Repeat(string(p.Mask), len([]rune(echo)))
	}

	p.out = p.Formatter(p.out)
	p.rl.Write([]byte(p.Indent + p.state + " " + p.prompt + p.InputResult(echo) + "\n"))

	return p.out, err
}

type defaultPainter struct {
	style StyleFn
}

func (p *defaultPainter) Paint(line []rune, _ int) []rune {
	return []rune(p.style(string(line)))
}

// Ask func is default prompt
func Ask(label, startString string) (string, error) {
	p := Prompt{
		BasicPrompt: BasicPrompt{
			Label:   label,
			Default: startString,
		},
	}
	return p.Run()
}

// AskMasked func is default masked prompt
func AskMasked(label, startString string) (string, error) {
	p := Prompt{
		BasicPrompt: BasicPrompt{
			Label:   label,
			Default: startString,
		},
		Mask: '*',
	}
	return p.Run()
}

// PromptAfterSelect func is predefined easy confirm promp
var PromptAfterSelect = func(label string, answers []string) (string, error) {
	s := Select{
		Label:   label,
		Items:   answers,
		Default: 0,
	}
	_, rs, err := s.Run()
	if err != nil {
		return rs, err
	}
	fmt.Print(upLine(1) + clearLine)
	p := Prompt{
		BasicPrompt: BasicPrompt{
			Label:   label,
			Default: rs,
		},
	}
	return p.Run()
}
