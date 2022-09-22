package clapper

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

/*----------------*/

// test unsupported flag
func TestUnsupportedAssignment(t *testing.T) {

	tests := []struct {
		flag    string
		options string
		err     error
	}{
		{"---version", "---version", ErrorUnsupportedFlag{}},
		{"---v", "---v=1.0.0", ErrorUnsupportedFlag{}},
		// Single-dash for long names was never supported, but now this is interpreted as:
		// -v -e -r -s -i -o -n; ergo, the error will be "unknown option '-e'"
		{"-e", "-version", ErrorUnknownFlag{}},
	}

	for _, test := range tests {
		reg := NewRegistry()
		root, _ := reg.Register("")
		root.AddFlag("version", "v", true, "")

		_, err := reg.Parse([]string{test.options})
		if err == nil {
			t.Errorf("(%s) expected parse error %T", test.options, test.err)
		} else {
			if fmt.Sprintf("%T", err) != fmt.Sprintf("%T", test.err) {
				t.Errorf("(%s) expected %T error -- got %T", test.options, test.err, err)
			}
			if !strings.Contains(err.Error(), test.flag) {
				t.Errorf("(%s) expected error %q -- got %q", test.options, test.flag, err)
			}
		}
	}
}

// test empty root command
func TestEmptyRootCommand(t *testing.T) {
	// command
	cmd := exec.Command("go", "run", "demo/cmd.go")

	// get output
	if output, err := cmd.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		lines := []string{
			`sub-command => ""`,
			`argument(output) => &clapper.Arg{Name:"output", IsVariadic:false, DefaultValue:"", Value:""}`,
			`flag(force) => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:""}`,
			`flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:""}`,
			`flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"", Value:""}`,
			`flag(dir) => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, IsInverted:false, DefaultValue:"/var/users", Value:""}`,
		}

		for _, line := range lines {
			if !strings.Contains(fmt.Sprintf("%s", output), line) {
				t.Fail()
			}
		}
	}
}

// test root command when not registered
func TestUnregisteredRootCommand(t *testing.T) {
	// command
	cmd := exec.Command("go", "run", "demo/cmd.go")
	cmd.Env = append(os.Environ(), "NO_ROOT=TRUE")

	// get output
	if output, err := cmd.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		lines := []string{
			`error => clapper.ErrorUnknownCommand{Name:""}`,
		}

		for _, line := range lines {
			if !strings.Contains(fmt.Sprintf("%s", output), line) {
				t.Fail()
			}
		}
	}
}

// test an unregistered flag
func TestUnregisteredFlag(t *testing.T) {

	// flags
	flags := map[string][]string{
		"-d":          {"-V", "1.0.1", "-v", "--force", "-d", "./sub/dir"},
		"--m":         {"-V", "1.0.1", "-v", "--force", "--m", "./sub/dir"},
		"--directory": {"-V", "1.0.1", "-v", "--force", "--directory", "./sub/dir"},
	}

	for flag, options := range flags {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf(`error => clapper.ErrorUnknownFlag{Name:"%s"}`, flag)) {
				t.Fail()
			}
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

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}`,
				`argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:""}`,
				`argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:""}`,
				`flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:""}`,
				`flag(output) => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:"./opt/dir"}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
				`flag(clean) => &clapper.Flag{Name:"clean", ShortName:"", IsBoolean:true, IsInverted:true, DefaultValue:"true", Value:"false"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fail()
				}
			}
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
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf(`error => clapper.ErrorUnknownFlag{Name:"%s"}`, flag)) {
				t.Fail()
			}
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
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}`,
				`argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:"thatisuday"}`,
				`argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:""}`,
				`flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:"2.0.0"}`,
				`flag(output) => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:""}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Errorf("expected\n  %s\ngot\n  %s", line, output)
				}
			}
		}
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
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}`,
				`argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:"thatisuday"}`,
				`argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:"math,science,physics"}`,
				`flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:""}`,
				`flag(output) => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:"./opt/dir"}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
				`flag(clean) => &clapper.Flag{Name:"clean", ShortName:"", IsBoolean:true, IsInverted:true, DefaultValue:"true", Value:"false"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fail()
				}
			}
		}
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
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			lines := []string{
				`sub-command => ""`,
				`argument(output) => &clapper.Arg{Name:"output", IsVariadic:false, DefaultValue:"", Value:"userinfo"}`,
				`flag(force) => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
				`flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"", Value:"1.0.1"}`,
				`flag(dir) => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, IsInverted:false, DefaultValue:"/var/users", Value:"./sub/dir"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fail()
				}
			}
		}
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
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}`,
				`argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:""}`,
				`argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:""}`,
				`flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:""}`,
				`flag(output) => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:"./opt/dir"}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
				`flag(clean) => &clapper.Flag{Name:"clean", ShortName:"", IsBoolean:true, IsInverted:true, DefaultValue:"true", Value:""}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fail()
				}
			}
		}
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
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}`,
				`argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:"thatisuday"}`,
				`argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:""}`,
				`flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:"2.0.0"}`,
				`flag(output) => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:""}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fail()
				}
			}
		}
	}
}

// test sub-command with valid and extra arguments
func TestCombinedFlags(t *testing.T) {

	// tests
	tests := []struct {
		Name           string
		Args           []string
		ExpectedNames  []string
		ExpectedValues []string
	}{
		{"one flag", []string{"-a"}, []string{"alpha"}, []string{"true"}},
		{"two flags", []string{"-ab"}, []string{"alpha", "bravo"}, []string{"true", "true"}},
		{"one flag & default", []string{"-a"}, []string{"alpha", "bravo"}, []string{"true", ""}},
		{"two flags & var", []string{"-abc", "value"}, []string{"alpha", "bravo", "charlie"}, []string{"true", "true", "value"}},
		// This is weird, but it was not an error to have non-boolean flags that are missing values, and to preserve
		// backwards compatability this was not changed.
		{"unset flag", []string{"-acb", "value"}, []string{"alpha", "bravo", "charlie"}, []string{"true", "true", ""}},
	}

	for _, test := range tests {
		reg := NewRegistry()
		root, _ := reg.Register("")
		root.AddFlag("alpha", "a", true, "")
		root.AddFlag("bravo", "b", true, "")
		root.AddFlag("charlie", "c", false, "none")

		cmd, err := reg.Parse(test.Args)
		if err != nil {
			fmt.Printf("parse error %s\n", err)
			if len(test.ExpectedNames) > 0 {
				t.Errorf("(%s) unexpected parse error: %s", test.Name, err)
			}
			// else, expected error
			continue
		}
		if len(test.ExpectedNames) == 0 && err == nil {
			t.Errorf("(%s) expected parse error; didn't get one", test.Name)
		}
		// Check that all expected arguments are there
		for i, n := range test.ExpectedNames {
			var found bool
			for _, a := range cmd.Flags {
				if a.Name == n {
					if a.Value != test.ExpectedValues[i] {
						t.Errorf("(%s) expected value %s for %s, got %s", test.Name, test.ExpectedValues[i], n, a.Value)
					}
					found = true
				}
			}
			if !found {
				t.Errorf("(%s) did not find expected argument %s", test.Name, n)
			}
		}
	}
}
