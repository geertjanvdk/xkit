// Copyright (c) 2021, Geert JM Vanderkelen

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/geertjanvdk/xkit/xprompt"
)

type EchoCommand struct {
	prompt *xprompt.Prompt
}

var _ xprompt.Commander = &EchoCommand{}

func (c EchoCommand) CommandName() string {
	return "echo"
}

func (c EchoCommand) CommandInput(s string) error {
	s = strings.Replace(c.CommandName(), s, "", 1)
	c.prompt.Printf("Echo: %s", strings.TrimSpace(s))
	return nil
}

func (c *EchoCommand) PromptReference(p *xprompt.Prompt) {
	c.prompt = p
}

func main() {
	prompt, _ := xprompt.NewPrompt(">> ", os.Stdin, os.Stdout)
	if err := prompt.RegisterCommand(&EchoCommand{}); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	prompt.Start()
}
