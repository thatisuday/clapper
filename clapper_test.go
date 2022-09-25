package clapper

import (
	"fmt"
	"strings"
	"testing"
)

// This entire suite would be far easier with testify

var tests []struct {
	subCommand string
	arg        string
	longName   string
	shortName  string
	isBool     bool
	defaultVal string
}

func setup(t *testing.T, args []string) (*CommandConfig, error) {
	tests = []struct {
		subCommand string
		arg        string
		longName   string
		shortName  string
		isBool     bool
		defaultVal string
	}{
		{"", "output", "", "", false, ""},
		{"", "", "force", "f", true, ""},
		{"", "", "verbose", "v", true, ""},
		{"", "", "version", "V", false, ""},
		{"", "", "dir", "", false, "/var/users"},
		{"info", "category", "", "", false, "manager"},
		{"info", "username", "", "", false, ""},
		{"info", "subjects...", "", "", false, ""},
		{"info", "", "verbose", "v", true, ""},
		{"info", "", "version", "V", false, "1.0.1"},
		{"info", "", "output", "o", false, "./"},
		{"info", "", "no-clean", "", true, ""},
		{"ghost", "", "", "", false, ""},
	}

	reg := NewRegistry()
	subs := make(map[string]*CommandConfig)
	for _, test := range tests {
		sub := subs[test.subCommand]
		if sub == nil {
			sub, _ = reg.Register(test.subCommand)
			subs[test.subCommand] = sub
		}
		if test.longName != "" {
			sub.AddFlag(test.longName, test.shortName, test.isBool, test.defaultVal)
		} else if test.arg != "" {
			sub.AddArg(test.arg, test.defaultVal)
		}
	}

	return reg.Parse(args)
}

/*----------------*/

// test unsupported flag
func TestUnsupportedAssignment(t *testing.T) {

	tests := []struct {
		flag     string
		expected string
		options  string
		err      error
	}{
		{"---version", "---version", "---version", ErrorUnsupportedFlag{}},
		{"---v", "---v", "---v=1.0.0", ErrorUnsupportedFlag{}},
		// Single-dash for long names was never supported, but now this is interpreted as:
		// -v -e -r -s -i -o -n; ergo, the error will be "unknown option '-e'"
		{"-e", "-version", "-version", ErrorUnsupportedFlag{}},
	}

	for _, test := range tests {
		reg := NewRegistry()
		root, _ := reg.Register("")
		root.AddFlag("version", "v", true, "")

		_, err := reg.Parse([]string{test.options})
		assertError(t, err, "(%s) %T", test.options, test.err)
		if err != nil {
			assertEqual(t, fmt.Sprintf("%T", err), fmt.Sprintf("%T", test.err))
			var name string
			switch e := err.(type) {
			case ErrorUnknownFlag:
				name = e.Name
			case ErrorUnsupportedFlag:
				name = e.Name
			default:
				assertError(t, err)
			}
			assertEqual(t, test.expected, name)
		}
	}
}

// test empty root command
func TestEmptyRootCommand(t *testing.T) {
	cmd, err := setup(t, []string{})

	assertNoError(t, err)
	assertEqual(t, "", cmd.Name)
	assertEqual(t, 1, len(cmd.Args))

	if cmd.Args["output"] == nil {
		for k := range cmd.Args {
			t.Errorf("expected one \"output\" argument; got %s", k)
		}
	}

}

func TestRootDefaults(t *testing.T) {
	cmd, err := setup(t, []string{})

	assertNoError(t, err)
	assertEqual(t, "", cmd.Args["output"].Value)

	for _, test := range tests {
		if test.longName != "" && test.subCommand == "" {
			f := cmd.Flags[test.longName]
			assertNotNil(t, f, "missing flag %q", test.longName)
			assertEqual(t, test.shortName, f.ShortName, "(%s)", test.longName)
			assertEqual(t, test.isBool, f.IsBoolean, "(%s)", test.longName)
			dv := test.defaultVal
			if test.isBool && dv == "" {
				dv = "false"
			}
			assertEqual(t, dv, f.DefaultValue, "(%+v %+v)", test, f)
		}
	}
}

// test root command when not registered
// REMOVED This was testing code in the demo program, which had nothing to do with the library.

// test an unregistered flag
func TestUnregisteredFlag(t *testing.T) {
	// flags
	flags := map[string][]string{
		"-d":          {"-V", "1.0.1", "-v", "--force", "-d", "./sub/dir"},
		"--m":         {"-V", "1.0.1", "-v", "--force", "--m", "./sub/dir"},
		"--directory": {"-V", "1.0.1", "-v", "--force", "--directory", "./sub/dir"},
	}

	for flag, options := range flags {
		_, err := setup(t, options)
		assertError(t, err)
		if e, ok := err.(ErrorUnknownFlag); !ok {
			t.Errorf("expected an ErrUnknownFlag; got %t", err)
		} else {
			assertEqual(t, flag, e.Name)
		}
	}
}

// test for valid inverted flag values
func TestValidInvertFlagValues(t *testing.T) {
	// options list
	optionsList := [][]string{
		{"info", "student", "-V", "-v", "--output", "./opt/dir", "--no-clean"},
		{"info", "student", "--version", "--no-clean", "--output", "./opt/dir", "--verbose"},
	}
	expecteds := map[string]string{
		"version": "",
		"clean":   "false",
		"output":  "./opt/dir",
		"verbose": "true",
	}

	for _, options := range optionsList {
		cmd, err := setup(t, options)
		assertNoError(t, err)
		assertEqual(t, options[0], cmd.Name)
		assertEqual(t, options[1], cmd.Args["category"].Value)
		assertEqual(t, "", cmd.Args["username"].Value)
		assertEqual(t, "", cmd.Args["subjects"].Value)
		for _, opt := range options {
			if strings.HasPrefix("--no", opt) && cmd.Flags[opt].IsInverted != true {
				t.Errorf("expected inverted flag for %s", opt)
			}
		}
		for k, v := range expecteds {
			assertEqual(t, v, cmd.Flags[k].Value)
		}
	}
}

// test for invalid flag error when an inverted flag is used without `--no-` prefix
func TestErrorUnknownFlagForInvertFlags(t *testing.T) {

	// options list
	optionsList := map[string][]string{
		"--clean":   {"info", "student", "-V", "-v", "--output", "./opt/dir", "--clean"},
		"--no-dump": {"info", "student", "--version", "--no-dump", "--output", "./opt/dir", "--verbose"},
	}

	for flag, options := range optionsList {
		_, err := setup(t, options)
		assertError(t, err)
		if e, ok := err.(ErrorUnknownFlag); !ok {
			t.Errorf("expected an ErrUnknownFlag; got %t", err)
		} else {
			assertEqual(t, flag, e.Name)
		}
	}
}

// test `--flag=value` syntax
func TestFlagAssignmentSyntax(t *testing.T) {
	// options list
	optionsList := [][]string{
		{"info", "student", "-v", "--version=2.0.0", "thatisuday"},
		{"info", "student", "thatisuday", "-v", "-V=2.0.0"},
	}

	for _, options := range optionsList {
		cmd, err := setup(t, options)
		assertNoError(t, err)
		assertEqual(t, options[0], cmd.Name)
		assertEqual(t, options[1], cmd.Args["category"].Value)
		assertEqual(t, "thatisuday", cmd.Args["username"].Value)
		assertEqual(t, "", cmd.Args["subjects"].Value)
		assertEqual(t, "2.0.0", cmd.Flags["version"].Value)
		assertEqual(t, "", cmd.Flags["output"].Value)
		assertEqual(t, "true", cmd.Flags["verbose"].Value)
	}
}

// test for valid variadic argument values
func TestValidVariadicArgumentValues(t *testing.T) {

	// options list
	optionsList := [][]string{
		{"info", "student", "thatisuday", "-V", "-v", "--output", "./opt/dir", "--no-clean", "math", "science", "physics"},
		{"info", "student", "--version", "--no-clean", "thatisuday", "--output", "./opt/dir", "math", "science", "--verbose", "physics"},
	}

	for _, options := range optionsList {
		cmd, err := setup(t, options)
		assertNoError(t, err)
		assertEqual(t, options[0], cmd.Name)
		assertEqual(t, options[1], cmd.Args["category"].Value)
		assertEqual(t, "thatisuday", cmd.Args["username"].Value)
		assertEqual(t, "math,science,physics", cmd.Args["subjects"].Value)
		assertEqual(t, "", cmd.Flags["version"].Value)
		assertEqual(t, "./opt/dir", cmd.Flags["output"].Value)
		assertEqual(t, "true", cmd.Flags["verbose"].Value)
		assertEqual(t, "false", cmd.Flags["clean"].Value)
		assertEqual(t, true, cmd.Flags["clean"].IsInverted)
	}
}

/*-------------------*/

// test root command with options
func TestRootCommandWithOptions(t *testing.T) {

	// options list
	optionsList := [][]string{
		{"userinfo", "-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir"},
		{"-V", "1.0.1", "--verbose", "--force", "userinfo", "--dir", "./sub/dir"},
		{"-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir", "userinfo"},
		{"--version", "1.0.1", "--verbose", "--force", "--dir", "./sub/dir", "userinfo"},
	}

	for _, options := range optionsList {
		cmd, err := setup(t, options)
		assertNoError(t, err)
		assertEqual(t, "", cmd.Name)
		assertEqual(t, "userinfo", cmd.Args["output"].Value)
		assertEqual(t, "true", cmd.Flags["force"].Value)
		assertEqual(t, false, cmd.Flags["force"].IsInverted)
		assertEqual(t, "true", cmd.Flags["verbose"].Value)
		assertEqual(t, "1.0.1", cmd.Flags["version"].Value)
		assertEqual(t, "./sub/dir", cmd.Flags["dir"].Value)
	}
}

// test sub-command with options
func TestSubCommandWithOptions(t *testing.T) {
	// options list
	optionsList := [][]string{
		{"info", "student", "-V", "-v", "--output", "./opt/dir"},
		{"info", "student", "--version", "--output", "./opt/dir", "--verbose"},
	}

	for _, options := range optionsList {
		cmd, err := setup(t, options)
		assertNoError(t, err)
		assertEqual(t, "info", cmd.Name)
		assertEqual(t, "student", cmd.Args["category"].Value)
		assertEqual(t, "", cmd.Args["username"].Value)
		assertEqual(t, "", cmd.Args["subjects"].Value)
		assertEqual(t, "", cmd.Flags["version"].Value)
		assertEqual(t, "./opt/dir", cmd.Flags["output"].Value)
		assertEqual(t, "true", cmd.Flags["verbose"].Value)
		assertEqual(t, "", cmd.Flags["clean"].Value)
	}
}

// test sub-command with valid and extra arguments
func TestSubCommandWithArguments(t *testing.T) {
	// options list
	optionsList := [][]string{
		{"info", "-v", "student", "-V", "2.0.0", "thatisuday"},
		{"info", "student", "-v", "thatisuday", "--version", "2.0.0"},
	}

	for _, options := range optionsList {
		cmd, err := setup(t, options)
		assertNoError(t, err)
		assertEqual(t, "info", cmd.Name)
		assertEqual(t, "student", cmd.Args["category"].Value)
		assertEqual(t, "thatisuday", cmd.Args["username"].Value)
		assertEqual(t, "", cmd.Args["subjects"].Value)
		assertEqual(t, "2.0.0", cmd.Flags["version"].Value)
		assertEqual(t, "", cmd.Flags["output"].Value)
		assertEqual(t, "true", cmd.Flags["verbose"].Value)
	}
}
