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

// Package clapper processes the command-line arguments of getopt(3) syntax.
// This package provides the ability to process the root command, sub commands,
// command-line arguments and command-line flags.
package clapper

import (
	"fmt"
	"strings"
)

/***********************************************
        PRIVATE FUNCTIONS AND VARIABLES
***********************************************/

// format command-line argument values
func formatCommandValues(values []string) (formatted []string) {

	formatted = make([]string, 0)

	// split a value by `=`
	for _, value := range values {
		if isFlag(value) {
			parts := strings.Split(value, "=")

			for _, part := range parts {
				if strings.Trim(part, " ") != "" {
					formatted = append(formatted, part)
				}
			}
		} else {
			formatted = append(formatted, value)
		}
	}

	return
}

// check if value is a flag
func isFlag(value string) bool {
	return len(value) >= 2 && strings.HasPrefix(value, "-")
}

// check if value is a short flag
func isShortFlag(value string) bool {
	return isFlag(value) && len(value) == 2 && !strings.HasPrefix(value, "--")
}

// check if value starts with `--no-` prefix
func isInvertFlag(value string) (bool, string) {
	if isFlag(value) && strings.HasPrefix(value, "--no-") {
		return true, strings.TrimLeft(value, "--no-") // trim `--no-` prefix
	}

	return false, ""
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

// ErrorUnknownCommand represents an error when command-line arguments contain an unregistered command.
type ErrorUnknownCommand struct {
	Name string
}

func (e ErrorUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command %s found in the arguments", e.Name)
}

// ErrorUnknownFlag represents an error when command-line arguments contain an unregistered flag.
type ErrorUnknownFlag struct {
	Name string
}

func (e ErrorUnknownFlag) Error() string {
	return fmt.Sprintf("unknown flag %s found in the arguments", e.Name)
}

// ErrorUnsupportedFlag represents an error when command-line arguments contain an unsupported flag.
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
// If "name" is an empty string, it is considered as a root command.
// If a command is already registered, the registered `*Carg` object is returned.
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

// Parse method parses command-line arguments and returns an appropriate "*Carg" object registered in the registry.
// If command is not registered, it return `ErrorUnknownCommand` error.
// If there is an error parsing a flag, it can return an `ErrorUnknownFlag` or `ErrorUnsupportedFlag` error.
func (registry Registry) Parse(values []string) (*Carg, error) {

	// command name
	var commandName string

	// command-line argument values to process
	valuesToProcess := values

	// check if command is a root command
	if isRootCommand(values, registry) {
		commandName = "" // root command name
	} else {
		commandName, valuesToProcess = nextValue(values)
	}

	// format command-line argument values
	valuesToProcess = formatCommandValues(valuesToProcess)

	// check for invalid flag structure
	for _, val := range valuesToProcess {
		if isFlag(val) && isUnsupportedFlag(val) {
			return nil, ErrorUnsupportedFlag{val}
		}
	}

	// if command is not registered, return `ErrorUnknownCommand` error
	if _, ok := registry[commandName]; !ok {
		return nil, ErrorUnknownCommand{commandName}
	}

	// get `Carg` object from the registry
	carg := registry[commandName]

	// process all command-line arguments (except command name)
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

			// get flag object stored in the `carg`
			var flag *Flag

			// check if flag is short or long
			if isShortFlag(value) {
				if _, ok := carg.flagsShort[name]; !ok {
					return nil, ErrorUnknownFlag{value}
				}

				// get long flag name
				flagName := carg.flagsShort[name]

				flag = carg.Flags[flagName]
			} else {

				// check if a flag is an invert flag
				if ok, flagName := isInvertFlag(value); ok {
					if _, ok := carg.Flags[flagName]; !ok {
						return nil, ErrorUnknownFlag{value}
					}

					flag = carg.Flags[flagName]
				} else {

					// flag should not registered as an invert flag
					if _flag, ok := carg.Flags[name]; !ok || _flag.IsInvert {
						return nil, ErrorUnknownFlag{value}
					}

					flag = carg.Flags[name]
				}
			}

			// set flag value
			if flag.IsBoolean {
				if flag.IsInvert {
					flag.Value = "false" // if flag is an invert flag, its value will be `false`
				} else {
					flag.Value = "true"
				}
			} else {
				if nextValue, nextValuesToProcess := nextValue(valuesToProcess); len(nextValue) != 0 && !isFlag(nextValue) {
					flag.Value = nextValue
					valuesToProcess = nextValuesToProcess
				}
			}
		} else {

			// process as argument
			for _, argName := range carg.argNames {

				// get argument object stored in the `carg`
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

// Carg type holds the structured information about the command-line arguments.
type Carg struct {

	// name of the sub-command ("" for the root command)
	Cmd string

	// command-line flags
	Flags map[string]*Flag

	// mapping of the short flag names with long flag names
	flagsShort map[string]string

	// registered command argument values
	Args map[string]*Arg

	// list of the argument names (for ordered iteration)
	argNames []string
}

// AddArg registers an argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `defaultValue` argument represents the default value of the argument.
// All arguments with a default value must be registered first.
// If an argument with given `name` is already registered, then argument registration is skipped and registered `*Carg` object returned.
func (carg *Carg) AddArg(name string, defaultValue string) *Carg {

	// return if argument is already registered
	if _, ok := carg.Args[name]; ok {
		return carg
	}

	// create Arg object
	arg := &Arg{
		Name:         name,
		DefaultValue: defaultValue,
	}

	// register argument
	carg.Args[name] = arg

	// store argument name
	carg.argNames = append(carg.argNames, name)

	return carg
}

// AddFlag method registers a command-line flag with the command.
// The `name` argument is the long-name of the flag and it should not start with `--` prefix.
// The `shortName` argument is the short-name of the flag and it should not start with `-` prefix.
// The `isBool` argument indicates whether the flag holds a boolean value.
// A boolean flag doesn't accept an input value such as `--flag=<value>` and its default value is "true".
// The `defaultValue` argument represents the default value of the flag.
// In case of a boolean flag, the `defaultValue` is redundant.
// If the `name` value starts with `no-` prefix, then it is considered as an invert flag.
// An invert flag is registered with the name `<flag>` produced by removing `no-` prefix from `no-<flag>` and its defaut value is "true".
// When command-line arguments contain `--no-<flag>`, the value of the `<flag>` becomes "false".
// If a flag with given `name` is already registered, then flag registration is skipped and registered `*Carg` object returned.
func (carg *Carg) AddFlag(name string, shortName string, isBool bool, defaultValue string) *Carg {

	// actual flag names and default value
	_name := name
	_shortName := shortName
	_defaultValue := defaultValue

	// invert flag is a boolean flag with `no-` prefix
	_isInvert := false

	// short flag name should be only one character long
	if _shortName != "" {
		_shortName = _shortName[:1]
	}

	// set up `Flag` field values
	if isBool {

		// check if flag name has `no-` prefix
		if strings.HasPrefix(name, "no-") {
			_isInvert = true                      // is an invert flag
			_name = strings.TrimLeft(name, "no-") // trim `no-` prefix
			_defaultValue = "true"                // default value of an invert flag is `true`
			_shortName = ""                       // no short flag name for invert flag
		} else {
			_defaultValue = "false" // default value of a boolean flag is `true`
		}
	} else {
		_defaultValue = defaultValue
	}

	// return `carg` if flag is already registered
	if _, ok := carg.Flags[_name]; ok {
		return carg
	}

	// create a `Flag` object
	flag := &Flag{
		Name:         _name,
		ShortName:    _shortName,
		IsBoolean:    isBool,
		IsInvert:     _isInvert,
		DefaultValue: _defaultValue,
	}

	// register a flag
	carg.Flags[_name] = flag

	// register short flag name
	if len(_shortName) > 0 {
		carg.flagsShort[_shortName] = _name
	}

	return carg
}

/*---------------------*/

// Flag type holds the structured information about a flag.
type Flag struct {

	// long name of the flag
	Name string

	// short name of the flag
	ShortName string

	// if the flag holds boolean value
	IsBoolean bool

	// if the flag is an invert flag (with `--no-` prefix)
	IsInvert bool

	// default value of the flag
	DefaultValue string

	// value of the flag (provided by the user)
	Value string
}

/*---------------------*/

// Arg type holds the structured information about an argument.
type Arg struct {
	// name of the argument
	Name string

	// default value of the argument
	DefaultValue string

	// value of the argument (provided by the user)
	Value string
}
