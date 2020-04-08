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

	// options
	options := map[string][]string{
		"---version": []string{"---version"},
		"---v":       []string{"---v=1.0.0"},
		"-version":   []string{"-version"},
	}

	for flag, options := range options {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf(`error => clapper.ErrorUnsupportedFlag{Name:"%s"}`, flag)) {
				t.Fail()
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
			`argument-value => &clapper.Arg{Name:"output", IsVariadic:false, DefaultValue:"", Value:""}`,
			`flag-value => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:""}`,
			`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:""}`,
			`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"", Value:""}`,
			`flag-value => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, IsInverted:false, DefaultValue:"/var/users", Value:""}`,
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
		"-d":          []string{"-V", "1.0.1", "-v", "--force", "-d", "./sub/dir"},
		"--m":         []string{"-V", "1.0.1", "-v", "--force", "--m", "./sub/dir"},
		"--directory": []string{"-V", "1.0.1", "-v", "--force", "--directory", "./sub/dir"},
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
		[]string{"info", "student", "-V", "-v", "--output", "./opt/dir", "--no-clean"},
		[]string{"info", "student", "--version", "--no-clean", "--output", "./opt/dir", "--verbose"},
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
				`argument-value => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}`,
				`argument-value => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:""}`,
				`argument-value => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:""}`,
				`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:""}`,
				`flag-value => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:"./opt/dir"}`,
				`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
				`flag-value => &clapper.Flag{Name:"clean", ShortName:"", IsBoolean:true, IsInverted:true, DefaultValue:"true", Value:"false"}`,
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
		"--clean":   []string{"info", "student", "-V", "-v", "--output", "./opt/dir", "--clean"},
		"--no-dump": []string{"info", "student", "--version", "--no-dump", "--output", "./opt/dir", "--verbose"},
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
		[]string{"info", "student", "-v", "--version=2.0.0", "thatisuday"},
		[]string{"info", "student", "thatisuday", "-v", "-V=2.0.0"},
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
				`argument-value => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}`,
				`argument-value => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:"thatisuday"}`,
				`argument-value => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:""}`,
				`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:"2.0.0"}`,
				`flag-value => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:""}`,
				`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fail()
				}
			}
		}
	}
}

// test for valid variadic argument values
func TestValidVariadicArgumentValues(t *testing.T) {

	// options list
	optionsList := [][]string{
		[]string{"info", "student", "thatisuday", "-V", "-v", "--output", "./opt/dir", "--no-clean", "math", "science", "physics"},
		[]string{"info", "student", "--version", "--no-clean", "thatisuday", "--output", "./opt/dir", "math", "science", "--verbose", "physics"},
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
				`argument-value => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}`,
				`argument-value => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:"thatisuday"}`,
				`argument-value => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:"math,science,physics"}`,
				`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:""}`,
				`flag-value => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:"./opt/dir"}`,
				`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
				`flag-value => &clapper.Flag{Name:"clean", ShortName:"", IsBoolean:true, IsInverted:true, DefaultValue:"true", Value:"false"}`,
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
		[]string{"userinfo", "-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir"},
		[]string{"-V", "1.0.1", "--verbose", "--force", "userinfo", "--dir", "./sub/dir"},
		[]string{"-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir", "userinfo"},
		[]string{"--version", "1.0.1", "--verbose", "--force", "--dir", "./sub/dir", "userinfo"},
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
				`argument-value => &clapper.Arg{Name:"output", IsVariadic:false, DefaultValue:"", Value:"userinfo"}`,
				`flag-value => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
				`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
				`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"", Value:"1.0.1"}`,
				`flag-value => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, IsInverted:false, DefaultValue:"/var/users", Value:"./sub/dir"}`,
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
		[]string{"info", "student", "-V", "-v", "--output", "./opt/dir"},
		[]string{"info", "student", "--version", "--output", "./opt/dir", "--verbose"},
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
				`argument-value => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}`,
				`argument-value => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:""}`,
				`argument-value => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:""}`,
				`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:""}`,
				`flag-value => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:"./opt/dir"}`,
				`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
				`flag-value => &clapper.Flag{Name:"clean", ShortName:"", IsBoolean:true, IsInverted:true, DefaultValue:"true", Value:""}`,
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
		[]string{"info", "-v", "student", "-V", "2.0.0", "thatisuday"},
		[]string{"info", "student", "-v", "thatisuday", "--version", "2.0.0"},
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
				`argument-value => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}`,
				`argument-value => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:"thatisuday"}`,
				`argument-value => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:""}`,
				`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:"2.0.0"}`,
				`flag-value => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:""}`,
				`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fail()
				}
			}
		}
	}
}
