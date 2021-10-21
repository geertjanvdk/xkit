// Copyright (c) 2021, Geert JM Vanderkelen

package xprompt

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/geertjanvdk/xkit/xansi"
	"github.com/geertjanvdk/xkit/xutil"
	"golang.org/x/crypto/ssh/terminal"
)

type history struct {
	history []string
	cursor  int
}

func newHistory() *history {
	return &history{}
}

func (h *history) previous() string {
	h.cursor++
	if h.cursor <= 0 {
		return h.history[h.cursor]
	}
	if h.cursor >= len(h.history) {
		return h.history[len(h.history)-1]
	}
	return h.history[h.cursor]
}

func (h *history) resetCursor() {
	h.cursor = len(h.history) - 1
}

type Prompt struct {
	prompt         string
	exitCommands   []string
	commands       []Commander
	defaultCommand Commander

	history *history

	in  os.File
	out os.File

	origIn  *terminal.State
	origOut *terminal.State
}

// NewPrompt instantiates a new Prompt object which will show prompt when waiting
// for input from in. Feedback is written to out.
func NewPrompt(prompt string, in os.File, out os.File) (*Prompt, error) {
	p := &Prompt{
		prompt:       prompt,
		exitCommands: []string{"exit", "quit"},
		history:      newHistory(),
	}

	p.in = int(in.Fd())
	p.out = int(out.Fd())

	var err error
	p.origIn, err = terminal.MakeRaw(p.in)
	if err != nil {
		return nil, err
	}

	p.origOut, err = terminal.MakeRaw(p.out)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Prompt) Stop(a ...interface{}) {
	_, _ = fmt.Fprint(p.out, a...)
	p.Printf("Bye!\n")
}

func (p *Prompt) readInput() (string, error) {

	defer func() {
		_ = terminal.Restore(inFd, origOutState)
	}()

	defer func() {
		_ = terminal.Restore(inFd, origInState)
		_ = terminal.Restore(outFd, origOutState)
	}()

	var input []byte

	for {
		buf := make([]byte, 4)
		read, err := in.Read(buf)
		if err != nil {
			return "", err
		}

		if read == 1 && buf[0] == 13 {
			break
		} else if read == 3 && buf[0] == 0x1b && buf[1] == 0x4f {
			switch buf[2] {
			case 0x41:
				// up
				continue
			case 0x42:
				// down
				continue
			case 0x43:
				// right
				_, _ = fmt.Fprintf(out, "\x1b[1C")
				continue
			case 0x44:
				// left
				_, _ = fmt.Fprintf(out, "\x1b[1D")
				continue
			}
		} else {
			_, _ = fmt.Fprintf(out, "%s", string(buf))
			input = append(input, buf...)
		}
	}

	return string(input), nil
}

func (p *Prompt) Start() {
	for {
		_, _ = fmt.Fprint(p.out, p.prompt)
		input, err := p.readInput()
		if err != nil {
			if err == io.EOF {
				p.Stop("\n")
				return
			}
			p.Errorf(err.Error())
			continue
		}

		switch {
		case xutil.HasString(p.exitCommands, input):
			p.Stop()
			return
		default:
			name := input
			i := strings.Index(input, " ")
			if i > 0 {
				name = input[0:i]
			}
			cmd, err := p.FindCommand(name)
			if err != nil && p.defaultCommand == nil {
				p.Errorf("command not available")
				continue
			} else if cmd != nil {
				if err := cmd.CommandInput(input); err != nil {
					p.Errorf(err)
				}
			} else {
				cmd = p.defaultCommand
			}

		}
	}
}

func (p Prompt) Errorf(format interface{}, a ...interface{}) {
	switch v := format.(type) {
	case string:
		_, _ = fmt.Fprintln(p.out, xansi.Render{xansi.Red, xansi.Bold}.Sprintf("Error: "+v, a...))
	case error:
		_, _ = fmt.Fprintln(p.out, xansi.Render{xansi.Red, xansi.Bold}.Sprint(
			append([]interface{}{"Error: " + v.Error()}, a...)))
	}
}

// SetDefaultCommand sets command as default command, which will be executed
// when no other command matches.
func (p *Prompt) SetDefaultCommand(command Commander) error {
	if p.defaultCommand != nil {
		return fmt.Errorf("default command already set")
	}
	command.PromptReference(p)
	p.defaultCommand = command
	return nil
}

func (p Prompt) Printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(p.out, strings.TrimSpace(format)+"\n", a...)
}

func (p *Prompt) FindCommand(s string) (Commander, error) {
	for _, cmd := range p.commands {
		if cmd.CommandName() == s {
			return cmd, nil
		}
	}
	return nil, fmt.Errorf("command '%s' not available", s)
}

func (p *Prompt) RegisterCommand(command Commander) error {
	if cmd, _ := p.FindCommand(command.CommandName()); cmd != nil {
		return fmt.Errorf("command %s already registered", command.CommandName())
	}

	command.PromptReference(p)
	p.commands = append(p.commands, command)
	return nil
}
