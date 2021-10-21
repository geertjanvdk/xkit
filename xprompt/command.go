// Copyright (c) 2021, Geert JM Vanderkelen

package xprompt

type Commander interface {
	CommandName() string
	CommandInput(string) error
	PromptReference(prompt *Prompt)
}
