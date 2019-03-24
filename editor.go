package promptui

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

var (
	// ErrorOpen describes open editor error
	ErrorOpen = errors.New("in promptui:Editor: Error open temp file")
	// ErrorWrite describes write editor error
	ErrorWrite = errors.New("in promptui:Editor: Error write temp file")
	// ErrorRead describes read editor error
	ErrorRead = errors.New("in promptui:Editor: Error read temp file")
	// ErrorClose describes close editor error
	ErrorClose = errors.New("in promptui:Editor: Error close temp file")
)

var (
	bom = []byte{0xef, 0xbb, 0xbf}
)

// Editor func opens default editor to edit `s string`.
// Returns edit result. Creates temp file to edit string
func Editor(editor, s string) (string, error) {
	if editor == "" {
		editor = getEditor()
	}
	tf, err := ioutil.TempFile("", "promptui.*.txt")
	if err != nil {
		return "", ErrorOpen
	}
	defer os.Remove(tf.Name()) // clean up

	if runtime.GOOS == "windows" {
		// write utf8 BOM header
		// The reason why we do this is because notepad.exe on Windows determines the
		// encoding of an "empty" text file by the locale, for example, GBK in China,
		// while golang string only handles utf8 well. However, a text file with utf8
		// BOM header is not considered "empty" on Windows, and the encoding will then
		// be determined utf8 by notepad.exe, instead of GBK or other encodings.
		if _, err := tf.Write(bom); err != nil {
			return "", err
		}
	}

	// write initial value
	if _, err := tf.WriteString(s); err != nil {
		return "", err
	}

	// close the fd to prevent the editor unable to save file
	if err := tf.Close(); err != nil {
		return "", err
	}

	args := []string{editor, tf.Name()}

	// open the editor
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// raw is a BOM-unstripped UTF8 byte slice
	raw, err := ioutil.ReadFile(tf.Name())
	if err != nil {
		return "", err
	}

	text := ""
	if runtime.GOOS == "windows" {
		// strip BOM header
		text = string(bytes.TrimPrefix(raw, bom))
	} else {
		text = string(raw)
	}

	return text, nil
}

func getEditor() string {
	if runtime.GOOS == "windows" {
		return "notepad"
	}
	if v := os.Getenv("VISUAL"); v != "" {
		return v
	} else if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	if runtime.GOOS == "linux" {
		return "editor"
	}
	return "vim"
}
