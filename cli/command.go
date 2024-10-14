package cli

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type HandlerFunc func(*Command) error
type ValidatorFunc func(*Command) error

type Helper interface {
	HelpString() string
}

type PositonalArg struct {
	Name       string
	Type       string
	Required   bool
	Help       string
	Validators []ValidatorFunc
}

func (a *PositonalArg) HelpString() string {
	required := ""
	if a.Required {
		required = ", required"
	}
	return fmt.Sprintf("Argument: %s (%s%s) %s", a.Name, a.Type, required, a.Help)
}

type FlagArg struct {
	Name       string
	ShortName  string
	Type       string
	Required   bool
	Help       string
	Validators []ValidatorFunc
}

func (f *FlagArg) HelpString() string {
	required := ""
	if f.Required {
		required = "required"
	}
	shortName := ""
	if f.ShortName != "" {
		shortName = ", -" + f.ShortName
	}

	return fmt.Sprintf("{--%s%s} (%s%s) %s", f.Name, shortName, f.Type, required, f.Help)
}

type Command struct {
	Name       string
	Args       map[string]interface{}
	Flags      map[string]interface{}
	ArgTypes   []PositonalArg
	FlagTypes  []FlagArg
	Handler    HandlerFunc
	Help       string
	flagSet    *flag.FlagSet
	Validators []ValidatorFunc
	Program    string
}

func (c *Command) Parse(args []string) error {
	c.flagSet = flag.NewFlagSet(c.Name, flag.ExitOnError)
	c.flagSet.Usage = func() {
		fmt.Println(c.HelpString())
	}
	flagsMap := c.createFlagsMap()
	err := c.flagSet.Parse(args)
	if err != nil {
		return err
	}

	err = c.validateNumRequiredArgs(c.flagSet.Args())
	if err != nil {
		return err
	}

	err = c.validateNumRequiredFlags()
	if err != nil {
		return err
	}

	argsMap, err := c.createArgsMap()
	if err != nil {
		return err
	}

	errs := c.Validate()
	if len(errs) > 0 {
		return fmt.Errorf("validation errors: %v\n\n%s", errs, c.HelpString())
	}

	c.Flags = flagsMap
	c.Args = argsMap
	return nil
}

func (c *Command) Run() error {
	return c.Handler(c)
}

func (c *Command) Validate() []error {
	var errs []error
	for _, validator := range c.Validators {
		err := validator(c)
		if err != nil {
			errs = append(errs, err)
		}
	}

	for _, arg := range c.ArgTypes {
		for _, validator := range arg.Validators {
			err := validator(c)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	for _, flag := range c.FlagTypes {
		for _, validator := range flag.Validators {
			err := validator(c)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func (c *Command) validateNumRequiredArgs(args []string) error {
	requiredArgs := 0
	for _, arg := range c.ArgTypes {
		if arg.Required {
			requiredArgs++
		}
	}

	if len(args) < requiredArgs  || len(args) > len(c.ArgTypes) {
		return fmt.Errorf(c.HelpString())
	}
	return nil
}

func (c *Command) validateNumRequiredFlags() error {
	requiredFlags := 0
	for _, flag := range c.FlagTypes {
		if flag.Required {
			requiredFlags++
		}
	}

	if c.flagSet.NFlag() < requiredFlags {
		return fmt.Errorf(c.HelpString())
	}

	return nil
}

func (c *Command) argHelpStrings() (string, string) {
	argStrings := []string{}
	for _, arg := range c.ArgTypes {
		argStrings = append(argStrings, arg.HelpString())
	}

	flagStrings := []string{}
	for _, flag := range c.FlagTypes {
		flagStrings = append(flagStrings, flag.HelpString())
	}
	return strings.Join(argStrings, "\n"), strings.Join(flagStrings, "\n")
}

func (c *Command) HelpString() string {
	argNames := []string{}
	for _, arg := range c.ArgTypes {
		argNames = append(argNames, arg.Name)
	}
	var argNamesString string
	if len(argNames) > 0 {
		argNamesString = fmt.Sprintf("[%s]", strings.Join(argNames, " "))
	}
	argsString, flagsString := c.argHelpStrings()
	return fmt.Sprintf("Usage: %s %s %s\n%s%s", c.Program, c.Name, argNamesString, argsString, flagsString)
}

func (c *Command) createArgsMap() (map[string]interface{}, error) {
	args := make(map[string]interface{})
	for i, aT := range c.ArgTypes {
		arg := c.flagSet.Args()[i]
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
			boolArg := c.flagSet.Bool(arg.Name, false, c.HelpString())
			if arg.ShortName != "" {
				c.flagSet.BoolVar(boolArg, arg.ShortName, false, c.HelpString())
			}
			flags[arg.Name] = boolArg
		} else if arg.Type == "int" {
			intArg := c.flagSet.Int(arg.Name, 0, c.HelpString())
			if arg.ShortName != "" {
				c.flagSet.IntVar(intArg, arg.ShortName, 0, c.HelpString())
			}
			flags[arg.Name] = intArg
		} else if arg.Type == "string" {
			stringArg := c.flagSet.String(arg.Name, "", c.HelpString())
			if arg.ShortName != "" {
				c.flagSet.StringVar(stringArg, arg.ShortName, "", c.HelpString())
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
