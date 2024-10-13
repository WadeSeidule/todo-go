package cli

import (
	"flag"
	"fmt"
	"strings"
)

type Commands struct {
	Commands map[string]*Command
}

func (c *Commands) Contains(cmd string) (bool, *Command) {
	cmdStruct, exists := c.Commands[cmd]
	return exists, cmdStruct
}

func (c *Commands) ListCommands() string {
	var cmdNames []string
	for name := range c.Commands {
		cmdNames = append(cmdNames, name)
	}
	return strings.Join(cmdNames, ", ")
}

func (c *Commands) Run(cmd string, args []string) error {
	if cmd == "help" {
		c.Help(cmd)
		return nil
	}
	exists, cmdStruct := c.Contains(cmd)
	if !exists {
		return fmt.Errorf("command '%s' not found. Expected one of '%s'", cmd, c.ListCommands())
	}
	return cmdStruct.Run(args)
}

func (c *Commands) AddCommand(cmd *Command) {
	if c.Commands == nil {
		c.Commands = make(map[string]*Command)
	}
	c.Commands[cmd.Name] = cmd
}

func (c *Commands) ParseArgs() (string, []string, error) {
	flag.Parse()
	if len(flag.Args()) == 0 {
		return "", nil, fmt.Errorf("expected a command. Expected one of '%s'", c.ListCommands())
	}
	return flag.Args()[0], flag.Args()[1:], nil
}

func (c *Commands) Help(cmd string) {
	for name, cmdStruct := range c.Commands {
		fmt.Printf("%s: %s\n", name, cmdStruct.Help)
	}
}
