// MIT License

// Copyright (c) 2020 Uday Hiwarale

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package clapper processes the command line arguments of getopt(3) syntax.
// This package provides the ability to process the root command, sub commands,
// command line arguments and command line flags.
package clapper

import (
	"fmt"
	"strings"
)

/***********************************************
        PRIVATE FUNCTIONS AND VARIABLES
***********************************************/

// check if value is a flag
func isFlag(value string) bool {
	return len(value) >= 2 && strings.HasPrefix(value, "-")
}

// check if value is a short flag
func isShortFlag(value string) bool {
	return isFlag(value) && len(value) == 2 && !strings.HasPrefix(value, "--")
}

// check if flag is unsupported
func isUnsupportedFlag(value string) bool {

	// a flag should be at least two characters log
	if len(value) >= 2 {

		// if short flag, it should start with `-` but not with `--`
		if len(value) == 2 {
			return !strings.HasPrefix(value, "-") || strings.HasPrefix(value, "--")
		}

		// if long flag, it should start with `--` and not with `---`
		return !strings.HasPrefix(value, "--") || strings.HasPrefix(value, "---")
	}

	return false
}

// check if values corresponds to the root command
func isRootCommand(values []string, registry Registry) bool {

	// FALSE: if the root command is not registered
	if _, ok := registry[""]; !ok {
		return false
	}

	// TRUE: if all `values` are empty or the first `value` is a flag
	if len(values) == 0 || isFlag(values[0]) {
		return true
	}

	// get root `Carg` value from the registry
	rootCarg := registry[""]

	// TRUE: if the first value is not a registered command
	// and some arguments are registered for the root command
	if _, ok := registry[values[0]]; len(rootCarg.Args) > 0 && !ok {
		return true
	}

	return false
}

// return next value and remaining values of a slice of strings
func nextValue(slice []string) (v string, newSlice []string) {

	if len(slice) == 0 {
		v, newSlice = "", make([]string, 0)
		return
	}

	v = slice[0]

	if len(slice) > 1 {
		newSlice = slice[1:]
	} else {
		newSlice = make([]string, 0)
	}

	return
}

/***********************************************/

// ErrorUnknownCommand represents an error when command line arguments contain an unregistered command.
type ErrorUnknownCommand struct {
	Name string
}

func (e ErrorUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command %s found in the arguments", e.Name)
}

// ErrorUnknownFlag represents an error when command line arguments contain an unregistered flag.
type ErrorUnknownFlag struct {
	Name    string
	IsShort bool
}

func (e ErrorUnknownFlag) Error() string {
	return fmt.Sprintf("unknown flag %s found in the arguments", e.Name)
}

// ErrorUnsupportedFlag represents an error when command line arguments contain an unsupported flag.
type ErrorUnsupportedFlag struct {
	Name string
}

func (e ErrorUnsupportedFlag) Error() string {
	return fmt.Sprintf("unsupported flag %s found in the arguments", e.Name)
}

/*---------------------*/

// Registry holds the configuration of the registered commands.
type Registry map[string]*Carg

// Register method registers a command.
// The "name" argument should be a simple string.
// If "name" is empty, it is considered as a root command.
// If a command is already registered, the registered command is returned.
func (registry Registry) Register(name string) *Carg {

	// check if command is already registered, if found, return existing entry
	if _carg, ok := registry[name]; ok {
		return _carg
	}

	// construct new `Carg` object
	carg := &Carg{
		Cmd:        name,
		Flags:      make(map[string]*Flag),
		flagsShort: make(map[string]string),
		Args:       make(map[string]*Arg),
		argNames:   make([]string, 0),
	}

	// add entry to the registry
	registry[name] = carg

	return carg
}

// Parse method parses command line arguments and returns an appropriate "Carg" object registered in the registry.
// If command is not registered, return `ErrorUnknownCommand` error
// If flag is not registered, return `ErrorUnknownFlag` error
func (registry Registry) Parse(values []string) (*Carg, error) {

	// command name
	var commandName string

	// command line argument values to process
	valuesToProcess := values

	// check for invalid flag structure
	for _, val := range values {
		if isFlag(val) && isUnsupportedFlag(val) {
			return nil, ErrorUnsupportedFlag{val}
		}
	}

	// check if command is a root command
	if isRootCommand(values, registry) {
		commandName = ""
	} else {
		commandName, valuesToProcess = nextValue(values)
	}

	// if command is not registered, return `ErrorUnknownCommand` error
	if _, ok := registry[commandName]; !ok {
		return nil, ErrorUnknownCommand{commandName}
	}

	// get `Carg` object from the registry
	carg := registry[commandName]

	// process all command line arguments (except command name)
	for {

		// get current command-line argument value
		var value string
		value, valuesToProcess = nextValue(valuesToProcess)

		// if `value` is empty, break the loop
		if len(value) == 0 {
			break
		}
		// check if `value` is a `flag` or an `argument`
		if isFlag(value) {

			// trim `-` characters from the `value`
			name := strings.TrimLeft(value, "-")

			// get flag object stored with the `carg`
			var flag *Flag

			if isShortFlag(value) {
				if _, ok := carg.flagsShort[name]; !ok {
					return nil, ErrorUnknownFlag{name, true}
				}

				flag = carg.Flags[carg.flagsShort[name]]
			} else {
				if _, ok := carg.Flags[name]; !ok {
					return nil, ErrorUnknownFlag{name, false}
				}

				flag = carg.Flags[name]
			}

			// set flag value
			if flag.IsBoolean {
				flag.Value = "true"
			} else {
				if nextValue, nextValuesToProcess := nextValue(valuesToProcess); len(nextValue) != 0 && !isFlag(nextValue) {
					flag.Value = nextValue
					valuesToProcess = nextValuesToProcess
				}
			}
		} else {

			// process as argument
			for _, argName := range carg.argNames {

				// get argument from the name of the argument
				arg := carg.Args[argName]

				// assign value if value of the argument is empty
				if len(arg.Value) == 0 {
					arg.Value = value
					break
				}
			}
		}

	}

	return carg, nil
}

// NewRegistry returns new instance of the "Registry"
func NewRegistry() Registry {
	return make(Registry)
}

/*---------------------*/

// Carg type holds the structured information about the command line arguments
type Carg struct {

	// name of the command executed
	Cmd string

	// command line flags
	Flags map[string]*Flag

	// mapping of the short flag names with long flag names
	flagsShort map[string]string

	// registered command argument values
	Args map[string]*Arg

	// list of the argument names (for ordered iteration)
	argNames []string
}

// AddFlag method registeres a "Flag" value
func (carg *Carg) AddFlag(name string, shortName string, isBool bool, defaultValue string) *Carg {

	// return if flag is already registered
	if _, ok := carg.Flags[name]; ok {
		return carg
	}

	// create a Flag object
	flag := &Flag{
		Name:      name,
		ShortName: shortName,
		IsBoolean: isBool,
	}

	// register flag
	carg.Flags[name] = flag

	// set default value
	if isBool {
		carg.Flags[name].DefaultValue = "false"
	} else {
		carg.Flags[name].DefaultValue = defaultValue
	}

	// store short flag name
	if len(shortName) > 0 {
		carg.flagsShort[shortName] = name
	}

	return carg
}

// AddArg registers a "Arg" value
func (carg *Carg) AddArg(name string) *Carg {

	// return if argument is already registered
	if _, ok := carg.Args[name]; ok {
		return carg
	}

	// create Arg object
	arg := &Arg{
		Name: name,
	}

	// register argument
	carg.Args[name] = arg

	// store argument name
	carg.argNames = append(carg.argNames, name)

	return carg
}

/*---------------------*/

// Flag type holds the structured information about the command line flag
type Flag struct {

	// long name of the flag
	Name string

	// short name of the flag
	ShortName string

	// if flag holds boolean value
	IsBoolean bool

	// default value of the flag
	DefaultValue string

	// value of the flag
	Value string
}

/*---------------------*/

// Arg type holds the structured information about the command line argument
type Arg struct {
	Name  string
	Value string
}
