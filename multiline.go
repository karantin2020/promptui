package promptui

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/karantin2020/readline"
)

// MultilinePrompt represents a single line text field input.
type MultilinePrompt struct {
	BasicPrompt

	// OnError func is called when prompt.Run() return error
	// alternative way to fix error
	OnError func(string) (string, error)

	// Editor is default editor to edit multiline
	Editor string
}

// Run func implements multiline prompt logic
func (mp *MultilinePrompt) Run() (string, error) {
	err := mp.Init()
	if err != nil {
		return "", err
	}

	mp.c.Stdin = ioutil.NopCloser(io.MultiReader(bytes.NewBuffer([]byte(mp.Default)), os.Stdin))

	mp.rl, err = readline.NewEx(mp.c)
	if err != nil {
		return "", err
	}
	mp.suggestedAnswer = " " + faint("Two empty lines to finish")
	mp.prompt = mp.LabelInitial(mp.Label) + mp.punctuation + mp.suggestedAnswer + " "
	mp.c.UniqueEditLine = false
	var (
		firstListen      = true
		numlines    uint = 1
		breaklines       = 0
		out         string
	)

	mp.rl.Write([]byte(mp.Indent + mp.state + " " + mp.PromptInitial(mp.prompt) + "\n"))
	mp.rl.SetPrompt("... ")
	mp.rl.Refresh()
	var multilineReader = func(line []rune, pos int, key rune) ([]rune, int, bool) {
		if firstListen {
			firstListen = false
		}
		return nil, 0, false
	}
	mp.c.SetListener(multilineReader)

	for {
		out, err = mp.rl.Readline()
		numlines++
		if out == "" {
			breaklines++
		} else {
			breaklines = 0
		}

		mp.out += out

		if breaklines > 1 {
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

		mp.rl.SetPrompt("... ")
		mp.rl.Refresh()
		mp.out += "\n"
	}

	if err != nil {
		if err.Error() == "Interrupt" {
			err = ErrInterrupt
		}
		return "", err
	}

	defer mp.rl.Close()

	var clearLines = func(numlines uint) {
		for ; numlines > 0; numlines-- {
			mp.rl.Write([]byte(upLine(1) + clearLine))
		}
	}

	for {
		msg, oerr := mp.formatAndValidate()

		clearLines(numlines)
		numlines = uint(len(strings.Split(mp.out, "\n")))
		mp.rl.Write([]byte(mp.Indent + mp.state + " " + mp.prompt + "\n" + mp.InputResult(mp.out) + "\n"))
		numlines++
		if oerr != nil {
			if mp.OnError != nil {
				return mp.OnError(mp.out)
			}
			mp.rl.Write([]byte(red("Error: ") + msg + "\n"))
			numlines++
			var yn string
			yn, oerr = Confirm("Open editor to edit input", "", true)
			numlines++
			if oerr != nil {
				return mp.out, oerr
			}
			if yn == "Y" {
				mp.out, oerr = Editor(mp.Editor, mp.out)
				if oerr != nil {
					return mp.out, oerr
				}
				continue
			} else {
				clearLines(2)
				break
			}
		}
		break
	}
	return mp.out, err
}

func (mp *MultilinePrompt) formatAndValidate() (msg string, oerr error) {
	mp.out = strings.Trim(mp.out, "\n\r")
	mp.out = mp.Formatter(mp.out)

	oerr = mp.validFn(mp.out)
	if oerr != nil {
		if verr, ok := oerr.(*ValidationError); ok {
			msg = verr.msg
		}
		mp.state = mp.IconBad
	} else {
		mp.state = mp.IconGood
	}
	return
}

// MultiLine func is default multiline prompt
func MultiLine(label, answer string) (string, error) {
	mp := MultilinePrompt{
		BasicPrompt: BasicPrompt{
			Label:   label,
			Default: answer,
		},
	}
	return mp.Run()
}
