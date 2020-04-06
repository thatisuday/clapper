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
		cmd := exec.Command("go", append([]string{"run", "tests/valid-registry.go"}, options...)...)

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
	cmd := exec.Command("go", "run", "tests/valid-registry.go")

	// get output
	if output, err := cmd.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		lines := []string{
			`sub-command => ""`,
			`argument-value => &clapper.Arg{Name:"output", DefaultValue:"", Value:""}`,
			`flag-value => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, DefaultValue:"false", Value:""}`,
			`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:""}`,
			`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"", Value:""}`,
			`flag-value => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, DefaultValue:"/var/users", Value:""}`,
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
	cmd := exec.Command("go", "run", "tests/valid-registry.go")
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
		cmd := exec.Command("go", append([]string{"run", "tests/valid-registry.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			lines := []string{
				`sub-command => ""`,
				`argument-value => &clapper.Arg{Name:"output", DefaultValue:"", Value:"userinfo"}`,
				`flag-value => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, DefaultValue:"false", Value:"true"}`,
				`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:"true"}`,
				`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"", Value:"1.0.1"}`,
				`flag-value => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, DefaultValue:"/var/users", Value:"./sub/dir"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fail()
				}
			}
		}
	}
}

// test unregistered flag
func TestUnregisteredFlag(t *testing.T) {

	// flags
	flags := map[[2]string][]string{
		[2]string{"d", "true"}:          []string{"-V", "1.0.1", "-v", "--force", "-d", "./sub/dir"},
		[2]string{"m", "false"}:         []string{"-V", "1.0.1", "-v", "--force", "--m", "./sub/dir"},
		[2]string{"directory", "false"}: []string{"-V", "1.0.1", "-v", "--force", "--directory", "./sub/dir"},
	}

	for key, options := range flags {
		// command
		cmd := exec.Command("go", append([]string{"run", "tests/valid-registry.go"}, options...)...)

		// key parts
		flag := key[0]
		isRequired := key[1]

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf(`error => clapper.ErrorUnknownFlag{Name:"%s", IsShort:%s}`, flag, isRequired)) {
				t.Fail()
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
		cmd := exec.Command("go", append([]string{"run", "tests/valid-registry.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument-value => &clapper.Arg{Name:"category", DefaultValue:"manager", Value:"student"}`,
				`argument-value => &clapper.Arg{Name:"username", DefaultValue:"", Value:""}`,
				`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"1.0.1", Value:""}`,
				`flag-value => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, DefaultValue:"./", Value:"./opt/dir"}`,
				`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:"true"}`,
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
		[]string{"info", "-v", "student", "-V", "2.0.0", "thatisuday", "extra"},
		[]string{"info", "student", "-v", "thatisuday", "--version", "2.0.0", "extra"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "tests/valid-registry.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument-value => &clapper.Arg{Name:"category", DefaultValue:"manager", Value:"student"}`,
				`argument-value => &clapper.Arg{Name:"username", DefaultValue:"", Value:"thatisuday"}`,
				`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"1.0.1", Value:"2.0.0"}`,
				`flag-value => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, DefaultValue:"./", Value:""}`,
				`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:"true"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fail()
				}
			}
		}
	}
}

// test `--flag=value` syntax
func TestFlagAssignment(t *testing.T) {

	// options list
	optionsList := [][]string{
		[]string{"info", "student", "-v", "--version=2.0.0", "thatisuday", "extra"},
		[]string{"info", "student", "thatisuday", "extra", "-v", "-V=2.0.0"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "tests/valid-registry.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument-value => &clapper.Arg{Name:"category", DefaultValue:"manager", Value:"student"}`,
				`argument-value => &clapper.Arg{Name:"username", DefaultValue:"", Value:"thatisuday"}`,
				`flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"1.0.1", Value:"2.0.0"}`,
				`flag-value => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, DefaultValue:"./", Value:""}`,
				`flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:"true"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fail()
				}
			}
		}
	}
}
