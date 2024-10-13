package cli

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

type PositonalArg struct {
	Name     string
	Type     string
	Required bool
	Help     string
}

type FlagArg struct {
	Name      string
	ShortName string
	Type      string
	Required  bool
	Help      string
}

type Command struct {
	Name      string
	Args      map[string]interface{}
	Flags     map[string]interface{}
	ArgTypes  []PositonalArg
	FlagTypes []FlagArg
	Handler   func(*Command) error
	Help      string
	FlagSet   *flag.FlagSet
}

func (c *Command) Run(args []string) error {
	flagsMap := c.createFlagsMap()
	err := c.FlagSet.Parse(args)
	if err != nil {
		return err
	}
	err = c.validateArgs(c.FlagSet.Args())
	if err != nil {
		return err
	}
	argsMap, err := c.createArgsMap()
	if err != nil {
		return err
	}
	c.Flags = flagsMap
	c.Args = argsMap
	return c.Handler(c)
}

func (c *Command) validateArgs(args []string) error {
	requiredArgs := 0
	optionalArgs := 0
	for _, arg := range c.ArgTypes {
		if arg.Required {
			requiredArgs++
		} else {
			optionalArgs++
		}
	}
	// Validate number of arguments
	if len(args) < requiredArgs || len(args) > requiredArgs+optionalArgs {
		return fmt.Errorf(c.HelpString())
	}
	return nil
}

func (c *Command) HelpString() string {
	helpString := fmt.Sprintf("%s\n", c.Help)
	for _, arg := range c.ArgTypes {
		helpString += fmt.Sprintf("Argument: %s\nType: %s\nRequired: %t\nHelp: %s", arg.Name, arg.Type, arg.Required, arg.Help)
	}
	return helpString
}

func (c *Command) createArgsMap() (map[string]interface{}, error) {
	args := make(map[string]interface{})
	for i, aT := range c.ArgTypes {
		arg := c.FlagSet.Args()[i]
		if aT.Type == "int" {
			int_arg, err := strconv.Atoi(arg)
			if err != nil {
				return nil, errors.New("Invalid type for positional argument: " + aT.Name + ". Expected int.")
			}
			args[aT.Name] = int_arg
		} else if aT.Type == "float" {
			float_arg, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				return nil, errors.New("Invalid type for positional argument: " + aT.Name + ". Expected float.")
			}
			args[aT.Name] = float_arg
		} else {
			args[aT.Name] = arg
		}
	}
	return args, nil
}

func (c *Command) createFlagsMap() map[string]interface{} {
	flags := make(map[string]interface{})
	for _, arg := range c.FlagTypes {
		if arg.Type == "bool" {
			boolArg := c.FlagSet.Bool(arg.Name, false, c.HelpString())
			if arg.ShortName != "" {
				c.FlagSet.BoolVar(boolArg, arg.ShortName, false, c.HelpString())
			}
			flags[arg.Name] = boolArg
		} else if arg.Type == "int" {
			intArg := c.FlagSet.Int(arg.Name, 0, c.HelpString())
			if arg.ShortName != "" {
				c.FlagSet.IntVar(intArg, arg.ShortName, 0, c.HelpString())
			}
			flags[arg.Name] = intArg
		} else if arg.Type == "string" {
			stringArg := c.FlagSet.String(arg.Name, "", c.HelpString())
			if arg.ShortName != "" {
				c.FlagSet.StringVar(stringArg, arg.ShortName, "", c.HelpString())
			}
			flags[arg.Name] = stringArg
		}
	}
	return flags
}

func (c *Command) GetStringArg(name string) string {
	return c.Args[name].(string)
}

func (c *Command) GetIntArg(name string) int {
	return c.Args[name].(int)
}

func (c *Command) GetFloatArg(name string) float64 {
	return c.Args[name].(float64)
}

func (c *Command) GetBoolFlag(name string) bool {
	return *(c.Flags[name].(*bool))
}

func (c *Command) GetIntFlag(name string) int {
	return *(c.Flags[name].(*int))
}

func (c *Command) GetStringFlag(name string) string {
	return *(c.Flags[name].(*string))
}

func (c *Command) GetFloatFlag(name string) float64 {
	return *(c.Flags[name].(*float64))
}
